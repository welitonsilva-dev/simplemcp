package native

import "simplemcp/internal/adapter/tools"

func init() {
	tools.GlobalRegistry().Register(&DoubleEcho{})

}

type DoubleEcho struct{}

func (e *DoubleEcho) Name() string {
	return "double_echo"
}

func (e *DoubleEcho) Description() string {
	return `
Prioridade de interpretação:

Duplicar texto  
Palavras associadas:
- duplicar
- duplique
- double
- duble
- duplicar mensagem

→ usar ferramenta "double_echo"

double_echo
Ferramenta que simplesmente duplica o valor do "message" fornecido.

Parâmetros:
- message (string)

Comportamento:
- Retorna exatamente o mesmo texto duplicado recebido no parâmetro "message".
- Não modifica, interpreta ou transforma o conteúdo.

Uso comum:
- Testar o sistema de execução de ferramentas
- Depuração de fluxo de agentes
- Validar comunicação entre componentes
	`
}

func (e *DoubleEcho) Execute(params map[string]interface{}) (interface{}, error) {
	dobre := params["message"].(string) + "," + params["message"].(string)
	return dobre, nil
}
