package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"net/http"
	"strings"
	"taletracker.com/internal/http/handler"
)

type TaleServer struct {
	Config *TaleConfig
	echo   *echo.Echo
}
type TaleConfig struct {
	DevelopmentMode bool
	AllowedOrigins  []string
}

func (t *TaleServer) Start() {
	t.echo = echo.New()
	t.echo.Renderer = &handler.Template{
		Templates: template.Must(template.ParseGlob("../../internal/view/*.html")),
	}
	t.echo.Pre(middleware.RemoveTrailingSlash())
	t.echo.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	// CORS config for non-API purposes
	t.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// Skip if the request is to the API
		Skipper: func(ctx echo.Context) bool {
			return isAPI(ctx)
		},
		AllowOrigins:     t.Config.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}))
	// CORS config for API purposes
	t.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// Skip if the request is not the API
		Skipper: func(ctx echo.Context) bool {
			return !isAPI(ctx)
		},
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: false,
	}))
	t.echo.GET("/home", handler.Home)
	t.echo.Logger.Fatal(t.echo.Start(":1323"))
}

func isAPI(c echo.Context) bool {
	path := c.Request().URL.Path
	if strings.Contains(path, "/api/") {
		return true
	}
	return false
}
