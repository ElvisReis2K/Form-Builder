# Contributing

## Commit style

Use Conventional Commits for every commit:

- `feat: add user registration`
- `fix: validate published form status`
- `chore: tidy backend go modules`
- `docs: explain local setup`
- `refactor: split form builder styles`

## Architecture rules

Keep logic, presentation, and styling separated.

### Frontend

- Route wiring belongs in `frontend/src/app`.
- Feature pages belong in `frontend/src/features/<feature>/pages`.
- Feature-specific styles belong in `frontend/src/features/<feature>/styles`.
- Shared app styles and theme belong in `frontend/src/styles`.
- Shared non-visual helpers belong in `frontend/src/lib`.
- API generated code belongs only in `frontend/src/api/generated`.
- Do not use inline `sx={{ ... }}` or `style={{ ... }}` in TSX files. Export named style objects from `*.styles.ts` files and use `sx={styles.name}`.
- Put reusable behavior in hooks or services before it grows inside a page component.

Run the architecture guard:

```bash
cd frontend
npm run lint:architecture
```

### Backend

- `cmd/server` wires commands and process lifecycle only.
- `internal/config` owns environment/config loading.
- `internal/database` owns PostgreSQL connection and migrations.
- Domain packages, such as `internal/auth`, own their model, repository, service, and HTTP handler.
- SQL migrations belong in `backend/migrations`.
