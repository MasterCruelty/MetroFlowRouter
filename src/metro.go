package main


/*
 * Design of the algorithm.
 *
 * Point 1.
 * The metro network can be described throught a graph as data structure.
 * Each station is a graph's node. Included the changing stations which allow to change line.
 * "Cadorna 1" is the node of line 1. "Cadorna 2" is the node of line 2. Both are linked.
 * The path between two near stations is represented by the graph's link. The weight of each link is 1(for now).
 * This is then a weighted graph with a common weight to each link, not oriented and strongly connected.
 * From each station I can reach every station in the network, I can go back too without any limitation of direction.
 *
 * 
 * Point 2.
 * We have the starting point and the destination node. So we can get the mininum path to the destination.
 * We are talking about the mininum lenght of the path between a S node and a D node, using a BFS of the graph.
 * So we are going to use S as origin node and we visit the adjacent nodes of S and thus so on until we reach D.
 * At the end of execution I'm going to compare the lenght of the paths and I return the mininum.
 *
 * Point 3.
 * The graph is implementedi by using adjacency list. This list is implemented using maps with station name as key and the near stations as value. 
*/

import (
		"fmt"
		"bufio"
		"os"
		"strings"
		"strconv"
)

type stazione struct{
	nome string
	linea int
	interscambio bool
        capolineaBiforcazione bool
}

type rete struct{
	adj map[string][]*stazione //Given a station's name, we have the adjacents.
	linee map[int][]*stazione  //Given the line number, we have the sorted stations of that line.
}

/*
 * Function that reads the subway graph data from a text file
 * and returns the graph as a network by populating it.
 * Complexity evaluation:
 * There is a first loop that performs a series of elementary operations for each line of the file.
 * Inside it, there is another nested loop that executes a number of elementary operations
 * dependent on the number of stations on each line of the file. If we denote R as the line of the file
 * and S as the number of stations, this first part has a complexity of O(R*S).
 * Subsequently, there are two more nested loops, the outer one performing S operations
 * and the inner one also performing S operations. This second part has a complexity of O(S^2).
 * Finally, there is one last for loop that performs S elementary operations, resulting in a cost of O(S).
 * O(R*S) + O(S^2) + O(S) = O(S^2)
*/

func leggiDati() rete{
	myFile, err := os.Open(os.Args[1])
	if err != nil{
		fmt.Println("file non trovato")
		os.Exit(0)
	}
	defer myFile.Close()

	scanner := bufio.NewScanner(myFile)
	var staz []stazione
	metro := &rete{adj:make(map[string][]*stazione),linee:make(map[int][]*stazione)}

	//I read the content of file
	for scanner.Scan() {
		//Format text fetched from file
		linea := scanner.Text()
		stazioniLinea := strings.Split(linea,": ")
		numLinea,_ := strconv.Atoi(strings.TrimPrefix(stazioniLinea[0], "Linea "))
		stazioni := strings.Split(stazioniLinea[1],"; ")
                var indexCapolinea int
		//Create a slice containing all the stations sorted
                indexCapolinea = 0
		for k,v := range stazioni{
			stazione := &stazione{nome: v,linea: numLinea}
                        //look for bifurcations stations
                        for i := k; i < len(stazioni); i++ {
                            for j := i + 1; j < len(stazioni); j++ {
                                if stazioni[i] == stazioni[j] {
                                    //keep in memory the index of the end of a branch
                                    indexCapolinea = len(staz) + (j-i) -1 
                                    fmt.Println(stazioni[i])
                                    break
                                }
                            }
                            if indexCapolinea != 0{
                                break
                            }
                        }
			staz = append(staz,*stazione)
			metro.linee[numLinea] = append(metro.linee[numLinea],stazione)
		}
                //I set the end of a branch to avoid wrong adj
                if indexCapolinea != 0{
                    staz[indexCapolinea].capolineaBiforcazione = true
                    //fmt.Println(staz[indexCapolinea].nome)
                }
	}
	//I look for changing stations and I manage the adjacency between them.
	//For example Cadorna on line 1 is adjacent to Cadorna on line 2
	for i:= 0; i < len(staz); i++{
		for j:= 0;j < len(staz); j++{
			if(i != j && staz[i].nome == staz[j].nome && staz[i].linea != staz[j].linea) {
				staz[j].interscambio = true
				staz[i].interscambio = true
				staz[j].nome = staz[j].nome + "-" + strconv.Itoa(staz[j].linea)
				staz[i].nome = staz[i].nome + "-" + strconv.Itoa(staz[i].linea)
				metro.adj[staz[j].nome] = append(metro.adj[staz[j].nome],&staz[i])
				metro.adj[staz[i].nome] = append(metro.adj[staz[i].nome],&staz[j])
			}
		}
	}
	//I fill the adjacent slice
	for i, stazione := range staz{
		if i > 0 && stazione.linea == staz[i-1].linea && staz[i-1].capolineaBiforcazione == false{
			metro.adj[stazione.nome] = append(metro.adj[stazione.nome],&staz[i-1])
		}
		if i < len(staz)-1 && stazione.linea == staz[i+1].linea && staz[i].capolineaBiforcazione == false{
			metro.adj[stazione.nome] = append(metro.adj[stazione.nome],&staz[i+1])
		}
	}
	return *metro
}

/*
 * Returns the slice with the stations of the line numLine in order.
 * Complexity evaluation:
 * O(1) - Directly retrieve the stations of the requested line from the map.
*/

