package router

import (
	"context"
	// "net/http"

	"github.com/omniful/go_commons/http"
	"github.com/varun-singhal-omniful/oms-service/controllers"
)

// type Server struct {
// 	Engine *gin.Engine
// 	Server *http.Server
// }

func Initialize(ctx context.Context, s *http.Server) (err error) {
	OrderV1 := s.Engine.Group("/api/v1")
	OrderV1.POST("/order/bulk", controllers.BulkOrders)
	// OrderV1.GET("/", controllers.GetOrders)
	// OrderV1.POST("/validate", controllers.ValidateOrders)
	return
}
