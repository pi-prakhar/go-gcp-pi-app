package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/metrics"
)

type UserMiddleware struct{}

func NewUserMiddleware() *UserMiddleware {
	return &UserMiddleware{}
}

func (m *UserMiddleware) PrometheusMiddleware(c *gin.Context) {
	start := time.Now()

	c.Next()

	duration := time.Since(start).Seconds()
	statusCode := c.Writer.Status()

	metrics.HttpRequestsTotal.WithLabelValues(
		strconv.Itoa(statusCode), c.Request.Method, c.FullPath(),
	).Inc()

	metrics.RequestDuration.WithLabelValues(
		c.Request.Method, c.FullPath(),
	).Observe(duration)
}
