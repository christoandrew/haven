package api

import (
	"github.com/christo-andrew/haven/docs"
	"github.com/christo-andrew/haven/internal/api/middleware"
	"github.com/christo-andrew/haven/pkg/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	app    *gin.Engine
	Config *config.ServerConfig
}

func NewApiServer(config *config.Config) *Server {
	app := gin.Default()
	return &Server{app: app, Config: &config.Server}
}

func (s *Server) SetupRouter(db *gorm.DB) *gin.Engine {
	s.app.Use(middleware.CorsMiddleware())

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

	UsersRouterV1(v1.Group("/users"), db)
	AuthRouterV1(v1.Group("/auth"), db)
	AccountsRouterV1(v1.Group("/accounts", middleware.WithAuthUser()), db)
	TransactionsRouterV1(v1.Group("/transactions", middleware.WithAuthUser()), db)
	CategoriesRouterV1(v1.Group("/categories", middleware.WithAuthUser()), db)
	DataRouterV1(v1.Group("/data", middleware.WithAuthUser()), db)
	BudgetsRouterV1(v1.Group("/budgets", middleware.WithAuthUser()), db)

	return s.app
}

func (s *Server) Run() error {
	err := s.app.Run(s.Config.GetAddress())
	if err != nil {
		return err
	}
	return nil
}
