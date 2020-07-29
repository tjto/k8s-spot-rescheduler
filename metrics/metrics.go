/*
Copyright 2017 Pusher Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s-spot-rescheduler/nodes"
)

const (
	reschedulerNamespace = "spot_rescheduler"
)

var (
	// nodePodsCount tracks how many pods are nodes by type and by node name.
	nodePodsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: reschedulerNamespace,
			Name:      "node_pods_count",
			Help:      "Number of pods on each node.",
		},
		[]string{"node_type", "node"})

	// nodesCount tracks the number of nodes in the cluster.
	nodesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: reschedulerNamespace,
			Name:      "nodes_count",
			Help:      "Number of nodes in cluster.",
		}, []string{"node_type"},
	)

	// nodeDrainCount counts the number of nodes drained by the rescheduler.
	nodeDrainCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: reschedulerNamespace,
			Name:      "node_drain_total",
			Help:      "Number of nodes drained by rescheduler.",
		}, []string{"drain_state", "node"},
	)

	// evictionsCount counts the number of pods evicted by the rescheduler
	evictionsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: reschedulerNamespace,
			Name:      "evicted_pods_total",
			Help:      "Number of pods evicted by the rescheduler.",
		},
	)
)

func init() {
	prometheus.MustRegister(nodePodsCount)
	prometheus.MustRegister(nodesCount)
	prometheus.MustRegister(nodeDrainCount)
	prometheus.MustRegister(evictionsCount)
}

// UpdateNodesMap updates the metrics calculated by the nodes map
func UpdateNodesMap(nm nodes.Map) {
	if nm == nil {
		return
	}
	nodesCount.WithLabelValues(nodes.OnDemandNodeLabel).Set(float64(len(nm[nodes.OnDemand])))
	nodesCount.WithLabelValues(nodes.SpotNodeLabel).Set(float64(len(nm[nodes.Spot])))

}

// UpdateNodePodsCount updates nodePodsCount for a given node
func UpdateNodePodsCount(nodeType string, nodeName string, numPods int) {
	nodePodsCount.WithLabelValues(nodeType, nodeName).Set(float64(numPods))
}

// UpdateEvictionsCount adds 1 to the evictions counter
func UpdateEvictionsCount() {
	evictionsCount.Add(1)
}

// UpdateNodeDrainCount updates the number drains and drain state for a node
func UpdateNodeDrainCount(state string, nodeName string) {
	nodeDrainCount.WithLabelValues(state, nodeName).Add(1)
}
