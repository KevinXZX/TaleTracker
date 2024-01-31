package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
	"taletracker.com/internal/http/handler"
	"taletracker.com/internal/http/handler/api"
	"taletracker.com/internal/taledb"
	"taletracker.com/internal/view"
)

type TaleServer struct {
	Config *TaleConfig
	echo   *echo.Echo
	Db     *taledb.TaleDatabase
}
type TaleConfig struct {
	DevelopmentMode bool
	AllowedOrigins  []string
}

func tdbMiddleware(tdb *taledb.TaleDatabase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("tdb", tdb)
			return next(c)
		}
	}
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
	t.echo.Use(tdbMiddleware(t.Db))
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
	t.echo.GET("/:user/list", handler.List)
	a := t.echo.Group("/api/v1")
	api.RegisterRoutes(a)
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
