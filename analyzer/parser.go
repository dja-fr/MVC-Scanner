package analyzer

import (
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// ExtractDependencies parses the file content and returns a list of unique absolute paths of dependencies.
func ExtractDependencies(f string, content []byte, ctx *ProjectContext) []string {
	ext := filepath.Ext(f)
	var lang *sitter.Language
	switch ext {
	case ".java":
		lang = java.GetLanguage()
	case ".js":
		lang = javascript.GetLanguage()
	case ".ts":
		lang = typescript.GetLanguage()
	case ".py":
		lang = python.GetLanguage()
	default:
		return nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	tree := parser.Parse(nil, content)

	var deps []string

	// Helper to find file by suffix
	findFileBySuffix := func(suffix string) string {
		for _, absPath := range ctx.FileLookup {
			normalizedAbs := strings.ReplaceAll(absPath, "\\", "/")
			if strings.HasSuffix(normalizedAbs, "/"+suffix) || normalizedAbs == suffix {
				return absPath
			}
		}
		return ""
	}

	var walk func(*sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		t := n.Type()

		// Java dependencies (imports)
		if ext == ".java" && t == "import_declaration" {
			for i := 0; i < int(n.NamedChildCount()); i++ {
				child := n.NamedChild(i)
				if child.Type() == "scoped_identifier" || child.Type() == "identifier" {
					imp := child.Content(content)
					relPath := strings.ReplaceAll(imp, ".", "/") + ".java"
					if targetFile := findFileBySuffix(relPath); targetFile != "" {
						deps = append(deps, targetFile)
					}
				}
			}
		}

		// JS/TS dependencies
		if (ext == ".js" || ext == ".ts") && t == "import_statement" {
			for i := 0; i < int(n.NamedChildCount()); i++ {
				child := n.NamedChild(i)
				if child.Type() == "string" {
					impPath := child.Content(content)
					impPath = strings.Trim(impPath, "\"'`")
					// Resolve relative path
					if strings.HasPrefix(impPath, ".") || strings.HasPrefix(impPath, "/") {
						dir := filepath.Dir(f)
						resolvedAbs := filepath.Clean(filepath.Join(dir, impPath))

						extensions := []string{"", ".js", ".ts", "/index.js", "/index.ts"}
						for _, e := range extensions {
							cand := resolvedAbs + e
							candRel, _ := filepath.Rel(ctx.AbsRoot, cand)
							candNorm := strings.ReplaceAll(candRel, "\\", "/")
							if targetFile, exists := ctx.FileLookup[candNorm]; exists {
								deps = append(deps, targetFile)
								break
							}
						}
					} else {
						// Fallback: try suffix
						extensions := []string{".js", ".ts", "/index.js", "/index.ts"}
						for _, e := range extensions {
							candSuffix := impPath + e
							if targetFile := findFileBySuffix(candSuffix); targetFile != "" {
								deps = append(deps, targetFile)
								break
							}
						}
					}
				}
			}
		}

		// Python dependencies
		if ext == ".py" && (t == "import_from_statement" || t == "import_statement") {
			for i := 0; i < int(n.NamedChildCount()); i++ {
				child := n.NamedChild(i)
				if child.Type() == "dotted_name" {
					imp := child.Content(content)
					relPath := strings.ReplaceAll(imp, ".", "/") + ".py"
					if targetFile := findFileBySuffix(relPath); targetFile != "" {
						deps = append(deps, targetFile)
					}
				}
			}
		}
		for i := 0; i < int(n.ChildCount()); i++ {
			walk(n.Child(i))
		}
	}
	walk(tree.RootNode())

	// Remove duplicates
	depSet := make(map[string]bool)
	var uniqueDeps []string
	for _, d := range deps {
		if !depSet[d] {
			depSet[d] = true
			uniqueDeps = append(uniqueDeps, d)
		}
	}

	return uniqueDeps
}
