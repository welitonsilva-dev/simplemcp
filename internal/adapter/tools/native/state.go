package native

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// hostMount é o ponto de montagem do disco do host no container.
// Definido aqui para ser compartilhado por todas as tools do pacote native.
//
// O disco do host é mapeado pelo cli.sh / cli.ps1 via:
//
//	-v /:/app/host          (Linux/Mac)
//	-v C:\:/app/host/c      (Windows)
const HostMount = "/app/host"

// CwdState mantém o diretório atual compartilhado entre todas as tools.
//
// É inicializado com CONTAINER_CWD, que é injetado em tempo de execução
// pelo cli.sh / cli.ps1 via flag: -e CONTAINER_CWD=/app/host/home/user/projetos
//
// Se CONTAINER_CWD não for injetado (ex: container subiu manualmente),
// usa hostMount como fallback.
var CwdState = &sharedCWD{
	cwd: func() string {
		if v := os.Getenv("CONTAINER_CWD"); v != "" {
			return v
		}
		return HostMount
	}(),
}

type sharedCWD struct {
	mu  sync.RWMutex
	cwd string
}

func (s *sharedCWD) Get() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cwd
}

func (s *sharedCWD) Set(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cwd = path
}

// ResolvePath resolve um path relativo ou absoluto a partir do cwd atual.
// Centralizado aqui para ser usado por todas as tools do pacote.
//
// Exemplos:
//
//	ResolvePath("/app/host/home/user/projetos", "../")   → /app/host/home/user
//	ResolvePath("/app/host/home/user/projetos", "repo1") → /app/host/home/user/projetos/repo1
//	ResolvePath("/app/host/home/user/projetos", "/tmp")  → /app/host/tmp
func ResolvePath(cwd, path string) string {
	// Já é um path dentro do container
	if strings.HasPrefix(path, HostMount) {
		return filepath.Clean(path)
	}

	// Path absoluto do host (ex: /home/user ou C:\Users) → converte para container
	if filepath.IsAbs(path) {
		return filepath.Clean(filepath.Join(HostMount, filepath.ToSlash(path)))
	}

	// Path relativo (ex: ../ ou subpasta)
	return filepath.Clean(filepath.Join(cwd, path))
}

// ToHostPath converte o path do container para o path legível do host.
//
// Ex: /app/host/home/user/projetos → /home/user/projetos
func ToHostPath(containerPath string) string {
	trimmed := strings.TrimPrefix(containerPath, HostMount)
	if trimmed == "" {
		return "/"
	}
	return trimmed
}
