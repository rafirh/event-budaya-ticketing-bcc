package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"event-budaya-ticketing-bcc/config"
	"event-budaya-ticketing-bcc/internal/handler"
	gormRepo "event-budaya-ticketing-bcc/internal/repository/gorm"
	"event-budaya-ticketing-bcc/internal/router"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/migrations"
	"event-budaya-ticketing-bcc/pkg/email"
	"event-budaya-ticketing-bcc/pkg/payment"
	"event-budaya-ticketing-bcc/pkg/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	seedFlag := flag.Bool("seed", false, "Run database seeders")
	freshFlag := flag.Bool("fresh", false, "Drop all tables and re-run migrations")

	flag.Parse()
	config.LoadConfig()
	config.InitDatabase()

	if *freshFlag {
		migrations.Fresh(config.DB)
		log.Println("Database refreshed successfully")
		return
	}

	if *migrateFlag {
		migrations.Migrate(config.DB)
		log.Println("Migrations completed successfully")
		return
	}

	if *seedFlag {
		migrations.Seed(config.DB)
		log.Println("Seeding completed successfully")
		return
	}

	userRepo := gormRepo.NewUserRepository(config.DB)
	tokenRepo := gormRepo.NewPersonalAccessTokenRepository(config.DB)
	emailVerificationRepo := gormRepo.NewEmailVerificationTokenRepository(config.DB)
	categoryRepo := gormRepo.NewEventCategoryRepository(config.DB)
	eventRepo := gormRepo.NewEventRepository(config.DB)
	eventCommentRepo := gormRepo.NewEventCommentRepository(config.DB)
	eventCreationPaymentRepo := gormRepo.NewEventCreationPaymentRepository(config.DB)
	orderRepo := gormRepo.NewOrderRepository(config.DB)
	ticketRepo := gormRepo.NewTicketRepository(config.DB)
	paymentRepo := gormRepo.NewPaymentRepository(config.DB)
	walletRepo := gormRepo.NewPromoterWalletRepository(config.DB)
	transactionRepo := gormRepo.NewWalletTransactionRepository(config.DB)
	feeRepo := gormRepo.NewFeeRepository(config.DB)
	adminWalletRepo := gormRepo.NewAdminWalletRepository(config.DB)
	promoterTransactionRepo := gormRepo.NewPromoterTransactionHistoryRepository(config.DB)
	midtransClient := payment.NewMidtransClient(config.AppConfig.MidtransServer, config.AppConfig.MidtransEnv)

	var uploader storage.Uploader
	mailSender := email.NewSMTPSender(email.SMTPConfig{
		Host:       config.AppConfig.MailHost,
		Port:       config.AppConfig.MailPort,
		Username:   config.AppConfig.MailUsername,
		Password:   config.AppConfig.MailPassword,
		Encryption: config.AppConfig.MailEncryption,
		FromAddr:   config.AppConfig.MailFromAddress,
		FromName:   config.AppConfig.MailFromName,
	})

	if config.AppConfig.S3Bucket != "" && config.AppConfig.S3Region != "" {
		s3Uploader, err := storage.NewS3Uploader(storage.S3Config{
			Region:        config.AppConfig.S3Region,
			Bucket:        config.AppConfig.S3Bucket,
			AccessKey:     config.AppConfig.S3Key,
			SecretKey:     config.AppConfig.S3Secret,
			PublicBaseURL: config.AppConfig.S3PublicBase,
		})
		if err != nil {
			log.Printf("Warning: failed to initialize S3 uploader: %v", err)
		} else {
			uploader = s3Uploader
		}
	}

	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo, emailVerificationRepo, mailSender, config.AppConfig.AppURL, uploader)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	eventUsecase := usecase.NewEventUsecase(eventRepo, eventCreationPaymentRepo, feeRepo, adminWalletRepo, promoterTransactionRepo, midtransClient)
	eventCommentUsecase := usecase.NewEventCommentUsecase(eventRepo, eventCommentRepo)
	orderUsecase := usecase.NewOrderUsecase(userRepo, eventRepo, orderRepo, ticketRepo, paymentRepo, walletRepo, transactionRepo, midtransClient, config.AppConfig.MidtransServer)
	ticketUsecase := usecase.NewTicketUsecase(ticketRepo)
	walletUsecase := usecase.NewWalletUsecase(walletRepo)
	authHandler := handler.NewAuthHandler(authUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	eventHandler := handler.NewEventHandler(eventUsecase)
	eventCommentHandler := handler.NewEventCommentHandler(eventCommentUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	ticketHandler := handler.NewTicketHandler(ticketUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase)

	app := fiber.New(fiber.Config{
		AppName:      config.AppConfig.AppName,
		ErrorHandler: customErrorHandler,
	})

	router.SetupRoutes(app, authHandler, categoryHandler, eventHandler, orderHandler, ticketHandler, eventCommentHandler, walletHandler, tokenRepo)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	port := config.AppConfig.AppPort
	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s", config.AppConfig.AppEnv)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
	})
}
