package health

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	grpc_health_v1.UnimplementedHealthServer
	status grpc_health_v1.HealthCheckResponse_ServingStatus
}

func NewServer() *Server {
	return &Server{
		status: grpc_health_v1.HealthCheckResponse_SERVING,
	}
}

func (s *Server) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: s.status}, nil
}

func (s *Server) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (s *Server) SetServingStatus(_ string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	s.status = status
}

func (s *Server) Shutdown() {
	s.status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
}
