package algorithm

import (
	"Capacitated/graph/types"
	"Capacitated/graph/impl"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	_LP_DIRECTORY         = "lp_files"
	_LP_COMMAND           = "C:\\Users\\dorharr\\lp_solve\\lp_solve.exe"
	_OBJECTIVE_VALUE_LINE = 1
	_VARIABLE_START_LINE  = 5
)

func NewCapacitatedVNFapLP() Algorithm {
	return &capacitatedVnfapLP{}
}

func NewRoundCapacitatedVNFapLP() Algorithm {
	return &roundCapacitatedVnfpapLp {}
}

type capacitatedVnfapLP struct {}

type roundCapacitatedVnfpapLp struct {
	capacitatedVnfapLP
}

func (v *roundCapacitatedVnfpapLp) Run(g types.Graph) int64 {
	lpFilename := v.createLpFile(g)
	variables := v.runLp(lpFilename)
	return v.createSolutionFromVariables(g, variables)
}

func (v *capacitatedVnfapLP ) Run(g types.Graph) int64 {
	lpFilename := v.createLpFile(g)
	return int64(v.runLp(lpFilename))
}

func (v *capacitatedVnfapLP) createLpFile(g types.Graph) string {
	clients := g.Clients()
	clientsVariables := v.createClientsVariables(clients)
	servers := g.Servers()
	clientServerVariables := v.createClientServersVariables(g, clients)
	colors := g.MaxColor()
	serversVariables := v.createServersVariables(servers, colors)
	objectiveFunction := v.createObjectiveFunction(clientsVariables)
	constraints := v.createConstraints(g, clients, servers, colors, clientsVariables, serversVariables,
		clientServerVariables)
	declerations := v.createDeclerations(clientsVariables, serversVariables,
		clientServerVariables)
	filename := v.writeLpFile(objectiveFunction, constraints, declerations)
	return filename
}

func (v *capacitatedVnfapLP) createClientsVariables(clients []types.Client) map[int]string {
	cVariables := map[int]string{}
	for _, c := range clients {
		cVariables[c.Vertex()] = v.clientVar(c)
	}
	return cVariables
}

func (v *capacitatedVnfapLP) clientVar(client types.Client) string {
	return fmt.Sprintf("c%d", client.Vertex())
}

func (v *capacitatedVnfapLP) createClientServersVariables(g types.Graph,
	clients []types.Client) map[int]map[int]string {
	ciVariables := map[int]map[int]string{}
	for _, client := range clients {
		clientServers := impl.GetClientServers(g, client)
		ciVariables[client.Vertex()] = v.createClientServerVariables(client, clientServers)
	}
	return ciVariables

}

func (v *capacitatedVnfapLP) createClientServerVariables(client types.Client, servers []types.Server) map[int]string {
	ciVariables := map[int]string{}
	for _, server := range servers {
		ciVariables[server.Vertex()] = clientServerVar(client, server)
	}
	return ciVariables
}

func clientServerVar(c types.Node, i types.Node) string {
	return fmt.Sprintf("x_%d_%d", c.Vertex(), i.Vertex())
}

func (v *capacitatedVnfapLP) createServersVariables(servers []types.Server, colors int) map[int]map[int]string {
	serversVariables := map[int]map[int]string{}
	for _, server := range servers {
		serversVariables[server.Vertex()] = v.createServerVariables(server, colors)
	}
	return serversVariables
}

func (v *capacitatedVnfapLP) createServerVariables(server types.Server, colors int) map[int]string {
	colorsVariables := map[int]string{}
	for j := 0; j < colors; j++ {
		colorsVariables[j] = v.serverColorVar(server, j)
	}
	return colorsVariables
}

func (v *capacitatedVnfapLP) serverColorVar(server types.Server, color int) string {
	return fmt.Sprintf("i%d_%d", server.Vertex(), color)
}

func (v *capacitatedVnfapLP) createObjectiveFunction(clientsVariables map[int]string) string {
	max_string := ""
	for _, variable := range clientsVariables {
		max_string = fmt.Sprintf("%s%s + ", max_string, variable)
	}
	max_string = strings.Trim(max_string, " +")
	return fmt.Sprintf("%v;", max_string)
}

func getFile(filePath string) *os.File {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err == nil {
			return file
		} else{
			fmt.Println(err)
		}
	}
	if file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0660); err == nil {
		return file
	}
	return nil
}

