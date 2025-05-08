package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
)

func task() {
	fmt.Println("Выполняется задача:", time.Now())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler, err := getScheduler()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := scheduler.Shutdown(); err != nil {
			log.Printf("failed to shutdown scheduler: %v", err)
		}
	}()

	err = addJobs()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Starting scheduler...")
	scheduler.Start()

	<-ctx.Done()
	log.Println("Shutting down scheduler...")
}

func handleShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Println("Received shutdown signal", "signal", s)
	cancel()
}

func getScheduler(cancel context.CancelFunc) (gocron.Scheduler, error) {
	go handleShutdown(cancel)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("not LoadLocation, message: %s", err.Error()))
	}

	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
		gocron.WithLimitConcurrentJobs(5, gocron.LimitModeWait),
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("not LoadLocation, message: %s", err.Error()))
	}

	return scheduler, nil
}

func addJobs(scheduler gocron.Scheduler) error {
	//get from GORM Mysql tasks
	_, err := scheduler.NewJob(
		gocron.CronJob(
			"1 * * * *",
			false,
		),
		gocron.NewTask(
			task,
		),
	)
	if err != nil {
		return errors.New(fmt.Sprintf("not LoadLocation, message: %s", err.Error()))
	}
	return nil
}
