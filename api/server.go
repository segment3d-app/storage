package api

import (
	"github.com/gin-gonic/gin"
	"github.com/segment3d-app/segment3d-storage/docs"
	"github.com/segment3d-app/segment3d-storage/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config util.Config
	router *gin.Engine
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewServer(config *util.Config) (*Server, error) {
	server := &Server{config: *config}
	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// configure swagger docs
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/:foldername/:filename", server.getFile)
	router.POST("/upload", server.uploadFile)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
