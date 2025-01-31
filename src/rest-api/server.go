package restapi

import (
	"net/http"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/rest-api/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/cloud-barista/cb-mcks/src/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Server() {
	// Echo instance
	e := echo.New()

	// Echo middleware func

	if *app.Config.LoglevelHTTP == true {
		e.Use(middleware.Logger()) // Setting logger
	}
	e.Use(middleware.Recover())                            // Recover from panics anywhere in the chain
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{ // CORS Middleware
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET(*app.Config.RootURL+"/healthy", router.Healthy)

	m := e.Group(*app.Config.RootURL + "/mcir/connections")

	m.GET("/:connection/specs", router.ListSpec)

	g := e.Group(*app.Config.RootURL+"/ns", validMiddlewareFunc())

	// Routes
	g.GET("/:namespace/clusters", router.ListCluster)
	g.POST("/:namespace/clusters", router.CreateCluster)
	g.GET("/:namespace/clusters/:cluster", router.GetCluster)
	g.DELETE("/:namespace/clusters/:cluster", router.DeleteCluster)

	g.GET("/:namespace/clusters/:cluster/nodes", router.ListNode)
	g.POST("/:namespace/clusters/:cluster/nodes", router.AddNode)
	g.GET("/:namespace/clusters/:cluster/nodes/:node", router.GetNode)
	g.DELETE("/:namespace/clusters/:cluster/nodes/:node", router.RemoveNode)

	// Start server
	e.Logger.Fatal(e.Start(":1470"))
}

func validMiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ns := c.Param("namespace")
			if ns == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "'Namespace' is a mandatory parameter.")
			}
			return next(c)
		}
	}
}
