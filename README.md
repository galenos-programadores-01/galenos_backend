# Galenos Pro — Appointments API

REST API en Go para el sistema hospitalario Galenos Pro, construida con **Arquitectura Hexagonal (Ports and Adapters)**. Expone endpoints para el frontend Angular y persiste en **SQL Server 2022**.

## Arquitectura

```
cmd/api/main.go                                    → composition root (inyección de dependencias)
internal/
├── domain/                                         → entidades y reglas de negocio (sin dependencias externas)
│   ├── appointment.go                              → Appointment, TimeSlot, reglas ("no en el pasado", etc.)
│   ├── patient.go                                  → Patient (modelo de lectura)
│   └── errors.go                                   → errores de dominio (sentinels)
├── ports/
│   ├── input/                                      → driving ports (casos de uso)
│   ├── output/                                     → driven ports (repositorios)
│   └── shared/pagination.go                        → PageRequest / PageResponse[T]
├── usecase/                                         → implementa los input ports orquestando dominio + output ports
├── adapters/
│   ├── input/http/                                  → driving adapter (Gin): handlers, router, DTOs, respuesta estándar
│   └── output/persistence/sqlserver/                 → driven adapter: database/sql + go-mssqldb
└── config/config.go                                 → lectura de variables de entorno
```

Reglas de dependencia: `domain` no importa nada de `adapters` ni de frameworks; `ports` define contratos; `usecase` conecta ambos lados; `adapters` implementa los contratos hacia HTTP y SQL Server. `main.go` es el único lugar donde se instancian las implementaciones concretas.

## Requisitos

- Go 1.25+
- SQL Server 2022 accesible desde la máquina donde corra la API
- Acceso de red/VPN al host de base de datos (interno)

## Configuración

Copia `.env.example` a `.env` y ajusta los valores (el `.env` real con credenciales ya está en `.gitignore`, nunca se sube a git):

| Variable | Descripción | Default |
|---|---|---|
| `SQLSERVER_DSN` | Cadena de conexión, formato `sqlserver://user:pass@host:port?database=db` | *(requerida)* |
| `SERVER_PORT` | Puerto HTTP de la API | `8080` |
| `ALLOWED_ORIGINS` | Orígenes permitidos por CORS (coma-separados), ej. la URL de Angular | `http://localhost:4200` |
| `DB_MAX_OPEN_CONNS` | Máximo de conexiones abiertas del pool | `25` |
| `DB_MAX_IDLE_CONNS` | Máximo de conexiones inactivas del pool | `10` |
| `DB_CONN_MAX_LIFETIME` | Tiempo de vida máximo de una conexión | `5m` |
| `JWT_SECRET` | Secreto usado para firmar y validar los JWT (HS256) | *(requerida)* |
| `JWT_EXPIRATION` | Tiempo de expiración del token emitido en el login | `24h` |
| `API_USERNAME` | Usuario válido para `POST /auth/login` | *(requerida)* |
| `API_PASSWORD` | Contraseña válida para `POST /auth/login` | *(requerida)* |

`main.go` carga `.env` automáticamente si existe (vía `godotenv`), útil solo para desarrollo local.

## Cómo correr

```bash
go mod tidy
go run ./cmd/api
```

Al arrancar, la API hace `PING` a SQL Server; si el host no es alcanzable (red/VPN) o la contraseña es incorrecta, falla rápido con el error correspondiente en el log.

### Formato y Linter (Validación de código)

Para revisar y formatear el código rápidamente antes de hacer un commit, ejecuta el script incluido en la raíz del proyecto según tu terminal:

```bash
# En Git Bash / Linux:
./fix.bat

# En PowerShell / Windows CMD:
.\fix.bat
```

Este comando ejecuta internamente `go fmt ./...` (para formatear el código) y `go vet ./...` (el linter oficial estricto de Go).
## Endpoints

Base: `http://localhost:8080/api/v1`

| Método | Ruta | Descripción |
|---|---|---|
| `POST` | `/appointments` | Agenda una nueva cita médica |
| `GET` | `/appointments/:id` | Obtiene una cita por id |
| `GET` | `/pacientes?page=1&pageSize=20` | Lista pacientes paginados (`NroDocumento`, `ApellidoPaterno`, `ApellidoMaterno`, `PrimerNombre`, `SegundoNombre`, `TercerNombre`) |
| `GET` | `/pacientes/buscar?documento=&hc=&paterno=&materno=&nombres=` | Busca pacientes por filtros (SP `usp_go_listarpacientes`) |
| `GET` | `/pacientes/:idOrDoc` | Detalle por id (SP `webPacientesListarIdPaciente`) si es numérico, o por documento (SP `sp_listarPaciente`) |
| `PUT` | `/pacientes/:id` | Modifica un paciente (SP `usp_go_ModificarPaciente`) |
| `GET` | `/etnias` | Catálogo de etnias (SP `ups_go_ListarEtnias`) |
| `GET` | `/idiomas` | Catálogo de idiomas (SP `ups_go_ListarIdiomas`) |
| `GET` | `/tipos-sexo` | Catálogo de sexos (SP `usp_go_ListarTiposSexos`) |
| `GET` | `/estados-civil` | Catálogo de estados civiles (SP `usp_go_ListarEstadosCivil`) |
| `GET` | `/grados-instruccion` | Catálogo de grados de instrucción (SP `usp_go_ListarGradoInstruccion`) |
| `GET` | `/ocupaciones` | Catálogo de ocupaciones (SP `usp_go_ListarOcupaciones`) |
| `GET` | `/tipos-documentos` | Catálogo de tipos de documento (SP `usp_go_ListarTiposDocumentos`) |
| `GET` | `/departamentos` | Catálogo de departamentos (SP `usp_go_ListarDepartamentos`) |
| `GET` | `/provincias/:id` | Provincias de un departamento (SP `usp_go_ListarProvincias`) |
| `GET` | `/distritos/:id` | Distritos de una provincia (SP `usp_go_ListarDistritos`) |
| `GET` | `/centros-poblados/:id` | Centros poblados de un distrito (SP `usp_go_ListarCentrosPoblados`) |
| `GET` | `/paises` | Catálogo de países (SP `usp_go_ListarPaises`) |
| `GET` | `/reniec/:nrodoc?operacion=completo` | Consulta RENIEC (SOAP externo, `basico` \| `completo`) |
| `GET` | `/health` | Estado de la API |

Todas las respuestas usan el sobre estándar:

```json
{ "success": true, "data": { } }
{ "success": false, "error": { "code": "VALIDATION_ERROR", "message": "..." } }
```

## Colección Postman

`postman/galenos-pro-api.postman_collection.json` — importar en Postman. Variable `baseUrl` apunta a `http://localhost:8080` por defecto. Ejecutar primero **Autenticación > Login** (usa las variables `username`/`password`): su test script guarda el JWT en la variable de colección `accessToken`, que el resto de peticiones reutiliza automáticamente vía Bearer Auth a nivel de colección.

## Notas

- `patient_repository.go` consulta directamente la tabla `pacientes`; si el esquema real cambia, ajustar la query en `internal/adapters/output/persistence/sqlserver/patient_repository.go`.
- El módulo de citas verifica disponibilidad del médico dentro de una transacción `SERIALIZABLE` con `UPDLOCK, HOLDLOCK` antes de insertar, para evitar dobles reservas concurrentes.
