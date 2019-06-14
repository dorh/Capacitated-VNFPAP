package graph_creator

import (
	"Capacitated/graph/impl"
	"Capacitated/graph/types"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"runtime"
	"path/filepath"
	"math/rand"
)

const (
	SERVER_DIRECTORY = "server_locations"
)

var (
	_SERVERS []int = []int{5, 10, 20, 50, 100, 200, 500}
)

func Create(servers int, clients int, colors int, connections int, radius bool, serverCapacity int) []types.Graph {
	if !isServerExists(servers) {
		return nil
	}
	rand.Seed(time.Now().UTC().UnixNano())
	serversNodes, boundaries := get_servers(servers, serverCapacity)
	clientsNodes := getClients(clients*servers, boundaries, colors, servers)
	edges := createEdges(clientsNodes, serversNodes, connections)
	vnfapGraph := impl.NewGraph(edges)
	if !radius {
		return []types.Graph{vnfapGraph}
	}
	r, m := findRadius(vnfapGraph, clientsNodes)
	connectionRadius := findConnectionRadius(connections, clientsNodes, serversNodes, r, m)
	radiusGraphEdges := getEdgesInRadius(clientsNodes, serversNodes, connectionRadius)
	radiusGraph := impl.NewGraph(radiusGraphEdges)
	radiusGraphClients := radiusGraph.Clients()
	radiusGraphServers := radiusGraph.Servers()
	vnfapGraphEdges := createEdges(radiusGraphClients, radiusGraphServers, connections)
	vnfapGraph = impl.NewGraph(vnfapGraphEdges)
	return []types.Graph{vnfapGraph, radiusGraph}
}

func get_servers(servers int, serverCapacity int) ([]types.Server, []float64) {
	filename := getServersFile(servers)
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil
	}
	r := bufio.NewReader(file)
	s, _ := readln(r)
	boundaries := extractBoundaries(s)
	Servers := []types.Server{}
	s, e := readln(r)
	for e == nil {
		Servers = append(Servers, buildServer(s, len(Servers), serverCapacity))
		s, e = readln(r)
	}
	file.Close()
	return Servers, boundaries
}

func extractBoundaries(s string) []float64 {
	boundariesString := strings.Split(s, " ")
	boundaries := []float64{}
	for _, val := range boundariesString {
		boundaries = append(boundaries, strToFloat(val))
	}
	return boundaries
}

func buildServer(location string, index int, serverCapacity int) types.Server {
	serverCoordinates := strings.Split(location, ",")
	latitude := strToFloat(serverCoordinates[0])
	longitude := strToFloat(serverCoordinates[1])
	return impl.NewServer(index, longitude, latitude, serverCapacity)
}

func strToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return float64(0)
	}
	return f
}

func getClients(clients int, boundaries []float64, colors int, min_index int) []types.Client {
	clients_list := []types.Client{}

	for i := 0; i < clients; i++ {
		latitude := randFloat(boundaries[0], boundaries[1])
		longitude := randFloat(boundaries[2], boundaries[3])
		client_color := rand.Intn(colors)

		clients_list = append(clients_list,
			impl.NewClient(min_index+len(clients_list), client_color , longitude, latitude))
	}
	return clients_list
}

func randFloat(min float64, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

func createEdges(clients []types.Client, servers []types.Server, connections int) []types.Edge {
	edges := []types.Edge{}
	for _, c := range clients {
		closestServers := getClosestServers(c, servers, connections)
		for _, i := range closestServers {
			edges = append(edges, impl.NewEdge(c, i, len(edges)))
		}
	}
	return edges
}

func getClosestServers(client types.Client, servers []types.Server, connections int) []types.Server {
	minimalDistances := []float64{}
	minimalservers := []types.Server{}
	for i := 0; i < connections; i++ {
		minimalDistances = append(minimalDistances, 9000.0+float64(i))
		minimalservers = append(minimalservers, nil)
	}
	for _, server := range servers {
		currentDistance := client.Distance(server)
		for i, distance := range minimalDistances {
			if currentDistance < distance {
				copy(minimalDistances[i+1:], minimalDistances[i:])
				minimalDistances[i] = currentDistance
				copy(minimalservers[i+1:], minimalservers[i:])
				minimalservers[i] = server
				break
			}
		}
	}
	return minimalservers
}


func findConnectionRadius(connections int, clients []types.Client, servers []types.Server, minimal_radius float64,
	maximal_radius float64) float64{
	min := minimal_radius
	max := maximal_radius
	for i:=0; i<15; i++{
		mid := (max + min)/2.0
		edges, clients_num  := getEdgesNumberInRadius(clients, servers, mid)
		if float64(edges)/float64(clients_num) > float64(connections){
			max = mid
		} else{
			min = mid
		}
	}
	return (max + min)/2
}

func findRadius(graph types.Graph, clients []types.Client) (float64, float64) {
	minimalDistanceSum := 0.0
	maxDist := 0.0
	for _, client := range clients {
		clientServer := impl.GetClientServers(graph, client)
		minDistance := client.Distance(clientServer[0])
		for _, i := range clientServer {
			currentDistance := client.Distance(i)
			if currentDistance < minDistance {
				minDistance = currentDistance
			}
			if currentDistance > maxDist {
				maxDist = currentDistance
			}

		}
		minimalDistanceSum += minDistance
	}
	return minimalDistanceSum / float64(len(clients)), maxDist
}

func getEdgesNumberInRadius(clients []types.Client, servers []types.Server, radius float64) (int, int) {
	edgesNumber := 0
	clientsNotConnected := 0
	for _, client := range clients {
		found := false
		for _, server := range servers {
			if client.Distance(server) < radius {
				edgesNumber += 1
				found = true
			}

		}
		if !found {
			clientsNotConnected++
		}
	}
	return edgesNumber, len(clients) - clientsNotConnected
}

func getEdgesInRadius(clients []types.Client, servers []types.Server, radius float64) ([]types.Edge) {
	edges := []types.Edge{}
	for _, client := range clients {
		for _, server := range servers {
			if client.Distance(server) < radius {
				edges = append(edges, impl.NewEdge(client, server, len(edges)))
			}
		}
	}
	return edges
}


func isServerExists(servers int) bool {
	for _, val := range _SERVERS {
		if servers == val {
			return true
		}
	}
	return false
}

func getServersFile(servers int) string {
	filename := fmt.Sprintf("%d.csv", servers)
	return filepath.Join(getFilePath(), SERVER_DIRECTORY, filename)
}

func getFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)

}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
