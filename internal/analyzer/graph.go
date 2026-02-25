package analyzer

// Node represents a Go package in the dependency graph.
type Node struct {
	ImportPath string
	Files      []string
	Imports    []string
}

// addImport appends an import path, skipping duplicates.
func (n *Node) addImport(importPath string) {
	for _, existing := range n.Imports {
		if existing == importPath {
			return
		}
	}
	n.Imports = append(n.Imports, importPath)
}

// DependencyGraph is an in-memory directed graph of Go package dependencies.
type DependencyGraph struct {
	nodes map[string]*Node
}

// NewDependencyGraph creates an empty dependency graph.
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{nodes: make(map[string]*Node)}
}

// AddNode adds a package to the graph. If the package is already present,
// the existing node is returned unchanged.
func (g *DependencyGraph) AddNode(importPath string, files []string) *Node {
	if n, ok := g.nodes[importPath]; ok {
		return n
	}
	n := &Node{ImportPath: importPath, Files: files}
	g.nodes[importPath] = n
	return n
}

// AddEdge records an import edge from one package to another.
// The source node must already exist in the graph.
func (g *DependencyGraph) AddEdge(from, to string) {
	if n, ok := g.nodes[from]; ok {
		n.addImport(to)
	}
}

// Nodes returns all nodes in the graph keyed by import path.
func (g *DependencyGraph) Nodes() map[string]*Node {
	return g.nodes
}

// FindCycles detects all cyclic import paths within the graph using
// iterative DFS with an explicit call stack.
// Each returned slice represents one cycle: [A, B, C, A].
func (g *DependencyGraph) FindCycles() [][]string {
	visited := make(map[string]bool)
	onStack := make(map[string]bool)
	var cycles [][]string

	var dfs func(path string, stack []string)
	dfs = func(path string, stack []string) {
		if onStack[path] {
			start := -1
			for i, p := range stack {
				if p == path {
					start = i
					break
				}
			}
			if start >= 0 {
				cycle := make([]string, len(stack)-start+1)
				copy(cycle, stack[start:])
				cycle[len(cycle)-1] = path
				cycles = append(cycles, cycle)
			}
			return
		}
		if visited[path] {
			return
		}

		visited[path] = true
		onStack[path] = true
		stack = append(stack, path)

		if n, ok := g.nodes[path]; ok {
			for _, imp := range n.Imports {
				dfs(imp, stack)
			}
		}

		onStack[path] = false
	}

	for path := range g.nodes {
		if !visited[path] {
			dfs(path, nil)
		}
	}

	return cycles
}
