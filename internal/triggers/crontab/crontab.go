package crontab

import (
	"context"

	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	"github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
)

var sched = quartz.NewStdScheduler()

func New(ctx context.Context) {
	sched.Start(ctx)
}

type CronTrigger struct{}

func AddCronTab(namespace, name, crontab string) error {
	log.WithFields(logrus.Fields{
		"crontab":   crontab,
		"namespace": namespace,
		"name":      name,
	}).Info("Registering crontab")

	cronTrigger, _ := quartz.NewCronTrigger(crontab)
	functionJob := job.NewFunctionJob(func(_ context.Context) (string, error) {
		log.WithFields(logrus.Fields{
			"namespace": namespace,
			"name":      name,
		}).Info("Crontab trigger refresh")

		_, err := triggers.Trigger(triggers.RefreshImage, namespace, name)
		return "", err
	})

	return sched.ScheduleJob(
		quartz.NewJobDetail(
			functionJob,
			quartz.NewJobKey(BuildKey(namespace, name)),
		), cronTrigger)
}

func RemoveJob(name string) error {
	jobs, err := sched.GetJobKeys()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if job.Name() == name {
			return sched.DeleteJob(job)
		}
	}

	return nil
}

func BuildKey(namespace, name string) string {
	return namespace + "-" + name
}

func IsExistingJob(name string) (bool, error) {
	jobs, err := sched.GetJobKeys()
	if err != nil {
		return false, err
	}
	for _, job := range jobs {
		if job.Name() == name {
			return true, nil
		}
	}

	return false, nil
}
