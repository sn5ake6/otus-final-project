//go:generate protoc AntiBruteForceService.proto --go_out=./pb/ --go-grpc_out=./pb/ --proto_path=../../../api

package internalgrpc

import (
	context "context"
	"fmt"
	"net"
	"time"

	"github.com/sn5ake6/otus-final-project/internal/server/grpc/pb"
	"github.com/sn5ake6/otus-final-project/internal/storage"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAntiBruteForceServiceServer
	addr       string
	app        Application
	logger     Logger
	grpcServer *grpc.Server
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
	LogGRPCRequest(r interface{}, method string, requestDuration time.Duration)
}

type Application interface {
	Authorize(ctx context.Context, authorize storage.Authorize) (bool, error)
	Reset(ctx context.Context, authorize storage.Authorize)
	AddToBlacklist(ctx context.Context, subnet string) error
	DeleteFromBlacklist(ctx context.Context, subnet string) error
	FindIPInBlacklist(ctx context.Context, ip string) (bool, error)
	AddToWhitelist(ctx context.Context, subnet string) error
	DeleteFromWhitelist(ctx context.Context, subnet string) error
	FindIPInWhitelist(ctx context.Context, ip string) (bool, error)
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{
		addr:   addr,
		app:    app,
		logger: logger,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingMiddleware(logger),
		),
	)

	s.grpcServer = grpcServer

	pb.RegisterAntiBruteForceServiceServer(s.grpcServer, s)

	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("GRPC server started: %s", s.addr))

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() error {
	s.logger.Info(fmt.Sprintf("GRPC server stopped: %s", s.addr))

	s.grpcServer.GracefulStop()

	return nil
}

func (s *Server) Authorize(ctx context.Context, r *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	authorize := storage.NewAuthorize(
		r.GetLogin(),
		r.GetPassword(),
		r.GetIp(),
	)

	res, err := s.app.Authorize(ctx, authorize)
	if err != nil {
		return &pb.AuthorizeResponse{Result: false}, err
	}

	return &pb.AuthorizeResponse{Result: res}, nil
}

func (s *Server) Reset(ctx context.Context, r *pb.ResetRequest) (*pb.ResultResponse, error) {
	authorize := storage.NewAuthorize(
		r.GetLogin(),
		r.GetPassword(),
		r.GetIp(),
	)

	s.app.Reset(ctx, authorize)

	return &pb.ResultResponse{}, nil
}

func (s *Server) AddToBlacklist(ctx context.Context, r *pb.BlacklistRequest) (*pb.ResultResponse, error) {
	if err := s.app.AddToBlacklist(ctx, r.GetSubnet()); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	return &pb.ResultResponse{}, nil
}

func (s *Server) DeleteFromBlacklist(ctx context.Context, r *pb.BlacklistRequest) (*pb.ResultResponse, error) {
	if err := s.app.DeleteFromBlacklist(ctx, r.GetSubnet()); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	return &pb.ResultResponse{}, nil
}

func (s *Server) AddToWhitelist(ctx context.Context, r *pb.WhitelistRequest) (*pb.ResultResponse, error) {
	if err := s.app.AddToWhitelist(ctx, r.GetSubnet()); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	return &pb.ResultResponse{}, nil
}

func (s *Server) DeleteFromWhitelist(ctx context.Context, r *pb.WhitelistRequest) (*pb.ResultResponse, error) {
	if err := s.app.DeleteFromBlacklist(ctx, r.GetSubnet()); err != nil {
		return &pb.ResultResponse{Error: err.Error()}, err
	}

	return &pb.ResultResponse{}, nil
}
