package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

// BucketConfig is struct that will receive configuration options via environment variables
type BucketConfig struct {
	Provider          string `env:"GCS_URL"`
	GCPCredentials    string `env:"GOOGLE_APPLICATION_CREDENTIALS"`
	PrivateBucket     string `env:"GOOGLE_STORAGE_BUCKET_PRIVATE"`
	PublicBucket      string `env:"GOOGLE_STORAGE_BUCKET_PUBLIC"`
	PrivateBucketList string `env:"PRIVATE_FOLDERS"`
	PublicBucketList  string `env:"PUBLIC_FOLDERS"`
}

var bucketConfig BucketConfig

// Get are responsible to load env and get data an return the struct
func GetBucketConfig() *BucketConfig {
	if err := env.Parse(&bucketConfig); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &bucketConfig
}
