package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"

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
			if err != nil {
				if err := crontab.RemoveJob(crontab.BuildKey(e.Data()["namespace"].(string), e.Data()["image"].(string))); err != nil {
					return err
				}
				return err
			}
			k.Image().Event(&image, corev1.EventTypeNormal, "Image update triggered", "")

			var auths kubeclient.K8sDockerRegistrySecretData

			if image.Spec.ImagePullSecrets != nil {
				auths, err = k.GetPullSecretsForImage(ctx, image)
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
				k.Image().Event(&image, corev1.EventTypeWarning, "Fetch image tags", fmt.Sprintf("Error fetching tags: %v", err))
				log.WithError(err).Error("Error fetching tags")
				return err
			}

			metrics.Tags().AvailableSum.WithLabelValues(image.Spec.Image).Observe(float64(len(tagsAvailable)))
			k.Image().Event(&image, corev1.EventTypeNormal, "Fetch image tags", fmt.Sprintf("Found %d tags", len(tagsAvailable)))

			log.Debugf("[RefreshImage] %d tags available for %s", len(tagsAvailable), image.Spec.Image)

			for _, rule := range image.Spec.Rules {
				r, err := rules.GetRule(rule.Type)
				if err != nil {
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
					k.Image().Event(&image, corev1.EventTypeWarning, "Evaluate rule", fmt.Sprintf("Error evaluating rule %s: %v", rule.Type, err))
					continue
				}

				k.Image().Event(&image, corev1.EventTypeNormal, "Evaluate rule", fmt.Sprintf("Rule %s evaluated", rule.Type))

				if match {
					for _, action := range rule.Actions {
						a, err := actions.GetActionWithUntypedName(action.Type)
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
							k.Image().Event(&image, corev1.EventTypeWarning, "Execute action", fmt.Sprintf("Error executing action %s: %v", action.Type, err))
							continue
						}
						k.Image().Event(&image, corev1.EventTypeNormal, "Execute action", fmt.Sprintf("Action %s executed", action.Type))
					}
					log.Debugf("[RefreshImage] Rule %s evaluated: %v -> %s", rule.Type, tag, newTag)
				}
			}

			return k.Image().Update(ctx, image)
		})

		// Prometheus metrics - Increment the counter for the events evaluated with error
		metrics.Events().TriggerdErrorTotal.Inc()
		return retryErr
	}), event.Normal)
}
