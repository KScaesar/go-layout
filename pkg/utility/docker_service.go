package utility

import (
	"context"
	"fmt"
	"io"
	"time"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DockerServiceConfig interface {
	SetHost(host string)
	SetPort(port string)
}

type DockerService func(compose tc.ComposeStack, ctx context.Context) (svc string, err error)

func NewRedisService(svc string, conf DockerServiceConfig, commands []string) DockerService {
	return func(compose tc.ComposeStack, ctx context.Context) (string, error) {
		container, err := compose.ServiceContainer(ctx, svc)
		if err != nil {
			return svc, err
		}

		const redisLogMessage_7_0 = "Ready to accept connections"
		waitStrategy := wait.ForAll(
			wait.ForLog(redisLogMessage_7_0).WithStartupTimeout(2 * time.Second),
		)
		err = waitStrategy.WaitUntilReady(ctx, container)
		if err != nil {
			return svc, fmt.Errorf("wait ready: %w", err)
		}

		for _, command := range commands {
			code, reader, err := container.Exec(ctx, []string{"bash", "-c", command})
			if err != nil {
				return svc, fmt.Errorf("container.Exec: %w", err)
			}
			const success = 0
			if code != success {
				errorMessage, _ := io.ReadAll(reader)
				return svc, fmt.Errorf("\n%v", string(errorMessage))
			}
		}

		host, err := container.Host(ctx)
		if err != nil {
			return svc, err
		}
		port, err := container.MappedPort(ctx, "6379")
		if err != nil {
			return svc, err
		}

		conf.SetHost(host)
		conf.SetPort(port.Port())

		return svc, nil
	}
}
