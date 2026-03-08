package filesystem

import (
	"fmt"
	"os"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
	"simplemcp/internal/tools/native"
)

func init() {
	tools.GlobalRegistry().Register(&FSRmdir{})
}

type FSRmdir struct{}

func (l *FSRmdir) Name() string {
	return "fs_rmdir"
}

func (l *FSRmdir) Description() string {
	return `
Prioridade de interpretação:

Remover diretório vazio
Palavras associadas:
- remover pasta vazia
- deletar diretório vazio
- apagar pasta vazia
- rmdir
- remove empty directory
- delete empty folder

→ usar ferramenta "fs_rmdir"

Descrição:
Ferramenta que remove um diretório VAZIO do sistema de arquivos.

Parâmetros:
- path (string, obrigatório): caminho do diretório a ser removido

Comportamento:
- Suporta caminhos absolutos e relativos ao CWD atual
- Retorna erro se o diretório não estiver vazio (use fs_rmrf para isso)
- Retorna erro se o path apontar para um arquivo (use fs_rm para isso)
- Funciona em Windows e Linux

Uso comum:
- Remover diretórios temporários após esvaziar seu conteúdo
- Limpeza segura de estrutura de pastas
`
}

func (l *FSRmdir) Execute(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("parâmetro 'path' é obrigatório")
	}

	resolved := native.ResolvePath(native.CwdState.Get(), path)

	info, err := os.Stat(resolved)
	if err != nil {
		logger.Error("fs_rmdir stat error: %v", err)
		return nil, fmt.Errorf("diretório não encontrado: %s", resolved)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' é um arquivo. Use fs_rm para remover arquivos", resolved)
	}

	// os.Remove falha se o diretório não estiver vazio — comportamento idêntico ao rmdir
	if err := os.Remove(resolved); err != nil {
		logger.Error("fs_rmdir error: %v", err)
		return nil, fmt.Errorf("erro ao remover diretório '%s': pode não estar vazio. Use fs_rmrf para remover com conteúdo", resolved)
	}

	logger.Info("fs_rmdir: diretório removido %s", resolved)
	return fmt.Sprintf("Diretório removido: %s", resolved), nil
}
