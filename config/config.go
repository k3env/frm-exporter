package config

import (
	"github.com/traefik/paerser/env"
	"github.com/traefik/paerser/file"
	"github.com/traefik/paerser/flag"
	"os"
	"slices"
	"time"
)

type Config struct {
	Web      WebConfig
	FRM      FRMConfig
	LogLevel string
}

type FRMConfig struct {
	Url            string
	ScrapeInterval time.Duration
}

type WebConfig struct {
	Address string
}

func defaultConfig() Config {
	return Config{
		Web: WebConfig{
			Address: ":8080",
		},
		FRM: FRMConfig{
			Url:            "http://localhost:8080",
			ScrapeInterval: time.Second * 1,
		},
		LogLevel: "INFO",
	}
}

func LoadConfig(path string, envPrefix string, ignoredArgs ...string) (*Config, error) {
	cfg := defaultConfig() // Default values
	var err error
	if path != "" {
		err = file.Decode(path, &cfg)
		if err != nil {
			return nil, err
		}
	}

	if envPrefix != "" {
		err = env.Decode(os.Environ(), envPrefix, &cfg)
		if err != nil {
			return nil, err
		}
	}

	args := os.Args[1:]
	for _, s := range ignoredArgs {
		idx := slices.Index(args, s)
		if idx != -1 {
			slices.Delete(args, idx, idx+1)
		}
	}

	err = flag.Decode(args, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
