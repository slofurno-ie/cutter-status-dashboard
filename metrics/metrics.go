package metrics

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

type Metrics struct {
	client          *monitoring.QueryClient
	cluster         string
	project         string
	nodeMemoryQuery string
	nodeCpuQuery    string
}

type NodeMetrics struct {
	MemoryUsage float64
	CpuUsage    float64
}

type NodeMetric struct {
	Node   string
	Values []float64
}

type Nodes struct {
	Memory []NodeMetric
	Cpu    []NodeMetric
}

func (n *Nodes) Healthy() bool {
	for _, cpu := range n.Cpu {
		ten := avg(cpu.Values[:min(len(cpu.Values), 10)])
		last := cpu.Values[0]
		if ten > 0.9 && last > 0.9 {
			return false
		}
	}

	for _, mem := range n.Memory {
		if mem.Values[0] > 0.9 {
			return false
		}
	}

	return true
}

func avg(xs []float64) float64 {
	var sum float64
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func IsHealthy(nm map[string]*NodeMetrics) bool {
	var maxCpu float64
	var maxMem float64

	for _, node := range nm {
		maxCpu = max(node.CpuUsage, maxCpu)
		maxMem = max(node.MemoryUsage, maxMem)
	}

	if maxCpu > .9 {
		return false
	}

	if maxMem > .9 {
		return false
	}

	return true
}

func New(project string, cluster string) (*Metrics, error) {
	ctx := context.Background()
	client, err := monitoring.NewQueryClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		client:          client,
		cluster:         cluster,
		project:         fmt.Sprintf("projects/%s", project),
		nodeMemoryQuery: fmt.Sprintf(nodeMemoryQuery, cluster),
		nodeCpuQuery:    fmt.Sprintf(nodeCpuQuery, cluster),
	}, nil
}

type Metric struct {
	Name  string
	Time  time.Time
	Value float64
}

var nodeMemoryQuery string = `
fetch k8s_node
| metric 'kubernetes.io/node/memory/allocatable_utilization'
| filter (resource.cluster_name == '%s')
| within 30m
| group_by 1m,
    [value_allocatable_utilization_mean: mean(value.allocatable_utilization)]
	| every 1m
	| group_by [resource.node_name],
	    [value_allocatable_utilization_mean_aggregate:
		       aggregate(value_allocatable_utilization_mean)]
`

var nodeCpuQuery string = `
fetch k8s_node
| metric 'kubernetes.io/node/cpu/allocatable_utilization'
| filter (resource.cluster_name == '%s')
| within 30m
| group_by 1m,
    [value_allocatable_utilization_mean: mean(value.allocatable_utilization)]
	| every 1m
	| group_by [resource.node_name],
	    [value_allocatable_utilization_mean_aggregate:
		       aggregate(value_allocatable_utilization_mean)]
`

func (m *Metrics) GetNodeMetrics(ctx context.Context) (*Nodes, error) {

	nodeMem, err := m.query(ctx, m.nodeMemoryQuery)
	if err != nil {
		return nil, err
	}

	nodeCpu, err := m.query(ctx, m.nodeCpuQuery)
	if err != nil {
		return nil, err
	}

	return &Nodes{
		Memory: nodeMem,
		Cpu:    nodeCpu,
	}, nil
}

func (m *Metrics) query(ctx context.Context, query string) ([]NodeMetric, error) {
	req := &monitoringpb.QueryTimeSeriesRequest{
		Name:  m.project,
		Query: query,
	}

	it := m.client.QueryTimeSeries(ctx, req)
	ret := []NodeMetric{}

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(resp.PointData) == 0 {
			continue
		}

		node := resp.LabelValues[0].GetStringValue()
		var values []float64

		for _, p := range resp.PointData {
			values = append(values, p.Values[0].GetDoubleValue())
		}

		ret = append(ret, NodeMetric{
			Node:   node,
			Values: values,
		})
	}

	return ret, nil
}
