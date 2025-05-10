package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"

	"database/sql"
	database "go-api-docker/internal/common/database"
	models "go-api-docker/internal/common/scheduler/models"

	gocron "github.com/go-co-op/gocron/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler, err := getScheduler(cancel)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := scheduler.Shutdown(); err != nil {
			log.Printf("failed to shutdown scheduler: %v", err)
		}
	}()

	err = addJobs(scheduler)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Starting scheduler...")
	scheduler.Start()

	<-ctx.Done()
	log.Println("Shutting down scheduler...")
}

func getScheduler(cancel context.CancelFunc) (gocron.Scheduler, error) {
	go handleShutdown(cancel)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("not LoadLocation, message: %s", err)
	}

	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
		gocron.WithLimitConcurrentJobs(5, gocron.LimitModeWait),
	)
	if err != nil {
		return nil, fmt.Errorf("not LoadLocation, message: %s", err)
	}

	return scheduler, nil
}

func handleShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Println("Received shutdown signal", "signal", s)
	cancel()
}

func addJobs(scheduler gocron.Scheduler) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	rows, err := getTasks(db)
	if err != nil {
		return err
	}
	for rows.Next() {
		var сronTask models.CronTask
		if err := db.ScanRows(rows, &сronTask); err != nil {
			return err
		}
		if err := addJob(&сronTask, scheduler); err != nil {
			return err
		}
	}
	return nil
}

func addJob(task *models.CronTask, scheduler gocron.Scheduler) error {
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
		return fmt.Errorf("not LoadLocation, message: %s", err)
	}
	return nil
}

func getDB() (*gorm.DB, error) {
	db, err := database.ProvideDBConnection()

	if err != nil {
		return nil, fmt.Errorf("error get db connection: %s", err)
	}

	return db, nil
}

func getTasks(db *gorm.DB) (*sql.Rows, error) {
	rows, err := db.Model(&models.CronTask{}).Rows()
	if err != nil {
		return nil, fmt.Errorf("query error: %s", err)
	}
	defer rows.Close()

	return rows, nil
}

func task() {
	fmt.Println("Выполняется задача:", time.Now())
}
