package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers/crontab"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

var (
	flagLogLevel string

	version = "dev"
)

func init() {
	flag.StringVar(&flagLogLevel, "loglevel", "info", "log level (debug, info, warn, error, fatal, panic)")
	flag.Parse()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	log.SetLevel(utils.ParseLogLevel(flagLogLevel))
	log.SetFormatter(&log.TextFormatter{})

	log.Infof("Starting kimup (version: %s)", version)

	// kubernetes golang library provide flag "kubeconfig" to specify the path to the kubeconfig file
	k, err := kubeclient.New(flag.Lookup("kubeconfig").Value.String())
	if err != nil {
		panic(err)
	}

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

		if err := k.SetImage(ctx, image); err != nil {
			log.Errorf("Error updating image: %v", err)
		}

		return nil
	}), event.Normal)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// Detect new images
			images, err := k.ListAllImages(ctx)
			if err != nil {
				log.Errorf("Error listing images: %v", err)
				continue
			}

			for _, image := range images.Items {
				an := annotations.New(ctx, &image)
				if !an.Action().IsNull() && an.Action().Is(annotations.ActionRefresh) {
					log.Infof("Image configration %s in namespace %s has changed", image.Name, image.Namespace)

					for _, trigger := range image.Spec.Triggers {
						switch trigger.Type {
						case v1alpha1.ImageTriggerTypeCrontab:
							if ok, err := crontab.IsExistingJob(crontab.BuildKey(image.Namespace, image.Name)); err != nil || ok {
								if err := crontab.RemoveJob(crontab.BuildKey(image.Namespace, image.Name)); err != nil {
									log.Errorf("Error removing cronjob: %v", err)
								}
							}
							if err := crontab.AddCronTab(image.Namespace, image.Name, trigger.Value); err != nil {
								log.Errorf("Error adding cronjob: %v", err)
							}
						case v1alpha1.ImageTriggerTypeWebhook:
							log.Info("Webhook trigger not implemented yet")
						}
					}
				}

				// Remove the annotation annotations.AnnotationActionKey in the map[string]string
				an.Remove(annotations.KeyAction)

				if err := k.SetImage(ctx, image); err != nil {
					log.Errorf("Error updating image: %v", err)
				}

				// * Triggers
				for _, trigger := range image.Spec.Triggers {
					switch trigger.Type {
					case v1alpha1.ImageTriggerTypeCrontab:
						if ok, err := crontab.IsExistingJob(crontab.BuildKey(image.Namespace, image.Name)); err != nil || !ok {
							if err := crontab.AddCronTab(image.Namespace, image.Name, trigger.Value); err != nil {
								log.Errorf("Error adding cronjob: %v", err)
							}
						}
					case v1alpha1.ImageTriggerTypeWebhook:
						log.Info("Webhook trigger not implemented yet")
					}
				}
			}
		}
	}()

	<-c
}
