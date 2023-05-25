package config

import (
    "fmt"
    "time"
    "github.com/ilyakaznacheev/cleanenv"
    log "github.com/sirupsen/logrus"
)

type (

    Config struct {
        HTTP `yaml:"http"`
        LOG `yaml:"log"`
        DB `yaml:"db"`
    }

    HTTP struct {
        Port uint16 `env-required:true yaml:"port" env:"HTTP_PORT"`
        Mode string `env-required:true yaml:"mode"`
        ReadTimeout time.Duration `env-required:true yaml:"read-timeout"`
        WriteTimeout time.Duration `env-required:true yaml:"write-timeout"`
    }

    LOG struct {
        Level log.Level `env-required:false env-default:"info" yaml:"level"`
        Console bool `env-required:false env-default:"true" yaml:"console"`
        File string `env-required:false yaml:"file"`
        Keep uint `env-required:false env-default:"30" yaml:"keep"`
    }

    DB struct {
        Host string `env-required:true yaml:"host"`
        Port uint16 `env-required:true yaml:"port"`
        User string `env-required:true yaml:"user"`
        Password string `env-required:true yaml:"password"`
        Database string `env-required:true yaml:"database"`
        Timeout string `env-required:true yaml:"timeout"`
        MaxOpenConns int `env-required:true yaml:"max-open-conns"`
        MaxIdleConns int `env-required:true yaml:"max-idle-conns"`
    }
)

func NewConfig() (*Config, error) {
    cfg := &Config{}

    err := cleanenv.ReadConfig("./config/server.yml", cfg)
    if err != nil {
        return nil, fmt.Errorf("file config error: %w", err)
    }

    err = cleanenv.ReadEnv(cfg)
    if err != nil {
        return nil, fmt.Errorf("env config error: %w", err)
    }

    return cfg, nil
}
