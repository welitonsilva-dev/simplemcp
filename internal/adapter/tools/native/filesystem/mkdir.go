package filesystem

import (
	"fmt"
	"os"

	"simplemcp/internal/adapter/tools"
	"simplemcp/internal/adapter/tools/native"
	"simplemcp/internal/infra/logger"
)

func init() {
	tools.GlobalRegistry().Register(&FSMkdir{})
}

type FSMkdir struct{}

func (l *FSMkdir) Name() string {
	return "fs_mkdir"
}

func (l *FSMkdir) Description() string {
	return `
Prioridade de interpretação:

Criar diretório/pasta
Palavras associadas:
- criar pasta
- criar diretório
- nova pasta
- mkdir
- make directory
- criar folder

→ usar ferramenta "fs_mkdir"

Descrição:
Ferramenta que cria um novo diretório no sistema de arquivos.

Parâmetros:
- path (string, obrigatório): caminho do diretório a ser criado
- parents (bool, opcional): se true, cria diretórios intermediários (como mkdir -p). Padrão: false

Comportamento:
- Suporta caminhos absolutos e relativos ao CWD atual
- Com parents=true, cria toda a árvore de diretórios necessária
- Retorna erro se o diretório já existir (quando parents=false)

Uso comum:
- Criar estrutura de pastas para projetos
- Organizar arquivos em subdiretórios
`
}

func (l *FSMkdir) Execute(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("parâmetro 'path' é obrigatório")
	}

	resolved := native.ResolvePath(native.CwdState.Get(), path)

	parents, _ := params["parents"].(bool)

	var err error
	if parents {
		err = os.MkdirAll(resolved, 0755)
	} else {
		err = os.Mkdir(resolved, 0755)
	}

	if err != nil {
		logger.Error("fs_mkdir error: %v", err)
		return nil, fmt.Errorf("erro ao criar diretório '%s': %v", resolved, err)
	}

	logger.Info("fs_mkdir: diretório criado em %s", resolved)
	return fmt.Sprintf("Diretório criado: %s", resolved), nil
}
