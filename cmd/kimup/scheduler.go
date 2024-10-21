package main

import (
	"context"
	"sync"
	"time"

	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/util/retry"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/actions"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/registry"
	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers/crontab"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

type locks map[string]*sync.RWMutex

func initScheduler(ctx context.Context, k kubeclient.Interface) {
	l := make(locks)

	// Start Crontab client
	crontab.New(ctx)
	// Add event lock
	event.On(triggers.RefreshImage.String(), event.ListenerFunc(func(e event.Event) (err error) {
		// Increment the counter for the events
		metrics.Events().TriggeredTotal.Inc()
		// Start the timer for the event execution
		timerEvents := metrics.Events().TriggeredDuration.NewTimer()
		defer timerEvents.ObserveDuration()

		if l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)] == nil {
			l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)] = &sync.RWMutex{}
		}

		// Lock the image to prevent concurrent refreshes
		l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)].Lock()
		defer l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)].Unlock()

		// Sleep for 1 second to prevent concurrent refreshes
		time.Sleep(1 * time.Second)

		retryErr := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			log.Infof("Refreshing image %s in namespace %s", e.Data()["image"], e.Data()["namespace"])

			ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			image, err := k.Image().Get(ctx, e.Data()["namespace"].(string), e.Data()["image"].(string))
			defer func() {
				// update the status of the image
				image.SetStatusTime(time.Now().Format(time.RFC3339))
				err := k.Image().UpdateStatus(ctx, image)
				if err != nil {
					log.WithError(err).
						WithFields(log.Fields{
							"Namespace": e.Data()["namespace"],
							"Image":     e.Data()["image"],
						}).Error("Error updating status of image")
				}
			}()
			if err != nil {
				image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorGetImage)
				if err := crontab.RemoveJob(crontab.BuildKey(e.Data()["namespace"].(string), e.Data()["image"].(string))); err != nil {
					return err
				}
				return err
			}

			var auths kubeclient.K8sDockerRegistrySecretData

			if image.Spec.ImagePullSecrets != nil {
				auths, err = k.GetPullSecretsForImage(ctx, image)
				image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorPullSecrets)
				if err != nil {
					return err
				}
			}

			i := utils.ImageParser(image.Spec.Image)

			// Prometheus metrics - Increment the counter for the registry
			metrics.Registry().RequestTotal.WithLabelValues(i.GetRegistry()).Inc()
			timerRegistry := metrics.Registry().RequestDuration.NewTimer(i.GetRegistry())

			re, err := registry.New(ctx, image.Spec.Image, registry.Settings{
				InsecureTLS: image.Spec.InsecureSkipTLSVerify,
				Username: func() string {
					if v, ok := auths.Auths[i.GetRegistry()]; ok {
						return v.Username
					}
					return ""
				}(),
				Password: func() string {
					if v, ok := auths.Auths[i.GetRegistry()]; ok {
						return v.Password
					}
					return ""
				}(),
			})
			timerRegistry.ObserveDuration()
			if err != nil {
				// Prometheus metrics - Increment the counter for the registry with error
				metrics.Registry().RequestErrorTotal.WithLabelValues(i.GetRegistry()).Inc()
				image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorRegistry)
				return err
			}

			// Prometheus metrics - Increment the counter for the tags
			metrics.Tags().RequestTotal.Inc()
			timerTags := metrics.Tags().RequestDuration.NewTimer()

			tagsAvailable, err := re.Tags()
			timerTags.ObserveDuration()
			if err != nil {
				// Prometheus metrics - Increment the counter for the tags with error
				metrics.Tags().RequestErrorTotal.Inc()
				return err
			}

			metrics.Tags().AvailableSum.WithLabelValues(image.Spec.Image).Observe(float64(len(tagsAvailable)))

			log.Debugf("[RefreshImage] %d tags available for %s", len(tagsAvailable), image.Spec.Image)

			for _, rule := range image.Spec.Rules {
				r, err := rules.GetRule(rule.Type)
				if err != nil {
					image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorGetRule)
					log.Errorf("Error getting rule: %v", err)
					continue
				}

				tag := image.Status.Tag
				if image.Status.Tag == "" {
					tag = image.Spec.BaseTag
				}

				r.Init(tag, tagsAvailable, rule.Value)

				// Prometheus metrics - Increment the counter for the rules
				metrics.Rules().EvaluatedTotal.Inc()
				timerRules := metrics.Rules().EvaluatedDuration.NewTimer()

				match, newTag, err := r.Evaluate()

				// Prometheus metrics - Observe the duration of the rule evaluation
				timerRules.ObserveDuration()

				if err != nil {
					// Prometheus metrics - Increment the counter for the evaluated rule with error
					metrics.Rules().EvaluatedErrorTotal.Inc()

					log.Errorf("Error evaluating rule: %v", err)
					continue
				}

				if match {
					for _, action := range rule.Actions {
						a, err := actions.GetActionWithUntypedName(action.Type)
						image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorAction)
						if err != nil {
							log.Errorf("Error getting action: %v", err)
							continue
						}

						a.Init(k, models.Tags{
							Actual:        tag,
							New:           newTag,
							AvailableTags: tagsAvailable,
						}, &image, action.Data)

						// Prometheus metrics - Increment the counter for the actions
						metrics.Actions().ExecutedTotal.Inc()
						timerActions := metrics.Actions().ExecutedDuration.NewTimer()

						err = a.Execute(ctx)

						// Prometheus metrics - Observe the duration of the action execution
						timerActions.ObserveDuration()

						if err != nil {
							// Prometheus metrics - Increment the counter for the executed action with error
							metrics.Actions().ExecutedErrorTotal.Inc()

							log.Errorf("Error executing action(%s): %v", action.Type, err)
							continue
						}
					}
					log.Debugf("[RefreshImage] Rule %s evaluated: %v -> %s", rule.Type, tag, newTag)
				}
			}

			image.SetStatusResult(v1alpha1.ImageStatusLastSyncSuccess)
			return k.Image().Update(ctx, image)
		})

		// Prometheus metrics - Increment the counter for the events evaluated with error
		metrics.Events().TriggerdErrorTotal.Inc()
		return retryErr
	}), event.Normal)
}
