package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers/crontab"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

var (
	version = "dev" // set by build script

	c = make(chan os.Signal, 1)
)

func init() {
	flag.String("loglevel", "info", "log level (debug, info, warn, error, fatal, panic)")
	flag.Parse()

	log.SetLevel(utils.ParseLogLevel(flag.Lookup("loglevel").Value.String()))
	log.SetFormatter(&log.TextFormatter{})
}

func main() {
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("Starting kimup (version: %s)", version)

	// kubernetes golang library provide flag "kubeconfig" to specify the path to the kubeconfig file
	k, err := kubeclient.New(flag.Lookup("kubeconfig").Value.String())
	if err != nil {
		log.Panicf("Error creating kubeclient: %v", err)
	}

	initScheduler(k)

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			images, err := k.ListAllImages(ctx)
			if err != nil {
				log.Errorf("Error listing images: %v", err)
				continue
			}

			for _, image := range images.Items {
				an := annotations.New(ctx, &image)
				if !an.Action().IsNull() && an.Action().Is(annotations.ActionRefresh) {
					// * Here is only if the yaml has been updated and the operator has detected it

					log.Infof("Image configuration %s in namespace %s has changed", image.Name, image.Namespace)

					for _, trigger := range image.Spec.Triggers {
						switch trigger.Type {
						case v1alpha1.ImageTriggerTypeCrontab:
							if ok, err := crontab.IsExistingJob(crontab.BuildKey(image.Namespace, image.Name)); err != nil || ok {
								if err := crontab.RemoveJob(crontab.BuildKey(image.Namespace, image.Name)); err != nil {
									log.Errorf("Error removing crontab: %v", err)
								}
							}
						case v1alpha1.ImageTriggerTypeWebhook:
							log.Info("Webhook trigger not implemented yet")
						}
					}

					// Remove the annotation annotations.AnnotationActionKey in the map[string]string
					an.Remove(annotations.KeyAction)
				} // * End refresh action

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

			time.Sleep(2 * time.Second)
		}
	}()

	<-c
}
