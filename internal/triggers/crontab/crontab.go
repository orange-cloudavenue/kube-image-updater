package crontab

import (
	"context"

	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	log "github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
)

var sched = quartz.NewStdScheduler()

func New(ctx context.Context) {
	sched.Start(ctx)
}

type CronTrigger struct{}

func AddCronTab(namespace, name, crontab string) error {
	log.Infof("Registering crontab (%s) for %s in namespace %s", crontab, name, namespace)
	cronTrigger, _ := quartz.NewCronTrigger(crontab)
	functionJob := job.NewFunctionJob(func(_ context.Context) (string, error) {
		log.Infof("Fire crontab refresh for %s in namespace %s", name, namespace)
		triggers.Trigger(triggers.RefreshImage, namespace, name)
		return "", nil
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
