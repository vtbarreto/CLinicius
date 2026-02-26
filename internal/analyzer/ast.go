package analyzer

import (
	"go/parser"
	"go/token"
	"strings"
)

// FileImports maps a Go source file path to its declared import paths.
type FileImports struct {
	File    string
	Imports []string
}

// ParseFileImports extracts the import declarations from a single Go source file.
// Only import paths are returned; aliases are ignored.
func ParseFileImports(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	imports := make([]string, 0, len(f.Imports))
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		imports = append(imports, path)
	}
	return imports, nil
}

// ParseDirImports walks every .go file in a directory (non-recursive) and
// returns the import declarations grouped by file path.
func ParseDirImports(dir string) ([]FileImports, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var result []FileImports
	for _, pkg := range pkgs {
		for filePath, astFile := range pkg.Files {
			fi := FileImports{File: filePath}
			for _, imp := range astFile.Imports {
				fi.Imports = append(fi.Imports, strings.Trim(imp.Path.Value, `"`))
			}
			result = append(result, fi)
		}
	}
	return result, nil
}
