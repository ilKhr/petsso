package health

import (
	"context"

	healthv1 "github.com/ilKhr/petprotos/gen/go/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Health interface {
	CheckHealth(ctx context.Context)
	WatchHealth()
}

type serverApi struct {
	healthv1.UnimplementedHealthServer
	health Health
}

func (s *serverApi) Check(ctx context.Context, in *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
	s.health.CheckHealth(ctx)
	return &healthv1.HealthCheckResponse{Status: healthv1.HealthCheckResponse_SERVING}, nil
}

func (s *serverApi) Watch(in *healthv1.HealthCheckRequest, _ healthv1.Health_WatchServer) error {
	s.health.WatchHealth()
	return status.Error(codes.Unimplemented, "unimplemented")
}

func Register(gRPC *grpc.Server, h Health) {
	healthv1.RegisterHealthServer(gRPC, &serverApi{health: h})
}
