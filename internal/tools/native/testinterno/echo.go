package native

import "simplemcp/internal/tools"

func init() {
	tools.GlobalRegistry.Register(&Echo{})
}

type Echo struct {
}

func (e *Echo) Name() string {
	return "echo"
}

func (e *Echo) Description() string {
	return `
Prioridade de interpretação:

Repetir texto  
Palavras associadas:
- repetir
- repita
- repeat
- copie
- copy
- echo
- repit

→ usar ferramenta "echo" 

echo  
Ferramenta que simplesmente repete o valor do "message" fornecido.

Parâmetros:
- message (string)

Comportamento:
- Retorna exatamente o mesmo texto recebido no parâmetro "message".
- Não modifica, interpreta ou transforma o conteúdo.

Uso comum:
- Testar o sistema de execução de ferramentas
- Depuração de fluxo de agentes
- Validar comunicação entre componentes
	`
}

func (e *Echo) Execute(params map[string]interface{}) (interface{}, error) {
	return params["message"], nil
}
