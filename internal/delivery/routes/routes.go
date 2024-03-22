package routes

import (
	"net/http"

	"github.com/Dzikuri/openidea-segokuning/internal/delivery/handler"
	"github.com/Dzikuri/openidea-segokuning/internal/delivery/middleware"
	"github.com/labstack/echo/v4"
)

type RoutesConfig struct {
	Echo       *echo.Echo
	Middleware middleware.Middleware
	Handler    handler.Handler
}

func (c *RoutesConfig) Setup() {

	c.Echo.GET("/health", func(c echo.Context) error {
		// Your existing health-check handler logic
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Ok",
		})
	})

	// Error Handling Page not found
	c.Echo.Any("/*", func(c echo.Context) error {
		var d struct{}
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not found", "data": d})
	})

	c.SetupRouteAuth()
	c.SetupRouteUser()
}

func (c *RoutesConfig) SetupRouteAuth() {
	c.Echo.POST("/v1/user/register", c.Handler.UserRegister)
	c.Echo.POST("/v1/user/login", c.Handler.UserLogin)
}

func (c *RoutesConfig) SetupRouteUser() {
	c.Echo.POST("/v1/user/link/email", c.Handler.UserLinkEmail, c.Middleware.Authentication(true))
	c.Echo.POST("/v1/user/link/phone", c.Handler.UserLinkPhone, c.Middleware.Authentication(true))
	c.Echo.PATCH("/v1/user", c.Handler.UserUpdateAccount, c.Middleware.Authentication(true))
}
