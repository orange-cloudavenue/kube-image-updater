package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/httpserver"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
)

var c = make(chan os.Signal, 1)

func init() {
	// Initialize the metrics
	metrics.Tags()
	metrics.Events()
	metrics.Actions()
	metrics.Rules()
	metrics.Registry()

	// TODO add namespace scope
	// Flag "loglevel" is set in log package
	flag.Parse()
}

func main() {
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	log.WithField("version", models.Version).Info("Starting kimup", models.Version)

	// kubernetes golang library provide flag "kubeconfig" to specify the path to the kubeconfig file
	k, err := kubeclient.New(flag.Lookup("kubeconfig").Value.String(), kubeclient.ComponentController)
	if err != nil {
		log.WithError(err).Panic("Error creating kubeclient")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// * Config the metrics and healthz server
	a, waitHTTP := httpserver.Init(ctx, httpserver.WithCustomHandlerForHealth(
		func() (bool, error) {
			// TODO improve
			_, err := net.DialTimeout("tcp", models.HealthzDefaultAddr, 5*time.Second)
			if err != nil {
				return false, err
			}
			return true, nil
		}))

	if err := a.Run(); err != nil {
		log.WithError(err).Error("Failed to start HTTP servers")
		// send signal to stop the program
		c <- syscall.SIGINT
	}

	initScheduler(ctx, k)

	go func() {
		x, err := k.Image().Watch(ctx)
		if err != nil {
			log.WithError(err).Panic("Error watching events")
		}

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-x:
				if !ok {
					return
				}

				ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				an := annotations.New(ctx, &event.Value)

				switch event.Type {
				case "ADDED":
					// Clean old action
					an.Remove(annotations.KeyAction)

					setupTriggers(&event.Value)
					refreshIfRequired(an, event.Value)
					if err := setTagIfNotExists(ctx, k, an, &event.Value); err != nil {
						log.WithError(err).Error("Error setting tag")
					}

					if err := k.Image().Update(ctx, event.Value); err != nil {
						log.WithError(err).Error("Error updating image")
					}

				case "MODIFIED":
					switch an.Action().Get() { //nolint:gocritic
					case annotations.ActionReload:

						// * Here is only if the yaml has been updated and the operator has detected it
						for _, trigger := range event.Value.Spec.Triggers {
							switch trigger.Type {
							case triggers.Crontab:
								cleanTriggers(&event.Value)
							case triggers.Webhook:
								log.Warn("Webhook trigger not implemented yet")
							}
						}

						refresh(event.Value)

						// Remove the annotation annotations.AnnotationActionKey in the map[string]string
						an.Remove(annotations.KeyAction)
					}

					refreshIfRequired(an, event.Value)

					if err := k.Image().Update(ctx, event.Value); err != nil {
						log.WithError(err).Error("Error updating image")
					}

					setupTriggers(&event.Value)

				case "DELETED":
					cleanTriggers(&event.Value)
				}
			}
		}
	}()

	<-c
	cancel()
	waitHTTP()
}
