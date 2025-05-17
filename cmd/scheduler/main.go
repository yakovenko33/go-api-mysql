package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"

	"database/sql"
	queue_tasks "go-api-docker/cmd/scheduler/queue"
	models "go-api-docker/cmd/scheduler/tasks"
	database "go-api-docker/internal/common/database"
	publisher "go-api-docker/internal/common/rabbitmq/publisher"

	gocron "github.com/go-co-op/gocron/v2"
	"github.com/rabbitmq/amqp091-go"
)

var (
	publisherInstance *publisher.Publisher
	dbInstance        *gorm.DB
)

func main() {
	limitConcurrentJobs := flag.Uint("limit_concurrent_jobs", 5, "limit_concurrent_jobs")
	channelPoolCount := flag.Int("channel_pool_count", 5, "channel_pool_count for queue")
	flag.Parse()

	if err := initConnetctions(channelPoolCount); err != nil {
		log.Printf("failed to shutdown scheduler: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler, err := getScheduler(cancel, *limitConcurrentJobs)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		defer publisherInstance.CloseChannelPool()
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

func initConnetctions(channelPoolCount *int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	dbInstance = db

	publisherValue, err := queue_tasks.CreateNewPublisher(*channelPoolCount)
	if err != nil {
		return err
	}
	publisherInstance = publisherValue

	return nil
}

func getScheduler(cancel context.CancelFunc, limitConcurrentJobs uint) (gocron.Scheduler, error) {
	go handleShutdown(cancel)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("not LoadLocation, message: %s", err)
	}

	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
		gocron.WithLimitConcurrentJobs(limitConcurrentJobs, gocron.LimitModeWait),
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
	rows, err := getTasks()
	if err != nil {
		return err
	}
	for rows.Next() {
		var сronTask models.CronTask
		if err := dbInstance.ScanRows(rows, &сronTask); err != nil {
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
			task.Cron,
			false,
		),
		gocron.NewTask(
			getTask(task),
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

func getTasks() (*sql.Rows, error) {
	rows, err := dbInstance.Model(&models.CronTask{}).Rows()
	if err != nil {
		return nil, fmt.Errorf("query error: %s", err)
	}
	defer rows.Close()

	return rows, nil
}

func getTask(task *models.CronTask) func() {
	return func() {
		publisherParams := publisher.PublisherParams{
			Exchange:   "cron_scheduler",
			RoutingKey: task.RoutingKey,
			RetryCount: task.MaxRetries,
			RetryDelay: task.RetryTTL * time.Second,
			Msg: amqp091.Publishing{
				ContentType: "application/json",
				Body:        []byte{},
			},
		}

		publisherInstance.Publish(publisherParams)
	}
}
