package filesystem

import (
	"fmt"
	"os"

	"humancli-server/internal/adapter/tools"
	"humancli-server/internal/adapter/tools/native"
	"humancli-server/internal/infra/logger"
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
- - confirmed (bool, obrigatório para execução): deve ser true para confirmar a remoção.
  REGRA ESTRITA: só coloque confirmed: true se o usuário escrever EXATAMENTE uma dessas frases:
  "eu permito", "pode apagar", "confirmo", "sim pode", "pode remover".
  Se o usuário apenas pediu para remover SEM usar essas frases, NÃO inclua confirmed. true — isso é uma ação destrutiva, exige consentimento explícito do usuário.

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

	// guarda de confirmação — exige consentimento explícito do usuário
	confirmed, _ := params["confirmed"].(bool)
	if !confirmed {
		return map[string]any{
			"requires_confirmation": true,
			"tool":                  "fs_rmdir",
			"message":               fmt.Sprintf("ação destrutiva detectada em '%s' — envie novamente com confirmação explícita para prosseguir", path),
		}, nil
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
	return map[string]any{
		"message": fmt.Sprintf("diretório removido: %s", native.ToHostPath(resolved)),
	}, nil
}
