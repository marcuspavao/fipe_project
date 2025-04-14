fipe_project/
├── internal/             # Código privado da aplicação
│   ├── database/         # Conexão com banco de dados
│   │   └── mongodb.go
│   ├── models/           # Definição de estruturas de dados
│   │   └── fipe.go       # Modelos relacionados à tabela FIPE
│   ├── handlers/         # Manipuladores HTTP
│   │   └── fipe.go       # Handlers para endpoints FIPE
│   └── services/         # Lógica de negócios
│       └── fipe.go       # Serviços relacionados à FIPE
├── frontend/             # Interface do usuário
│   ├── static/           # Arquivos estáticos
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └──                   # Templates HTML
├── configs/              # Arquivos de configuração
├── docs/                 # Documentação
├── go.mod                # Dependências Go
└── go.sum
└── main.go       # Arquivo principal
