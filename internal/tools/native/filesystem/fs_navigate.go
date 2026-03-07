package filesystem

import (
	"fmt"
	"os"
	"strings"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
	"simplemcp/internal/tools/native"
)

func init() {
	tools.GlobalRegistry.Register(&FSNavigate{})
}

type FSNavigate struct{}

func (n *FSNavigate) Name() string {
	return "fs_navigate"
}

func (n *FSNavigate) Description() string {
	return `
Prioridade de interpretação:

Navegar pelo sistema de arquivos
Palavras associadas:
- entrar na pasta
- ir para
- voltar
- navegar
- cd
- mudar diretório
- change directory
- onde estou
- diretório atual
- pwd
- trocar disco
- ir pro disco

→ usar ferramenta "fs_navigate"

Descrição:
Ferramenta que navega pelo sistema de arquivos do host mapeado no container.
O host é acessível em /app/host (Linux/Mac) ou /app/host/<letra> (Windows, ex: /app/host/c).

Parâmetros:
- action (string, obrigatório): "cd" | "pwd" | "drives"
- path (string, opcional): usado na action "cd"

Comportamento:
- pwd    : retorna o diretório atual
- cd     : navega para o path informado (relativo ou absoluto)
- drives : lista os discos disponíveis mapeados em /app/host

Uso comum:
- Navegar entre pastas e discos
- Identificar onde o usuário está
- Preparar contexto para outras ferramentas (fs_list, etc.)
`
}

func (n *FSNavigate) Execute(params map[string]interface{}) (interface{}, error) {

	action, _ := params["action"].(string)
	action = strings.ToLower(strings.TrimSpace(action))

	path, _ := params["path"].(string)

	switch action {

	case "pwd":
		cwd := native.CwdState.Get()
		return fmt.Sprintf("Diretório atual: %s\n(container: %s)", native.ToHostPath(cwd), cwd), nil

	case "cd":
		return navigateCD(path)

	case "drives":
		return listDrives()

	default:
		return nil, fmt.Errorf("action inválida: %q — use: cd, pwd, drives", action)
	}
}

func navigateCD(path string) (string, error) {
	if path == "" {
		logger.Error("fs_navigate error: path vazio para cd")
		return "", fmt.Errorf("informe um path. Ex: cd ../ ou cd /home/user/projetos")
	}

	resolved := native.ResolvePath(native.CwdState.Get(), path)

	info, err := os.Stat(resolved)
	if err != nil {
		logger.Error("fs_navigate error: diretório não encontrado: %s", resolved)
		return "", fmt.Errorf("diretório não encontrado: %s", resolved)
	}
	if !info.IsDir() {
		logger.Error("fs_navigate error: %s não é um diretório", resolved)
		return "", fmt.Errorf("%s não é um diretório", resolved)
	}

	native.CwdState.Set(resolved)

	return fmt.Sprintf("✔ Navegou para: %s", resolved), nil
}

func listDrives() (string, error) {
	entries, err := os.ReadDir(native.HostMount)
	if err != nil {
		logger.Error("fs_navigate error: falha ao listar discos em %s: %v", native.HostMount, err)
		return "", fmt.Errorf("erro ao listar discos: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("Discos disponíveis:\n")

	found := false
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) == 1 {
			sb.WriteString(fmt.Sprintf("  %s:\\ → /app/host/%s\n", strings.ToUpper(e.Name()), e.Name()))
			found = true
		}
	}

	if !found {
		sb.WriteString("  / → /app/host\n")
	}

	return sb.String(), nil
}
