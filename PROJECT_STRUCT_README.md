.
├── Dockerfile
├── Dockerfile.standalone
├── FLOW_README.md
├── PLUGIN_README.md
├── README.md
├── RUN_README.md
├── build.sh
├── cli.ps1
├── cli.sh
├── cmd
│   ├── logs
│   │   └── simplemcp.log
│   └── main.go // Excuta o servidor para o usuário
├── docker
│   ├── docker-clean.md
│   ├── docker-clean.sh
│   └── ollama
│       ├── entrypoint.md
│       └── entrypoint.sh
├── docker-compose.yml
├── exemple.env
├── generate.go
├── go.mod
├── internal
│   ├── agent
│   │   ├── agent.go
│   │   └── planner.go
│   ├── config
│   │   └── config.go
│   ├── llm
│   │   ├── Client.go
│   │   └── prompt.go
│   ├── logger
│   │   └── logger.go
│   ├── pipeline
│   │   ├── cleaner.go
│   │   ├── optimizer.go
│   │   ├── reduce.go
│   │   └── sanitizer.go
│   ├── protocol
│   │   ├── request.go
│   │   └── response.go
│   ├── server
│   │   └── handler.go
│   └── tools
│       ├── native
│       │   ├── filesystem
│       │   │   ├── fs_cd.go
│       │   │   ├── fs_list.go
│       │   │   ├── fs_mkdir.go
│       │   │   ├── fs_rm.go
│       │   │   ├── fs_rmdir.go
│       │   │   ├── fs_rmrf.go
│       │   │   └── fs_touch.go
│       │   ├── state.go
│       │   ├── testinterno
│       │   │   ├── double_echo.go
│       │   │   └── echo.go
│       │   └── tool_list.go
│       ├── registry.go
│       └── tool.go
├── logs
│   └── simplemcp.log
├── scripts
│   ├── genplugins
│   │   └── main.go
│   └── gentools
│       └── main.go
├── sdk
│   └── sdk.go
└── utils
    └── json.go