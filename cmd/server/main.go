package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/api"
	grpcHandler "github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/api/grpc"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/api/handlers"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/config"
	infra_rmq "github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/infrastructure/rabbitmq"
	infra_redis "github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/infrastructure/redis"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/service"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/telemetry"
	"google.golang.org/grpc"
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
	// 4. Init Redis
	redisClient, err := infra_redis.NewClient(cfg.RedisURL, "", 0)
	if err != nil {
		log.Printf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	// 5. Init RabbitMQ
	rabbitClient, err := infra_rmq.NewClient(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("failed to connect to rabbitmq: %v", err)
	}
	defer rabbitClient.Close()

	// 6. Init Services
	droneService := service.NewDroneService(repo, redisClient)
	orderService := service.NewOrderService(repo, rabbitClient)
	dispatcherService := service.NewDispatcherService(repo, repo)

	// Worker (Async)
	orderWorker := service.NewOrderDispatcherWorker(rabbitClient, dispatcherService, repo)
	if err := orderWorker.Start(); err != nil {
		log.Printf("Failed to start order worker: %v", err)
	}

	// Recovery Handler (Observer)
	recoveryHandler := service.NewRecoveryHandler(repo)
	droneService.AddObserver(recoveryHandler)

	// Heartbeat Monitor (Async)
	heartbeatMonitor := service.NewHeartbeatMonitor(repo, droneService, redisClient)
	go heartbeatMonitor.Start(ctx)

	// 5. Init Handlers
	droneHandler := handlers.NewDroneHandler(droneService, dispatcherService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// 6. Init Router
	r := api.SetupRouter(droneHandler, orderHandler)

	// 7. Start servers
	// HTTP Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("HTTP Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// gRPC Server
	lis, err := net.Listen("tcp", ":50051") // TODO: Move port to config
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	droneGrpcServer := grpcHandler.NewDroneServer(droneService)
	grpcHandler.Register(grpcServer, droneGrpcServer) // Ensure Register function exists in grpc package

	go func() {
		log.Printf("gRPC Server starting on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
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
