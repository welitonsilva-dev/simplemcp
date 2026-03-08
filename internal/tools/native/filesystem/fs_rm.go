package filesystem

import (
	"fmt"
	"os"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
	"simplemcp/internal/tools/native"
)

func init() {
	tools.GlobalRegistry.Register(&FSRm{})
}

type FSRm struct{}

func (l *FSRm) Name() string {
	return "fs_rm"
}

func (l *FSRm) Description() string {
	return `
Prioridade de interpretação:

Remover arquivo
Palavras associadas:
- remover arquivo
- deletar arquivo
- apagar arquivo
- excluir arquivo
- rm
- delete file
- remove file

→ usar ferramenta "fs_rm"

Descrição:
Ferramenta que remove um arquivo do sistema de arquivos.

Parâmetros:
- path (string, obrigatório): caminho do arquivo a ser removido

Comportamento:
- Suporta caminhos absolutos e relativos ao CWD atual
- Remove APENAS arquivos (não diretórios)
- Para remover diretórios, use fs_rmdir ou fs_rmrf
- Retorna erro se o path apontar para um diretório
- Funciona em Windows e Linux

Uso comum:
- Remover arquivos desnecessários
- Limpar arquivos temporários
`
}

func (l *FSRm) Execute(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("parâmetro 'path' é obrigatório")
	}

	resolved := native.ResolvePath(native.CwdState.Get(), path)

	info, err := os.Stat(resolved)
	if err != nil {
		logger.Error("fs_rm stat error: %v", err)
		return nil, fmt.Errorf("arquivo não encontrado: %s", resolved)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("'%s' é um diretório. Use fs_rmdir ou fs_rmrf para remover diretórios", resolved)
	}

	if err := os.Remove(resolved); err != nil {
		logger.Error("fs_rm error: %v", err)
		return nil, fmt.Errorf("erro ao remover arquivo '%s': %v", resolved, err)
	}

	logger.Info("fs_rm: arquivo removido %s", resolved)
	return fmt.Sprintf("Arquivo removido: %s", resolved), nil
}
