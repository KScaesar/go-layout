package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func TestMain(m *testing.M) {
	defer pkg.Shutdown().Notify(nil)

	// wlogger := wlog.NewStderrLoggerWhenNormal(false)
	// wlogger := wlog.NewStderrLoggerWhenDebug()
	wlogger := wlog.NewStderrLoggerWhenTesting()
	pkg.Logger().PointToNew(wlogger)

	code := m.Run()
	os.Exit(code)
}

func testHttpResponseJsonBody(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	bBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	return string(bBody)
}

func testExpectedHttpResponse(t *testing.T, appResponse any) string {
	bBody, err := json.Marshal(appResponse)
	assert.NoError(t, err)
	return string(bBody)
}
