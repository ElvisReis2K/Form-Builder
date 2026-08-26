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
- PostgreSQL 16+ ou Docker com Docker Compose para subir o PostgreSQL localmente.

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
│       ├── app/
│       ├── api/
│       ├── features/
│       ├── lib/
│       └── styles/
├── scripts/
├── docker-compose.yml
├── Makefile
├── CONTRIBUTING.md
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
- `GOOGLE_REDIRECT_URL`: callback autorizado no Google Cloud. Em local, use `http://localhost:8080/api/auth/google/callback`.

## Banco de dados

Configure `DATABASE_URL` apontando para uma instancia PostgreSQL. Voce pode usar um PostgreSQL ja instalado na maquina ou subir o servico auxiliar via Docker Compose.

Opcao com Docker Compose:

```bash
make db-up
```

Opcao com PostgreSQL local:

```bash
createdb form_builder
```

Depois ajuste `.env`:

```bash
DATABASE_URL=postgres://<usuario>:<senha>@localhost:5432/form_builder?sslmode=disable
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

## Autenticacao com Google

O login com Google usa OAuth 2.0 Authorization Code no backend Go. O frontend apenas redireciona para `GET /api/auth/google`; o callback valida `state`, consulta o perfil no Google e cria a mesma sessao HTTP-only usada pelo login por e-mail/senha.

No Google Cloud, crie um OAuth Client do tipo Web application e configure o redirect URI:

```text
http://localhost:8080/api/auth/google/callback
```

Depois preencha no `.env`:

```bash
GOOGLE_CLIENT_ID=<client-id>
GOOGLE_CLIENT_SECRET=<client-secret>
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

Rotas:

- `GET /api/auth/google`: inicia login e redireciona para o Google.
- `GET /api/auth/google/callback`: recebe `code` e `state`, cria/associa o usuario e redireciona para `/admin`.

## CRUD de formularios

O backend ja expõe CRUD autenticado para formularios. Campos sao persistidos junto com o formulario e substituidos em bloco no `PUT`, mantendo a camada de regra no backend. Cada formulario tambem possui metadados de privacidade usados no aviso publico: e-mail do controlador, finalidade do tratamento e politica de retencao.

Rotas autenticadas:

- `GET /api/forms`: lista formularios do usuario autenticado.
- `POST /api/forms`: cria formulario em rascunho.
- `GET /api/forms/{formId}`: consulta formulario do usuario autenticado.
- `PUT /api/forms/{formId}`: atualiza titulo, descricao e campos.
- `DELETE /api/forms/{formId}`: remove formulario.
- `POST /api/forms/{formId}/publish`: publica e gera `publicSlug`.
- `POST /api/forms/{formId}/unpublish`: volta o formulario para rascunho.

Rota publica:

- `GET /api/public/forms/{slug}`: consulta formulario publicado para preenchimento.
- `POST /api/public/forms/{slug}/responses`: envia uma resposta publica para formulario publicado.

Respostas:

- `GET /api/forms/{formId}/responses`: lista respostas recebidas por formulario do usuario autenticado.
- `GET /api/forms/{formId}/responses/export`: exporta respostas em JSON para atendimento administrativo.
- `DELETE /api/forms/{formId}/responses/{responseId}`: exclui uma resposta especifica do formulario.

Criar formulario:

```bash
curl -i http://localhost:8080/api/forms \
  --cookie "form_builder_session=<cookie-value>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Pesquisa de satisfacao",
    "description": "Formulario publico para clientes",
    "controllerEmail": "privacidade@example.com",
    "privacyPurpose": "Coletar feedback para melhoria do atendimento.",
    "retentionPolicy": "As respostas serao mantidas por ate 90 dias apos o envio.",
    "fields": [
      {
        "type": "text",
        "label": "Nome",
        "required": true,
        "placeholder": "Digite seu nome"
      },
      {
        "type": "select",
        "label": "Plano",
        "required": true,
        "options": ["Basic", "Pro", "Enterprise"]
      }
    ]
  }'
```

Publicar formulario:

```bash
curl -i -X POST http://localhost:8080/api/forms/<form-id>/publish \
  --cookie "form_builder_session=<cookie-value>"
```

Consultar formulario publicado:

```bash
curl -i http://localhost:8080/api/public/forms/<public-slug>
```

Enviar resposta publica:

```bash
curl -i http://localhost:8080/api/public/forms/<public-slug>/responses \
  -H "Content-Type: application/json" \
  -d '{
    "privacyAcknowledged": true,
    "answers": {
      "<field-id>": "Ada Lovelace"
    }
  }'
```

Listar respostas no admin:

```bash
curl -i http://localhost:8080/api/forms/<form-id>/responses \
  --cookie "form_builder_session=<cookie-value>"
```

Exportar respostas:

```bash
curl -i http://localhost:8080/api/forms/<form-id>/responses/export \
  --cookie "form_builder_session=<cookie-value>"
```

Excluir uma resposta:

```bash
curl -i -X DELETE http://localhost:8080/api/forms/<form-id>/responses/<response-id> \
  --cookie "form_builder_session=<cookie-value>"
```

## Pacote minimo LGPD

O projeto inclui controles basicos para apoiar transparencia e atendimento de direitos do titular, sem substituir revisao juridica do controlador.

