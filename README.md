# Form Builder

Aplicação full stack para criação, publicação e resposta de formulários. O projeto foi feito para um teste prático de vaga Full Stack Developer e os requisitos do PDF são tratados como regra do projeto.

## O que o sistema faz

- Permite criar conta e entrar como administrador por e-mail/senha.
- Cria um usuário administrador padrão para facilitar teste local, quando o seed estiver ativado.
- Permite entrar com conta Google usando OAuth 2.0.
- Exibe, após login, uma tela de formulários salvos.
- Mantém a área administrativa protegida por autenticação.
- Permite criar, editar e excluir formulários.
- Permite adicionar e configurar campos.
- Permite publicar um formulário e gerar uma URL pública.
- Permite que qualquer pessoa com a URL pública responda sem login.
- Armazena respostas no banco de dados.
- Permite ao administrador consultar, exportar e excluir respostas.
- Inclui política de privacidade, aviso público e registros mínimos de ciência LGPD.

## Tecnologias

- Backend: Go 1.25+ como serviço independente.
- Banco de dados: PostgreSQL.
- Frontend: React, TypeScript e Vite.
- UI: MUI.
- Rotas: React Router.
- Requisições/cache: TanStack Query.
- Contrato da API: OpenAPI 3.
- Client TypeScript: gerado a partir do OpenAPI.
- Docker: usado apenas como conveniência para subir o PostgreSQL localmente.

## Estrutura do projeto

```text
.
├── backend/
│   ├── cmd/server/              # Entrada do backend Go
│   ├── internal/                # Regras, handlers, repositórios e serviços
│   ├── migrations/              # Migrations SQL executadas pelo backend
│   └── openapi/                 # Especificação OpenAPI gerada
├── frontend/
│   └── src/
│       ├── api/                 # Client e tipos gerados pelo OpenAPI
│       ├── app/                 # Rotas e providers
│       ├── features/            # Funcionalidades por domínio
│       ├── lib/                 # Helpers sem UI
│       └── styles/              # Tema e estilos compartilhados
├── scripts/
├── docker-compose.yml
├── Makefile
├── .env.example
└── README.md
```

## Pré-requisitos

Antes de começar, instale:

- Git.
- Go 1.25 ou superior.
- Node.js 20 ou superior.
- PostgreSQL 16+ ou Docker Desktop.

Você pode rodar o sistema com Docker ou sem Docker.

Com Docker, o Docker sobe somente o PostgreSQL. O backend e o frontend continuam rodando direto na sua máquina.

Sem Docker, você precisa ter o PostgreSQL instalado e criar o banco manualmente.

## Passo a passo rápido com Docker

Este é o caminho mais simples para rodar localmente.

### 1. Clone o repositório

```powershell
git clone https://github.com/ElvisReis2K/Form-Builder.git
cd Form-Builder
```

Se você já tem a pasta do projeto, entre nela:

```powershell
cd "C:\Users\ElviZ\Documents\ChatGPT\Teste Prático Falqon"
```

### 2. Crie o arquivo `.env`

Na raiz do projeto, copie o exemplo:

```powershell
Copy-Item .env.example .env
```

Abra o arquivo:

```powershell
notepad .env
```

Para usar o PostgreSQL do Docker, deixe o `DATABASE_URL` assim:

```env
DATABASE_URL=postgres://form_builder:form_builder@localhost:5432/form_builder?sslmode=disable
```

O arquivo `.env` deve ficar na raiz do projeto, no mesmo nível de `README.md`, `backend`, `frontend` e `docker-compose.yml`.

Nunca envie o `.env` para o GitHub. Ele fica ignorado pelo `.gitignore`.

Com o `.env.example`, o projeto já vem preparado para criar um usuário padrão de desenvolvimento:

```text
e-mail: admin@gmail.com
senha: 12345678
```

Esse usuário é criado quando você roda as migrations com `go run ./cmd/server migrate up`.

### 3. Suba o PostgreSQL

Abra o Docker Desktop e espere ele iniciar. Depois rode:

```powershell
docker compose up -d postgres
```

Esse comando cria um banco PostgreSQL em:

```text
host: localhost
porta: 5432
banco: form_builder
usuário: form_builder
senha: form_builder
```

