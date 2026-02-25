package analyzer_test

import (
	"testing"

	"github.com/vtbarreto/CLinicius/internal/analyzer"
)

func TestDependencyGraph_AddNode(t *testing.T) {
	g := analyzer.NewDependencyGraph()

	node := g.AddNode("myapp/pkg/a", []string{"a.go"})
	if node == nil {
		t.Fatal("AddNode returned nil")
	}
	if node.ImportPath != "myapp/pkg/a" {
		t.Errorf("ImportPath = %q, want %q", node.ImportPath, "myapp/pkg/a")
	}

	// Adding same node twice must return the existing one.
	same := g.AddNode("myapp/pkg/a", []string{"other.go"})
	if same != node {
		t.Error("AddNode should return existing node on duplicate")
	}

	if len(g.Nodes()) != 1 {
		t.Errorf("expected 1 node, got %d", len(g.Nodes()))
	}
}

func TestDependencyGraph_AddEdge(t *testing.T) {
	g := analyzer.NewDependencyGraph()
	g.AddNode("myapp/handler", nil)
	g.AddEdge("myapp/handler", "myapp/repository")
	g.AddEdge("myapp/handler", "myapp/repository") // duplicate must be ignored

	node := g.Nodes()["myapp/handler"]
	if len(node.Imports) != 1 {
		t.Errorf("expected 1 import after dedup, got %d", len(node.Imports))
	}
	if node.Imports[0] != "myapp/repository" {
		t.Errorf("Imports[0] = %q, want %q", node.Imports[0], "myapp/repository")
	}
}

func TestDependencyGraph_FindCycles_NoCycle(t *testing.T) {
	g := analyzer.NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", nil)
	g.AddNode("c", nil)
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")

	cycles := g.FindCycles()
	if len(cycles) != 0 {
		t.Errorf("expected no cycles, got %v", cycles)
	}
}

func TestDependencyGraph_FindCycles_DirectCycle(t *testing.T) {
	g := analyzer.NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", nil)
	g.AddEdge("a", "b")
	g.AddEdge("b", "a")

	cycles := g.FindCycles()
	if len(cycles) == 0 {
		t.Fatal("expected at least one cycle, got none")
	}
}

func TestDependencyGraph_FindCycles_LongCycle(t *testing.T) {
	g := analyzer.NewDependencyGraph()
	for _, p := range []string{"a", "b", "c", "d"} {
		g.AddNode(p, nil)
	}
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.AddEdge("c", "d")
	g.AddEdge("d", "a")

	cycles := g.FindCycles()
	if len(cycles) == 0 {
		t.Fatal("expected a cycle in a→b→c→d→a, got none")
	}
}

func TestDependencyGraph_FindCycles_SelfLoop(t *testing.T) {
	g := analyzer.NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddEdge("a", "a")

	cycles := g.FindCycles()
	if len(cycles) == 0 {
		t.Fatal("expected self-loop to be detected as a cycle")
	}
}
