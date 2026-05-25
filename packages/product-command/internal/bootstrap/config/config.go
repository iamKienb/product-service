package config

import (
	"time"

	configx "github.com/iamKienb/go-core/config"
)

type TemporalConfig struct {
	Address            string        `env:"TEMPORAL_ADDRESS"`
	Namespace          string        `env:"TEMPORAL_NAMESPACE"`
	ProductTaskQueue   string        `env:"TEMPORAL_PRODUCT_TASK_QUEUE"`
	InventoryTaskQueue string        `env:"TEMPORAL_INVENTORY_TASK_QUEUE"`
	ActivityTimeout    time.Duration `env:"TEMPORAL_ACTIVITY_TIMEOUT"`
	RollbackTimeout    time.Duration `env:"TEMPORAL_ROLLBACK_TIMEOUT"`
}

type ApiConfig struct {
	APIGatewayAddr string `env:"API_GATEWAY_SERVICE_ADDR"`
}

type ProductCommandConfig struct {
	Postgres configx.PostgresConfig `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Redis    configx.RedisConfig    `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Server   configx.Server         `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Temporal TemporalConfig         `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Gateway  ApiConfig              `envPrefix:"PRODUCT_COMMAND_SERVICE"`
}
