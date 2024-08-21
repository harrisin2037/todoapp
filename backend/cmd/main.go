package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/handlers"
	"github.com/harrisin2037/todoapp/internal/middlewares"
	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/repository"
	"github.com/harrisin2037/todoapp/internal/service"
)

func main() {
	var (
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASS")
		dbHost = os.Getenv("DB_HOST")
		dbPort = os.Getenv("DB_PORT")
		dbName = os.Getenv("DB_NAME")
		dsn    = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.Todo{}, &models.User{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	weekLater := time.Now().AddDate(0, 0, 7)

	_ = db.Create(&models.Todo{
		Name:        "Create a new todo",
		Description: "This is a new todo",
		DueDate:     &weekLater,
	})

	admin := &models.User{
		Username: "admin",
		Email:    "admin",
		Role:     models.RoleAdmin,
	}
	admin.SetPassword("admin123")
	_ = db.Create(admin)

	var (
		todoRepo    = repository.NewTodoRepository(db)
		todoService = service.NewTodoService(todoRepo)
		todoHandler = handlers.NewTodoHandler(todoService)
		userRepo    = repository.NewUserRepository(db)
		userService = service.NewUserService(userRepo)
		userHandler = handlers.NewUserHandler(userService)
		router      = gin.Default()
	)

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	// userRouter := router.Group("/")
	// userRouter.Use(middlewares.AuthMiddleware())
	// {
	// 	userRouter.POST("/todos", todoHandler.CreateTodo)
	// 	userRouter.GET("/todos", todoHandler.GetTodos)
	// 	userRouter.GET("/todos/:id", todoHandler.GetTodo)
	// 	userRouter.PUT("/todos/:id", todoHandler.UpdateTodo)
	// 	userRouter.DELETE("/todos/:id", todoHandler.DeleteTodo)
	// }

	router.POST("/todos", todoHandler.CreateTodo)
	router.GET("/todos", todoHandler.GetTodos)
	router.GET("/todos/:id", todoHandler.GetTodo)
	router.PUT("/todos/:id", todoHandler.UpdateTodo)
	router.DELETE("/todos/:id", todoHandler.DeleteTodo)

	adminRouter := router.Group("/admin")
	adminRouter.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		adminRouter.PUT("/role", userHandler.ChangeRole)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
