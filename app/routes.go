package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/mnizarzr/dot-test/config"
	"github.com/mnizarzr/dot-test/db"
	"github.com/mnizarzr/dot-test/middleware"
	"github.com/mnizarzr/dot-test/modules/auth"
	"github.com/mnizarzr/dot-test/modules/project"
	"github.com/mnizarzr/dot-test/modules/task"
	"github.com/mnizarzr/dot-test/modules/user"
	"gorm.io/gorm"
)

// Dependencies holds all the application dependencies
type Dependencies struct {
	Config    *config.Config
	DB        *gorm.DB
	Redis     *db.RedisClient
	JobClient *asynq.Client
}

// BuildHandler creates and configures all route handlers with dependency injection
func BuildHandler(config *config.Config, router *gin.Engine, database *gorm.DB, redisClient *db.RedisClient) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
	}
	jobClient := asynq.NewClient(redisOpt)

	deps := &Dependencies{
		Config:    config,
		DB:        database,
		Redis:     redisClient,
		JobClient: jobClient,
	}

	setupRoutesV1(router, deps)
}

// setupRoutesV1 configures all application routes
func setupRoutesV1(router *gin.Engine, deps *Dependencies) {
	api := router.Group("/api/v1")

	//  middleware to inject user information for GORM hooks
	api.Use(middleware.AuditResourceContext())

	setupAuthRoutes(api, deps)

	setupUserRoutes(api, deps)
	setupProjectRoutes(api, deps)
	setupTaskRoutes(api, deps)
}

// setupAuthRoutes configures auth module routes with dependency injection
func setupAuthRoutes(api *gin.RouterGroup, deps *Dependencies) {
	userRepo := user.NewRepository(deps.DB, deps.Redis)
	authService := auth.NewService(userRepo, deps.JobClient, deps.Config.JWTSecret)
	authHandler := auth.NewHandler(authService)

	authGroup := api.Group("/auth")
	{
		authGroup.Use(middleware.AuditMiddleware(deps.DB))
		// Registration allows both authenticated (admin) and unauthenticated users
		authGroup.POST("/register", middleware.OptionalJWTAuth(deps.Config.JWTSecret), authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}
}

// setupUserRoutes configures user module routes with dependency injection
func setupUserRoutes(api *gin.RouterGroup, deps *Dependencies) {
	userRepo := user.NewRepository(deps.DB, deps.Redis)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	userGroup := api.Group("/user")
	userGroup.Use(middleware.JWTAuth(deps.Config.JWTSecret))
	{
		userGroup.GET("/me", userHandler.GetProfile)
	}
}

// setupProjectRoutes configures project module routes with dependency injection
func setupProjectRoutes(api *gin.RouterGroup, deps *Dependencies) {
	projectRepo := project.NewRepository(deps.DB)
	projectService := project.NewService(projectRepo)
	projectHandler := project.NewHandler(projectService)

	projectGroup := api.Group("/projects")
	projectGroup.Use(middleware.JWTAuth(deps.Config.JWTSecret))
	{
		projectGroup.POST("", projectHandler.CreateProject)
		projectGroup.GET("", projectHandler.GetAllProjects)
		projectGroup.GET("/:id", projectHandler.GetProject)
		projectGroup.PUT("/:id", projectHandler.UpdateProject)
		projectGroup.DELETE("/:id", projectHandler.DeleteProject)
	}
}

// setupTaskRoutes configures task module routes with dependency injection
func setupTaskRoutes(api *gin.RouterGroup, deps *Dependencies) {
	userRepo := user.NewRepository(deps.DB, deps.Redis)
	userService := user.NewService(userRepo)

	projectRepo := project.NewRepository(deps.DB)
	projectService := project.NewService(projectRepo)

	taskRepo := task.NewRepository(deps.DB)
	taskService := task.NewService(taskRepo, projectService, userService)
	taskHandler := task.NewHandler(taskService)

	taskGroup := api.Group("/tasks")
	taskGroup.Use(middleware.JWTAuth(deps.Config.JWTSecret))
	{
		taskGroup.POST("", taskHandler.CreateTask)
		taskGroup.GET("", taskHandler.GetTasksWithFilters)
		taskGroup.GET("/:id", taskHandler.GetTask)
		taskGroup.PUT("/:id", taskHandler.UpdateTask)
		taskGroup.PUT("/:id/assign", taskHandler.AssignTask)
		taskGroup.DELETE("/:id", taskHandler.DeleteTask)
	}
}
