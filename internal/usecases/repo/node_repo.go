package repo

import (
	"fmt"
	"io"
	"sync"
)

type NodeRepo struct {
	nodes map[string]io.ReadWriter
	mutex *sync.Mutex
}

func NewNodeRepo() *NodeRepo {
	return &NodeRepo{
		nodes: make(map[string]io.ReadWriter),
		mutex: &sync.Mutex{},
	}
}

func (r *NodeRepo) Get(id string) (io.ReadWriter, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if node, ok := r.nodes[id]; ok {
		return node, nil
	}
	return nil, fmt.Errorf("node not found")
}

func (r *NodeRepo) Set(node NodeInterface) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.nodes[node.GetId()] = node.GetConn()
}

func (r *NodeRepo) List() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var nodes []string
	for k := range r.nodes {
		nodes = append(nodes, k)
	}
	return nodes
}

func (r *NodeRepo) Remove(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.nodes[id]; !ok {
		return fmt.Errorf("node not found")
	}
	delete(r.nodes, id)
	return nil
}
