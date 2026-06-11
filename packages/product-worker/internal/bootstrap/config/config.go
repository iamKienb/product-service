package config

import configx "github.com/iamKienb/go-core/config"

type ProductWorkerConfig struct {
	ES       configx.ElasticSearchConfig `envPrefix:"PRODUCT_WORKER_SERVICE"`
	Redis    configx.RedisConfig         `envPrefix:"PRODUCT_WORKER_SERVICE"`
	Kafka    configx.KafkaConfig         `envPrefix:"PRODUCT_WORKER_SERVICE"`
	Consumer configx.ConsumerConfig      `envPrefix:"PRODUCT_WORKER_SERVICE"`
}
