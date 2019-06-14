package algorithm

import "Capacitated/graph/types"

type Algorithm interface {
	Run(types.Graph) int64
}
