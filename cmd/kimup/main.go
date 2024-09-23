package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
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

	// kubernetes golang libraru provide flag "kubeconfig" to specify the path to the kubeconfig file
	k, err := kubeclient.New(flag.Lookup("kubeconfig").Value.String())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			images, err := k.ListAllImages(ctx)
			if err != nil {
				log.Errorf("Error listing images: %v", err)
				continue
			}

			for _, image := range images.Items {
				action, err := image.GetAnnotationAction()
				if err != nil {
					continue
				}
				if action == annotations.ActionRefresh {
					log.Infof("Image %s needs to be refreshed", image.Name)
				}

				// Remove the annotation annotations.AnnotationActionKey in the map[string]string
				delete(image.Annotations, annotations.AnnotationActionKey)

				if err := k.SetImage(ctx, image); err != nil {
					log.Errorf("Error updating image: %v", err)
				}
			}
		}
	}()

	// images, err := k.ListAllImages(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, image := range images.Items {
	// 	log.Infof("Image: %s", image.Name)
	// }

	<-c
}
