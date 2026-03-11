package filesystem

import (
	"fmt"
	"os"

	"simplemcp/internal/adapter/tools"
	"simplemcp/internal/adapter/tools/native"
	"simplemcp/internal/infra/logger"
)

func init() {
	tools.GlobalRegistry().Register(&FSRm{})
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

ATENÇÃO — escolha correta de ferramenta:
- O alvo é um ARQUIVO (ex: foto.jpg, dados.txt, script.sh)? → use fs_rm
- O alvo é uma PASTA/DIRETÓRIO VAZIO? → use fs_rmdir
- O alvo é uma PASTA/DIRETÓRIO COM CONTEÚDO? → use fs_rmrf

Como identificar se é arquivo ou diretório:
- Tem extensão (ex: .txt, .go, .json, .sh)? → é um ARQUIVO → use fs_rm
- Não tem extensão e parece um nome de pasta? → pode ser diretório → use fs_rmdir ou fs_rmrf

Descrição:
Ferramenta que remove um arquivo do sistema de arquivos.

Parâmetros:
- path (string, obrigatório): caminho do arquivo a ser removido
- confirmed (bool, obrigatório para execução): deve ser true para confirmar a remoção
  Só inclua confirmed: true se o usuário deixar explicitamente claro que permite a ação.
  Palavras como "pode", "permite", "confirmo", "eu permito", "pode apagar" indicam confirmação.

Comportamento:
- Suporta caminhos absolutos e relativos ao CWD atual
- Remove APENAS arquivos (não diretórios)
- Para remover diretórios vazios use fs_rmdir
- Para remover diretórios com conteúdo use fs_rmrf
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

	// guarda de confirmação — exige consentimento explícito do usuário
	confirmed, _ := params["confirmed"].(bool)
	if !confirmed {
		return map[string]any{
			"requires_confirmation": true,
			"tool":                  "fs_rm",
			"message":               fmt.Sprintf("ação destrutiva detectada em '%s' — envie novamente com confirmação explícita para prosseguir", path),
		}, nil
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
	return map[string]any{
		"message": fmt.Sprintf("arquivo removido: %s", native.ToHostPath(resolved)),
	}, nil
}
