package config

import (
	"bytes"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gopkg.in/go-playground/validator.v9"
)

type (
	Config struct {
		Logger        Logger        `mapstructure:"logger" validate:"required"`
		Nats          Nats          `mapstructure:"nats" validate:"required"`
		NatsStreaming NatsStreaming `mapstructure:"nats-streaming" validate:"required"`
		Elasticsearch Elasticsearch `mapstructure:"elasticsearch" validate:"required"`
	}

	Nats struct {
		Servers []string `mapstructure:"servers" validate:"required"`
	}

	NatsStreaming struct {
		Servers []string `mapstructure:"servers" validate:"required"`
	}

	Elasticsearch struct {
		Servers []string `mapstructure:"servers" validate:"required"`
	}

	Logger struct {
		Level string `mapstructure:"level" validate:"required"`
	}
)

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

func Init(path string) Config {
	var cfg Config

	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader([]byte(Default))); err != nil {
		logrus.Panicf("error loading default configs: %s", err.Error())
	}

	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.SetEnvPrefix(Namespace)
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err := v.MergeInConfig()
	if err != nil {
		logrus.Warn("no config file found. Using defaults and environment variables")
	}

	if err := v.UnmarshalExact(&cfg); err != nil {
		logrus.Fatalf("invalid configuration: %s", err)
	}

	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("invalid configuration: %s", err)
	}

	return cfg
}
