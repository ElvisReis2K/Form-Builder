# Form Builder

Aplicacao full stack para criacao, publicacao e resposta de formularios. Este repositorio foi criado para o desafio tecnico de Full Stack Developer.

## Stack escolhida

- Backend: Go como servico independente.
- Banco: PostgreSQL.
- Frontend: React, TypeScript e Vite.
- UI: MUI.
- Rotas: React Router.
- Estado remoto: TanStack Query.
- Contrato de API: OpenAPI 3.
- Client TypeScript: gerado a partir da especificacao OpenAPI.

## Requisitos locais

- Go 1.25+.
- Node.js 20+.
- Docker com Docker Compose.

## Estrutura

```text
.
├── backend/
│   ├── cmd/server/
│   ├── internal/
│   ├── migrations/
│   └── openapi/
├── frontend/
│   └── src/
├── docker-compose.yml
├── Makefile
└── README.md
```

## Configuracao

Copie `.env.example` para `.env` e ajuste os valores conforme necessario.

Variaveis principais:

- `ADDRESS`: endereco do backend.
- `DATABASE_URL`: conexao do PostgreSQL.
- `FRONTEND_URL`: origem permitida para o frontend.
- `SESSION_SECRET`: segredo usado para gerar o hash dos tokens de sessao.
- `SESSION_TTL_HOURS`: duracao da sessao em horas.
- `COOKIE_SECURE`: use `true` quando a API estiver atras de HTTPS.
- `GOOGLE_CLIENT_ID` e `GOOGLE_CLIENT_SECRET`: credenciais OAuth do Google.

## Banco de dados

Suba o PostgreSQL local:

```bash
make db-up
```

As migrations ficam em `backend/migrations`.

```bash
make migrate-up
```

Para desfazer a ultima migration aplicada:

```bash
make migrate-down
```

## Backend

O backend roda como um servico independente em Go:

```bash
make backend-dev
```

Health check:

```bash
curl http://localhost:8080/healthz
```

Readiness com ping no PostgreSQL:

```bash
curl http://localhost:8080/readyz
```

Comandos do binario:

```bash
cd backend
go run ./cmd/server run
go run ./cmd/server migrate up
go run ./cmd/server migrate down
go run ./cmd/server openapi
```

## Autenticacao por e-mail e senha

A base atual possui cadastro, login, logout e consulta do usuario autenticado. A sessao usa cookie HTTP-only (`form_builder_session`) e persiste apenas o hash HMAC do token no banco.

Cadastro:

```bash
curl -i http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Admin","email":"admin@example.com","password":"password123"}'
```

Login:

```bash
curl -i http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'
```

Usuario autenticado:

```bash
curl -i http://localhost:8080/api/auth/me \
  --cookie "form_builder_session=<cookie-value>"
```

## Frontend

Instale as dependencias e rode o frontend:

```bash
cd frontend
npm install
npm run dev
```

## OpenAPI e client TypeScript

A especificacao OpenAPI sera mantida como parte do fluxo de build do backend.

```bash
make openapi
```

O client TypeScript do frontend sera gerado a partir da especificacao:

```bash
make client
```

Arquivos gerados ficam em `frontend/src/api/generated` e nao devem ser editados manualmente.

## Decisoes iniciais

- Monorepo para simplificar setup local, revisao e DX.
- PostgreSQL por ser multiplataforma, robusto e simples de subir via Docker Compose.
- Vite em vez de Next.js para deixar claro que o backend Go e responsavel pela API, regras de negocio e persistencia.
- Formularios, campos e respostas serao modelados em tabelas relacionais, com `jsonb` para configuracoes e respostas flexiveis.
- Sessoes em banco com token opaco em cookie HTTP-only. O banco armazena apenas o hash HMAC do token, reduzindo impacto se os dados de sessao vazarem.
- Migrations embutidas no binario Go para manter o setup local reproduzivel sem depender de uma CLI externa.

## Status

Ja existe a base real do backend com conexao PostgreSQL, migrations executaveis, modelo de usuarios/sessoes e autenticacao por e-mail/senha. As proximas etapas sao autenticar com Google, implementar CRUD de formularios, publicacao, validacao de respostas e integrar o frontend ao client TypeScript gerado.
