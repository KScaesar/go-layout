package pkg

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func MustLoadConfig() *Config {
	const defaultPath = "./configs/config.yml"

	filePath := flag.String("conf", defaultPath, "Path to the configuration file")
	flag.Parse()

	logger := Logger().Slog()

	conf, err := utility.LoadLocalConfigFromMultiSource[Config](yaml.Unmarshal, *filePath, logger)
	if err != nil {
		logger.Error("load config fail", slog.Any("err", err))
		os.Exit(1)
	}
	return conf
}

type Config struct {
	AppId_      string       `yaml:"AppId"`
	Hack        utility.Hack `yaml:"Hack"`
	ShowErrCode bool         `yaml:"ShowErrCode"`

	Filepath Filepath `yaml:"Filepath"`
	Http     Http     `yaml:"Http"`
	MySql    MySql    `yaml:"MySql"`
	Redis    Redis    `yaml:"Redis"`

	O11Y   utility.O11YConfig `yaml:"O11Y"`
	Logger wlog.Config        `yaml:"Logger"`
}

func (c *Config) AppId() string {
	if c.AppId_ == "" {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		DefaultAppId := hostname
		c.AppId_ = DefaultAppId
	}
	return c.AppId_
}

type Filepath struct {
	Logger string `yaml:"Logger"` // output to stderr if empty
	Event  string `yaml:"Event"`  // output to stdout if empty
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
