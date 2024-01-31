package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"taletracker.com/internal/model"
	"taletracker.com/internal/taledb"
)

func RegisterRoutes(g *echo.Group) string {
	g.GET("/tale", GetTales)
	g.POST("/tale", AddTale)
	g.POST("/tale/:taleID/review", ReviewTale)
	g.POST("/tale/:taleID/tag", AddTags)
	g.DELETE("/tale/:taleID/tag/:tagID", DeleteTags)
	return "api"
}

func GetTales(c echo.Context) error {
	tdb := c.Get("tdb").(*taledb.TaleDatabase)
	tales, err := tdb.GetTales()
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

func ReviewTale(c echo.Context) error {
	tdb := c.Get("tdb").(*taledb.TaleDatabase)
	var review model.Review
	err := c.Bind(&review)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	q := c.Param("taleID")
	taleID, err := strconv.Atoi(q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = tdb.ReviewTale(taleID, review)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, review)
}
func AddTags(c echo.Context) error {
	tdb := c.Get("tdb").(*taledb.TaleDatabase)
	var tags = struct {
		Tags []model.Tag `json:"tags"`
	}{}
	err := c.Bind(&tags)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	q := c.Param("taleID")
	taleID, err := strconv.Atoi(q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = tdb.AddTagToTale(taleID, tags.Tags...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, tags)
}
func DeleteTags(c echo.Context) error {
	tdb := c.Get("tdb").(*taledb.TaleDatabase)
	q := c.Param("taleID")
	taleID, err := strconv.Atoi(q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	tagID, err := strconv.Atoi(c.Param("tagID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = tdb.RemoveTagFromTale(taleID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "")
}
