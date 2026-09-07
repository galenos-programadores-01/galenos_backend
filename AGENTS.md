# AGENTS.md

Galenos Pro Appointments API — REST API en Go (Gin + SQL Server 2022), arquitectura hexagonal (Ports and Adapters). Las reglas de arquitectura estrictas están en `ARCHITECTURE.md`; el catálogo de endpoints y variables está en `README.md`.

## Comandos

```bash
go mod tidy
go run ./cmd/api        # arranca la API
go build ./...          # verificar compilación
./fix.bat               # Git Bash / Linux: formatea el código (go fmt) y revisa errores estáticos (go vet)
.\fix.bat               # PowerShell / Windows CMD: formatea código y revisa errores estáticos
```

- **No hay** tests, linter ni CI. Verifica con `go build ./...`. Formatea con `gofmt -w .`.
- `go.mod` declara `go 1.27rc2` (toolchain release candidate); si `go build` falla por versión, revisa hay que tener instalada esa toolchain, no bajar el go.mod.
- **Swagger**: `docs/docs.go` es **generado** por swaggo ("DO NOT EDIT"). Al cambiar anotaciones `@Summary`/`@Router`/etc. de los handlers, regenéralo con `swag init` (los handlers ya usan `@Router` en comentarios). UI en `/swagger/index.html`.
- El código y los comentarios **están en español** (identificadores en inglés). Mantén ese estilo.
- No hacer push ni commit salvo que se pida. Se trabaja en branches de feature (`feature/gustavo`, `feature/carlos`); `main` puede tener cambios staged sin commitear.

## Arquitectura Hexagonal — dirección de dependencias (ver ARCHITECTURE.md)

La regla central (verificada): `internal/domain` no importa nada de `adapters` ni frameworks.

- `domain` → entidades y errores sentinel (solo stdlib).
- `ports/{input,output}` → contratos. `ports/shared/pagination.go` → `PageRequest`/`PageResponse[T]`.
- `usecase` → implementa los input ports; recibe repos/clients por interfaz, nunca implementaciones concretas.
- `adapters/input/http` → Gin handlers, router, DTOs, `response.go` (wrapper estándar).
- `adapters/output/persistence/sqlserver` → `database/sql` + `microsoft/go-mssqldb`.
- `adapters/output/{reniec,sis}` → clients SOAP de servicios externos.
- `cmd/api/main.go` → **único** composition root; solo aquí se instancian implementaciones concretas.

Nota: `ARCHITECTURE.md` pide que los structs de `domain` no tengan tags de framework, pero el código real sí usa `json:` en varios. Sigue el código existente.

## Configuración / entorno

Requeridas (si faltan, `config.Load` corta la app): `SQLSERVER_DSN`, `JWT_SECRET`, `API_USERNAME`, `API_PASSWORD`.

- `main.go` carga `.env` solo si existe (copia `.env.example → .env`). Los valores reales nunca se suben a git.
- Autenticación: `POST /api/v1/auth/login` devuelve un JWT HS256 (secret `JWT_SECRET`, TTL por `JWT_TTL`, default 24h). Exigen Bearer (`RequireBearerToken`): `/pacientes`, `/appointments`, `/evoluciones`, `/ordenes`, `/resultados`, `/interconsultas`, `/sintomas` y `/auth/menus`; el resto (catálogos, `/reniec`, `/sis`, `/triaje`…) es público. Ojo: el README menciona `JWT_EXPIRATION`, pero el código lee **`JWT_TTL`**.
- Otras vars: `SERVER_PORT` (8080), `SERVER_HOST` (IP autodetectado de la máquina si vacío; la API escucha en todas las interfaces), `ALLOWED_ORIGINS` (por defecto añade el propio host+localhost para Swagger/CORS), `DB_MAX_OPEN_CONNS` (25), `DB_MAX_IDLE_CONNS` (10), `DB_CONN_MAX_LIFETIME` (5m).
- Servicios externos SOAP (con timeout, pueden no estar alcanzables en dev): RENIEC (`RENIEC_APP`, `RENIEC_USUARIO`, `RENIEC_CLAVE`, `RENIEC_URL`, `RENIEC_TIMEOUT`) y SIS (`SIS_USUARIO`, `SIS_CLAVE`, `SIS_URL`, `SIS_DNI_AUTORIZADO`, `SIS_TIMEOUT`).
- Al arrancar `main.go` abre el pool y hace **PING con timeout de 5s**: sin red/VPN al host SQL Server o credenciales mal, la API falla al inicio (quick fail).

## SQL / persistencia

- Los repositorios llaman stored procedures (`sp_listarPaciente`, `usp_go_*`, etc.) y algunos consultan tablas directo (ej. `pacientes` en `patient_repository.go`). Si el esquema cambia, ajusta query y columnas que espera el SP.
- `appointment_repository.go` usa transacción `SERIALIZABLE` con `WITH (UPDLOCK, HOLDLOCK)` para evitar dobles reservas concurrentes. No cambies ese patrón sin revisarlo.

## Router / HTTP

- La lista fiable de rutas es `router.go` (README quedó desactualizado: le faltan `/triaje`, `/sis`, `/evoluciones`, `/ordenes`, `/resultados`, `/interconsultas`, `/sintomas`, `/auth/menus`, etc.).
- Base path: `http://localhost:8080/api/v1`. Respuestas con el wrapper estándar `{ "success": true, "data": <...> }` / `{ "success": false, "error": { "code": "...", "message": "..." } }` (ver `response.go`).
- El middleware `RequireBearerToken` pone en el contexto de Gin el `idEmpleado` del claim del JWT (`c.Set("idEmpleado", ...)`). Los handlers protegidos lo leen con `c.GetInt("idEmpleado")` para registrar quién ejecuta la acción; si es 0, responden error.
- Nuevos handlers: inspírate en los existentes de `internal/adapters/input/http/` (patrón DTO → servicio por port → DTO de respuesta, anotaciones swag en comentarios, y registrar la ruta en `router.go` + cablear en `main.go`).
