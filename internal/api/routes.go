package api

import (
	"github.com/christo-andrew/haven/internal/api/handlers"
	"github.com/christo-andrew/haven/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AccountsRouterV1(router *gin.RouterGroup, db *gorm.DB) {
	router.GET("", func(ctx *gin.Context) {
		handlers.GetAllAccountsHandler(ctx, db)
	})

	router.GET("/:id", func(ctx *gin.Context) {
		handlers.GetAccountHandler(ctx, db)
	})

	router.GET("/:id/statistics", func(ctx *gin.Context) {
		handlers.AccountStatisticsHandler(ctx, db)
	})

	router.GET("/:id/transactions", func(ctx *gin.Context) {
		handlers.GetAccountTransactionsHandler(ctx, db)
	})

	router.POST("/:id/transactions/create", func(ctx *gin.Context) {
		handlers.CreateAccountTransactionHandler(ctx, db)
	})

	router.POST("/create", func(ctx *gin.Context) {
		handlers.CreateAccountHandler(ctx, db)
	})

	router.GET("/:id/transactions/recent", func(ctx *gin.Context) {
		handlers.GetRecentTransactionsHandler(ctx, db)
	})

	router.POST("/:id/transactions/upload", func(ctx *gin.Context) {
		handlers.UploadAccountTransactionsHandler(ctx, db)
	})

	router.GET("/:id/transactions/percentage", func(ctx *gin.Context) {
		handlers.PercentageOfTotalAmountByTransactionHandler(ctx, db)
	})
}

func TransactionsRouterV1(router *gin.RouterGroup, db *gorm.DB) {
	router.GET("", func(ctx *gin.Context) {
		handlers.GetAllTransactionsHandler(ctx, db)
	})

	router.GET("/:id", func(ctx *gin.Context) {
		handlers.GetTransactionHandler(ctx, db)
	})

	router.POST("/create", func(ctx *gin.Context) {
		handlers.CreateAccountTransactionHandler(ctx, db)
	})

	router.GET("/recent", func(ctx *gin.Context) {
		handlers.GetRecentTransactionsHandler(ctx, db)
	})

	router.POST("/:id/tags", func(ctx *gin.Context) {
		handlers.AddTransactionTagHandler(ctx, db)
	})

	router.GET("/:id/tags", func(ctx *gin.Context) {
		handlers.GetTransactionTagsHandler(ctx, db)
	})

	router.GET("/schemas", func(ctx *gin.Context) {
		handlers.GetTransactionSchemasHandler(ctx)
	})
}

func CategoriesRouterV1(router *gin.RouterGroup, db *gorm.DB) {
	router.GET("", func(ctx *gin.Context) {
		handlers.GetAllCategoriesHandler(ctx, db)
	})

	router.GET("/:id", func(ctx *gin.Context) {
		handlers.GetCategoryHandler(ctx, db)
	})

	router.POST("/create", func(ctx *gin.Context) {
		handlers.CreateCategoryHandler(ctx, db)
	})

	router.GET("/types", func(ctx *gin.Context) {
		handlers.GetCategoryByContextAndContextTypeHandler(ctx, db)
	})
}

func UsersRouterV1(router *gin.RouterGroup, db *gorm.DB) {
	router.GET("/", middleware.WithAuthUser(), func(ctx *gin.Context) {
		handlers.GetAllUsersHandler(ctx, db)
	})

	router.POST("/create", func(ctx *gin.Context) {
		handlers.CreateUserHandler(ctx, db)
	})

	router.GET("/:id", middleware.WithAuthUser(), func(ctx *gin.Context) {
		handlers.GetUserHandler(ctx, db)
	})
}

func AuthRouterV1(router *gin.RouterGroup, db *gorm.DB) {
	router.POST("/login", func(ctx *gin.Context) {
		handlers.LoginHandler(ctx, db)
	})
}

func DataRouterV1(router *gin.RouterGroup, db *gorm.DB) {
	router.GET("/:account_id/transactions/histogram", func(ctx *gin.Context) {
		handlers.TransactionsHistogramHandler(ctx, db)
	})
	router.GET("/:account_id/transactions/summary", func(ctx *gin.Context) {
		handlers.TransactionsSummaryHandler(ctx, db)
	})
}

func BudgetsRouterV1(router *gin.RouterGroup, db *gorm.DB) {
	router.POST("/create", func(ctx *gin.Context) {
		handlers.CreateBudgetHandler(ctx, db)
	})

	router.GET("/:id", func(ctx *gin.Context) {
		handlers.GetBudgetHandler(ctx, db)
	})

	router.PUT("/:id/update", func(ctx *gin.Context) {
		handlers.UpdateBudgetHandler(ctx, db)
	})

	router.GET("", func(ctx *gin.Context) {
		handlers.GetBudgetsHandler(ctx, db)
	})

	//router.POST(":id/tags", func(ctx *gin.Context) {
	//	handlers.AddBudgetTagHandler(ctx, db)
	//})
}
