package config

import (
	"time"

	configx "github.com/iamKienb/go-core/config"
)

type TemporalConfig struct {
	Address          string        `env:"_TEMPORAL_ADDRESS"`
	Namespace        string        `env:"_TEMPORAL_NAMESPACE"`
	ProductTaskQueue string        `env:"_TEMPORAL_PRODUCT_TASK_QUEUE"`
	ActivityTimeout  time.Duration `env:"_TEMPORAL_ACTIVITY_TIMEOUT"`
	RollbackTimeout  time.Duration `env:"_TEMPORAL_ROLLBACK_TIMEOUT"`
}

type UpstreamConfig struct {
	ShopCommandURL      string `env:"_SHOP_COMMAND_URL" envDefault:"http://localhost:8002"`
	InventoryCommandURL string `env:"_INVENTORY_COMMAND_URL" envDefault:"http://localhost:8004"`
}

type ProductCommandConfig struct {
	Postgres configx.PostgresConfig `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Redis    configx.RedisConfig    `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Server   configx.Server         `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Temporal TemporalConfig         `envPrefix:"PRODUCT_COMMAND_SERVICE"`
	Upstream UpstreamConfig         `envPrefix:"PRODUCT_COMMAND_SERVICE"`
}
