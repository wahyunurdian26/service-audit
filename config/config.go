package config

import (
	"github.com/wahyunurdian26/util/config"
	"github.com/wahyunurdian26/util/constanta"
)

type Config struct {
	MessageBrokerConfig MessageBrokerConfig
	DBConfiguration     DBConfiguration
}

type MessageBrokerConfig struct {
	RabbitMQUrl string
}

type DBConfiguration struct {
	DatabaseUrl string
}

func LoadConfigs() Config {
	return Config{
		MessageBrokerConfig: MessageBrokerConfig{
			RabbitMQUrl: config.Get(constanta.RabbitMqUrl, "amqp://guest:guest@localhost:5672/"),
		},
		DBConfiguration: DBConfiguration{
			DatabaseUrl: config.Get(constanta.DatabaseUrl, "postgres://postgres:postgres@localhost:5432/omnipay_db?sslmode=disable"),
		},
	}
}
