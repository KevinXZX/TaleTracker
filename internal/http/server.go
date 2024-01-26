package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
	"taletracker.com/internal/http/handler"
	"taletracker.com/internal/view"
)

type TaleServer struct {
	Config *TaleConfig
	echo   *echo.Echo
}
type TaleConfig struct {
	DevelopmentMode bool
	AllowedOrigins  []string
}

func (t *TaleServer) Start() error {
	t.echo = echo.New()
	templates, err := view.ParseTemplates()
	if err != nil {
		return err
	}
	t.echo.Renderer = &handler.Template{
		Templates: templates,
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
	//TODO: One Handler for UI (including login)
	//      one handler for API. This handler includes an auth middleware which either checks for cookies or a bearer token
	t.echo.GET("/home", handler.Home)
	t.echo.GET("/:user/list", handler.List)
	t.echo.Logger.Fatal(t.echo.Start(":1323"))
	return nil
}

func isAPI(c echo.Context) bool {
	path := c.Request().URL.Path
	if strings.Contains(path, "/api/") {
		return true
	}
	return false
}
