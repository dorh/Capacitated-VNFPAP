package algorithm

import (
	"Capacitated/graph/types"
	"Capacitated/graph/impl"
	"math/rand"
	"time"
)

func NewSmartRandom() Algorithm {
	return &smartRandom{}
}

type smartRandom struct{}

func (sr *smartRandom) Run(g types.Graph) int64 {
	rand.Seed(time.Now().UTC().UnixNano())
	servers := g.Servers()
	coloringMap := map[int]int{}
	for _, server := range servers {
		serverClients := impl.GetServerClients(g, server)
		serverColors := []int{}
		for _, client := range serverClients {
			serverColors = append(serverColors, client.Color())
		}
		coloringMap[server.Vertex()] = serverColors[rand.Intn(len(serverColors))]
	}
	return impl.CalculateMaxAssignment(g, coloringMap)
}
