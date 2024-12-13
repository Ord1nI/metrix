package grpcserv

import (
	"context"
	"net"

	"github.com/Ord1nI/metrix/internal/server"
	"github.com/Ord1nI/metrix/internal/server/grpcserv/interceptors"
	"google.golang.org/grpc"

	"github.com/Ord1nI/metrix/internal/repo/metrics"

	pb "github.com/Ord1nI/metrix/internal/proto"
)

type GrpcServer struct {
	*server.Server
	pb.UnimplementedMetrixServerServer
	GServer *grpc.Server
	Interceptors []grpc.UnaryServerInterceptor
}

func Default() (*GrpcServer, error) {
	serv, err := new()
	if err != nil {
		return nil, err;
	}
	serv.Add(interceptors.LoggerInterceptor(serv.Logger))

	if serv.Config.Key != "" {
		serv.Add(interceptors.SignInterceptor(serv.Logger, []byte(serv.Config.Key)))
	}

	if serv.Config.TrustedSubnet != "" {
		serv.Add(interceptors.CheckSubnetInterceptor(serv.Logger, net.ParseIP(serv.Config.TrustedSubnet)))
	}

	serv.GServer = grpc.NewServer(grpc.ChainUnaryInterceptor(serv.Interceptors...))

	return serv, nil
}


func new() (*GrpcServer, error) {
	mServer, err := server.New()

	if err != nil {
		return nil, err
	}

	s := &GrpcServer{
		Server: mServer,
	}


	return s, nil

}

func (g *GrpcServer) Add(i ...grpc.UnaryServerInterceptor) {
	g.Interceptors = append(g.Interceptors, i...)
}

func (g *GrpcServer) Run(<-chan struct{}) error {
	listen, err := net.Listen("tcp", g.Config.Address)
	if err != nil {
        g.Logger.Fatal(err)
    }


	pb.RegisterMetrixServerServer(g.GServer, g)

	if err := g.GServer.Serve(listen); err != nil {
        g.Logger.Fatal(err)
    }

	return nil
}

func (g *GrpcServer) SendMetric(ctx context.Context, in *pb.Metric) (*pb.Error, error) {

	m := metrics.Metric{
		ID: in.ID,
		Delta: &in.Delta,
		Value: &in.Value,
		MType: in.MType,
	}

	err := g.Repo.Add(m.ID, m)

	if err != nil {
		return &pb.Error{Err:"error adding metric"}, err;
	}

	return &pb.Error{Err:"All good"}, nil;

}

func (g *GrpcServer) GetMetric(ctx context.Context, in *pb.MetricName) (*pb.Metric, error) {
	var m metrics.Metric
	err := g.Repo.Get(in.ID, &m)

	if err != nil {
		return nil, err;
	}

	return nil, nil
}

func (g *GrpcServer) SendMetrics(ctx context.Context, in *pb.Metrics) (*pb.Error, error) {
	for _, i := range in.Metrics {

		m := metrics.Metric{
			ID: i.ID,
			Delta: &i.Delta,
			Value: &i.Value,
			MType: i.MType,
		}

		err := g.Repo.Add(m.ID, m)

		if err != nil {
			return &pb.Error{Err:"error adding metric"}, err;
		}
	}

	return &pb.Error{Err:"All good"}, nil;
}
