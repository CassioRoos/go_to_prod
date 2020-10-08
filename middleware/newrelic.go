package middleware

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"io/ioutil"
	"os"
	"strconv"
)

func ConfigureNewRelic(echoInstance *echo.Echo) *newrelic.Application {
	//log := logger.GetLogger(nil)
	//defer log.Sync()

	newRelicEnvVar := os.Getenv("ENABLE_NEW_RELIC")
	newRelicEnable, err := strconv.ParseBool(newRelicEnvVar)

	if err == nil && newRelicEnable {
		licenseKeyEnvVar := os.Getenv("NEW_RELIC_LICENSE_KEY")
		//proxyURL, _ := url.Parse(os.Getenv("NEW_RELIC_PROXY_URL"))
		//useProxy := false
		//if proxyURL.String() == "" {
		//	fmt.Println(fmt.Sprintf("Error when parsing new relic proxy url - %e", err ))
		//	//log.Error("Error when parsing new relic proxy url", zap.Error(err))
		//} else {
		//	useProxy = true
		//}
		newRelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName("Go to Prod"),
			newrelic.ConfigLicense(licenseKeyEnvVar),
			//func(config *newrelic.Config) {
			//	if useProxy {
			//		config.Transport = &http.Transport{
			//			Proxy: http.ProxyURL(proxyURL),
			//		}
			//	}
			//},
		)
		if err == nil {
			//log.Info("New Relic ENABLED")
			return newRelicApp
		}
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("Error when creating new relic - %e", err ))
		//log.Error("Error when creating new relic", zap.Error(err))
	}
	fmt.Println("New Relic DISABLED")
	//log.Info("New Relic DISABLED")
	return nil
}

// NewRelicMiddleware is a middleware to send request info to New Relic.
func NewRelicMiddleware(app *newrelic.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if app != nil {
				responseWriter := ctx.Response().Writer
				// Copy struct request to remove body.
				requestCopy := *ctx.Request()
				// Set body empty to send to New Relic
				requestCopy.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
				txn := app.StartTransaction(requestCopy.URL.Path)
				txn.SetWebResponse(responseWriter)
				txn.SetWebRequestHTTP(&requestCopy)
				defer txn.End()
				ctx.SetRequest(ctx.Request().WithContext(newrelic.NewContext(ctx.Request().Context(), txn)))
			}
			return next(ctx)
		}
	}
}
