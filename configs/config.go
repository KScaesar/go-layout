package configs

import (
	"encoding/json"
	"fmt"

	"github.com/KScaesar/go-layout/pkg/utility"
)

func NewConfig(path string) (*Config, error) {
	return utility.NewConfigFromLocal[Config](json.Unmarshal, path)
}

type Config struct {
	Http  Server
	MySql MySql
	Redis Redis
	O11Y  Observability
}

type Server struct {
	Port  int
	Debug bool
}

type MySql struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Debug    bool
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
	User     string
	Password string
	Host     string
	Port     string
}

func (conf *Redis) Address() string {
	return fmt.Sprintf("%v:%v",
		conf.Host,
		conf.Port,
	)
}

type Observability struct {
	Enable bool
}