Se você já tiver PostgreSQL instalado usando a porta `5432`, o Docker pode dar conflito ou o backend pode conectar no banco errado. Nesse caso, pare o serviço local do PostgreSQL ou rode o projeto sem Docker usando os dados corretos do seu banco.

### 4. Baixe dependências do backend

```powershell
cd backend
go mod tidy
```

### 5. Rode as migrations

Ainda dentro de `backend`, rode:

```powershell
go run ./cmd/server migrate up
```

Mensagem esperada:

```text
migrate up completed
```

Se `SEED_DEFAULT_ADMIN=true` estiver no `.env`, o backend também confere/cria o usuário padrão:

```text
e-mail: admin@gmail.com
senha: 12345678
```

### 6. Inicie o backend

Ainda dentro de `backend`, rode:

```powershell
go run ./cmd/server run
```

Deixe esse terminal aberto.

Mensagem esperada:

```text
server listening on http://localhost:8080
```

Para testar no navegador, abra:

```text
http://localhost:8080/healthz
```

Também existe o teste de conexão com o banco:

```text
http://localhost:8080/readyz
```

### 7. Instale dependências do frontend

Abra um segundo terminal na raiz do projeto e rode:

```powershell
cd frontend
npm install
```

### 8. Inicie o frontend

Ainda dentro de `frontend`, rode:

```powershell
npm run dev
```

Deixe esse terminal aberto.

Mensagem esperada:

```text
Local: http://localhost:5173/
```

Abra no navegador:

```text
http://localhost:5173
```

## Passo a passo sem Docker

Use esta opção se você prefere usar um PostgreSQL instalado diretamente na máquina.

### 1. Instale e inicie o PostgreSQL

Instale PostgreSQL 16 ou superior. Depois confirme que o serviço está rodando.

No Windows, você pode verificar no PowerShell:

```powershell
Get-Service *postgres*
```

### 2. Crie usuário e banco

Entre no `psql` com um usuário administrador, por exemplo `postgres`:

```powershell
psql -U postgres
```

Crie usuário e banco:

```sql
CREATE USER form_builder WITH PASSWORD 'form_builder';
CREATE DATABASE form_builder OWNER form_builder;
```

Saia do `psql`:

```sql
\q
```

### 3. Configure o `.env`

Na raiz do projeto:

```powershell
Copy-Item .env.example .env
notepad .env
```

Se você criou o banco igual ao exemplo acima, use:

```env
DATABASE_URL=postgres://form_builder:form_builder@localhost:5432/form_builder?sslmode=disable
```

Se você usa outro usuário, senha, porta ou nome de banco, ajuste o valor:

```env
DATABASE_URL=postgres://seu_usuario:sua_senha@localhost:5432/seu_banco?sslmode=disable
```

Depois siga os mesmos passos do modo com Docker:

```powershell
cd backend
go mod tidy
go run ./cmd/server migrate up
go run ./cmd/server run
```

Em outro terminal:

```powershell
cd frontend
npm install
npm run dev
```

## Configuração do `.env`

O backend lê variáveis de ambiente e também carrega o arquivo `.env` da raiz do projeto.

Exemplo completo para desenvolvimento local:

