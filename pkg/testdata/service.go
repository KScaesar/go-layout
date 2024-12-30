package wgin

import (
	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func NewDockerServices(conf *pkg.Config) []utility.DockerService {
	return []utility.DockerService{
		utility.NewRedisService("redis1", &conf.Redis, []string{
			"cat /testdata_db0 && cat /testdata_db0 | redis-cli --pipe -n 0",
			"cat /testdata_db1 && cat /testdata_db1 | redis-cli --pipe -n 1",
		}),
	}
}
