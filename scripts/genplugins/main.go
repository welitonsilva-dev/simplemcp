//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	mainFile      = "cmd/main.go"
	markerPlugins = "pacotes de ferramentas externas/plugins"
)

func root() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

func main() {
	r := root()

	main_ := filepath.Join(r, mainFile)
	gomod := filepath.Join(r, "go.mod")
	pluginsDir := filepath.Join(r, "..", "simplemcpplugins")

	if _, err := os.Stat(main_); err != nil {
		fatalf("erro: %s não encontrado.", main_)
	}

	// se a pasta de plugins não existir, encerra silenciosamente
	if !dirExists(pluginsDir) {
		return
	}

	module := readPluginsModule(filepath.Join(pluginsDir, "go.mod"))
	if module == "" {
		// fallback: tenta ler do go.mod do simplemcp via replace
		module = readPluginsModuleFromReplace(gomod)
	}
	if module == "" {
		fatalf("erro: não foi possível detectar o module name do simplemcpplugins.")
	}

	lines, err := readLines(main_)
	if err != nil {
		fatalf("erro ao ler %s: %v", main_, err)
	}

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		fatalf("erro ao ler %s: %v", pluginsDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		goFiles, _ := filepath.Glob(filepath.Join(pluginsDir, entry.Name(), "*.go"))
		if len(goFiles) == 0 {
			continue
		}

		importPath := fmt.Sprintf(`"%s/%s"`, module, entry.Name())
		importLine := fmt.Sprintf("\t_ %s", importPath)

		if containsLine(lines, importPath) {
			continue
		}

		lines = insertAfterMarker(lines, markerPlugins, importLine)
		fmt.Printf("adicionado: %s\n", importPath)
	}

	if err := writeLines(main_, lines); err != nil {
		fatalf("erro ao gravar %s: %v", main_, err)
	}
}

func insertAfterMarker(lines []string, marker, line string) []string {
	for i, l := range lines {
		if strings.Contains(l, marker) {
			result := make([]string, 0, len(lines)+1)
			result = append(result, lines[:i+1]...)
			result = append(result, line)
			result = append(result, lines[i+1:]...)
			return result
		}
	}
	return lines
}

func containsLine(lines []string, substr string) bool {
	for _, l := range lines {
		if strings.Contains(l, substr) {
			return true
		}
	}
	return false
}

// readPluginsModule lê o module name do go.mod do simplemcpplugins
func readPluginsModule(gomod string) string {
	return readModuleName(gomod)
}

// readPluginsModuleFromReplace lê o module name do replace no go.mod do simplemcp
func readPluginsModuleFromReplace(gomod string) string {
	f, err := os.Open(gomod)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "require") {
			continue
		}
		// ex: require github.com/weliton/simplemcpplugins v0.0.0...
		if strings.Contains(line, "simplemcpplugins") && !strings.HasPrefix(line, "replace") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[len(fields)-2]
			}
		}
	}
	return ""
}

func readModuleName(gomod string) string {
	f, err := os.Open(gomod)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
