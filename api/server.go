package api

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"manager-service/config"
	"manager-service/db"
	docs "manager-service/docs"
	"manager-service/proto"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config     config.Config
	store      *db.Store
	grpcClient proto.AuthServiceClient
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config config.Config, store *db.Store, grpcClient proto.AuthServiceClient) (*Server, error) {

	gin.SetMode(config.GinMode)
	router := gin.Default()

	server := &Server{
		config:     config,
		store:      store,
		grpcClient: grpcClient,
	}

	// Setup routing for server
	v1 := router.Group("v1")
	{
		v1.GET("/users/:id", server.GetUserByID)
		v1.GET("/users", server.GetAllUsers)
		v1.POST("/users", server.CreateUser)
		v1.POST("/login", server.LoginUser)
	}

	// Setup health check routes
	health := router.Group("health")
	{
		health.GET("/live", server.Live)
		health.GET("/ready", server.Ready)
	}

	docs.SwaggerInfo.BasePath = "/manager-service"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// TODO: Setup metrics routes

	server.router = router
	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
