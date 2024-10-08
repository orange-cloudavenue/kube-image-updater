package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/actions"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers/crontab"
)

func setupTriggers(x *v1alpha1.Image) {
	// * Triggers
	for _, trigger := range x.Spec.Triggers {
		switch trigger.Type {
		case triggers.Crontab:
			if ok, err := crontab.IsExistingJob(crontab.BuildKey(x.Namespace, x.Name)); err != nil || !ok {
				if err := crontab.AddCronTab(x.Namespace, x.Name, trigger.Value); err != nil {
					log.Errorf("Error adding cronjob: %v", err)
				}
			}
		case triggers.Webhook:
			log.Info("Webhook trigger not implemented yet")
		}
	}
}

func cleanTriggers(x *v1alpha1.Image) {
	for _, trigger := range x.Spec.Triggers {
		switch trigger.Type {
		case triggers.Crontab:
			if ok, err := crontab.IsExistingJob(crontab.BuildKey(x.Namespace, x.Name)); err != nil || ok {
				if err := crontab.RemoveJob(crontab.BuildKey(x.Namespace, x.Name)); err != nil {
					log.Errorf("Error removing crontab: %v", err)
				}
			}
		case triggers.Webhook:
			log.Info("Webhook trigger not implemented yet")
		}
	}
}

func refreshIfRequired(an annotations.Annotation, image v1alpha1.Image) {
	if an.Action().Get() == annotations.ActionRefresh {
		// * Here is only if the image has annotations.ActionRefresh
		log.Infof("[Fire] Annotation trigger refresh for image %s in namespace %s", image.Name, image.Namespace)
		_, err := triggers.Trigger(triggers.RefreshImage, image.Namespace, image.Name)
		if err != nil {
			log.Errorf("Error triggering event: %v", err)
		}
		an.Remove(annotations.KeyAction)
	}
}

func setTagIfNotExists(ctx context.Context, kubeClient kubeclient.Interface, an annotations.Annotation, image *v1alpha1.Image) error {
	if an.Tag().IsNull() {
		a, err := actions.GetAction(actions.Apply)
		if err != nil {
			return err
		}

		a.Init(kubeClient, models.Tags{
			Actual: image.Status.Tag,
			New:    image.Spec.BaseTag,
		}, image, v1alpha1.ValueOrValueFrom{})

		if err := a.Execute(ctx); err != nil {
			return err
		}
	}

	return nil
}
