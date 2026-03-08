package filesystem

import (
	"fmt"
	"os"
	"time"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
	"simplemcp/internal/tools/native"
)

func init() {
	tools.GlobalRegistry().Register(&FSTouch{})
}

type FSTouch struct{}

func (l *FSTouch) Name() string {
	return "fs_touch"
}

func (l *FSTouch) Description() string {
	return `
Prioridade de interpretação:

Criar arquivo vazio ou atualizar timestamp
Palavras associadas:
- criar arquivo
- novo arquivo
- arquivo vazio
- touch
- create file
- criar file

→ usar ferramenta "fs_touch"

Descrição:
Ferramenta que cria um arquivo vazio ou atualiza o timestamp de um arquivo existente.

Parâmetros:
- path (string, obrigatório): caminho do arquivo a ser criado ou atualizado

Comportamento:
- Suporta caminhos absolutos e relativos ao CWD atual
- Se o arquivo não existir, cria um arquivo vazio
- Se o arquivo já existir, atualiza seu timestamp (access e modification time)
- Funciona em Windows e Linux

Uso comum:
- Criar arquivos placeholder
- Inicializar arquivos de configuração
- Atualizar timestamps de arquivos existentes
`
}

func (l *FSTouch) Execute(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("parâmetro 'path' é obrigatório")
	}

	resolved := native.ResolvePath(native.CwdState.Get(), path)

	// Tenta abrir ou criar o arquivo
	file, err := os.OpenFile(resolved, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("fs_touch error: %v", err)
		return nil, fmt.Errorf("erro ao criar arquivo '%s': %v", resolved, err)
	}
	file.Close()

	// Atualiza o timestamp (equivalente ao touch em sistemas Unix)
	now := time.Now()
	if err := os.Chtimes(resolved, now, now); err != nil {
		logger.Error("fs_touch chtimes error: %v", err)
		return nil, fmt.Errorf("erro ao atualizar timestamp de '%s': %v", resolved, err)
	}

	logger.Info("fs_touch: arquivo criado/atualizado em %s", resolved)
	return fmt.Sprintf("Arquivo criado/atualizado: %s", resolved), nil
}
