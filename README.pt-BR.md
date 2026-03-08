# CLinicius

<p align="center">
  <img src="assets/logo.png" alt="CLinicius logo" width="160" style="border-radius: 50%;"/>
</p>

<p align="center">
  <a href="README.md">🇺🇸 English</a>
</p>

> CLI de governança arquitetural para projetos Go.

CLinicius valida fronteiras de camadas e detecta dependências cíclicas em codebases Go — exatamente o que a maioria dos linters silenciosamente ignora.

```bash
clinicius check ./...
```

```
❌ Violação Arquitetural
  Camada:  handler
  Pacote:  myapp/internal/handler
  Importa: myapp/internal/repository
  Regra:   layer-boundary
  Detalhe: handler layer cannot depend on internal/repository

1 violação encontrada.
```

---

## Instalação

### Pré-requisitos

Go 1.22+ instalado. Download em [go.dev/dl](https://go.dev/dl).

### Build a partir do código fonte

```bash
git clone https://github.com/vtbarreto/CLinicius.git
cd CLinicius
go build -o clinicius .
```

Em seguida mova o binário para um diretório no seu `$PATH` (veja abaixo).

---

### Adicionando o binário ao PATH

#### Linux e macOS

Mova para `/usr/local/bin` (requer sudo):

```bash
sudo mv clinicius /usr/local/bin/
```

Ou instale direto no diretório Go bin (sem sudo):

```bash
go build -o ~/go/bin/clinicius .
```

Certifique-se que `~/go/bin` está no PATH. Adicione ao `~/.bashrc` (bash) ou `~/.zshrc` (zsh):

```bash
export PATH="$PATH:$HOME/go/bin"
```

Recarregue o shell:

```bash
source ~/.bashrc   # ou source ~/.zshrc
```

Verifique:

```bash
clinicius --help
```

#### Windows

Build do binário:

```powershell
go build -o clinicius.exe .
```

Mova para um local permanente, por exemplo `C:\tools\`:

```powershell
mkdir C:\tools
move clinicius.exe C:\tools\
```

Adicione `C:\tools` ao PATH do sistema:

1. Abra o **Menu Iniciar** → pesquise por **"Variáveis de Ambiente"**
2. Clique em **"Editar as variáveis de ambiente do sistema"**
3. Clique em **Variáveis de Ambiente...**
4. Em **Variáveis do sistema**, selecione **Path** e clique em **Editar**
5. Clique em **Novo** e adicione `C:\tools`
6. Clique em **OK** em todas as janelas

Abra um novo terminal e verifique:

```powershell
clinicius --help
```

---

## Uso

### `clinicius init` — gerar configuração automaticamente

Execute primeiro em qualquer projeto Go. O CLinicius escaneia a árvore de diretórios, detecta as pastas de camadas pelo nome e gera um `clinicius.yaml` com as regras certas:

```bash
clinicius init
```

```
✅ Generated clinicius.yaml with 4 layer(s) detected:

  handler       internal/handler
                forbids: internal/repository
                forbids: internal/infra
  domain        internal/domain
                forbids: internal/infra
                forbids: internal/repository
                forbids: internal/handler
  repository    internal/repository
                forbids: internal/handler
  infra         internal/infra

--------------------------------------------------
Review the generated file, then run:
  clinicius check ./...
```

Nomes de pastas reconhecidos:

| Tipo de camada | Nomes de pasta |
|---|---|
| handler | `handler`, `handlers`, `controller`, `controllers`, `http`, `rest`, `grpc`, `api` |
| domain | `domain`, `core`, `model`, `entity`, `entities` |
| usecase | `usecase`, `usecases`, `service`, `services`, `application` |
| repository | `repository`, `repositories`, `repo`, `store`, `storage` |
| infra | `infra`, `infrastructure`, `database`, `db`, `cache`, `queue` |

Flags adicionais:

```bash
clinicius init --dry-run                     # visualiza sem escrever
clinicius init --output configs/rules.yaml   # caminho de saída customizado
clinicius init --force                       # sobrescreve arquivo existente
```

---

### `clinicius check` — executar a análise

```bash
# Verifica todos os pacotes do módulo atual
clinicius check ./...

# Usa um arquivo de configuração customizado
clinicius check ./... --config caminho/para/clinicius.yaml

# Modo CI: sai com código 1 se houver violações
clinicius check ./... --ci

# Saída em JSON para integração com ferramentas
clinicius check ./... --json
```

### `--lang` — idioma da saída

O CLinicius detecta automaticamente seu idioma pela variável `LANG` do sistema. Você pode sobrescrever com flag ou variável de ambiente:

```bash
# Via flag (maior prioridade)
clinicius --lang=pt-BR check ./...

# Via variável de ambiente (sessão atual)
CLINICIUS_LANG=pt-BR clinicius check ./...

# Auto-detectado pelo SO (ex: LANG=pt_BR.UTF-8)
clinicius check ./...
```

Para configurar `CLINICIUS_LANG` permanentemente:

**Linux / macOS** — adicione ao `~/.bashrc` ou `~/.zshrc`:
```bash
export CLINICIUS_LANG=pt-BR
```
Depois recarregue: `source ~/.bashrc`

**Windows** (PowerShell):
```powershell
[System.Environment]::SetEnvironmentVariable("CLINICIUS_LANG", "pt-BR", "User")
```

Para verificar o idioma que o SO está reportando:
```bash
echo $LANG   # ex: pt_BR.UTF-8 → CLinicius usa pt-BR automaticamente
```

Idiomas suportados: `en-US` (padrão), `pt-BR`.
Para adicionar um novo idioma, crie `internal/i18n/locales/<tag>.json` — sem alterações de código.

---

## Configuração

O `clinicius.yaml` pode ser gerado com `clinicius init` ou escrito manualmente.

Cada camada mapeia um **prefixo de caminho** para uma lista de **imports proibidos**. Qualquer pacote cujo import path contenha `path` que importe algo que contenha um valor de `forbid` é reportado como violação.

```yaml
layers:
  - name: domain
    path: internal/domain
    forbid:
      - internal/infra
      - internal/repository

  - name: handler
    path: internal/handler
    forbid:
      - internal/repository
```

---

## Por que CLinicius?

A erosão arquitetural é silenciosa. Com o tempo:

- Handlers começam a importar repositórios diretamente
- A lógica de domínio cria dependências com infraestrutura
- Imports cíclicos entre pacotes se acumulam

Linters convencionais (`golangci-lint`, `staticcheck`) não detectam isso. O CLinicius detecta.

---

## Exemplos

### ✅ Estrutura correta

Cada camada importa apenas o que é permitido:

```
internal/
├── handler/
│   └── user_handler.go   → importa apenas domain
├── domain/
│   └── user_service.go   → não importa nada abaixo
├── repository/
│   └── user_repo.go      → importa apenas domain
└── infra/
    └── database.go       → sem imports internos
```

```bash
$ clinicius check ./...

✅ Nenhuma violação arquitetural encontrada.
```

---

### ❌ Estrutura com problemas

`domain` acessa `infra` e `repository` diretamente. `handler` bypassa o domain e fala direto com `repository`:

```
internal/
├── handler/
│   └── user_handler.go   → importa repository ⚠
├── domain/
│   └── user_service.go   → importa infra ⚠ e repository ⚠
├── repository/
│   └── user_repo.go
└── infra/
    └── database.go
```

```bash
$ clinicius --lang=pt-BR check ./...

❌ Violação Arquitetural
  Camada:  domain
  Pacote:  myapp/internal/domain
  Importa: myapp/internal/infra
  Regra:   layer-boundary
  Detalhe: domain layer cannot depend on internal/infra

❌ Violação Arquitetural
  Camada:  domain
  Pacote:  myapp/internal/domain
  Importa: myapp/internal/repository
  Regra:   layer-boundary
  Detalhe: domain layer cannot depend on internal/repository

❌ Violação Arquitetural
  Camada:  handler
  Pacote:  myapp/internal/handler
  Importa: myapp/internal/repository
  Regra:   layer-boundary
  Detalhe: handler layer cannot depend on internal/repository

3 violações encontradas.
```

Um exemplo funcional com esse projeto quebrado está disponível em [`examples/myapp`](./examples/myapp).

---

## Funcionalidades

| Funcionalidade | Status |
|---|---|
| Auto-descoberta de camadas (`clinicius init`) | ✅ |
| Validação de fronteiras entre camadas | ✅ |
| Detecção de dependências cíclicas | ✅ |
| Regras configuráveis via YAML | ✅ |
| Saída em JSON | ✅ |
| Exit codes para CI | ✅ |
| Saída multilíngue (`en-US`, `pt-BR`) | ✅ |
| Exportação de grafo DOT | 🔜 |
| Relatório HTML | 🔜 |
| Sistema de plugins | 🔜 |

---

## Como funciona

1. Carrega pacotes Go usando `golang.org/x/tools/go/packages` (module-aware)
2. Faz parse dos imports por arquivo usando `go/ast`
3. Constrói um grafo dirigido em memória
4. Executa o motor de regras sobre o grafo
5. Reporta violações no stdout (console ou JSON)

O motor de regras é construído em torno de uma interface simples, facilitando a adição de novas regras:

```go
type Rule interface {
    Name() string
    Validate(graph *DependencyGraph, cfg *config.Config) []Violation
}
```

---

## Roadmap

- [ ] Exportação de grafo DOT
- [ ] Relatório HTML
- [ ] Sistema de plugins
- [ ] Análise incremental baseada em diff
- [ ] Benchmarks de performance

---

## Licença

MIT — feito por [Vinicius Teixeira](https://github.com/vtbarreto).
