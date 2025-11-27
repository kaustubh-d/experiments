package main

import (
	"log"
	"os"
	"path/filepath"

	// Import app package to use its types and handlers
	"github.com/kaustubh-d/experiments/inventory/app"
	"github.com/labstack/echo/v4"
)

func main() {
	// data directory can be overridden by DATA_DIR env var, default "data"
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		// used in development and testing
		dataDir = filepath.Join("..", "data")
	}

	ds := app.NewDataStore(dataDir)

	e := echo.New()

	// GET /inventory/apps -> return enabled-app-list.yaml contents
	e.GET("/inventory/apps", func(c echo.Context) error {
		log.Println("Called /inventory/apps.")
		c.Set("ds", ds)
		return app.GetEnabledApps(c)
	})

	// GET /inventory/apps/:app -> list envs (files in data/<app>/)
	e.GET("/inventory/apps/:app", func(c echo.Context) error {
		log.Printf("Called /inventory/apps/%s.\n", c.Param("app"))
		c.Set("ds", ds)
		return app.GetAppEnvs(c)
	})

	// GET /inventory/apps/:app/:env -> read and return data/<app>/<env>.yaml
	e.GET("/inventory/apps/:app/:env", func(c echo.Context) error {
		log.Printf("Called /inventory/apps/%s/%s.\n",
			c.Param("app"), c.Param("env"))
		c.Set("ds", ds)
		return app.GetAppEnvDetails(c)
	})

	e.Logger.Fatal(e.Start(":9080"))
}
