package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/handlers"
)

type Router struct {
	Engine  *gin.Engine
	handler *handlers.UserHandler
}

func NewRouter(handler *handlers.UserHandler) *Router {
	engine := gin.Default()
	router := &Router{
		Engine:  engine,
		handler: handler,
	}
	router.initRoutes()
	return router
}

func (r *Router) initRoutes() {

	api := r.Engine.Group("/api/v1/users")
	{
		api.GET("/get/:mail", r.handler.GetUser)
		api.GET("/create", r.handler.CreateUser)
	}
}
