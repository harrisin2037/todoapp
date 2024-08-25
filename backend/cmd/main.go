package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/handlers"
	"github.com/harrisin2037/todoapp/internal/middlewares"
	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/repository"
	"github.com/harrisin2037/todoapp/internal/service"
	"github.com/harrisin2037/todoapp/internal/websocket"
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

	admin := &models.User{
		Username: "admin",
		Email:    "admin",
		Role:     models.RoleAdmin,
	}
	admin.SetPassword("admin123")
	_ = db.Create(admin)

	hub := websocket.NewHub()
	go hub.Run()

	var (
		router      = gin.Default()
		userRepo    = repository.NewUserRepository(db)
		userService = service.NewUserService(userRepo)
		todoRepo    = repository.NewTodoRepository(db)
		todoService = service.NewTodoService(todoRepo)
		userHandler = handlers.NewUserHandler(userService)
		todoHandler = handlers.NewTodoHandler(todoService, userService, hub)
	)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	userRouter := router.Group("/")
	userRouter.Use(middlewares.AuthMiddleware())
	{
		userRouter.GET("/roles", userHandler.CheckRoles)
		userRouter.POST("/todos", todoHandler.CreateTodo)
		userRouter.GET("/todos", todoHandler.GetTodos)
		userRouter.GET("/todos/:id", todoHandler.GetTodo)
		userRouter.PUT("/todos/:id", todoHandler.UpdateTodo)
		userRouter.DELETE("/todos/:id", todoHandler.DeleteTodo)

		router.GET("/users", userHandler.GetAllUsers)
	}

	adminRouter := router.Group("/admin")
	adminRouter.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		adminRouter.PUT("/role", userHandler.ChangeRole)
		adminRouter.PUT("/users/:id", userHandler.UpdateUser)
		adminRouter.GET("/users", userHandler.GetAllUsers)
		adminRouter.POST("/users", userHandler.CreateUser)
		adminRouter.DELETE("/users/:id", userHandler.DeleteUser)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.GET("/ws", func(c *gin.Context) {
		hub.HandleWebSocket(c)
	})

	// online status websocket
	router.GET("/online", func(c *gin.Context) {
		// hub.HandleOnlineStatus(c)
	})

	go func() {
		hub.Broadcast <- []byte(`{"message": "Welcome to the chat room!"}`)
	}()

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
