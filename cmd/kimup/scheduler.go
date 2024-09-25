package main

import (
	"context"
	"time"

	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/util/retry"

	"github.com/orange-cloudavenue/kube-image-updater/internal/actions"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/registry"
	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers/crontab"
)

func initScheduler(k *kubeclient.Client) {
	// Start Crontab client
	crontab.New(context.Background())

	event.On(triggers.RefreshImage.String(), event.ListenerFunc(func(e event.Event) error {
		retryErr := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			// TODO: implement image refresh
			log.Infof("Refreshing image %s in namespace %s", e.Data()["image"], e.Data()["namespace"])

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			image, err := k.GetImage(ctx, e.Data()["namespace"].(string), e.Data()["image"].(string))
			if err != nil {
				if err := crontab.RemoveJob(crontab.BuildKey(e.Data()["namespace"].(string), e.Data()["image"].(string))); err != nil {
					return err
				}
				return err
			}

			// an := annotations.New(ctx, &image)
			// TODO add last refresh annotation

			re, err := registry.New(ctx, image.Spec.Image)
			if err != nil {
				return err
			}

			tagsAvailable, err := re.Tags()
			if err != nil {
				return err
			}

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
				match, newTag, err := r.Evaluate()
				if err != nil {
					log.Errorf("Error evaluating rule: %v", err)
					continue
				}

				if match {
					for _, action := range rule.Actions {
						a, err := actions.GetAction(action.Type)
						if err != nil {
							log.Errorf("Error getting action: %v", err)
							continue
						}

						a.Init(tag, newTag, &image)
						if err := a.Execute(ctx); err != nil {
							log.Errorf("Error executing action: %v", err)
							continue
						}
					}

					log.Debugf("[RefreshImage] Rule %s evaluated: %v -> %s", rule.Type, tag, newTag)
				}
			}

			if err := k.SetImage(ctx, image); err != nil {
				log.Errorf("Error updating image: %v", err)
			}

			return nil
		})

		return retryErr
	}), event.Normal)
}
