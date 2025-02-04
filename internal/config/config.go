package config

import (
    "os"
    "fmt"
    "time"
    "github.com/ilyakaznacheev/cleanenv"
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/utils"
)

type (

    Config struct {
        HTTP `yaml:"http"`
        LOG `yaml:"log"`
        DB `yaml:"db"`
        CHAINS map[string]CHAIN `yaml:"chains"`
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

    ContractAddress struct {
        Poop string `yaml:"Poop"`
        PoopRouter string `yaml:"PoopRouter"`
        WBNBRouter string `yaml:"WBNBRouter"`
    }

    CHAIN struct {
        ChainName string `yaml:"chain-name"`
        ChainId uint `yaml:"chain-id"`
        Interval time.Duration `yaml:"interval"`
        StartBlock uint64 `yaml:"start-block"`
        DelayBlock uint64 `yaml:"delay-block"`
        MaxBlock uint64 `yaml:"max-block"`
        Urls []string `yaml:"urls"`
        Contracts []string `yaml:"contracts"`
        Addresses ContractAddress `yaml:"addresses"`
        Abi map[string]string `yaml:"abi"`
    }
)

var chainIdMap = map[uint]*string{}

func NewConfig() (*Config, error) {
    cfg := &Config{}
    fileName := os.Getenv("CONFIG_FILE")
    if fileName == "" {
        fileName = "./config/config.yml"
    }
    fmt.Printf("config file: %s\n", fileName)

    err := cleanenv.ReadConfig(fileName, cfg)
    if err != nil {
        return nil, fmt.Errorf("file config error: %w", err)
    }
    for _, value := range cfg.CHAINS {
        chainIdMap[value.ChainId] = utils.PTR(value.ChainName)
    }

    /*
    err = cleanenv.ReadEnv(cfg)
    if err != nil {
        return nil, fmt.Errorf("env config error: %w", err)
    }
    */
    return cfg, nil
}

func GetChainNameByChainId(chainId uint) *string {
    return chainIdMap[chainId]
}
