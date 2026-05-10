# 📚 Biblioteca API

CRUD de livros em **Go + PostgreSQL** com pipeline completo de CI/CD em **GitHub Actions** e publicação no **Docker Hub**.

> Trabalho da disciplina **Infraestrutura de TIC** — Prof. Gabriel Castellani de Oliveira (FURB).

---

## 👥 Integrantes
<!-- Substitua pelos nomes do grupo -->
- Fulano de Tal
- Beltrano de Tal
- Sicrano de Tal

---

## 🐳 Imagem no Docker Hub
<!-- Substitua "GuiNaumann" pelo seu usuário real do Docker Hub -->
**`GuiNaumann/biblioteca-api:latest`** → https://hub.docker.com/r/GuiNaumann/biblioteca-api

---

## 🧱 Stack

| Camada | Tecnologia |
|--------|-----------|
| Linguagem | Go 1.22 (stdlib + lib/pq) |
| Banco | PostgreSQL 16 |
| Container | Docker (multi-stage build) |
| Orquestração | Docker Compose |
| CI/CD | GitHub Actions |
| Registry | Docker Hub |

---

## 📁 Estrutura do projeto

```
seminario-infra/
├── .github/workflows/ci.yml    # Pipeline CI/CD
├── handlers/
│   ├── book_handler.go         # HTTP handlers
│   └── book_handler_test.go    # 8 testes unitários
├── models/
│   └── book.go                 # Struct Book
├── repository/
│   └── book_repo.go            # Queries SQL no PostgreSQL
├── main.go                     # Setup: DB + rotas + servidor
├── init.sql                    # Cria tabela + seed
├── Dockerfile                  # Multi-stage build
├── docker-compose.yml          # API + PostgreSQL
├── .gitignore
├── go.mod / go.sum
└── README.md
```

---

## 🚀 Como rodar

### 1) Com Docker Compose (recomendado)
```bash
docker compose up --build
```
A API sobe em `http://localhost:8080` e o Postgres em `localhost:5432`.

### 2) Com a imagem oficial do Docker Hub
```bash
DOCKER_IMAGE=GuiNaumann/biblioteca-api:latest docker compose up
```

### 3) Sem Docker (PostgreSQL local)
```bash
export DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=biblioteca sslmode=disable"
go run .
```

---

## 🔌 Endpoints

| Método | Rota         | Descrição              |
|--------|--------------|------------------------|
| GET    | /health      | Health check           |
| GET    | /books       | Listar todos os livros |
| POST   | /books       | Criar livro            |
| GET    | /books/{id}  | Buscar por ID          |
| PUT    | /books/{id}  | Atualizar              |
| DELETE | /books/{id}  | Remover                |

### Exemplos
```bash
# Listar
curl http://localhost:8080/books

# Criar
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"DDD","author":"Eric Evans","year":2003,"available":true}'

# Buscar por ID
curl http://localhost:8080/books/1

# Atualizar
curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Go","author":"Alan Donovan","year":2015,"available":false}'

# Deletar
curl -X DELETE http://localhost:8080/books/2
```

---

## 🧪 Testes

```bash
go test -v -cover ./...
```

São **8 testes unitários** em `handlers/book_handler_test.go` usando um mock da interface `BookStore`:
1. `TestListBooks_OK` — listar livros
2. `TestGetBook_Found` — buscar livro existente
3. `TestGetBook_NotFound` — 404 para livro inexistente
4. `TestCreateBook_OK` — criar livro
5. `TestCreateBook_InvalidJSON` — 400 para body inválido
6. `TestUpdateBook_OK` — atualizar livro
7. `TestDeleteBook_OK` — deletar livro
8. `TestBookHandler_InvalidID` — 400 para id não numérico

---

## ⚙️ Pipeline CI/CD

