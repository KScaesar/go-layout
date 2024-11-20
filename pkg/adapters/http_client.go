package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/KScaesar/go-layout/pkg"
)

func NewHttpClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConnsPerHost: 5,
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

//

// ConvertErrorFromHttpClient 還不知道會有怎樣的 http error, 所以只先定義預設錯誤
func ConvertErrorFromHttpClient(err error) error {
	switch {
	default:
		return fmt.Errorf("server invoke server issue: %w", pkg.ErrSystem)
	}
}

//

func GetHttpJsonBodyByType[T any](
	// dependency
	client *http.Client,
	logger *slog.Logger,

	// parameter
	ctx context.Context,
	endpoint string,
) (view T, Err error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Error(err.Error(), slog.Any("cause", http.NewRequestWithContext))
		Err = ConvertErrorFromHttpClient(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error(), slog.Any("cause", client.Do))
		Err = ConvertErrorFromHttpClient(err)
		return
	}
	defer response.Body.Close()

	code := response.StatusCode
	if code != http.StatusOK {
		Err = fmt.Errorf("call external service but http status code != 200: %w", pkg.ErrSystem)
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&view)
	if err != nil {
		logger.Error(err.Error(), slog.Any("cause", decoder.Decode))
		Err = fmt.Errorf("json.Decode issue when http.client.Do: %w", pkg.ErrSystem)
		return
	}

	return
}
