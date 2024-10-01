package main

import (
	"context"
	"sync"
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
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

type locks map[string]*sync.RWMutex

func initScheduler(ctx context.Context, k *kubeclient.Client) {
	l := make(locks)

	// Start Crontab client
	crontab.New(ctx)
	event.On(triggers.RefreshImage.String(), event.ListenerFunc(func(e event.Event) (err error) {
		if l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)] == nil {
			l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)] = &sync.RWMutex{}
		}

		// Lock the image to prevent concurrent refreshes
		l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)].Lock()
		defer l[e.Data()["namespace"].(string)+"/"+e.Data()["image"].(string)].Unlock()

		retryErr := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			// TODO: implement image refresh
			log.Infof("Refreshing image %s in namespace %s", e.Data()["image"], e.Data()["namespace"])

			ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
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

			var auths kubeclient.K8sDockerRegistrySecretData

			if image.Spec.ImagePullSecrets != nil {
				auths, err = k.GetPullSecretsForImage(ctx, image)
				if err != nil {
					return err
				}
			}

			i := utils.ImageParser(image.Spec.Image)

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