func linea(metro rete,numLinea int) []*stazione{
	stazioni := metro.linee[numLinea]
	return stazioni
}

/*
 * Returns the stations adjacent to the input station.
 * Complexity evaluation:
 * The initial assignment is executed in O(1).
 * The for loop performs at most a few iterations to populate the string slice with neighboring stations,
 * resulting in a complexity of O(S). It's important to note that, in a subway network, the number of adjacent stations (S) for a given station is typically not large.
*/

func stazioniVicine(metro rete, s string) []string{
	vicine := metro.adj[s]
	result := make([]string,0)
	for i:= 0;i< len(vicine);i++{
		result = append(result,vicine[i].nome)
	}
	return result
}

/*
 * Returns the slice containing the names of interchange stations.
 * Complexity evaluation:
 * The first loop iterates through the elements of the metro.adj map.
 * The second loop iterates through the stations slice connected to the map as a value.
 * Let R be the number of keys in the map and S be the number of stations per key.
 * Final complexity: O(R*S)
*/

func interscambio(metro rete) []string{
    trovatoScambio := make(map[string]bool)
    interscambi := make([]string, 0)

    for _, stazioni := range metro.adj {
        for _, stazione := range stazioni {
            if stazione.interscambio {
                nomeStazione := stazione.nome
                if !trovatoScambio[nomeStazione] {
                    trovatoScambio[nomeStazione] = true
                    interscambi = append(interscambi, nomeStazione)
                }
            }
        }
    }

    return interscambi
}

/*
 * Returns true if two stations are on the same line, otherwise false.
 * Complexity evaluation:
 * O(1) - All operations are elementary, directly retrieved from the map.
*/
func stessaLinea(metro rete,s1 string,s2 string) bool{
	stazione1 := metro.adj[s1][len(metro.adj[s1])-1].linea
	stazione2 := metro.adj[s2][len(metro.adj[s2])-1].linea
	if stazione1 == stazione2{
		return true
	}else{
		return false
	}
}


/*
 * Returns the minimum time to reach the destination given the departure station.
 * This represents the minimum number of stations crossed, considering a fixed weight on each node.
 * The first loop iterates until the queue is empty, executed S times where S is the number of stations.
 * The second loop is executed a number of times equal to the current level of breadth-first search.
 * It iterates over all stations at the current level, which can be at most S.
 * The third loop iterates over all adjacent stations, a number R of adjacent stations.
 * Final complexity: O(V) * O(V) + O(E) = O(V^2 + E)
*/

func tempo(metro rete, partenza string, arrivo string) ([]string,int) {
    coda := []string{partenza}
    aux := make(map[string]bool)
    aux[partenza] = true
    result := 0
	percorso := make(map[string][]string)
	percorso[partenza] = []string{partenza}

    for len(coda) > 0 {
        livello := len(coda)
        for i := 0; i < livello; i++ {
            partenza := coda[0]
            coda = coda[1:]
            //fmt.Println("Sto visitando i vicini di: " + partenza)
            for _, vicino := range metro.adj[partenza] {
                if !aux[vicino.nome] {
                    if vicino.nome == arrivo {
        		percorso[vicino.nome] = append(percorso[partenza],vicino.nome)
			return percorso[vicino.nome], result + 1
                    }
                    coda = append(coda, vicino.nome)
                    aux[vicino.nome] = true
                    //To avoid loss of information in the case we have more than one near station not visited
                    //I create distint copies of the path for each near station not visited and I updated them separately
                    nuovoPercorso := make([]string,len(percorso[partenza]))
                    copy(nuovoPercorso,percorso[partenza])
                    percorso[vicino.nome] = append(nuovoPercorso,vicino.nome)
                }
            }
        }
        result++
    }
	return nil,-1
}

//main di test per le funzioni
func main() {
	metro := leggiDati()

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1) Adjacent stations (type 2 stations to check if they're adjacent)")
		fmt.Println("2) Show all the stations of a line")
		fmt.Println("3) Show all changing stations")
		fmt.Println("4) Check if two stations are on the same line")
		fmt.Println("5) Find minimum path and time between two stations")
		fmt.Println("6) Exit")

		var choice int
		fmt.Print("Enter your choice (1-6): ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("Which station you wanna know if they are adjacent?")
			stazione := ""
			fmt.Scan(&stazione)
			vicine := stazioniVicine(metro, stazione)
			fmt.Println(vicine)

		case 2:
			fmt.Println("Which line you wanna know all the stations?")
			var line int
			fmt.Scan(&line)
			stazioni := linea(metro, line)
			for i := 0; i < len(stazioni); i++ {
				fmt.Println(stazioni[i].nome)
			}

		case 3:
			fmt.Println("\nAll the changing stations:")
			interscambi := interscambio(metro)
			fmt.Println(interscambi)

		case 4:
			fmt.Println("Which stations you wanna check if they are on the same line?")
			var s1, s2 string
			fmt.Scan(&s1, &s2)
			ok := stessaLinea(metro, s1, s2)
			fmt.Println(ok)

		case 5:
			fmt.Println("Starting point and destination:")
			var s1, s2 string
			fmt.Scan(&s1, &s2)
			fmt.Println()
			percorso, time := tempo(metro, s1, s2)
			fmt.Println("Minimum path: " + strconv.Itoa(time))
			fmt.Println("Route to be covered: ")
			fmt.Println(percorso)

		case 6:
			fmt.Println("Exiting the program.")
			return

		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 6.")
		}
	}
}

