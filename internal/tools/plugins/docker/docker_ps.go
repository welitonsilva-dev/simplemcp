package docker

import (
	"bytes"
	"os/exec"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
)

func init() {
	tools.GlobalRegistry.Register(&DockerPS{})
}

type DockerPS struct{}

func (d *DockerPS) Name() string {
	return "docker_ps"
}

func (d *DockerPS) Description() string {
	return `
Prioridade de interpretação:

Listar containers Docker
Palavras associadas:
- docker ps
- listar containers
- containers rodando
- containers ativos
- mostrar containers
- ver containers
- containers docker
- ps docker

→ usar ferramenta "docker_ps"

Descrição:
Ferramenta que lista os containers Docker em execução.

Parâmetros:
- all (bool, opcional): se true, lista todos os containers incluindo os parados.

Comportamento:
- Por padrão lista apenas containers em execução.
- Se all=true, executa "docker ps -a".

Uso comum:
- Verificar containers ativos
- Debug de ambiente Docker
- Inspecionar estado dos containers
`
}

func (d *DockerPS) Execute(params map[string]interface{}) (interface{}, error) {

	args := []string{"ps"}

	if all, ok := params["all"].(bool); ok && all {
		args = append(args, "-a")
	}

	cmd := exec.Command("docker", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Error("docker_ps error: falha ao executar docker ps: %v, stderr: %s", err, stderr.String())
		return stderr.String(), err
	}

	return out.String(), nil
}
