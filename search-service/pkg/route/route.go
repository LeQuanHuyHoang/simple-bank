package route

import (
	"Go_Learn/conf"
	"Go_Learn/pkg/handler"
	"Go_Learn/pkg/repo"
	srv "Go_Learn/pkg/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config conf.Config
	pg     *gorm.DB
	router *gin.Engine
}

func NewServer(config conf.Config, pg *gorm.DB) (*Server, error) {
	server := &Server{
		config: config,
		pg:     pg,
	}

	userRepo := repo.NewRepo(pg)
	userService := srv.NewUserService(userRepo, config)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	v1Api := router.Group("/api/v1")

	v1Api.POST("/login", userHandler.Login)
	v1Api.POST("/search-file", userHandler.SearchFile)

	migrate := handler.NewMigrationHandler(pg)
	router.POST("/internal/migrate", migrate.Migrate)
	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
