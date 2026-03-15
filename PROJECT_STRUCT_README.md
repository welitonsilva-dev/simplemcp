
```
├── Dockerfile
├── Dockerfile.standalone
├── FLOW_README.md
├── PLUGIN_README.md
├── PROJECT_STRUCT_README.md
├── README.md
├── RUN_README.md
├── build.sh
├── cli.ps1
├── cli.sh
├── cmd
│   ├── logs
│   │   └── humancli-server.log
│   └── server
│       └── main.go
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
│   ├── adapter
│   │   ├── llm
│   │   │   ├── Client.go
│   │   │   ├── parser.go
│   │   │   └── prompt.go
│   │   ├── pipeline
│   │   │   ├── cleaner.go
│   │   │   ├── normalize.go
│   │   │   ├── optimizer.go
│   │   │   ├── pipeline.go
│   │   │   ├── sanitizer.go
│   │   │   └── validator.go
│   │   └── tools
│   │       ├── native
│   │       │   ├── echo
│   │       │   │   ├── double_echo.go
│   │       │   │   └── echo.go
│   │       │   ├── filesystem
│   │       │   │   ├── cd.go
│   │       │   │   ├── list.go
│   │       │   │   ├── mkdir.go
│   │       │   │   ├── mr.go
│   │       │   │   ├── rmdir.go
│   │       │   │   └── touch.go
│   │       │   ├── state.go
│   │       │   └── tool_list.go
│   │       └── registry.go
│   ├── domain
│   │   ├── message
│   │   │   └── message.go
│   │   ├── plan
│   │   │   └── plan.go
│   │   └── tool
│   │       ├── registry.go
│   │       └── tool.go
│   ├── infra
│   │   ├── config
│   │   │   └── config.go
│   │   ├── logger
│   │   │   └── logger.go
│   │   └── server
│   │       ├── handler.go
│   │       ├── middleware.go
│   │       ├── ratelimit.go
│   │       ├── server.go
│   │       └── timeout.go
│   └── usecase
│       └── agent
│           └── agent.go
├── logs
│   └── humancli-server.log
├── scripts
│   ├── genplugins
│   │   └── main.go
│   └── gentools
│       └── main.go
└── sdk
    └── sdk.go
```