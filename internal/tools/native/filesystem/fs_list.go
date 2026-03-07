package filesystem

import (
	"bytes"
	"os/exec"
	"runtime"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
	"simplemcp/internal/tools/native"
)

func init() {
	tools.GlobalRegistry.Register(&FSList{})
}

type FSList struct{}

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

	// Usa o CWD compartilhado como padrão
	// ResolvePath e native.CwdState estão definidos em state.go (mesmo pacote)
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

	err := cmd.Run()

	if err != nil {
		logger.Error("fs_list error: %v, stderr: %s", err, stderr.String())
		return stderr.String(), err
	}

	return out.String(), nil
}