- Pagina publica `/privacidade` com politica de privacidade do sistema.
- Formulario publicado exige e exibe e-mail do controlador, finalidade do tratamento e politica de retencao.
- Publicacao bloqueada quando os metadados de privacidade obrigatorios nao foram preenchidos.
- Envio publico exige `privacyAcknowledged: true` e grava `privacyAcknowledgedAt` na resposta.
- Painel administrativo permite exportar respostas em JSON e excluir uma resposta especifica.
- Cookies usados hoje sao essenciais: sessao HTTP-only e estado temporario do Google OAuth.

## Frontend

Instale as dependencias e rode o frontend:

```bash
cd frontend
npm install
npm run dev
```

O frontend usa `http://localhost:8080` como API por padrao. Para mudar, copie `frontend/.env.example` para `frontend/.env` e ajuste `VITE_API_URL`.

Guarda de arquitetura do frontend:

```bash
cd frontend
npm run lint:architecture
```

Regra do projeto: componentes e paginas nao devem misturar logica com estetica. Estilos MUI ficam em arquivos `*.styles.ts`; rotas e providers ficam em `src/app`; codigo por dominio fica em `src/features`; helpers sem UI ficam em `src/lib`.

Depois do login, `/admin` exibe apenas a lista de formularios salvos. A area administrativa completa fica em `/admin/workspace` e consome a API real para criar, editar, configurar campos, publicar, despublicar e excluir formularios. A pagina publica `/f/:slug` carrega o formulario publicado por slug e envia respostas. A rota `/admin/forms/:formId/responses` lista as respostas recebidas.

## Validacao manual local

Com PostgreSQL, backend e frontend rodando, valide o fluxo principal no navegador:

```text
http://localhost:5173
```

Checklist:

- Criar uma conta por e-mail e senha.
- Entrar e conferir a lista de formularios salvos em `/admin`.
- Acessar a area administrativa em `/admin/workspace`.
- Criar um formulario.
- Preencher e-mail do controlador, finalidade do tratamento e retencao das respostas.
- Adicionar campos ao formulario.
- Salvar o formulario.
- Publicar o formulario.
- Abrir o link publico gerado em `/f/:slug`.
- Confirmar ciencia do aviso de privacidade.
- Enviar uma resposta sem autenticacao.
- Voltar ao admin e abrir `/admin/forms/:formId/responses`.
- Confirmar que a resposta enviada aparece na listagem.
- Exportar respostas em JSON e excluir uma resposta de teste.

Sinais esperados:

- Backend respondendo em `http://localhost:8080/healthz`.
- Readiness respondendo em `http://localhost:8080/readyz`.
- Frontend servindo em `http://localhost:5173`.
- Cookie HTTP-only `form_builder_session` criado apos login.

Se houver um PostgreSQL local usando a porta `5432`, pare esse servico antes de subir o PostgreSQL do Docker ou altere a porta publicada no `docker-compose.yml` e ajuste `DATABASE_URL`.

## OpenAPI e client TypeScript

Sempre que uma rota, payload, schema ou status code da API mudar no backend, atualize primeiro a especificacao OpenAPI gerada:

```bash
make openapi
```

Esse comando executa o backend Go e regrava `backend/openapi/openapi.json`. Depois gere novamente os tipos TypeScript consumidos pelo frontend:

```bash
make client
```

Esse comando roda `npm run generate:api` dentro de `frontend` e atualiza `frontend/src/api/generated/schema.ts` e `frontend/src/api/generated/client.ts` a partir de `backend/openapi/openapi.json`.

Fluxo obrigatorio para mudancas de contrato:

```bash
make openapi && make client
```

Em ambientes sem GNU Make, rode os comandos equivalentes:

```bash
cd backend
go run ./cmd/server openapi > ./openapi/openapi.json
cd ../frontend
npm run generate:api
```

Arquivos gerados ficam em `frontend/src/api/generated`, sao versionados para revisao e nao devem ser editados manualmente. O build do frontend tambem executa `npm run generate:api` antes do typecheck. Os tipos usados nas features de auth, formularios e respostas sao aliases derivados de `frontend/src/api/generated/schema.ts`, e as paginas/features consomem a API por meio das funcoes geradas em `frontend/src/api/generated/client.ts`.

## Decisoes iniciais

- Monorepo para simplificar setup local, revisao e DX.
- PostgreSQL por ser multiplataforma, robusto e adequado para persistencia relacional. O Docker Compose existe apenas como conveniencia para desenvolvimento local.
- Vite em vez de Next.js para deixar claro que o backend Go e responsavel pela API, regras de negocio e persistencia.
- Frontend organizado por `app`, `features`, `styles`, `lib` e `api`, mantendo estilos fora dos componentes e logica fora da camada visual.
- Formularios, campos e respostas serao modelados em tabelas relacionais, com `jsonb` para configuracoes e respostas flexiveis.
- Sessoes em banco com token opaco em cookie HTTP-only. O banco armazena apenas o hash HMAC do token, reduzindo impacto se os dados de sessao vazarem.
- Migrations embutidas no binario Go para manter o setup local reproduzivel sem depender de uma CLI externa.

## Status

Ja existe a base real do backend com conexao PostgreSQL, migrations executaveis, modelo de usuarios/sessoes, autenticacao por e-mail/senha, autenticacao com Google, CRUD autenticado de formularios com publicacao, envio/listagem de respostas e frontend conectado aos fluxos principais.
