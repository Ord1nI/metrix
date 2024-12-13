package grpcagent

import (
	"context"

	"github.com/Ord1nI/metrix/internal/agent"
	pb "github.com/Ord1nI/metrix/internal/proto"
	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcAgent struct {
	*agent.Agent
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

	conn, err := grpc.Dial(agent.Config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	agent.Client = pb.NewMetrixServerClient(conn)

	agent.Logger.Infoln("Agent inited successfuly")

	return &agent, nil
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

	res, err := g.Client.SendMetric(context.Background(), gMetric)
	g.Logger.Infoln(res.GetErr)

	return err
}


func (g *GrpcAgent) Run() chan struct{} {
	end := make(chan struct{})
	// if (g.Config.PublicKeyFile != "") {
	// 	g.StartWorkers(g.TaskPoll(end, g.StartMetricCollector(end)), g.SendMetricJSONwithEncryption(g.Config.PublicKeyFile))
	// } else {
	// 	a.StartWorkers(g.TaskPoll(end, g.StartMetricCollector(end)), g.SendMetricJSON)
	// }
	g.StartWorkers(g.TaskPoll(end, g.StartMetricCollector(end)), g.SendMetric)

	return end
}
