package displayGoroutine

import (
	"os"
  "fmt"
   "io/ioutil"
   "strings"
	"github.com/awalterschulze/gographviz"
)

func make_tree(all_edges [][]string, all_nodes []string, tree_num int) {
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

	for _, node := range all_nodes {
		if err := g.AddNode("G", node, nodeAttrs); err != nil {
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
	file, err := os.Create("graph/" + strings.Split(string(os.Args[1]), ".")[0] + ".dot")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(s))
}

func main() {
  file_name := os.Args[1]
  f, err := os.Open(file_name)
  // 読み取り時の例外処理
  if err != nil{
      fmt.Println("error reading file: " + file_name)
  }
  file_data, err := ioutil.ReadAll(f)
  lines := strings.Split(string(file_data), "\n")
  defer f.Close()
  all_edges := [][]string{}
  all_nodes := []string{}
  already_found_node := map[string]bool{}
  already_found_edge := map[string]bool{}
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
  make_tree(all_edges, all_nodes, 0)
}
