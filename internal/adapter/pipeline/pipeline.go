package pipeline

// Pipeline processa o input do usuário antes de enviar para a LLM.
type Pipeline struct{}

// New retorna uma nova instância do Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

type step func(string) (string, error)

// Process executa os passos em ordem:
// 1. validate  — rejeita entradas inválidas
// 2. sanitize  — bloqueia intenções perigosas
// 3. clean     — remove ruído
// 4. normalize — padroniza formato
// 5. optimize  — melhora clareza para a LLM
func (p *Pipeline) Process(input string) (string, error) {
	steps := []step{
		validate,
		sanitize,
		clean,
		normalize,
		optimize,
	}

	text := input
	for _, s := range steps {
		result, err := s(text)
		if err != nil {
			return "", err
		}
		text = result
	}

	return text, nil
}
