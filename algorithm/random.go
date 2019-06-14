package algorithm

import (
	"Capacitated/graph/types"
	"Capacitated/graph/impl"
	"math/rand"
	"time"
)

func NewRandom() Algorithm {
	return &random{}
}

type random struct{}

func (r *random) Run(g types.Graph) int64 {
	rand.Seed(time.Now().UTC().UnixNano())
	servers := g.Servers()
	colors := g.MaxColor()
	coloringMap := map[int]int{}
	for _, server := range servers {
		coloringMap[server.Vertex()] = rand.Intn(colors)
	}
	return impl.CalculateMaxAssignment(g, coloringMap)
}