func (v *capacitatedVnfapLP) writeLpFile(obecjtiveFunction string, constraints []string, declerations []string) string {
	filename := fmt.Sprintf("%d.%f.lp", time.Now().UTC().Unix(), rand.Float64())
	lpFilename := filepath.Join(v.lpDirectory(), filename)
	f, _ := os.Create(lpFilename)
	w := bufio.NewWriter(f)
	v.writeString(w, obecjtiveFunction)
	for _, constraint := range constraints {
		v.writeString(w, constraint)
	}
	for _, decleration := range declerations {
		v.writeString(w, decleration)
	}
	f.Close()
	return lpFilename
}

func (v *capacitatedVnfapLP) lpDirectory() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), _LP_DIRECTORY)
}

func (v *capacitatedVnfapLP) writeString(writer *bufio.Writer, s string) {
	writer.WriteString(fmt.Sprintf("%s\n", s))
	writer.Flush()
}

func (v *capacitatedVnfapLP) createConstraints(g types.Graph, clients []types.Client, servers []types.Server,
	colors int, clientsVariables map[int]string, serversVariables map[int]map[int]string,
		clientServersVariables map[int]map[int]string) []string {
	clientConstraints := v.createClientConstraints(g, clients, clientsVariables, clientServersVariables)
	serverConstraints := v.createServersConstraints(servers, colors, serversVariables)
	constraints := append(clientConstraints, serverConstraints...)
	colorConstraints := v.createColorConstraints(g, servers, colors, serversVariables,
		clientServersVariables)
	return append(constraints, colorConstraints...)
}

func (v *capacitatedVnfapLP) createClientConstraints(g types.Graph, clients []types.Client,
	clientsVariables map[int]string, clientServersVariables map[int]map[int]string) []string {
	constraints := []string{}
	for _, client := range clients {
		clientVar := clientsVariables[client.Vertex()]
		clientVarConstraint := fmt.Sprintf("0.0 <= %v <= 1.0;", clientVar)
		constraints = append(constraints, clientVarConstraint)

		clientServers := impl.GetClientServers(g, client)
		clientServerConstraint := fmt.Sprint(clientVar)

		for _, server := range clientServers {
			clientServerConstraint = fmt.Sprintf("%s - %s",
				clientServerConstraint, clientServersVariables[client.Vertex()][server.Vertex()])
		}
		clientServerConstraint = fmt.Sprintf("%s <= 0;", clientServerConstraint)
		constraints = append(constraints, clientServerConstraint)
	}

	return constraints
}

func (v *capacitatedVnfapLP) createServersConstraints(servers []types.Server, colors int,
	serversVariables map[int]map[int]string) []string {
	constraints := []string{}
	for _, server := range servers {
		serverConstraint := fmt.Sprint(serversVariables[server.Vertex()][0])
		for j := 1; j < colors; j++ {
			currVariable := serversVariables[server.Vertex()][j]
			serverConstraint = fmt.Sprintf("%s + %s", serverConstraint, currVariable)
		}
		serverConstraint = fmt.Sprintf("%s <= 1;", serverConstraint)
		constraints = append(constraints, serverConstraint)

		for j := 0; j < colors; j++ {
			currVariable := serversVariables[server.Vertex()][j]
			serverColorVarConstraint := fmt.Sprintf("0.0 <= %v <= 1.0;",currVariable)
			constraints = append(constraints, serverColorVarConstraint)
		}

	}
	return constraints
}

func (v *capacitatedVnfapLP) createDeclerations(clientsVariables map[int]string, serversVariables map[int]map[int]string,
	clientServersVariables map[int]map[int]string) []string {
	declerations := []string{}
	clientsDecleration := "sec"
	for _, cVar := range clientsVariables {
		clientsDecleration = fmt.Sprintf("%s %s,", clientsDecleration, cVar)
	}
	clientsDecleration = strings.TrimRight(clientsDecleration, ", ")
	declerations = append(declerations, fmt.Sprintf("%s;", clientsDecleration))
	for _, serverVariables := range serversVariables {
		serverDecleration := "sec"
		for _, iVar := range serverVariables {
			serverDecleration = fmt.Sprintf("%s %s,", serverDecleration, iVar)
		}
		serverDecleration = strings.TrimRight(serverDecleration, ", ")
		declerations = append(declerations, fmt.Sprintf("%s;", serverDecleration))
	}
	for _, clientServerVariable := range clientServersVariables {
		clientColorDecleration := "sec"
		for _, cVar := range clientServerVariable {
			clientColorDecleration  = fmt.Sprintf("%s %s,", clientColorDecleration , cVar)
		}
		clientColorDecleration  = strings.TrimRight(clientColorDecleration , ", ")
		declerations = append(declerations, fmt.Sprintf("%s;", clientColorDecleration ))
	}
	return declerations
}

