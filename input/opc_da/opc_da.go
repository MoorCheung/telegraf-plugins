package opc_da

// simple.go

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/konimarti/opc"
	"time"
)

type OpcDa struct {
	Server string
	Host   string
	Items  []string
	client opc.Connection
}

func (s *OpcDa) Description() string {
	return "a demo plugin"
}

func (s *OpcDa) SampleConfig() string {
	return `
  ## Indicate if everything is fine
  server = "Graybox.Simulator"
  host = "localhost"
  items = ["numeric.sin.int64"]
`
}

func (s *OpcDa) Gather(acc telegraf.Accumulator) error {
	if len(s.Items) == 0 {
		return nil
	}

	if s.client == nil {
		s.client = opc.NewConnection(
			s.Server,         // ProgId
			[]string{s.Host}, //  OPC servers nodes
			s.Items,          // slice of OPC tags
		)
	}

	result := s.client.Read()
	now := time.Now()
	for key := range result {
		tags := map[string]string{
			"tagCode": key,
		}
		acc.AddFields("opc", map[string]interface{}{"value": result[key]}, tags, now)
	}

	return nil
}

func init() {
	inputs.Add("simple", func() telegraf.Input { return &OpcDa{} })
}
