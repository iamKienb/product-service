package config

import configx "github.com/iamKienb/go-core/config"

type ProductQueryConfig struct {
	ES     configx.ElasticSearchConfig `envPrefix:"PRODUCT_QUERY_SERVICE"`
	Server configx.Server              `envPrefix:"PRODUCT_QUERY_SERVICE"`
}
