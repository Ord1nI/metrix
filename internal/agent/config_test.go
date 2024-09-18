package agent

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
    type want struct{
        Address string
        PollInterval int64 
        ReportInterval int64
        Key string
        RateLimit int
    }
    tests := []struct{
        want want
        Address string
        PollInterval string
        ReportInterval string
        Key string
        RateLimit string
    }{
        {
            want: want{
                Address: "10.5",
                PollInterval:15,
                ReportInterval:15,
                Key: "slkfj",
                RateLimit: 20,
            },
            Address:"10.5",
            PollInterval:"15",
            ReportInterval:"15",
            Key: "slkfj",
            RateLimit: "20",
        },
    }

    for i, test := range tests{
        t.Run(fmt.Sprintf("Test %d",i),func(t *testing.T){
            if test.Address != ""{
                os.Setenv("ADDRESS", test.Address)
            }
            if test.PollInterval != "" {
                os.Setenv("POLL_INTERVAL", test.PollInterval)
            }
            if test.ReportInterval != "" {
                os.Setenv("REPORT_INTERVAL", test.ReportInterval)
            }
            if test.Key != "" {
                os.Setenv("KEY",test.Key)
            }
            if test.RateLimit != "" {
                os.Setenv("RATE_LIMIT", test.RateLimit)
            }

            agent, _ := New()
            assert.Equal(t, test.want.Address, agent.Config.Address)
            assert.Equal(t, test.want.PollInterval, agent.Config.PollInterval)
            assert.Equal(t, test.want.ReportInterval, agent.Config.ReportInterval)
            assert.Equal(t, test.want.Key, agent.Config.Key)
            assert.Equal(t, test.want.RateLimit, agent.Config.RateLimit)
        })
    }
}
