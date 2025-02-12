package utility

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var ErrDecodeConf = errors.New("decode config")

// LoadLocalConfigFromMultiSource attempts to load a local configuration file from multiple sources.
// It tries to locate the configuration file using the following sources in order:
//
// 1. The specified FilePath parameter
// 2. The current directory
// 3. The user's home directory
// 4. The path defined by the environment variable CONF_PATH
//
// Parameters:
// - decode: A function used to decode the content of the configuration file into an instance of type T.
// - FilePath: The explicit path to the configuration file.
func LoadLocalConfigFromMultiSource[T any](
	decode Unmarshal,
	FilePath string,
	logger *slog.Logger,
) (
	conf *T,
	err error,
) {
	const (
		byNormal int = iota + 1
		byCurrentDir
		byHomeDir
		byEnvironmentVariable
		stop
	)

	const DefaultEnvName = "CONF_PATH"

	fileName := ""
	if FilePath != "" {
		fileName = filepath.Base(FilePath)
	}

	pathSources := map[int]func() (string, error){
		byNormal: func() (string, error) {
			return FilePath, nil
		},
		byCurrentDir: func() (string, error) {
			currentDir, err := os.Getwd()
			if err != nil {
				return "", err
			}
			return filepath.Join(currentDir, fileName), nil
		},
		byHomeDir: func() (string, error) {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return filepath.Join(homeDir, fileName), nil
		},
		byEnvironmentVariable: func() (string, error) {
			byEnvPath, ok := os.LookupEnv(DefaultEnvName)
			if !ok {
				return "", errors.New("ENV ${CONF_PATH} not set")
			}
			return byEnvPath, nil
		},
	}

	var path string
	for source := byNormal; source < stop; source++ {
		path, err = pathSources[source]()
		if err != nil {
			continue
		}

		conf, err = LoadLocalFile[T](decode, path)
		if err == nil {
			logger.Info("load config", slog.String("path", path))
			return conf, nil
		}

		if errors.Is(err, ErrDecodeConf) {
			return nil, err
		}

		logger.Warn("try load config", slog.String("path", path))
	}
	return nil, err
}

func LoadLocalFile[T any](decode Unmarshal, FilePath string) (*T, error) {
	if FilePath == "" {
		return nil, fmt.Errorf("FilePath cannot be empty")
	}

	file, err := os.Open(FilePath)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}

	bData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read byte: %w", err)
	}

	var conf T
	err = decode(bData, &conf)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeConf, err)
	}

	return &conf, nil
}
