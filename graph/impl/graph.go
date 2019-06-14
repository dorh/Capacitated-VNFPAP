package impl

import (
	"Capacitated/graph/types"
)

func NewGraph(edges []types.Edge) types.Graph {
	clientsMap := map[int]bool{}
	clients := []types.Client{}
	serversMap := map[int]bool{}
	servers := []types.Server{}
	max_color := 0
	for _, e := range edges {
		client := e.Client()
		server := e.Server()

		_, ok := clientsMap[client.Vertex()]
		if !ok {
			clientsMap[client.Vertex()] = true
			clients = append(clients, client)
		}

		_, ok = serversMap[server.Vertex()]
		if !ok {
			serversMap[server.Vertex()] = true
			servers = append(servers, server)
		}

		if  client.Color() > max_color{
			max_color = client.Color()
		}
	}


	return &graphImpl{edges, clients, servers, max_color + 1}
}

type graphImpl struct {
	edges           []types.Edge
	clients         []types.Client
	infrastructures []types.Server
	colors          int
}

func (g *graphImpl) Edges() []types.Edge {
	return g.edges
}

func (g *graphImpl) Clients() []types.Client {
	return g.clients
}

func (g *graphImpl) Servers() []types.Server {
	return g.infrastructures
}

func (g *graphImpl) MaxColor() int {
	return g.colors
}