```
push / PR
    │
    ├──────────┬──────────┐
    ▼          ▼          
 ┌──────┐  ┌──────┐       
 │build │  │ test │   ← jobs PARALELOS
 │go1.21│  │  8   │
 │go1.22│  │testes│
 └───┬──┘  └───┬──┘
     │         │
     └────┬────┘
          ▼  (só na branch main)
     ┌─────────┐
     │ docker  │  ← Build + push para Docker Hub
     │  build  │     tags: latest + ${{ github.sha }}
     │ + push  │
     └─────────┘
```

### Jobs do workflow
- **build** → matriz `[1.21, 1.22]` + `actions/upload-artifact` (binário)
- **test** → roda os 8 testes em paralelo
- **docker** → `needs: [build, test]`, login + push no Docker Hub (latest + SHA)

### Secrets configurados no repositório
| Secret | Uso |
|--------|-----|
| `DOCKERHUB_USERNAME` | Usuário do Docker Hub |
| `DOCKERHUB_TOKEN` | Access Token (nunca a senha!) |
| `API_KEY` | Exemplo de variável sensível |

---

## 📝 Respostas das perguntas do trabalho

### Tarefa 3 — O que acontece se um teste falhar propositalmente?
Se quebrarmos um `t.Fatalf` em qualquer teste, o `go test` retorna exit code ≠ 0, o step "Rodar testes" do job `test` falha, **o job inteiro fica vermelho**, e por consequência o job `docker` (que tem `needs: [build, test]`) **nem chega a rodar**. Resultado: nenhum push de imagem nova é feito. Foi exatamente isso que vimos quando alteramos o `TestGetBook_Found` pra esperar `Title` errado: pipeline ficou vermelho e a CD não disparou.

### Tarefa 4 — Em que cenário real o `upload-artifact` é útil?
Para guardar o binário compilado, relatórios de cobertura, screenshots de testes E2E, dumps de logs em caso de falha, etc. Útil quando o build é caro/lento e a equipe de QA precisa baixar exatamente o mesmo binário que rodou no CI, sem ter que recompilar localmente. Também serve como rastreabilidade — dá pra baixar o binário de qualquer commit antigo direto do GitHub.

### Tarefa 5 — Por que nunca commitar credenciais no código?
Porque o git tem **histórico imutável**: mesmo apagando o segredo num commit posterior, ele continua disponível em qualquer clone do repo. Repositórios públicos são varridos por bots em segundos — token vazado vira mineração de cripto, vazamento de dados, conta bloqueada por uso indevido. Secrets do GitHub Actions são criptografados, expostos só em runtime e mascarados nos logs.

### Tarefa 6 — Diferença de comportamento entre versões
No Go 1.21 e 1.22 o código compila igual; nada muda no comportamento da nossa API. Se usássemos a nova `range int` (introduzida no 1.22) só compilaria na 1.22. A matriz é importante para detectar essas regressões antes de produção.

### Tarefa 8 — Por que paralelismo importa em CI?
Reduz o **tempo total do pipeline**. Build e testes não dependem um do outro — rodar em série gasta o dobro de tempo. Pipeline rápido = feedback rápido = menos contexto perdido pelo dev.

### Tarefa 9 — Diferença entre tag `latest` e tag por SHA?
- **`latest`**: ponteiro mutável, sempre aponta para a última versão publicada. Usar em ambientes de desenvolvimento/teste, ou quando quer "sempre o mais recente". Perigoso em produção porque a imagem muda sob seus pés.
- **`${{ github.sha }}`**: imutável, identifica unicamente o commit. Use em produção, em rollbacks, em auditoria — você sabe exatamente qual código está rodando.

---

## 🎬 Roteiro de apresentação (10–15 min)
1. **Projeto** (1 min) — CRUD de livros em Go + PostgreSQL
2. **Walkthrough do código** (3 min) — `main.go` (setup) → `handlers` → `repository` → `models`
3. **Demonstração ao vivo** (5 min):
   - `docker compose up` → mostrar API rodando + curl nos endpoints
   - Push em `develop` → mostrar pipeline rodando no GitHub
   - Mostrar imagem publicada no Docker Hub
4. **Explicar cada job/step do `ci.yml`** (3 min)
5. **Aprendizados e dificuldades** (2 min)