func (v *capacitatedVnfapLP) createColorConstraints(g types.Graph, servers []types.Server, colors int,
	serversVariables map[int]map[int]string, clientServersVariables map[int]map[int]string) []string {
	constraints := []string{}
	for _, server := range servers {
		serverClients := impl.GetServerClients(g, server)
		for _, client := range serverClients {
			currVariable := clientServersVariables[client.Vertex()][server.Vertex()]
			serverVarConstraint := fmt.Sprintf("0.0 <= %v <= 1.0;",currVariable)
			constraints = append(constraints, serverVarConstraint)
		}

		for j:=0; j< colors; j++{
			serverColorConstraint := ""
			found := false
			for _, client := range serverClients {
				if client.Color() == j{
					currVariable:= clientServersVariables[client.Vertex()][server.Vertex()]
					serverColorConstraint = fmt.Sprintf("%s%s + ",
						serverColorConstraint, currVariable)
					found = true

					curServerColorVar := serversVariables[server.Vertex()][j]
					constraint := fmt.Sprintf("%s - %s <= 0.0;", currVariable, curServerColorVar)
					constraints = append(constraints, constraint)
				}
			}
			if found{
				serverColorConstraint = strings.Trim(serverColorConstraint, " +")
				curServerColorVar := serversVariables[server.Vertex()][j]
				serverColorConstraint = fmt.Sprintf("%s - %v*%s <= 0.0;",
					serverColorConstraint, server.Capacity() ,curServerColorVar)
				constraints = append(constraints, serverColorConstraint)
			}
		}
	}
	return constraints
}

func (v *capacitatedVnfapLP) runLp(lpFilename string) float64 {
	cmd := exec.Command(_LP_COMMAND, lpFilename)
	o, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Run()
	objectiveValue := v.analyzeResults(string(o))
	os.Remove(lpFilename)
	return objectiveValue
}

func (v *capacitatedVnfapLP) analyzeResults(results string) float64 {
	lines := strings.Split(results, "\n")
	objectiveValue :=0.0
	for _, line := range lines{
		if strings.HasPrefix(line, "Value of objective"){
			splittedObjective := strings.Fields(line)
			objectiveValue, _ = strconv.ParseFloat(splittedObjective[len(splittedObjective)-1], 64)
			return objectiveValue
		}
	}
	return objectiveValue
}


func (v *roundCapacitatedVnfpapLp) runLp(lpFilename string) map[string]float64 {
	cmd := exec.Command(_LP_COMMAND, lpFilename)
	o, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Run()
	variableValues := v.analyzeResults(string(o))
	os.Remove(lpFilename)
	return variableValues
}

func (v *roundCapacitatedVnfpapLp) analyzeResults(results string) map[string]float64{
	lines := strings.Split(results, "\n")
	start_line := 0
	for i, line := range lines{
		if strings.HasPrefix(line, "Actual values"){
			start_line = i+1
			break
		}
	}
	variableValues := map[string]float64{}
	for lineNumber := start_line; lineNumber < len(lines)-1; lineNumber++ {
		splittedS := strings.Fields(lines[lineNumber])
		variableValue, _ := strconv.ParseFloat(splittedS[1], 64)
		variable := splittedS[0]
		variableValues[variable] = variableValue
	}
	return variableValues
}

func (v *roundCapacitatedVnfpapLp) createSolutionFromVariables(g types.Graph, variables map[string]float64) int64 {
	servers := g.Servers()
	placementFunction := map[int]int{}
	colors := g.MaxColor()
	for _, server := range servers {
		colorProbability := map[int]float64{}
		for c := 0; c < colors; c++ {
			iVar := v.serverColorVar(server, c)
			colorProbability[c] = variables[iVar]
		}
		x := rand.Float64()
		for c, prob := range colorProbability {

			x = x - prob
			if x <= 0 {
				placementFunction[server.Vertex()] = c
				break
			}
		}
	}
	return impl.CalculateMaxAssignment(g, placementFunction)
}