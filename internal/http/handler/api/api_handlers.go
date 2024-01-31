package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"taletracker.com/internal/model"
	"taletracker.com/internal/taledb"
)

func RegisterRoutes(g *echo.Group) string {
	g.GET("/tales", GetTales)
	g.POST("/tales", AddTale)
	return "api"
}

func GetTales(c echo.Context) error {
	tdb := c.Get("tdb").(*taledb.TaleDatabase)
	tales, err := tdb.GetTales("1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, tales)
}

func AddTale(c echo.Context) error {
	tdb := c.Get("tdb").(*taledb.TaleDatabase)
	var tale model.Tale
	err := c.Bind(&tale)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = tdb.AddTale(tale)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, tale)
}
