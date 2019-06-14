package algorithm

import (
	"Capacitated/graph/types"
	"Capacitated/graph/impl"
)

func NewGreedy() Algorithm {
	return &greedy{}
}

type greedy struct{}

func (gr *greedy) Run(g types.Graph) int64 {
	servers := g.Servers()
	servedClients := []int{}

	for _, server := range servers {
		serverClients := impl.GetServerClients(g, server)
		serverColors := map[int]int{}
		for _, client := range serverClients {
			if clientIsServed(client.Vertex(), servedClients) {
				continue
			}
			_, ok := serverColors[client.Color()]
			if !ok {
				serverColors[client.Color()] = 0
			}
			serverColors[client.Color()] += 1
		}

		bestColor := findBestColor(serverColors)
		if bestColor == -1 {
			continue
		}

		leftToServe := server.Capacity()
		for _, client := range serverClients {
			if leftToServe == 0 || clientIsServed(client.Vertex(), servedClients) {
				continue
			}
			if client.Color() == bestColor{
				servedClients = append(servedClients, client.Vertex())
				leftToServe--
			}
		}
	}

	return int64(len(servedClients))
}
func findBestColor(serverColors map[int]int) int {
	maxColor := -1
	max := -1
	for color, amount := range serverColors{
		if amount > max{
			max = amount
			maxColor = color
		}
	}

	return maxColor
}

func clientIsServed(vertex int, clients_served []int) bool{
	for _, client := range clients_served{
		if client == vertex{
			return true
		}
	}
	return false
}
