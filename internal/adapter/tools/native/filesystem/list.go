package filesystem

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"simplemcp/internal/adapter/tools"
	"simplemcp/internal/adapter/tools/native"
	"simplemcp/internal/infra/logger"
)

func init() {
	tools.GlobalRegistry().Register(&FSList{})
}

type FSList struct{}

// ListResult é o output estruturado da tool fs_list.
type ListResult struct {
	Message string   `json:"message"`
	Items   []string `json:"items"`
}

func (l *FSList) Name() string {
	return "fs_list"
}

func (l *FSList) Description() string {
	return `
Prioridade de interpretação:

Listar arquivos e pastas
Palavras associadas:
- listar arquivos
- listar pasta
- mostrar arquivos
- list files
- ls
- dir
- ver pasta
- arquivos ocultos

→ usar ferramenta "fs_list"

Descrição:
Ferramenta que lista arquivos e diretórios do sistema, incluindo arquivos ocultos.

Parâmetros:
- path (string, opcional)

Comportamento:
- Se nenhum path for informado, usa o diretório atual do usuário (CONTAINER_CWD).
- Linux: executa "ls -a [path]"
- Windows: executa "dir /a"

Uso comum:
- Inspecionar diretórios
- Debug de ambiente
- Verificar arquivos criados por ferramentas
`
}

func (l *FSList) Execute(params map[string]interface{}) (interface{}, error) {
	path := native.CwdState.Get()

	if p, ok := params["path"].(string); ok && p != "" {
		path = native.ResolvePath(native.CwdState.Get(), p)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "dir", "/a", path)
	} else {
		cmd = exec.Command("ls", "-a", path)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Error("fs_list error: %v, stderr: %s", err, stderr.String())
		return nil, fmt.Errorf("erro ao listar '%s': %s", path, stderr.String())
	}

	// transforma string com \n em slice, remove . e ..
	var items []string
	for _, item := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if item == "." || item == ".." || item == "" {
			continue
		}
		items = append(items, item)
	}

	hostPath := native.ToHostPath(path)
	return ListResult{
		Message: fmt.Sprintf("Encontrei %d itens em %s", len(items), hostPath),
		Items:   items,
	}, nil
}
