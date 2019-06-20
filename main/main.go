package main

import (
	"os"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"bufio"
	"Capacitated/graph/impl/graph_creator"
	"Capacitated/algorithm"
	"Capacitated/graph/types"
	"time"
)


type setting struct {
	servers	 	int
	clients     int
	colors      int
	connections int
}

const RESULT_DIRECTORY = "Results"

func main() {
	ProccesorsNumber := runtime.NumCPU() - 1
	settingsChannel := make(chan setting, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < ProccesorsNumber; i++ {
		wg.Add(1)
		go func() {
			for s := range settingsChannel {
				runSettings(s)
			}
			wg.Done()
		}()
	}

	servers := []int{20}
	clients := []int{2,3,5,10} // 2, 5,
	colors := []int{3}//,5,10} // 5,
	connections := []int{3,5}// 2,
	numberOfRuns := 15

	settingNumber := 0

	for _, i := range servers {
		for _, c := range clients {
			for _, col := range colors {
				for _, con := range connections {
					s := setting{i, c, col, con}
					settingNumber += 1
					fmt.Println(settingNumber, s)

					for j := 0; j < numberOfRuns; j++ {
						settingsChannel<- s
						fmt.Println(j)
					}
				}
			}
		}
	}

	close(settingsChannel)
	wg.Wait()
}

type algResult struct{
	res int64
	duration int64

}

func timeAlgRun(alg algorithm.Algorithm, g types.Graph) algResult{
	start := time.Now()
	res := alg.Run(g)
	dur := time.Since(start)
	return algResult{res,int64(dur/time.Nanosecond)}
}

func max(a,b int) int {
	if a > b{
		return a
	}
	return b
}

func runSettings(s setting) {
	graphs := graph_creator.Create(s.servers, s.clients, s.colors, s.connections, true, max(s.connections, s.clients))

	random := algorithm.NewRandom()
	smartRandom := algorithm.NewSmartRandom()
	greedy := algorithm.NewGreedy()
	multilLinear := algorithm.NewMultiLinearCalculation(s.servers*s.colors*s.servers*s.colors)
	lpRound := algorithm.NewRoundCapacitatedVNFapLP()
	lpMax := algorithm.NewCapacitatedVNFapLP()


	radius_client_loss := s.clients*s.servers-len(graphs[0].Clients())
	results := []algResult{}
	for _, g := range graphs {
		results = append(results, timeAlgRun(random, g))
		results = append(results, timeAlgRun(smartRandom, g))
		results = append(results, timeAlgRun(greedy, g))
		results = append(results, timeAlgRun(multilLinear, g))
		results = append(results, timeAlgRun(lpRound, g))
		results = append(results, timeAlgRun(lpMax, g))
	}

	writeResults(s, results, radius_client_loss)
}

func writeResults(s setting, res []algResult, radius_client_loss int) {
	filePath := buildFileName(s)

	file := getFile(filePath)
	if file != nil {
		b := bufio.NewWriter(file)
		b.WriteString(createResultString(res, radius_client_loss))
		b.Flush()
		file.Close()
	}
}
func createResultString(results []algResult, radius_client_loss int) string {
	resultString := fmt.Sprint(radius_client_loss)
	for _, result := range results {
		resultString = fmt.Sprintf("%s, %v", resultString, result)
	}
	resultString = fmt.Sprintf("%s\n", resultString)
	return resultString
}

func getFile(filePath string) *os.File {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if file, err := os.Create(filePath); err == nil {
			return file
		}
	}
	if file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0660); err == nil {
		return file
	}
	return nil
}

func buildFileName(s setting) string {
	filename :=  fmt.Sprintf("res_s%d_cl%d_col%d_con%d.csv", s.servers, s.clients, s.colors, s.connections)
	return filepath.Join(getDirectoryName(), filename)
}

func getDirectoryName() string{
	_,b,_,_ := runtime.Caller(0)
	results_directory := filepath.Dir(filepath.Dir(b))
	return filepath.Join(results_directory, RESULT_DIRECTORY)
}