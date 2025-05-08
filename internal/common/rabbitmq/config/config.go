package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CronTask struct {
	Name       string `yaml:"name"`
	RoutingKey string `yaml:"routing_key"`
	Queue      string `yaml:"queue"`
	Cron       string `yaml:"cron"`
	RetryTTL   int    `yaml:"retry_ttl"`
	MaxRetries int    `yaml:"max_retries"`
	DLXEnabled bool   `yaml:"dlx_enabled"`
}

type CronConfig struct {
	Exchange     string     `yaml:"exchange"`
	ExchangeType string     `yaml:"exchange_type"`
	Tasks        []CronTask `yaml:"tasks"`
}

func LoadCronConfig(path string) (*CronConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg CronConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func DeclareTasks(ch *amqp.Channel, cfg *CronConfig) error {
	err := ch.ExchangeDeclare(
		cfg.Exchange,
		cfg.ExchangeType,
		true, false, false, false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, task := range cfg.Tasks {
		args := amqp.Table{}
		if task.DLXEnabled {
			args["x-dead-letter-exchange"] = cfg.Exchange + ".dlx"
			args["x-dead-letter-routing-key"] = task.RoutingKey + ".retry"
			args["x-message-ttl"] = int32(task.RetryTTL)
		}
		_, err := ch.QueueDeclare(
			task.Queue,
			true, false, false, false,
			args,
		)
		if err != nil {
			return err
		}

		err = ch.QueueBind(
			task.Queue,
			task.RoutingKey,
			cfg.Exchange,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
