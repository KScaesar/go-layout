//go:build intg

package utility

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

// https://golang.testcontainers.org/features/docker_compose/

// UpDocker 進行 integration test 需要啟動 docker
//
// 必要的環境變數：
//
//   - WORK_DIR:
//     為了找到 `docker-compose.yml`，需要指定工作目錄
//   - CGO_ENABLED: 設定為 `1`
//     因為在 macOS 中，需要啟用 CGO 才能正確運行 fsevents， https://github.com/fsnotify/fsevents
func UpDocker(removeData bool, services []DockerService) (DownDocker func() error) {
	var Err error
	defer func() {
		if Err != nil {
			panic(Err)
		}
	}()

	compose, err := newDockerCompose()
	if err != nil {
		Err = err
		return
	}

	ctx := context.Background()

	err = compose.Up(ctx, tc.RemoveOrphans(true))
	if err != nil {
		Err = fmt.Errorf("compose.Up: %w", err)
		return
	}

	down := func() error {
		err := compose.Down(ctx, tc.RemoveVolumes(removeData), tc.RemoveOrphans(true))
		if err != nil {
			return fmt.Errorf("compose.Down: %w", err)
		}
		return nil
	}
	defer func() {
		if Err != nil {
			downErr := down()
			if downErr != nil {
				panic(downErr)
			}
		}
	}()

	mqError := make(chan error, len(services))
	wg := sync.WaitGroup{}
	for _, setup := range services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc, err := setup(compose, ctx)
			if err != nil {
				mqError <- fmt.Errorf("svc=%v: setup service: %w", svc, err)
			}
		}()
	}

	go func() {
		wg.Wait()
		mqError <- nil
	}()

	Err = <-mqError
	return down
}

// 嘗試找正確的 docker-compose.yml 路徑
func newDockerCompose() (tc.ComposeStack, error) {
	workDirs := []func() (string, error){
		func() (string, error) {
			wd, err := os.Getwd()
			if err != nil {
				return "", fmt.Errorf("os.Getwd: %w", err)
			}
			return wd, nil
		},

		func() (string, error) {
			wd, ok := os.LookupEnv("WORK_DIR")
			if !ok {
				return "", errors.New("ENV ${WORK_DIR} not set")
			}
			return wd, nil
		},
	}

	var Err error
	for i := 0; i < len(workDirs); i++ {
		workDir, err := workDirs[i]()
		if err != nil {
			Err = err
			continue
		}

		files := []string{
			filepath.Join(workDir, "pkg", "testdata", "docker-compose.yml"),
			filepath.Join(workDir, "docker-compose.yml"),
		}
		for _, f := range files {
			_, err := os.Stat(f)
			if err != nil {
				Err = err
				continue
			}
			compose, err := tc.NewDockerComposeWith(tc.WithStackFiles(f))
			if err != nil {
				Err = err
				continue
			}
			return compose, nil
		}
	}

	return nil, Err
}
