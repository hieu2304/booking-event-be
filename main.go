package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-booking-be/internal/config"
	"event-booking-be/internal/handler"
	"event-booking-be/internal/models"
	"event-booking-be/internal/repository"
	"event-booking-be/internal/routes"
	"event-booking-be/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize GORM database
	db, err := initGormDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient, err := initRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	eventRepo := repository.NewEventRepository(db)
	userRepo := repository.NewUserRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// Initialize services
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, db, cfg.BookingTimeoutMinutes)

	// Initialize handlers
	eventHandler := handler.NewEventHandler(eventService)
	userHandler := handler.NewUserHandler(userService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Initialize router
	router := routes.NewRouter(userHandler, eventHandler, bookingHandler)

	app := fiber.New(fiber.Config{
		AppName:      "Event Booking API",
		ErrorHandler: customErrorHandler,
	})

	setupMiddlewares(app)
	
	app.Get("/health", healthCheckHandler(db, redisClient))
	
	router.Setup(app)

	// Start background worker for expired bookings
	go startBookingWorker(bookingService)

	go gracefulShutdown(app, db, redisClient)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on http://localhost%s", addr)
	
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initGormDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	if err := db.AutoMigrate(&models.Event{}, &models.User{}, &models.Booking{}); err != nil {
		return nil, err
	}

	log.Println("Database connected (GORM)")
	return db, nil
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Redis connected")
	return client, nil
}

func setupMiddlewares(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}

func startBookingWorker(bookingService service.BookingService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Booking expiration worker started")

	for range ticker.C {
		ctx := context.Background()
		if err := bookingService.ProcessExpiredBookings(ctx); err != nil {
			log.Printf("Error processing expired bookings: %v", err)
		}
	}
}

func healthCheckHandler(db *gorm.DB, redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		dbStatus := "ok"
		if err := db.PingContext(ctx); err != nil {
			dbStatus = "error"
		}

		redisStatus := "ok"
		if err := redisClient.Ping(ctx).Err(); err != nil {
			redisStatus = "error"
		}

		return c.JSON(fiber.Map{
			"status": "ok",
			"database": dbStatus,
			"redis": redisStatus,
		})
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("Error: %v | %s %s", err, c.Method(), c.Path())

	return c.Status(code).JSON(fiber.Map{
		"error":   message,
		"status":  code,
	})
}

func gracefulShutdown(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down...")

	app.Shutdown()
	db.Close()
	redisClient.Close()

	os.Exit(0)
}