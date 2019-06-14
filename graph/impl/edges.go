package impl

import (
	"Capacitated/graph/types"
)

func NewEdge(client types.Client, server types.Server, edge_id int) types.Edge {
	return &edge{client, server, edge_id}
}

type edge struct {
	client types.Client
	server types.Server
	edge_id    int
}

func (e *edge) Client() types.Client {
	return e.client
}

func (e *edge) Server() types.Server {
	return e.server
}

func (e *edge) Edge() int {
	return e.edge_id
}