```env
ADDRESS=localhost:8080
DATABASE_URL=postgres://form_builder:form_builder@localhost:5432/form_builder?sslmode=disable
FRONTEND_URL=http://localhost:5173
SESSION_SECRET=dev-session-secret-change-me-before-production
SESSION_TTL_HOURS=168
COOKIE_SECURE=false
SEED_DEFAULT_ADMIN=true
DEFAULT_ADMIN_NAME=Administrador
DEFAULT_ADMIN_EMAIL=admin@gmail.com
DEFAULT_ADMIN_PASSWORD=12345678
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

Significado de cada variável:

- `ADDRESS`: endereço onde o backend Go vai rodar.
- `DATABASE_URL`: conexão do PostgreSQL.
- `FRONTEND_URL`: endereço do frontend autorizado pelo backend.
- `SESSION_SECRET`: segredo usado para proteger tokens de sessão.
- `SESSION_TTL_HOURS`: tempo de validade da sessão.
- `COOKIE_SECURE`: use `false` em local com HTTP e `true` em produção com HTTPS.
- `SEED_DEFAULT_ADMIN`: quando `true`, cria o usuário padrão ao rodar `migrate up`.
- `DEFAULT_ADMIN_NAME`: nome do usuário padrão.
- `DEFAULT_ADMIN_EMAIL`: e-mail do usuário padrão.
- `DEFAULT_ADMIN_PASSWORD`: senha do usuário padrão.
- `GOOGLE_CLIENT_ID`: Client ID do OAuth Google.
- `GOOGLE_CLIENT_SECRET`: Client Secret do mesmo OAuth Client.
- `GOOGLE_REDIRECT_URL`: URL de callback cadastrada no Google Cloud.

Não coloque aspas ao redor dos valores e não coloque espaços antes ou depois do `=`.

Sempre que alterar o `.env`, pare e inicie o backend novamente.

Em produção ou em qualquer ambiente público, altere a senha padrão ou defina:

```env
SEED_DEFAULT_ADMIN=false
```

## Login com Google localmente

O código já tem login com Google, mas as credenciais do Google não podem ir para o GitHub. Por isso, cada pessoa que clonar o projeto precisa configurar o próprio OAuth Client no Google Cloud.

Isso é uma regra de segurança: se o `GOOGLE_CLIENT_SECRET` fosse versionado, qualquer pessoa poderia usar a credencial privada do projeto.

### Criar credenciais no Google Cloud

1. Acesse o [Google Cloud Console](https://console.cloud.google.com/).
2. Crie ou selecione um projeto.
3. Abra **APIs & Services** ou **Google Auth Platform**.
4. Configure a tela de consentimento OAuth.
5. Informe nome do app e e-mail de suporte.
6. Para teste local, você pode deixar o app em modo de teste e adicionar seu e-mail em **Test users**.
7. Se quiser que qualquer conta Google consiga entrar sem estar em **Test users**, publique a tela de consentimento em produção.
8. Vá em **Clients** ou **Credentials**.
9. Crie um OAuth Client do tipo **Web application**.
10. Em **Authorized JavaScript origins**, adicione exatamente:

```text
http://localhost:5173
```

11. Em **Authorized redirect URIs**, adicione exatamente:

```text
http://localhost:8080/api/auth/google/callback
```

12. Salve e copie o **Client ID** e o **Client secret**.

### Preencher no `.env`

Abra o `.env` da raiz:

```powershell
notepad .env
```

Preencha:

```env
GOOGLE_CLIENT_ID=cole_aqui_o_client_id_completo
GOOGLE_CLIENT_SECRET=cole_aqui_o_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

O `GOOGLE_CLIENT_ID` precisa ser o valor completo que termina com:

```text
.apps.googleusercontent.com
```

Exemplo do formato esperado:

```env
GOOGLE_CLIENT_ID=1234567890-abcdefg.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-exemplo_de_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

Não use API Key, Project ID, Service Account ou Client Secret no campo `GOOGLE_CLIENT_ID`.

Depois de salvar o `.env`, reinicie o backend:

```powershell
cd backend
go run ./cmd/server run
```

Abra o frontend e clique em **Continuar com Google**:

```text
http://localhost:5173
```

### Erros comuns no Google

`Erro 401: invalid_client`:

- O `GOOGLE_CLIENT_ID` está errado, incompleto ou não termina em `.apps.googleusercontent.com`.
- O `GOOGLE_CLIENT_SECRET` não pertence ao mesmo OAuth Client.
- O OAuth Client não é do tipo **Web application**.
- O backend não foi reiniciado depois da alteração no `.env`.

`redirect_uri_mismatch`:

- A URL cadastrada em **Authorized redirect URIs** precisa ser exatamente `http://localhost:8080/api/auth/google/callback`.
- Não coloque barra final extra.
- Não troque `localhost` por `127.0.0.1` se o `.env` está usando `localhost`.

Acesso bloqueado ou app não verificado:

- Confira a tela de consentimento.
- Se estiver em modo de teste, adicione seu e-mail em **Test users**.
- Alterações no Google Cloud podem levar alguns minutos para refletir.

Referências oficiais:

