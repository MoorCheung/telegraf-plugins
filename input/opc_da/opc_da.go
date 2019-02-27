package opc_da

// opc_da.go

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/konimarti/opc"
)

type OpcDa struct {
	Server string
	Host   string
	Items  map[string]string
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
  [[inputs.simple.items]]
    test001="numeric.sin.double"
    test002="numeric.sin.int16"
    test003="numeric.sin.float"
    test004="numeric.square.double"
    test005="numeric.square.double"
`
}

func (s *OpcDa) Gather(acc telegraf.Accumulator) error {
	if len(s.Items) == 0 {
		return nil
	}

	if s.client == nil {
		var points []string
		for v := range s.Items {
			points = append(points, s.Items[v])
		}
		items := RemoveRepeatedElement(points)
		s.client = opc.NewConnection(
			s.Server,         // ProgId
			[]string{s.Host}, //  OPC servers nodes
			items,            // slice of OPC tags
		)
	}

	temp := s.client.Read()
	result := ConvertMap(temp)
	for k, code := range s.Items {
		if _, ok := result[code]; ok {
			tags := map[string]string{
				"tagCode": k,
				"quality": "good",
			}
			fields := map[string]interface{}{
				"value": result[code],
			}
			acc.AddFields("live", fields, tags)
		}
	}

	return nil
}

func ConvertMap(result map[string]interface{}) (newMap map[string]float64) {
	newMap=make(map[string]float64)
	for key := range result {
		switch result[key].(type) {
		case int:
			newMap[key] = float64(result[key].(int))
			break
		case int16:
			newMap[key] = float64(result[key].(int16))
			break
		case int32:
			newMap[key] = float64(result[key].(int32))
			break
		case float32:
			newMap[key] = float64(result[key].(float32))
			break
		case float64:
			newMap[key] = result[key].(float64)
			break
		default:
			break
		}
	}
	return newMap
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func init() {
	inputs.Add("opc_da", func() telegraf.Input { return &OpcDa{} })
}
