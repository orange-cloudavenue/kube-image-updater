package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/actions"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
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
					log.
						WithError(err).
						WithFields(logrus.Fields{
							"crontab":   trigger.Value,
							"namespace": x.Namespace,
							"name":      x.Name,
						}).Error("Error adding cronjob")
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
					log.
						WithError(err).
						WithFields(logrus.Fields{
							"namespace": x.Namespace,
							"name":      x.Name,
						}).Error("Error removing crontab")
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
		log.
			WithFields(logrus.Fields{
				"namespace": image.Namespace,
				"name":      image.Name,
			}).Info("Annotation trigger refresh")
		_, err := triggers.Trigger(triggers.RefreshImage, image.Namespace, image.Name)
		if err != nil {
			log.
				WithFields(logrus.Fields{
					"namespace": image.Namespace,
					"name":      image.Name,
				}).
				Error("Error triggering event")
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
