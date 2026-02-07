package grpc

import (
	"io"
	"log"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/service"
	pb "github.com/MohamedDenta/Drone-Delivery-Management-Backend/proto/drone"
	"google.golang.org/grpc"
)

type DroneServer struct {
	pb.UnimplementedDroneServiceServer
	droneService *service.DroneService
}

func NewDroneServer(droneService *service.DroneService) *DroneServer {
	return &DroneServer{droneService: droneService}
}

func (s *DroneServer) ReportLocation(stream pb.DroneService_ReportLocationServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Update Location in Service (which handles DB + Redis + Observers)
		err = s.droneService.UpdateLocation(req.DroneId, req.Latitude, req.Longitude)
		if err != nil {
			log.Printf("Failed to update location for drone %s: %v", req.DroneId, err)
			continue
		}

		// Optional: Send Ack back
		/*
			if err := stream.Send(&pb.LocationResponse{Message: "Ack"}); err != nil {
				return err
			}
		*/
	}
}

func Register(s *grpc.Server, srv *DroneServer) {
	pb.RegisterDroneServiceServer(s, srv)
}
