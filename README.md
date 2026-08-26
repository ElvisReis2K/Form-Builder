# Form Builder

Aplicação full stack para criação, publicação e resposta de formulários. Este repositório foi criado para o desafio técnico de Full Stack Developer.

## Stack escolhida

- Backend: Go como serviço independente.
- Banco: PostgreSQL.
- Frontend: React, TypeScript e Vite.
- UI: MUI.
- Rotas: React Router.
- Estado remoto: TanStack Query.
- Contrato de API: OpenAPI 3.
- Client TypeScript: gerado a partir da especificação OpenAPI.

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

## Configuração

Copie `.env.example` para `.env` e ajuste os valores conforme necessário.

Variáveis principais:

- `ADDRESS`: endereço do backend.
- `DATABASE_URL`: conexão do PostgreSQL.
- `FRONTEND_URL`: origem permitida para o frontend.
- `SESSION_SECRET`: segredo usado para gerar o hash dos tokens de sessão.
- `SESSION_TTL_HOURS`: duração da sessão em horas.
- `COOKIE_SECURE`: use `true` quando a API estiver atrás de HTTPS.
- `GOOGLE_CLIENT_ID` e `GOOGLE_CLIENT_SECRET`: credenciais OAuth do Google.
- `GOOGLE_REDIRECT_URL`: callback autorizado no Google Cloud. Em local, use `http://localhost:8080/api/auth/google/callback`.

## Banco de dados

Configure `DATABASE_URL` apontando para uma instância PostgreSQL. Você pode usar um PostgreSQL já instalado na máquina ou subir o serviço auxiliar via Docker Compose.

Opção com Docker Compose:

```bash
make db-up
```

Opção com PostgreSQL local:

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

Para desfazer a última migration aplicada:

```bash
make migrate-down
```

## Backend

O backend roda como um serviço independente em Go:

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

Comandos do binário:

```bash
cd backend
go run ./cmd/server run
go run ./cmd/server migrate up
go run ./cmd/server migrate down
go run ./cmd/server openapi
```

## Autenticação por e-mail e senha

A base atual possui cadastro, login, logout e consulta do usuário autenticado. A sessão usa cookie HTTP-only (`form_builder_session`) e persiste apenas o hash HMAC do token no banco.

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

Usuário autenticado:

```bash
curl -i http://localhost:8080/api/auth/me \
  --cookie "form_builder_session=<cookie-value>"
```

## Autenticação com Google

O login com Google usa OAuth 2.0 Authorization Code no backend Go. O frontend apenas redireciona para `GET /api/auth/google`; o callback valida `state`, consulta o perfil no Google e cria a mesma sessão HTTP-only usada pelo login por e-mail/senha.

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
- `GET /api/auth/google/callback`: recebe `code` e `state`, cria/associa o usuário e redireciona para `/admin`.

## CRUD de formulários

O backend já expõe CRUD autenticado para formulários. Campos são persistidos junto com o formulário e substituídos em bloco no `PUT`, mantendo a camada de regra no backend. Cada formulário também possui metadados de privacidade usados no aviso público: e-mail do controlador, finalidade do tratamento e política de retenção.

Rotas autenticadas:

- `GET /api/forms`: lista formulários do usuário autenticado.
- `POST /api/forms`: cria formulário em rascunho.
- `GET /api/forms/{formId}`: consulta formulário do usuário autenticado.
- `PUT /api/forms/{formId}`: atualiza título, descrição e campos.
- `DELETE /api/forms/{formId}`: remove formulário.
- `POST /api/forms/{formId}/publish`: publica e gera `publicSlug`.
- `POST /api/forms/{formId}/unpublish`: volta o formulário para rascunho.

Rota pública:

- `GET /api/public/forms/{slug}`: consulta formulário publicado para preenchimento.
- `POST /api/public/forms/{slug}/responses`: envia uma resposta pública para formulário publicado.

Respostas:

- `GET /api/forms/{formId}/responses`: lista respostas recebidas por formulário do usuário autenticado.
- `GET /api/forms/{formId}/responses/export`: exporta respostas em JSON para atendimento administrativo.
- `DELETE /api/forms/{formId}/responses/{responseId}`: exclui uma resposta específica do formulário.

O frontend também permite exportar a listagem de respostas em PDF e Excel a partir do painel administrativo.

Ao editar um formulário, os IDs dos campos existentes são preservados para manter as respostas já enviadas associadas às colunas corretas. Caso existam respostas antigas gravadas com campos que já foram removidos, a tela e as exportações exibem esses valores em "Outros dados salvos".

Tipos de campo suportados:

- Texto curto, texto longo, e-mail, número, telefone, seleção e caixa de seleção.
- Campos de e-mail usam validação de formato.
- Campos de número exigem valor numérico válido.
- Campos de telefone exigem 12 dígitos numéricos no padrão brasileiro usado neste projeto: 3 dígitos de DDD e 9 dígitos de telefone.

Criar formulário:

```bash
curl -i http://localhost:8080/api/forms \
  --cookie "form_builder_session=<cookie-value>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Pesquisa de satisfação",
    "description": "Formulário público para clientes",
    "controllerEmail": "privacidade@example.com",
    "privacyPurpose": "Coletar feedback para melhoria do atendimento.",
    "retentionPolicy": "As respostas serão mantidas por até 90 dias após o envio.",
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

Publicar formulário:

```bash
curl -i -X POST http://localhost:8080/api/forms/<form-id>/publish \
  --cookie "form_builder_session=<cookie-value>"
