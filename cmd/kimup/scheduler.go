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

		var (
			namespaceName = e.Data()["namespace"].(string)
			imageName     = e.Data()["image"].(string)
		)

		if l[namespaceName+"/"+imageName] == nil {
			l[namespaceName+"/"+imageName] = &sync.RWMutex{}
		}

		// Lock the image to prevent concurrent refreshes
		l[namespaceName+"/"+imageName].Lock()
		defer l[namespaceName+"/"+imageName].Unlock()

		// Sleep for 1 second to prevent concurrent refreshes
		time.Sleep(1 * time.Second)

		retryErr := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			log.Infof("Refreshing image %s in namespace %s", imageName, namespaceName)

			ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			image, err := k.Image().Get(ctx, namespaceName, imageName)
			defer func() {
				// Add delay to avoid conflicts
				time.Sleep(500 * time.Millisecond)

				// update the status of the image
				image.SetStatusTime(time.Now().Format(time.RFC3339))

				// Need to get image again to avoid conflicts
				imageRefreshed, err := k.Image().Get(ctx, namespaceName, imageName)
				if err != nil {
					log.WithError(err).
						WithFields(log.Fields{
							"Namespace": namespaceName,
							"Image":     imageName,
						}).Error("Error getting image")
					return
				}
				imageRefreshed.Status = image.Status

				if err := k.Image().UpdateStatus(ctx, imageRefreshed); err != nil {
					log.WithError(err).
						WithFields(log.Fields{
							"Namespace": e.Data()["namespace"],
							"Image":     e.Data()["image"],
						}).Error("Error updating status of image")
				}
			}()
			if err != nil {
				image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorGetImage)
				if err := crontab.RemoveJob(crontab.BuildKey(namespaceName, imageName)); err != nil {
					return err
				}
				return err
			}
			k.Image().Event(&image, corev1.EventTypeNormal, "Image update triggered", "")

			// Set Status to Scheduled permit in the execution of the refresh if the image have a error or not
			image.SetStatusResult(v1alpha1.ImageStatusLastSyncScheduled)

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
				metrics.Registry().RequestErrorTotal.WithLabelValues(i.GetRegistry()).Inc()
				image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorRegistry)
				k.Image().Event(&image, corev1.EventTypeWarning, "Fetch image", fmt.Sprintf("Error fetching image: %v", err))
				log.WithError(err).Error("Error fetching image")
				return err
			}

			// Prometheus metrics - Increment the counter for the tags
			metrics.Tags().RequestTotal.Inc()
			timerTags := metrics.Tags().RequestDuration.NewTimer()

			tagsAvailable, err := re.Tags()
			timerTags.ObserveDuration()
			if err != nil {
				metrics.Tags().RequestErrorTotal.Inc()
				image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorTags)
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
					image.SetStatusResult(v1alpha1.ImageStatusLastSyncError)
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
							image.SetStatusResult(v1alpha1.ImageStatusLastSyncErrorAction)
							k.Image().Event(&image, corev1.EventTypeWarning, "Get action", fmt.Sprintf("Error getting action %s: %v", action.Type, err))
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
							image.SetStatusResult(v1alpha1.ImageStatusLastSyncError)
							log.Errorf("Error executing action(%s): %v", action.Type, err)
							k.Image().Event(&image, corev1.EventTypeWarning, "Execute action", fmt.Sprintf("Error executing action %s: %v", action.Type, err))
							continue
						}
						k.Image().Event(&image, corev1.EventTypeNormal, "Execute action", fmt.Sprintf("Action %s executed", action.Type))
					}
					log.Debugf("[RefreshImage] Rule %s evaluated: %v -> %s", rule.Type, tag, newTag)
				}
			}

			if image.Status.Result == v1alpha1.ImageStatusLastSyncScheduled {
				image.SetStatusResult(v1alpha1.ImageStatusLastSyncSuccess)
			}
			return k.Image().Update(ctx, image)
		})

		// Prometheus metrics - Increment the counter for the events evaluated with error
		metrics.Events().TriggerdErrorTotal.Inc()
		return retryErr
	}), event.Normal)
}
