package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/api"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/api/handlers"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/config"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/service"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/telemetry"
)

func main() {
	// 1. Load Config
	cfg := config.Load()

	// 2. Init Telemetry
	telemetry.InitLogger()
	tp, err := telemetry.InitTracer("drone-backend", cfg.OTelCollector)
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	mp, err := telemetry.InitMeter()
	if err != nil {
		log.Fatalf("failed to init meter: %v", err)
	}

	// Graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer telemetry.Shutdown(context.Background(), tp, mp)

	// 3. Init Database
	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()

	// 4. Init Services
	droneService := service.NewDroneService(repo)
	orderService := service.NewOrderService(repo)
	dispatcherService := service.NewDispatcherService(repo, repo)

	// Recovery Handler (Observer)
	recoveryHandler := service.NewRecoveryHandler(repo)
	droneService.AddObserver(recoveryHandler)

	// 5. Init Handlers
	droneHandler := handlers.NewDroneHandler(droneService, dispatcherService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// 6. Init Router
	r := api.SetupRouter(droneHandler, orderHandler)

	// 7. Start Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 8. Wait for Interrupt
	<-ctx.Done()
	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	// 9. Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
