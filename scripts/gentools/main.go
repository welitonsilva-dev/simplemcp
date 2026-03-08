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
	mainFile     = "cmd/main.go"
	nativeDir    = "internal/tools/native"
	markerNative = "pacotes de ferramentas nativas"
)

// root retorna a raiz do projeto independente de onde o go generate foi chamado
func root() string {
	_, file, _, _ := runtime.Caller(0)
	// file = .../simplemcp/cmd/gentools/main.go
	// sobe dois níveis: gentools/ → cmd/ → raiz
	return filepath.Join(filepath.Dir(file), "..", "..")
}

func main() {
	r := root()

	main_ := filepath.Join(r, mainFile)
	native := filepath.Join(r, nativeDir)
	gomod  := filepath.Join(r, "go.mod")

	if _, err := os.Stat(main_); err != nil {
		fatalf("erro: %s não encontrado.", main_)
	}

	if !dirExists(native) {
		fatalf("erro: diretório não encontrado (%s).", native)
	}

	module := readModuleName(gomod)
	if module == "" {
		fatalf("erro: não foi possível detectar o module name no go.mod.")
	}

	lines, err := readLines(main_)
	if err != nil {
		fatalf("erro ao ler %s: %v", main_, err)
	}

	entries, err := os.ReadDir(native)
	if err != nil {
		fatalf("erro ao ler %s: %v", native, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subdir := filepath.Join(nativeDir, entry.Name())
		goFiles, _ := filepath.Glob(filepath.Join(r, subdir, "*.go"))
		if len(goFiles) == 0 {
			continue
		}

		importPath := fmt.Sprintf(`"%s/%s"`, module, filepath.ToSlash(subdir))
		importLine := fmt.Sprintf("\t_ %s", importPath)

		if containsLine(lines, importPath) {
			continue
		}

		lines = insertAfterMarker(lines, markerNative, importLine)
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
