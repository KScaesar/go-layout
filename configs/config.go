package configs

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility"
	"gopkg.in/yaml.v3"
)

func MustLoadConfig(filePath string) *Config {
	logger := pkg.DefaultLogger().Logger

	conf, err := utility.LoadLocalConfigFromMultiSource[Config](
		yaml.Unmarshal,
		filePath,
		"local.yml",
		logger,
	)
	if err != nil {
		logger.Error("load config fail", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Debug("print config", slog.Any("conf", conf))
	return conf
}

type Config struct {
	ServiceId_ string       `yaml:"ServiceId"`
	Hack       utility.Hack `yaml:"Hack"`

	Http  Http  `yaml:"Http"`
	MySql MySql `yaml:"MySql"`
	Redis Redis `yaml:"Redis"`

	O11Y   utility.O11YConfig   `yaml:"O11Y"`
	Logger utility.LoggerConfig `yaml:"Logger"`
}

func (c *Config) ServiceId() string {
	if c.ServiceId_ == "" {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		DefaultServiceId := hostname
		c.ServiceId_ = DefaultServiceId
	}
	return c.ServiceId_
}

type Http struct {
	Port  string `yaml:"Port"`
	Debug bool   `yaml:"Debug"`
}

type MySql struct {
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Database string `yaml:"Database"`
	Debug    bool   `yaml:"Debug"`
}

func (conf *MySql) DSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)
}

type Redis struct {
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
}

func (conf Redis) Address() string {
	return fmt.Sprintf("%v:%v", conf.Host, conf.Port)
}
