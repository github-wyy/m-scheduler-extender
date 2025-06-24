package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	schedulerextapi "k8s.io/kube-scheduler/extender/v1"
)

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	klog.Info("filterHandler")

	var args schedulerextapi.ExtenderArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	klog.Infof("doFilter: pod: %v", args.Pod.Name)
	filteredNodes := doFilter(args.Pod, args.Nodes)
	result := &schedulerextapi.ExtenderFilterResult{
		Nodes: filteredNodes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func scoreHandler(w http.ResponseWriter, r *http.Request) {
	klog.Info("scoreHandler")

	var args schedulerextapi.ExtenderArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	klog.Infof("doScore: pod: %v", args.Pod.Name)
	priorityList := make(schedulerextapi.HostPriorityList, 0, len(args.Nodes.Items))
	for _, node := range args.Nodes.Items {
		priorityList = append(priorityList, schedulerextapi.HostPriority{
			Host:  node.Name,
			Score: 100,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(priorityList); err != nil {
		klog.Errorf("Failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func doFilter(pod *v1.Pod, nodes *v1.NodeList) *v1.NodeList {
	var filteredNodes []v1.Node
	for _, node := range nodes.Items {
		filteredNodes = append(filteredNodes, node)
	}
	return &v1.NodeList{Items: filteredNodes}
}
