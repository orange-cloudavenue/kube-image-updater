package main

import (
	"context"
	"time"

	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/internal/actions"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers/crontab"
)

func initScheduler(k *kubeclient.Client) {
	// Start Crontab client
	crontab.New(context.Background())

	event.On(triggers.RefreshImage.String(), event.ListenerFunc(func(e event.Event) error {
		// TODO: implement image refresh
		log.Infof("Refreshing image %s in namespace %s", e.Data()["name"], e.Data()["namespace"])

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		image, err := k.GetImage(ctx, e.Data()["namespace"].(string), e.Data()["name"].(string))
		if err != nil {
			if err := crontab.RemoveJob(crontab.BuildKey(e.Data()["namespace"].(string), e.Data()["nginx"].(string))); err != nil {
				return err
			}
			return err
		}

		an := annotations.New(ctx, &image)
		an.Tag().Set(image.Spec.BaseTag)
		// TODO: Implement logic here

		// Add tags refresh

		for _, rule := range image.Spec.Rules {
			r, err := rules.GetRuleWithUntypedName(string(rule.Type))
			if err != nil {
				log.Errorf("Error getting rule: %v", err)
				continue
			}

			// TODO: missing tags
			r.Init(image.Status.Tag, make([]string, 0), rule.Value)
			match, newTag, err := r.Evaluate()
			if err != nil {
				log.Errorf("Error evaluating rule: %v", err)
				continue
			}

			if match {
				for _, action := range rule.Actions {
					a, err := actions.GetActionWithUntypedName(string(action.Type))
					if err != nil {
						log.Errorf("Error getting action: %v", err)
						continue
					}

					a.Init(image.Status.Tag, newTag, &image)
					if err := a.Execute(ctx); err != nil {
						log.Errorf("Error executing action: %v", err)
						continue
					}
				}
			}
		}

		if err := k.SetImage(ctx, image); err != nil {
			log.Errorf("Error updating image: %v", err)
		}

		return nil
	}), event.Normal)
}
