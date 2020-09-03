package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GRPConfig struct {
	Host string `yaml:"host"`
}

type GRPCServer struct {
	host     string
	listener net.Listener
	grpc     *grpc.Server
	services []GRPCService
	tracer   trace.Tracer
	shutdown chan os.Signal
	logger   *logrus.Logger
}

func NewServer(opts ...GRPCServerOption) *GRPCServer {
	server := new(GRPCServer)
	server.defaults()

	for _, o := range opts {
		o(server)
	}

	server.creteServer()

	for _, service := range server.services {
		service.RegisterService(server.grpc)
	}

	return server
}

func (s *GRPCServer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ln, err := net.Listen("tcp", s.host)

	if err != nil {
		return err
	}

	s.listener = ln
	defer s.listener.Close()

	signal.Notify(s.shutdown, os.Interrupt)

	go func() {
		select {
		case sig := <-s.shutdown:
			s.logger.WithField("sig", sig.String()).
				Info("starting server graceful shutdown")
			s.grpc.GracefulStop()
		case <-ctx.Done():
		}
	}()

	s.logger.WithField("addr", s.listener.Addr().String()).
		Info("starting grpc server")

	if err := s.grpc.Serve(s.listener); err != nil {
		return err
	}

	return nil
}

func (s *GRPCServer) creteServer() {
	entry := logrus.NewEntry(s.logger)

	s.grpc = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(entry),
			grpctrace.UnaryServerInterceptor(s.tracer),
		),

		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(entry),
			grpctrace.StreamServerInterceptor(s.tracer),
		),

		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Second * 10,
			PermitWithoutStream: true,
		}),

		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Second * 30,
			Timeout: time.Second * 10,
		}),
	)
}

func (s *GRPCServer) defaults() {
	s.logger = logrus.New()
	s.tracer = global.Tracer("grpc")
	s.services = []GRPCService{}
	s.shutdown = make(chan os.Signal)
}
