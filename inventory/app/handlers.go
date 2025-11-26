package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// GET /inventory/apps -> return enabled-app-list.yaml contents
func getEnabledApps(c echo.Context) error {
	ds := c.Get("ds").(*DataStore)

	apps, err := ds.loadEnabledApps()
	if err != nil {
		if os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound,
				map[string]string{"error": "enabled-app-list.yaml not found"})
		}
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, apps)
}

// GET /inventory/apps/:app -> list envs (files in data/<app>/)
func getAppEnvs(c echo.Context) error {
	ds := c.Get("ds").(*DataStore)

	app := c.Param("app")
	log.Println("App param:", app)
	if app == "" {
		return c.JSON(http.StatusBadRequest,
			map[string]string{"error": "app name required"})
	}
	envs, err := ds.listEnvs(app)
	if err != nil {
		if os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound,
				map[string]string{"error": "app not found"})
		}
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
	}
	log.Println("Envs found:", envs)
	return c.JSON(http.StatusOK, AppEnvListResponse{Environments: envs})
}

// GET /inventory/apps/:app/:env -> read and return data/<app>/<env>.yaml
func getAppEnvDetails(c echo.Context) error {
	ds := c.Get("ds").(*DataStore)

	app := c.Param("app")
	env := c.Param("env")
	if app == "" || env == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "app and env required"})
	}
	data, err := ds.loadAppEnv(app, env)
	if err != nil {
		if os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "env file not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}
