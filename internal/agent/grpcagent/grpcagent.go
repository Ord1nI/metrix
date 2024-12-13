package grpcagent

import (
	"context"

	"github.com/Ord1nI/metrix/internal/agent"
	"github.com/Ord1nI/metrix/internal/agent/grpcagent/interceptors"
	pb "github.com/Ord1nI/metrix/internal/proto"
	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GrpcAgent struct {
	*agent.Agent
	Interceptors []grpc.UnaryClientInterceptor
	Client pb.MetrixServerClient
}

func New() (*GrpcAgent, error) {
	mAgent, err := agent.New()

	if err != nil {
		return nil, err
	}

	agent := GrpcAgent{
		Agent: mAgent,
	}

	if agent.Config.Key != "" {
		agent.Add(interceptors.SignInterceptor(agent.Logger, []byte(agent.Config.Key)))
	}
	if agent.Config.IP != "" {
		agent.Add(interceptors.AddIPInterceptro(agent.Logger, agent.Config.IP))
	}

	conn, err := grpc.Dial(agent.Config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(agent.Interceptors...))

	if err != nil {
		return nil, err
	}

	agent.Client = pb.NewMetrixServerClient(conn)

	agent.Logger.Infoln("Agent inited successfuly")

	return &agent, nil
}

func (g *GrpcAgent) Add(i ...grpc.UnaryClientInterceptor) {
	g.Interceptors = append(g.Interceptors, i...)
}

func (g *GrpcAgent) SendMetric(m metrics.Metric) error{
	var gMetric *pb.Metric
	switch m.MType {
	case "gauge":
		gMetric = &pb.Metric{
			ID: m.ID,
			Value: *m.Value,
			MType: m.MType,
		}
	case "counter":
		gMetric = &pb.Metric{
			ID: m.ID,
			Delta: *m.Delta,
			MType: m.MType,
		}
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{}))
	res, err := g.Client.SendMetric(ctx, gMetric)
	g.Logger.Infoln(res.GetErr)

	return err
}


func (g *GrpcAgent) Run() chan struct{} {
	end := make(chan struct{})
	g.StartWorkers(g.TaskPoll(end, g.StartMetricCollector(end)), g.SendMetric)

	return end
}
