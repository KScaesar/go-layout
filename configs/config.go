package configs

import (
	"fmt"
	"os"

	"github.com/KScaesar/go-layout/pkg/utility"
	"gopkg.in/yaml.v3"
)

func MustLoadConfig(filePath string) *Config {
	conf, err := utility.LoadLocalConfigFromMultiSource[Config](yaml.Unmarshal, filePath, "conf.yml")
	if err != nil {
		panic(err)
	}
	return conf
}

type Config struct {
	Biz       Business      `yaml:"Business"`
	ServiceId string        `yaml:"ServiceId"`
	Hack      utility.Hack  `yaml:"Hack"`
	Http      Http          `yaml:"Http"`
	MySql     MySql         `yaml:"MySql"`
	Redis     Redis         `yaml:"Redis"`
	O11Y      Observability `yaml:"O11Y"`
	Logger    Logger        `yaml:"Logger"`
}

func (c *Config) ServiceId_() string {
	if c.ServiceId == "" {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		DefaultServiceId := hostname
		c.ServiceId = DefaultServiceId
	}
	return c.ServiceId
}

type Business struct {
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

func (conf *Redis) Address() string {
	return fmt.Sprintf("%v:%v",
		conf.Host,
		conf.Port,
	)
}

type Observability struct {
	Enable     bool   `yaml:"Enable"`
	MetricPort string `yaml:"MetricPort"`
}

func (o *Observability) MetricPort_() string {
	if o.MetricPort == "" {
		const DefaultMetricPort = "2112"
		o.MetricPort = DefaultMetricPort
	}
	return o.MetricPort
}

type Logger struct {
	AddSource  bool `yaml:"AddSource"`
	JsonFormat bool `yaml:"JsonFormat"`
}
