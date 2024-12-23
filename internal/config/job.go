package config

import (
	"context"
	"simple-crud-clean-architecture/internal/usecase"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type job struct {
	Log          *logrus.Logger
	JobMethod    *gocron.Scheduler
	CronTime     time.Time
	CronInterval int
}

type Job interface {
	RunCron(discountUseCase *usecase.DiscountUseCase)
}

func NewCronJob(viper *viper.Viper, log *logrus.Logger) Job {
	cronTimeStr := viper.GetString("cron.schedule_time")
	cronInterval := viper.GetInt("cron.interval_minutes")

	if cronTimeStr == "" {
		cronTimeStr = "00:00"
	}

	if cronInterval == 0 {
		cronInterval = 10
	}

	log.Infof("Current cron time is" + cronTimeStr + "and current interval is")
	parsedCronTime, err := time.Parse("15:04", cronTimeStr) // 15:04 is the format for HH:mm
	if err != nil {
		log.Fatalf("Failed to parse cron time: %v", err)
	}

	sampleJob := gocron.NewScheduler(time.UTC)

	parsedCronTime = parsedCronTime.UTC()

	return &job{
		Log:          log,
		JobMethod:    sampleJob,
		CronTime:     parsedCronTime,
		CronInterval: cronInterval,
	}
}

func (j *job) RunCron(discountUseCase *usecase.DiscountUseCase) {

	j.Log.Infof("Starting jobs...")

	_, err := j.JobMethod.Every(1).Day().At(j.CronTime.Format("15:04")).Do(func() {
		ctxActivate, cancelActivate := context.WithCancel(context.Background())
		defer cancelActivate()
		errActivate := discountUseCase.ActivateDiscount(ctxActivate)
		if errActivate != nil {
			j.Log.Error("Error activating discount:", errActivate)
		}

		ctxDeactivate, cancelDeactivate := context.WithCancel(context.Background())
		defer cancelDeactivate()
		errDeactivate := discountUseCase.DeactivateDiscount(ctxDeactivate)
		if errDeactivate != nil {
			j.Log.Error("Error deactivating discount:", errDeactivate)
		}
	})

	if err != nil {
		j.Log.Error("Error setting daily cron job: ", err)
	}

	_, err2 := j.JobMethod.Every(1).Minutes().Do(func() {
		ctxActivate, cancelActivate := context.WithCancel(context.Background())
		defer cancelActivate()
		errActivate := discountUseCase.ActivateDiscount(ctxActivate)
		if errActivate != nil {
			j.Log.Error("Error activating discount:", errActivate)
		}

		ctxDeactivate, cancelDeactivate := context.WithCancel(context.Background())
		defer cancelDeactivate()
		errDeactivate := discountUseCase.DeactivateDiscount(ctxDeactivate)
		if errDeactivate != nil {
			j.Log.Error("Error deactivating discount:", errDeactivate)
		}

	})
	if err2 != nil {
		j.Log.Error("Error setting interval cron job: ", err2)
	}

	j.Log.Infof("Scheduler has started")

	j.JobMethod.StartAsync()
}
