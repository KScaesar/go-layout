package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

var testConfig pkg.Config

func TestMain(m *testing.M) {
	// wlogger := wlog.NewStderrLoggerWhenNormal(false)
	// wlogger := wlog.NewStderrLoggerWhenDebug()
	wlogger := wlog.NewDiscardLogger()
	pkg.Logger().PointToNew(wlogger)

	// DownDocker := testdata.UpDocker(true, &testConfig)

	goleak.VerifyTestMain(m,
		goleak.IgnoreCurrent(),
		goleak.Cleanup(func(code int) {
			pkg.Shutdown().Notify(nil)
			<-pkg.Shutdown().WaitChannel()
			// DownDocker()
			os.Exit(code)
		}),
	)
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
