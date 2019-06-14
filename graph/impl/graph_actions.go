package impl

import "Capacitated/graph/types"
import "github.com/yourbasic/graph"

func GetClientServers(g types.Graph, client types.Client) []types.Server{
	clientId := client.Vertex()
	clientServers := []types.Server{}
	for _, edge := range g.Edges() {
		if clientId == edge.Client().Vertex() {
			clientServers = append(clientServers, edge.Server())
		}
	}
	return clientServers
}

func GetServerClients(g types.Graph, server types.Server) []types.Client{
	serverId := server.Vertex()
	serverClients := []types.Client{}
	for _, edge := range g.Edges() {
		if serverId == edge.Server().Vertex() {
			serverClients = append(serverClients , edge.Client())
		}
	}
	return serverClients
}

func GetServerColoredClients(g types.Graph, server types.Server, color int) []types.Client{
	serverId := server.Vertex()
	serverColoredClients := []types.Client{}
	for _, edge := range g.Edges() {
		if serverId == edge.Server().Vertex() && edge.Client().Color() == color {
			serverColoredClients = append(serverColoredClients, edge.Client())
		}
	}
	return serverColoredClients
}

func CalculateMaxAssignment(g types.Graph, placementFunction map[int]int) int64{
	res := int64(0)
	for color:=0; color<g.MaxColor(); color++ {
		res += CalculateColorValue(g, placementFunction, color)
	}

	return res
}

func CalculateColorValue(g types.Graph, placementFunction map[int]int, color int) int64 {
	total_vertices := getGraphMaxIndex(g)+3
	source := total_vertices - 2
	dest := total_vertices -1
	base_graph := graph.New(total_vertices)
	edgess := map[int][]int{}
	for _, server := range g.Servers(){
		if placementFunction[server.Vertex()] != color {
			continue
		}
		base_graph.AddCost(server.Vertex(), dest, int64(server.Capacity()))
		edgess[server.Vertex()] = []int{}
		coloredClients := GetServerColoredClients(g, server, color)
		for _, client := range coloredClients{
			base_graph.AddCost(client.Vertex(), server.Vertex(), 1)
			edgess[server.Vertex()] = append(edgess[server.Vertex()], client.Vertex())
		}
	}

	cclients := []int{}
	for _, client := range g.Clients(){
		if client.Color() == color {
			base_graph.AddCost(source, client.Vertex(), 1)
			cclients = append(cclients, client.Vertex())
		}
	}

	flow_value, _ := graph.MaxFlow(base_graph, source, dest)

	return flow_value
}

func getGraphMaxIndex(g types.Graph) int{
	max := -1
	for _, client := range g.Clients(){
		if client.Vertex() > max{
			max = client.Vertex()
		}
	}
	return max
}