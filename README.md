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
- `SESSION_SECRET`: segredo usado nas sessoes.
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

## Backend

O backend roda como um servico independente em Go:

```bash
make backend-dev
```

Health check:

```bash
curl http://localhost:8080/healthz
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

## Status

Este commit inicial prepara a base do projeto. As proximas etapas sao implementar autenticacao, CRUD de formularios, publicacao, validacao de respostas, OpenAPI gerado pelo backend e client TypeScript gerado no frontend.
