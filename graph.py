from graphviz import Digraph
import sys

def make_tree(edges, tree_num):
    g = Digraph(format='png')
    g.attr('node', shape='circle')
    for edge in edges:
        g.edge(edge[1], edge[0])
    g.render(f'graph/{tree_num}')

path = sys.argv[1]
with open(path) as f:
    ls = f.readlines()
    edges = []
    tree_num = 0
    for l in ls:
        l = l.replace('\n', '')
        if 'start' in l:
            edges = []
        elif 'end' in l:
            make_tree(edges, tree_num)
            tree_num += 1
        else:
            edges.append(list(l.split(' ')))