- [Google OAuth 2.0 for Web Server Applications](https://developers.google.com/identity/protocols/oauth2/web-server)
- [Google Cloud - Manage OAuth Clients](https://support.google.com/cloud/answer/15549257)

## Como usar o sistema no navegador

Com banco, backend e frontend rodando, abra:

```text
http://localhost:5173
```

### Criar conta por e-mail e senha

1. Na tela inicial, clique em **Criar conta**.
2. Informe nome, e-mail e senha.
3. Clique em **Criar**.
4. Depois do login, você será enviado para a tela de **Formulários salvos**.

### Entrar por e-mail e senha

1. Na tela inicial, informe e-mail e senha.
2. Clique em **Entrar**.
3. Após o login, a aplicação abre a tela de **Formulários salvos**.

Se você acabou de rodar as migrations com o seed ativado, pode entrar com:

```text
e-mail: admin@gmail.com
senha: 12345678
```

### Entrar com Google

1. Configure o Google OAuth no `.env`.
2. Reinicie o backend.
3. Abra `http://localhost:5173`.
4. Clique em **Continuar com Google**.
5. Escolha sua conta Google.
6. Após o retorno do Google, a aplicação abre a tela de **Formulários salvos**.

### Sessão e segurança de acesso

A área administrativa só pode ser acessada depois do login.

Por regra definida no projeto, se o usuário recarregar uma página protegida ou voltar para a tela de login, a aplicação encerra a sessão no backend e exige novo login. Isso impede que alguém avance para a área administrativa apenas pelo histórico do navegador.

Use o botão **Sair** no cabeçalho para encerrar a sessão manualmente.

## Fluxo completo do formulário

### 1. Abrir formulários salvos

Depois do login, a primeira tela protegida é:

```text
http://localhost:5173/admin
```

Essa tela mostra os formulários já criados.

### 2. Acessar a área administrativa

No menu, clique em **Administração**.

Rota:

```text
http://localhost:5173/admin/workspace
```

Essa área permite criar, editar, configurar campos, salvar, publicar e despublicar formulários.

### 3. Criar ou editar formulário

Na área administrativa:

1. Informe o título do formulário.
2. Informe uma descrição, se quiser.
3. Preencha o e-mail do controlador dos dados.
4. Preencha a finalidade do tratamento dos dados.
5. Preencha a política de retenção das respostas.
6. Adicione os campos necessários.
7. Configure tipo, nome, obrigatoriedade, placeholder e opções.
8. Clique em **Salvar**.

Os dados de privacidade são obrigatórios para publicar o formulário.

### 4. Tipos de campo

Tipos disponíveis:

- Texto curto.
- Texto longo.
- E-mail.
- Número.
- Telefone.
- Seleção.
- Caixa de seleção.

Regras de validação:

- Campo obrigatório precisa ser preenchido.
- E-mail precisa ter formato válido.
- Número precisa ser numérico.
- Telefone precisa ter 12 dígitos numéricos no padrão usado neste projeto: 3 dígitos de DDD e 9 dígitos de telefone.
- Seleção precisa ter opções cadastradas.

Para telefone, informe somente números. Exemplo:

```text
011912345678
```

### 5. Publicar formulário

Depois de salvar e preencher os dados LGPD mínimos, clique em **Publicar**.

A aplicação gera uma URL pública no formato:

```text
http://localhost:5173/f/slug-do-formulario
```

Qualquer pessoa com essa URL pode preencher o formulário sem login.

### 6. Responder formulário público

Abra a URL pública gerada.

Na página pública:

1. Leia o aviso de privacidade.
2. Preencha os campos.
3. Confirme a ciência do aviso LGPD.
4. Envie a resposta.

O backend valida os dados usando a definição do formulário publicado.

### 7. Consultar respostas

Entre como administrador, vá para **Formulários salvos** e abra as respostas do formulário.

Rota:

```text
http://localhost:5173/admin/forms/:formId/responses
```

Na tela de respostas, o administrador pode:

- Ver data de envio.
- Ver ciência LGPD.
- Ver respostas por campo.
- Exportar JSON.
- Exportar PDF.
- Exportar Excel.
- Excluir uma resposta específica.

Se um campo for removido depois de já ter recebido respostas, os dados antigos aparecem em **Outros dados salvos** para não perder histórico.

## Rotas principais do frontend

- `/`: login e cadastro.
- `/admin`: formulários salvos, somente após login.
- `/admin/workspace`: área administrativa, somente após login.
- `/admin/forms/:formId/responses`: respostas de um formulário, somente após login.
- `/f/:slug`: formulário público publicado.
- `/privacidade`: política de privacidade.

## Rotas principais da API

Autenticação:

- `POST /api/auth/register`: cria conta.
- `POST /api/auth/login`: entra por e-mail e senha.
- `POST /api/auth/logout`: sai da sessão atual.
- `GET /api/auth/me`: consulta usuário autenticado.
- `GET /api/auth/google`: inicia login Google.
- `GET /api/auth/google/callback`: recebe callback do Google.

Formulários autenticados:

- `GET /api/forms`: lista formulários do administrador.
- `POST /api/forms`: cria formulário.
- `GET /api/forms/{formId}`: consulta formulário.
- `PUT /api/forms/{formId}`: atualiza formulário.
- `DELETE /api/forms/{formId}`: remove formulário.
- `POST /api/forms/{formId}/publish`: publica formulário.
- `POST /api/forms/{formId}/unpublish`: despublica formulário.

Formulários públicos:

- `GET /api/public/forms/{slug}`: carrega formulário publicado.
- `POST /api/public/forms/{slug}/responses`: envia resposta pública.

Respostas autenticadas:

- `GET /api/forms/{formId}/responses`: lista respostas.
- `GET /api/forms/{formId}/responses/export`: exporta respostas em JSON.
- `DELETE /api/forms/{formId}/responses/{responseId}`: exclui resposta.

## LGPD mínima implementada

O projeto possui um pacote mínimo de apoio à LGPD. Ele não substitui revisão jurídica, mas atende ao fluxo básico pedido no teste.

- Página `/privacidade` com política de privacidade.
- Aviso de privacidade no formulário público.
- Campo de e-mail do controlador por formulário.
- Campo de finalidade do tratamento por formulário.
- Campo de retenção das respostas por formulário.
- Bloqueio de publicação quando os dados de privacidade obrigatórios estão vazios.
- Confirmação de ciência LGPD no envio público.
- Registro de `privacyAcknowledgedAt` na resposta.
- Exportação administrativa das respostas.
- Exclusão administrativa de respostas.
- Cookies usados apenas para sessão HTTP-only e estado temporário do Google OAuth.

## OpenAPI e client TypeScript gerado

O backend é a fonte do contrato da API. Sempre que mudar rota, payload, schema ou status code, atualize o OpenAPI e gere o client TypeScript.

Com GNU Make:

```powershell
make openapi
make client
```

Ou em um único fluxo:

```powershell
make openapi && make client
```

Sem GNU Make, no PowerShell:

```powershell
cd backend
go run ./cmd/server openapi | Set-Content -Encoding utf8 .\openapi\openapi.json
cd ..\frontend
npm run generate:api
```

Arquivos gerados:

- `backend/openapi/openapi.json`
- `frontend/src/api/generated/schema.ts`
- `frontend/src/api/generated/client.ts`

Os arquivos em `frontend/src/api/generated` são versionados para revisão, mas não devem ser editados manualmente.

## Testes e validações

Backend:

```powershell
cd backend
go test ./...
```

Frontend:

```powershell
cd frontend
npm run lint:architecture
npm run build
```

O build do frontend executa `npm run generate:api` antes do typecheck.

## Regra de organização do frontend

O projeto deve separar lógica de estética.

- Páginas e componentes ficam dentro de `frontend/src/features`.
- Estilos MUI ficam em arquivos `*.styles.ts`.
- Rotas e providers ficam em `frontend/src/app`.
- Código de API fica em `frontend/src/api`.
- Helpers sem UI ficam em `frontend/src/lib`.
- Tipos de auth, forms e responses devem vir do OpenAPI gerado.
- Wrappers e chamadas HTTP devem respeitar o client gerado.

Para checar essa regra:

```powershell
cd frontend
npm run lint:architecture
```

## Comandos úteis com Make

Se você tiver GNU Make instalado, pode usar estes atalhos na raiz do projeto:

```powershell
make db-up
make migrate-up
make backend-dev
make frontend-dev
make openapi
make client
```

O que cada comando faz:

- `make db-up`: sobe o PostgreSQL via Docker Compose.
- `make db-down`: derruba o PostgreSQL do Docker Compose.
- `make migrate-up`: aplica migrations.
- `make migrate-down`: desfaz a última migration aplicada.
- `make backend-dev`: roda o backend.
- `make frontend-dev`: roda o frontend.
- `make openapi`: gera `backend/openapi/openapi.json`.
- `make client`: gera o client TypeScript do frontend.

## Como parar o sistema

Para parar backend e frontend, pressione `Ctrl+C` nos terminais onde eles estão rodando.

Para parar o PostgreSQL do Docker:

```powershell
docker compose down
```

Se quiser apagar também os dados do banco criado pelo Docker:

```powershell
docker compose down -v
```

Use `docker compose down -v` com cuidado, porque ele remove o volume do PostgreSQL e apaga os dados locais.

## Problemas comuns

`Failed to fetch` no frontend:

- Confirme que o backend está rodando em `http://localhost:8080`.
- Abra `http://localhost:8080/healthz`.
- Confira se `frontend/.env`, caso exista, usa `VITE_API_URL=http://localhost:8080`.
- Confira se `FRONTEND_URL=http://localhost:5173` no `.env` da raiz.

Erro de senha do PostgreSQL:

- Confirme se o backend está conectando no banco certo.
- Se usar Docker, o usuário e senha padrão são `form_builder`.
- Se houver PostgreSQL local na porta `5432`, ele pode estar recebendo a conexão no lugar do container Docker.
- Ajuste `DATABASE_URL` ou pare o serviço local conflitante.

`go: command not found` ou `go não é reconhecido`:

- Instale Go 1.25+.
- Feche e abra o terminal depois da instalação.
- Rode `go version` para confirmar.

Erro em `npm install`:

- Confirme Node.js 20+ com `node -v`.
- Confirme npm com `npm -v`.
- Verifique sua conexão com a internet.

Backend não carrega mudança no `.env`:

- Pare o backend com `Ctrl+C`.
- Inicie novamente com `go run ./cmd/server run`.

Google OAuth não funciona:

- Confira `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` e `GOOGLE_REDIRECT_URL`.
- Confirme que as URLs cadastradas no Google Cloud são exatamente as mesmas do README.
- Reinicie o backend após editar o `.env`.
- Aguarde alguns minutos caso tenha acabado de criar ou alterar o OAuth Client.

Usuário padrão não entra:

- Confira se `SEED_DEFAULT_ADMIN=true` está no `.env`.
- Rode novamente `cd backend` e depois `go run ./cmd/server migrate up`.
- Use exatamente `admin@gmail.com` e `12345678`, a menos que tenha alterado `DEFAULT_ADMIN_EMAIL` ou `DEFAULT_ADMIN_PASSWORD`.
- Se esse e-mail já existia no banco com outra senha, o seed não sobrescreve a senha existente.

## Checklist para saber se tudo está funcionando

- `docker compose up -d postgres` rodou sem erro, se você escolheu Docker.
- `go run ./cmd/server migrate up` terminou com sucesso.
- Backend mostra `server listening on http://localhost:8080`.
- `http://localhost:8080/healthz` abre no navegador.
- `http://localhost:8080/readyz` abre no navegador.
- Frontend mostra `Local: http://localhost:5173/`.
- `http://localhost:5173` abre a tela de login.
- Login com `admin@gmail.com` e `12345678` funciona quando o seed padrão está ativado.
- Cadastro por e-mail e senha funciona.
- Login por Google funciona depois de configurar OAuth.
- Após login, a primeira tela é **Formulários salvos**.
- A área administrativa aparece apenas para usuário logado.
- Um formulário pode ser criado, salvo e publicado.
- A URL pública `/f/:slug` abre sem login.
- Resposta pública é salva no banco.
- Respostas aparecem para o administrador.
- Exportação JSON, PDF e Excel funciona.
- Botão **Sair** encerra a sessão e volta para login.
- Recarregar página protegida exige login novamente.

## Status atual

O projeto já possui a base real exigida para o fluxo principal:

- Backend Go independente.
- PostgreSQL com migrations.
- Usuários e sessões.
- Usuário administrador padrão para desenvolvimento local.
- Login por e-mail/senha.
- Login com Google OAuth.
- CRUD autenticado de formulários.
- Campos configuráveis.
- Publicação de formulário com URL pública.
- Envio público de respostas.
- Listagem e exportação de respostas.
- Pacote mínimo LGPD.
- OpenAPI 3 e client TypeScript gerado.
