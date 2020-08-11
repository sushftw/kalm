package server

import (
	"flag"
	"github.com/go-playground/validator/v10"
	"github.com/kalmhq/kalm/api/log"
	"net"
	"net/http"
	"os"

	"github.com/kalmhq/kalm/api/config"
	"github.com/kalmhq/kalm/api/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func isTest() bool {
	testFlag := flag.Lookup("test.v")

	if testFlag == nil {
		return false
	}

	return testFlag.Value.String() == "true"
}

// Only trust envoy external address or tcp remote address
func getClientIP(req *http.Request) string {
	if req.Header.Get("X-Envoy-External-Address") != "" {
		return req.Header.Get("X-Envoy-External-Address")
	}

	ra, _, _ := net.SplitHostPort(req.RemoteAddr)

	return ra
}

func NewEchoInstance() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Gzip())
	e.Use(middlewareLogging)

	e.IPExtractor = getClientIP

	e.Pre(middleware.RemoveTrailingSlash())

	// TODO, only enabled cors on dev env
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	if !isTest() {
		// TODO is our safe to ignore CSRF protection?
		//e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		//	CookiePath:     "/",
		//	CookieHTTPOnly: true,
		//}))
	}

	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewEchoServer(runningConfig *config.Config) *echo.Echo {
	e := NewEchoInstance()

	// in production docker build, all things are in a single docker
	// golang api server is charge of return frontend files to users
	// If the STATIC_FILE_ROOT is set, add extra routes to handle static files
	staticFileRoot := os.Getenv("STATIC_FILE_ROOT")

	if staticFileRoot != "" {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  staticFileRoot,
			HTML5: true,
		}))
	}

	e.Validator = &CustomValidator{validator: validator.New()}

	return e
}

func middlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c != nil {
			log.Info("receive request", "method", c.Request().Method, "uri", c.Request().URL.String(), "ip", c.RealIP())
		} else {
			log.Info("receive request bad request")
		}

		return next(c)
	}
}
