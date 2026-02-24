# CLnicius

> Opinionated architecture governance CLI for Go projects.\
> Created by Vinicius Teixeira.

**CLnicius** (CLI + Vinicius) is a static analysis tool designed to
enforce architectural integrity in Go codebases.

It focuses on something most linters ignore:

**Architectural boundaries and dependency governance.**

------------------------------------------------------------------------

## 🎯 Vision

As systems evolve, architectural erosion happens:

-   Handlers start importing repositories directly
-   Domain logic depends on infrastructure
-   Cycles appear between packages
-   Forbidden libraries leak into core layers

CLnicius exists to:

-   Make architectural violations visible
-   Enforce boundaries automatically
-   Integrate into CI pipelines
-   Scale with growing backend systems

------------------------------------------------------------------------

## 🚀 What It Does

CLnicius:

-   Loads Go packages using module-aware resolution
-   Builds a dependency graph
-   Parses AST for accurate import detection
-   Applies rule engine validations
-   Reports structured architectural violations

------------------------------------------------------------------------

## 📦 Features

-   Layer boundary enforcement
-   Forbidden dependency detection
-   Cyclic dependency detection
-   YAML-configurable rules
-   CI-friendly exit codes
-   JSON output mode
-   Dependency graph export (DOT)

------------------------------------------------------------------------

## 🔧 Example Configuration

`clnicius.yaml`

``` yaml
layers:
  - name: domain
    path: internal/domain
    forbid:
      - internal/infra
      - internal/repository

  - name: handler
    path: internal/handler
    forbid:
      - internal/repository
```

------------------------------------------------------------------------

## ▶ Running

``` bash
clnicius check ./...
```

### Output Example

``` text
❌ Architectural Violation

Layer: handler
File: internal/handler/user.go
Imports: internal/repository/user_repo.go
Rule: handler layer cannot depend on repository layer
```

CI mode:

``` bash
clnicius check ./... --ci
```

Returns exit code `1` if violations exist.

------------------------------------------------------------------------

## 🏗 Internal Architecture

    cmd/
        root.go
        check.go

    internal/
        analyzer/
            loader.go
            graph.go
            ast.go
        rules/
            engine.go
            layer_rule.go
            cycle_rule.go
        reporter/
            console.go
            json.go

------------------------------------------------------------------------

## 🧠 Core Components

### Package Loader

Uses:

-   `golang.org/x/tools/go/packages`

Provides module-aware dependency resolution.

------------------------------------------------------------------------

### AST Analyzer

Uses:

-   `go/parser`
-   `go/ast`
-   `go/token`

Responsible for precise import detection and extensibility.

------------------------------------------------------------------------

### Dependency Graph

In-memory directed graph.

Supports:

-   DFS cycle detection
-   Layer boundary validation
-   Rule-based traversal

------------------------------------------------------------------------

### Rule Engine

Extensible rule interface:

``` go
type Rule interface {
    Name() string
    Validate(graph *DependencyGraph) []Violation
}
```

Allows:

-   Custom rule injection
-   Plugin evolution
-   Clean separation of concerns

------------------------------------------------------------------------

## 🛠 Tech Stack

-   Go 1.22+
-   Cobra (CLI framework)
-   Go AST (`go/ast`)
-   go/packages
-   YAML config parsing
-   Table-driven tests
-   GitHub Actions CI

------------------------------------------------------------------------

## 📊 Design Philosophy

CLnicius is built around:

-   Deterministic analysis
-   Extensibility over rigidity
-   Architectural governance as code
-   CI-first mindset

It is not a style linter.\
It is an architecture guardrail.

------------------------------------------------------------------------

## 🔮 Roadmap

-   [ ] DOT graph export
-   [ ] HTML report
-   [ ] Plugin system
-   [ ] Performance benchmarks
-   [ ] Incremental diff-based analysis

------------------------------------------------------------------------

## 🧪 Testing Strategy

-   Table-driven tests
-   Golden file output validation
-   Dependency fixture projects
-   Benchmark analysis

------------------------------------------------------------------------

## 📄 License

MIT
