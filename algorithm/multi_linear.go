package algorithm

import (
	"Capacitated/graph/types"
	"math/rand"
	"Capacitated/graph/impl"
)

func NewMultiLinearCalculation(iterations int) Algorithm {
	return &multiLinearCalculation {iterations }
}

type multiLinearCalculation struct{
	iterations int
}


func (mlc *multiLinearCalculation ) Run(g types.Graph) int64 {
	addedValue := 1.0/float64(mlc.iterations)
	multi_linear_extension := make([][]float64, len(g.Servers()))
	for i := range multi_linear_extension{
		multi_linear_extension[i] = make([]float64, g.MaxColor())
		for j := range multi_linear_extension[i]{
			multi_linear_extension[i][j] = 0.0
		}
	}

	for i:= 0; i<mlc.iterations; i++{

		currentIterationBest := map[int]int{}
		for server := range g.Servers(){
			currentIterationBest[server] = mlc.findCurrentColorToServer(g, multi_linear_extension, server)
		}


		for server,color := range currentIterationBest{
			multi_linear_extension[server][color] += addedValue
		}
	}

	placementFunction := randomPlacementFunction(multi_linear_extension)
	return impl.CalculateMaxAssignment(g, placementFunction)
}
func randomPlacementFunction(multi_linear_extension [][]float64) map[int]int {
	placementFunction := map[int]int{}
	for server, serverColors := range 	multi_linear_extension{
		x := rand.Float64()
		for c, prob := range serverColors {
			x = x - prob
			if x <= 0 {
				placementFunction[server] = c
				break
			}
		}
		if _, ok := placementFunction[server]; !ok{
			placementFunction[server] = rand.Intn(len(serverColors))
		}
	}

	return placementFunction
}
func (mlc *multiLinearCalculation) findCurrentColorToServer(g types.Graph, multi_linear_extension [][]float64,
	server int) int {
	max := int64(-1)
	maxColor := -1
	for color := 0; color < g.MaxColor(); color++{
		currentValue := mlc.calculateMultiLinearValue(g, multi_linear_extension, server, color)
		if currentValue > max{
			max = currentValue
			maxColor = color
		}
	}

	return maxColor
}

func (mlc *multiLinearCalculation) calculateMultiLinearValue(g types.Graph, multi_linear_extension [][]float64,
	server int, color int) int64 {
	sum := int64(0)
	numOfIterations := len(g.Servers())*g.MaxColor()*g.MaxColor()
	for i:=0; i< numOfIterations; i++{
		sum += mlc.calculateOneTime(g, multi_linear_extension, server, color)
	}
	return sum
}

func (mlc *multiLinearCalculation) calculateOneTime(g types.Graph, multi_linear_extension [][]float64, server int,
	color int) int64 {
	addedValue := 1.0/float64(mlc.iterations)

	sum := (int64)(0)
	for c := 0; c<g.MaxColor(); c++ {
		placementFunction := map[int]int{}
		for i:= range g.Servers(){
			prob := multi_linear_extension[i][c]
			if i == server && c == color{
				prob += addedValue
			}
			if rand.Float64() < prob{
				placementFunction[i] = c
			} else {
				placementFunction[i] = -1
			}
		}

		sum += impl.CalculateColorValue(g, placementFunction, c)

	}
	return sum

}