package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/handlers"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	Engine     *gin.Engine
	handler    *handlers.UserHandler
	middleware *middleware.UserMiddleware
}

func NewRouter(handler *handlers.UserHandler, middleware *middleware.UserMiddleware) *Router {
	engine := gin.Default()
	engine.Use(middleware.PrometheusMiddleware)

	router := &Router{
		Engine:     engine,
		handler:    handler,
		middleware: middleware,
	}
	router.initRoutes()
	return router
}

func (r *Router) initRoutes() {

	api := r.Engine.Group("/api/v1/users")
	{
		api.GET("/metrics", gin.WrapH(promhttp.Handler()))
		api.GET("/get/:mail", r.handler.GetUserByEmail)
		api.GET("/get", r.handler.GetUsers)
		api.POST("/create", r.handler.CreateUser)
		api.PUT("/update", r.handler.UpdateUser)
		api.DELETE("/delete/:mail", r.handler.DeleteUser)

	}
}
