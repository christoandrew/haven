package api

import (
	"github.com/christo-andrew/haven/internal/api/middleware"
	"github.com/gin-contrib/cors"

	"github.com/christo-andrew/haven/docs"
	"github.com/christo-andrew/haven/internal/api/routes"
	"github.com/christo-andrew/haven/pkg/config"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	app    *gin.Engine
	config *config.ServerConfig
}

func NewServer(config config.ServerConfig) *Server {
	app := gin.Default()
	return &Server{app: app, config: &config}
}

func (s *Server) SetupRouter(db *gorm.DB) *gin.Engine {
	allowedHeaders := []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
		"Access-Control-Allow-Origin",
	}

	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     allowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	docs.SwaggerInfo.BasePath = "/api/v1"
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Haven API",
		})
	})

	s.app.GET("/provider/auth", func(ctx *gin.Context) {
		// Redirect to the provider's auth page
		ctx.JSON(200, gin.H{
			"message": "Welcome to the Haven API",
		})
	})

	routes.UsersRouterV1(v1.Group("/users"), db)
	routes.AuthRouterV1(v1.Group("/auth"), db)
	routes.AccountsRouterV1(v1.Group("/accounts", middleware.AuthMiddleware()), db)
	routes.TransactionsRouterV1(v1.Group("/transactions", middleware.AuthMiddleware()), db)
	routes.CategoriesRouterV1(v1.Group("/categories", middleware.AuthMiddleware()), db)
	routes.DataRouterV1(v1.Group("/data", middleware.AuthMiddleware()), db)
	routes.BudgetsRouterV1(v1.Group("/budgets", middleware.AuthMiddleware()), db)

	return s.app
}

func (s *Server) Run() error {
	err := s.app.Run(s.config.Host + ":" + s.config.Port)
	if err != nil {
		return err
	}
	return nil
}
