package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"simplemcp/internal/adapter/tools"
	"simplemcp/internal/adapter/tools/native"
	"simplemcp/internal/infra/logger"
)

func init() {
	tools.GlobalRegistry().Register(&FSCd{})
}

type FSCd struct{}

func (l *FSCd) Name() string {
	return "fs_cd"
}

func (l *FSCd) Description() string {
	return `
Prioridade de interpretação:

Mudar de diretório
Palavras associadas:
- mudar de diretório
- mudar pasta
- entrar na pasta
- navegar para
- cd
- change directory
- ir para

→ usar ferramenta "fs_cd"

Descrição:
Ferramenta que muda o diretório de trabalho atual do usuário.

Parâmetros:
- path (string, obrigatório): caminho do diretório de destino

Comportamento:
- Suporta caminhos absolutos (ex: /home/user ou C:\Users\user)
- Suporta caminhos relativos (ex: ../pasta ou ..\pasta)
- Suporta "~" para o diretório home do usuário
- Atualiza o CWD compartilhado (CONTAINER_CWD) após a mudança
- Retorna erro se o diretório não existir ou não for acessível
- Compatível com Windows, Linux e macOS

Uso comum:
- Navegar entre diretórios
- Definir contexto para outras ferramentas (fs_list, fs_move, etc.)
`
}

// expandHome substitui "~" pelo diretório home do usuário,
// compatível com Windows, Linux e macOS.
func expandHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("não foi possível determinar o diretório home: %v", err)
	}

	// Substitui apenas o "~" inicial, preservando o restante do path
	return filepath.Join(home, path[1:]), nil
}

func (l *FSCd) Execute(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("parâmetro 'path' é obrigatório")
	}

	// Expande "~" para o home do usuário (cross-platform)
	expanded, err := expandHome(path)
	if err != nil {
		return nil, err
	}

	// Resolve o path relativo ao CWD atual
	resolved := native.ResolvePath(native.CwdState.Get(), expanded)

	// Normaliza separadores de path para o OS atual (\ no Windows, / no Unix)
	resolved = filepath.Clean(resolved)

	// Verifica se o diretório existe e é acessível
	info, err := os.Stat(resolved)
	if err != nil {
		logger.Error("fs_cd error: %v", err)
		return nil, fmt.Errorf("diretório não encontrado: %s", resolved)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("o caminho informado não é um diretório: %s", resolved)
	}

	// Atualiza o CWD compartilhado
	native.CwdState.Set(resolved)

	logger.Info("fs_cd: diretório alterado para %s", resolved)
	return fmt.Sprintf("Diretório alterado para: %s", resolved), nil
}
