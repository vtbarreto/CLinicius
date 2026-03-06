package discovery

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// LayerType represents the architectural role of a layer.
type LayerType string

const (
	LayerTypeHandler    LayerType = "handler"
	LayerTypeDomain     LayerType = "domain"
	LayerTypeUsecase    LayerType = "usecase"
	LayerTypeRepository LayerType = "repository"
	LayerTypeInfra      LayerType = "infra"
)

// forbidMatrix defines which layer types are forbidden from importing which others.
// Based on Clean Architecture / DDD principles.
var forbidMatrix = map[LayerType][]LayerType{
	LayerTypeHandler:    {LayerTypeRepository, LayerTypeInfra},
	LayerTypeDomain:     {LayerTypeInfra, LayerTypeRepository, LayerTypeHandler},
	LayerTypeUsecase:    {LayerTypeInfra},
	LayerTypeRepository: {LayerTypeHandler},
	LayerTypeInfra:      {},
}

// aliasToType maps well-known folder names to their architectural layer type.
// Covers common naming conventions across Go projects.
var aliasToType = map[string]LayerType{
	// Handler / delivery layer
	"handler":     LayerTypeHandler,
	"handlers":    LayerTypeHandler,
	"controller":  LayerTypeHandler,
	"controllers": LayerTypeHandler,
	"http":        LayerTypeHandler,
	"rest":        LayerTypeHandler,
	"grpc":        LayerTypeHandler,
	"api":         LayerTypeHandler,

	// Domain / core layer
	"domain":   LayerTypeDomain,
	"core":     LayerTypeDomain,
	"model":    LayerTypeDomain,
	"models":   LayerTypeDomain,
	"entity":   LayerTypeDomain,
	"entities": LayerTypeDomain,

	// Use case / application layer
	"usecase":     LayerTypeUsecase,
	"usecases":    LayerTypeUsecase,
	"service":     LayerTypeUsecase,
	"services":    LayerTypeUsecase,
	"application": LayerTypeUsecase,
	"app":         LayerTypeUsecase,

	// Repository / persistence layer
	"repository":   LayerTypeRepository,
	"repositories": LayerTypeRepository,
	"repo":         LayerTypeRepository,
	"repos":        LayerTypeRepository,
	"store":        LayerTypeRepository,
	"storage":      LayerTypeRepository,

	// Infrastructure layer
	"infra":          LayerTypeInfra,
	"infrastructure": LayerTypeInfra,
	"database":       LayerTypeInfra,
	"db":             LayerTypeInfra,
	"cache":          LayerTypeInfra,
	"queue":          LayerTypeInfra,
}

// foundLayer holds raw discovery data before rules are computed.
type foundLayer struct {
	name     string
	path     string
	layerType LayerType
}

// DiscoveredLayer is a layer ready to be written to clinicius.yaml.
type DiscoveredLayer struct {
	Name   string
	Path   string
	Type   LayerType
	Forbid []string
}

// Result holds all layers found during a discovery scan.
type Result struct {
	Layers []DiscoveredLayer
}

// Discover walks root looking for directories whose names match known layer
// patterns and returns a Result with computed forbidden-import rules.
//
// Skips hidden directories, vendor/ and testdata/.
func Discover(root string) (*Result, error) {
	found, err := walkLayers(root)
	if err != nil {
		return nil, err
	}

	// Build a map from LayerType → []relative paths for forbidden resolution.
	typeToPath := make(map[LayerType][]string)
	for _, l := range found {
		typeToPath[l.layerType] = append(typeToPath[l.layerType], l.path)
	}

	seen := make(map[string]bool)
	var layers []DiscoveredLayer

	for _, l := range found {
		if seen[l.path] {
			continue
		}
		seen[l.path] = true

		var forbid []string
		for _, forbiddenType := range forbidMatrix[l.layerType] {
			forbid = append(forbid, typeToPath[forbiddenType]...)
		}

		layers = append(layers, DiscoveredLayer{
			Name:   string(l.layerType),
			Path:   l.path,
			Type:   l.layerType,
			Forbid: forbid,
		})
	}

	return &Result{Layers: layers}, nil
}

func walkLayers(root string) ([]foundLayer, error) {
	var found []foundLayer

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		name := d.Name()
		if shouldSkip(name) {
			return filepath.SkipDir
		}

		layerType, ok := aliasToType[strings.ToLower(name)]
		if !ok {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		found = append(found, foundLayer{
			name:     name,
			path:     rel,
			layerType: layerType,
		})
		return nil
	})

	return found, err
}

func shouldSkip(name string) bool {
	return strings.HasPrefix(name, ".") ||
		name == "vendor" ||
		name == "testdata" ||
		name == "node_modules"
}
