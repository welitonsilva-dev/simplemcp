package filesystem

import (
	"fmt"
	"os"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
	"simplemcp/internal/tools/native"
)

func init() {
	tools.GlobalRegistry.Register(&FSRmRf{})
}

type FSRmRf struct{}

func (l *FSRmRf) Name() string {
	return "fs_rmrf"
}

func (l *FSRmRf) Description() string {
	return `
Prioridade de interpretação:

Remover diretório recursivamente (com todo o conteúdo)
Palavras associadas:
- remover pasta com conteúdo
- deletar pasta inteira
- apagar diretório recursivo
- rm -rf
- remover tudo dentro de
- delete folder recursively
- force remove directory

→ usar ferramenta "fs_rmrf"

Descrição:
Ferramenta que remove um diretório e todo o seu conteúdo recursivamente.

Parâmetros:
- path (string, obrigatório): caminho do diretório a ser removido

ATENÇÃO:
- Esta operação é IRREVERSÍVEL
- Remove o diretório e TODOS os arquivos e subdiretórios dentro dele
- Suporta caminhos absolutos e relativos ao CWD atual
- Funciona em Windows e Linux

Uso comum:
- Remover projetos ou pastas inteiras
- Limpar diretórios temporários com conteúdo
`
}

func (l *FSRmRf) Execute(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("parâmetro 'path' é obrigatório")
	}

	resolved := native.ResolvePath(native.CwdState.Get(), path)

	info, err := os.Stat(resolved)
	if err != nil {
		logger.Error("fs_rmrf stat error: %v", err)
		return nil, fmt.Errorf("diretório não encontrado: %s", resolved)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' é um arquivo. Use fs_rm para remover arquivos", resolved)
	}

	if err := os.RemoveAll(resolved); err != nil {
		logger.Error("fs_rmrf error: %v", err)
		return nil, fmt.Errorf("erro ao remover diretório '%s': %v", resolved, err)
	}

	logger.Info("fs_rmrf: diretório removido recursivamente %s", resolved)
	return fmt.Sprintf("Diretório removido recursivamente: %s", resolved), nil
}
