package displayGoroutine

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

func make_tree(all_edges [][]string, all_nodes []string, active_nodes []string, treeName string) {
	g := gographviz.NewGraph()
	if err := g.SetName("G"); err != nil {
		panic(err)
	}
	// 有向グラフか
	if err := g.SetDir(true); err != nil {
		panic(err)
	}
	// Node設定
	nodeAttrs := make(map[string]string)
	activeNodeAttrs := make(map[string]string)
	activeNodeAttrs["fillcolor"] = "\"#00FF00\""
	activeNodeAttrs["style"] = "\"solid,filled\""

	for _, node := range all_nodes {
		if err := g.AddNode("G", node, nodeAttrs); err != nil {
			panic(err)
		}
	}

	for _, node := range active_nodes {
		if err := g.AddNode("G", node, activeNodeAttrs); err != nil {
			panic(err)
		}
	}

	edgeAttrs := make(map[string]string)
	for _, edge := range all_edges {
		if err := g.AddEdge(edge[0], edge[1], true, edgeAttrs); err != nil {
			panic(err)
		}
	}
	s := g.String()
	file, err := os.Create(treeName + ".dot")
	if err != nil {
		fmt.Println("in file panic")
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(s))
}

func Show(goroutineName string) {
	fileName := getFileName()
	folderName := fileName + "_" + goroutineName
	f, err := os.Open(folderName + "/tree_data.txt")
	// 読み取り時の例外処理
	if err != nil {
		fmt.Println("error reading file: " + fileName)
	}
	file_data, err := ioutil.ReadAll(f)
	lines := strings.Split(string(file_data), "\n")
	defer f.Close()
	all_edges := [][]string{}
	all_nodes := []string{}
	already_found_node := map[string]bool{}
	already_found_edge := map[string]bool{}
	fmt.Printf("hi")
	for _, line := range lines {
		line = strings.Replace(line, "\n", "", -1)
		if !strings.HasPrefix(line, "start") && !strings.HasPrefix(line, "end") && len(line) != 0 {
			child := strings.Split(line, " ")[0]
			parent := strings.Split(line, " ")[1]
			edge := []string{parent, child}
			key := parent + "$" + child
			// _, found := already_found_edge[key]
			if _, found := already_found_edge[key]; !found {
				all_edges = append(all_edges, edge)
				already_found_edge[key] = true
			}
			if _, found := already_found_node[child]; !found {
				all_nodes = append(all_nodes, child)
				already_found_node[child] = true
			}
			if _, found := already_found_node[parent]; !found {
				all_nodes = append(all_nodes, parent)
				already_found_node[parent] = true
			}
		}
	}
	treeCnt := 0
	active_nodes := []string{}
	for _, line := range lines {
		line = strings.Replace(line, "\n", "", -1)
		if strings.HasPrefix(line, "start") {
			active_nodes = []string{}
		} else if strings.HasPrefix(line, "end") {
			make_tree(all_edges, all_nodes, active_nodes, folderName + "/tree" + strconv.Itoa(treeCnt))
			treeCnt += 1
		} else if len(line) != 0 {
			child := strings.Split(line, " ")[0]
			active_nodes = append(active_nodes, child)
		}
	}
}

// func main() {
// 	fileName := os.Args[1]
// 	f, err := os.Open(fileName)
// 	// 読み取り時の例外処理
// 	if err != nil {
// 		fmt.Println("error reading file: " + fileName)
// 	}
// 	file_data, err := ioutil.ReadAll(f)
// 	lines := strings.Split(string(file_data), "\n")
// 	defer f.Close()
// 	all_edges := [][]string{}
// 	all_nodes := []string{}
// 	already_found_node := map[string]bool{}
// 	already_found_edge := map[string]bool{}
// 	for _, line := range lines {
// 		line = strings.Replace(line, "\n", "", -1)
// 		if !strings.HasPrefix(line, "start") && !strings.HasPrefix(line, "end") && len(line) != 0 {
// 			child := strings.Split(line, " ")[0]
// 			parent := strings.Split(line, " ")[1]
// 			edge := []string{parent, child}
// 			key := parent + "$" + child
// 			// _, found := already_found_edge[key]
// 			if _, found := already_found_edge[key]; !found {
// 				all_edges = append(all_edges, edge)
// 				already_found_edge[key] = true
// 			}
// 			if _, found := already_found_node[child]; !found {
// 				all_nodes = append(all_nodes, child)
// 				already_found_node[child] = true
// 			}
// 			if _, found := already_found_node[parent]; !found {
// 				all_nodes = append(all_nodes, parent)
// 				already_found_node[parent] = true
// 			}
// 		}
// 	}
// 	// make_tree(all_edges, all_nodes, 0)
// }