```

Consultar formulário publicado:

```bash
curl -i http://localhost:8080/api/public/forms/<public-slug>
```

Enviar resposta pública:

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

## Pacote mínimo LGPD

O projeto inclui controles básicos para apoiar transparência e atendimento de direitos do titular, sem substituir revisão jurídica do controlador.

- Página pública `/privacidade` com política de privacidade do sistema.
- Formulário publicado exige e exibe e-mail do controlador, finalidade do tratamento e política de retenção.
- Publicação bloqueada quando os metadados de privacidade obrigatórios não foram preenchidos.
- Envio público exige `privacyAcknowledged: true` e grava `privacyAcknowledgedAt` na resposta.
- Painel administrativo permite exportar respostas em JSON, PDF e Excel, além de excluir uma resposta específica.
- Cookies usados hoje são essenciais: sessão HTTP-only e estado temporário do Google OAuth.

## Frontend

Instale as dependências e rode o frontend:

```bash
cd frontend
npm install
npm run dev
```

O frontend usa `http://localhost:8080` como API por padrão. Para mudar, copie `frontend/.env.example` para `frontend/.env` e ajuste `VITE_API_URL`.

Guarda de arquitetura do frontend:

```bash
cd frontend
npm run lint:architecture
```

Regra do projeto: componentes e páginas não devem misturar lógica com estética. Estilos MUI ficam em arquivos `*.styles.ts`; rotas e providers ficam em `src/app`; código por domínio fica em `src/features`; helpers sem UI ficam em `src/lib`.

Depois do login, `/admin` exibe apenas a lista de formulários salvos. A área administrativa completa fica em `/admin/workspace` e consome a API real para criar, editar, configurar campos, publicar, despublicar e excluir formulários. A página pública `/f/:slug` carrega o formulário publicado por slug e envia respostas. A rota `/admin/forms/:formId/responses` lista as respostas recebidas.

## Validação manual local

Com PostgreSQL, backend e frontend rodando, valide o fluxo principal no navegador:

```text
http://localhost:5173
```

Checklist:

- Criar uma conta por e-mail e senha.
- Entrar e conferir a lista de formulários salvos em `/admin`.
- Acessar a área administrativa em `/admin/workspace`.
- Criar um formulário.
- Preencher e-mail do controlador, finalidade do tratamento e retenção das respostas.
- Adicionar campos ao formulário, incluindo telefone para validar o padrão brasileiro.
- Salvar o formulário.
- Publicar o formulário.
- Abrir o link público gerado em `/f/:slug`.
- Confirmar ciência do aviso de privacidade.
- Enviar uma resposta sem autenticação.
- Voltar ao admin e abrir `/admin/forms/:formId/responses`.
- Confirmar que a resposta enviada aparece na listagem.
- Exportar respostas em JSON, PDF e Excel e excluir uma resposta de teste.

Sinais esperados:

- Backend respondendo em `http://localhost:8080/healthz`.
- Readiness respondendo em `http://localhost:8080/readyz`.
- Frontend servindo em `http://localhost:5173`.
- Cookie HTTP-only `form_builder_session` criado após login.

Se houver um PostgreSQL local usando a porta `5432`, pare esse serviço antes de subir o PostgreSQL do Docker ou altere a porta publicada no `docker-compose.yml` e ajuste `DATABASE_URL`.

## OpenAPI e client TypeScript

Sempre que uma rota, payload, schema ou status code da API mudar no backend, atualize primeiro a especificação OpenAPI gerada:

```bash
make openapi
```

Esse comando executa o backend Go e regrava `backend/openapi/openapi.json`. Depois gere novamente os tipos TypeScript consumidos pelo frontend:

```bash
make client
```

Esse comando roda `npm run generate:api` dentro de `frontend` e atualiza `frontend/src/api/generated/schema.ts` e `frontend/src/api/generated/client.ts` a partir de `backend/openapi/openapi.json`.

Fluxo obrigatório para mudanças de contrato:

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

Arquivos gerados ficam em `frontend/src/api/generated`, são versionados para revisão e não devem ser editados manualmente. O build do frontend também executa `npm run generate:api` antes do typecheck. Os tipos usados nas features de auth, formulários e respostas são aliases derivados de `frontend/src/api/generated/schema.ts`, e as páginas/features consomem a API por meio das funções geradas em `frontend/src/api/generated/client.ts`.

## Decisões iniciais

- Monorepo para simplificar setup local, revisão e DX.
- PostgreSQL por ser multiplataforma, robusto e adequado para persistência relacional. O Docker Compose existe apenas como conveniência para desenvolvimento local.
- Vite em vez de Next.js para deixar claro que o backend Go é responsável pela API, regras de negócio e persistência.
- Frontend organizado por `app`, `features`, `styles`, `lib` e `api`, mantendo estilos fora dos componentes e lógica fora da camada visual.
- Formulários, campos e respostas serão modelados em tabelas relacionais, com `jsonb` para configurações e respostas flexíveis.
- Sessões em banco com token opaco em cookie HTTP-only. O banco armazena apenas o hash HMAC do token, reduzindo impacto se os dados de sessão vazarem.
- Migrations embutidas no binário Go para manter o setup local reproduzível sem depender de uma CLI externa.

## Status

Já existe a base real do backend com conexão PostgreSQL, migrations executáveis, modelo de usuários/sessões, autenticação por e-mail/senha, autenticação com Google, CRUD autenticado de formulários com publicação, envio/listagem de respostas e frontend conectado aos fluxos principais.
