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

## Cómo compilar y ejecutar

> [!IMPORTANT]
> Si te encuentras en la raíz del repositorio, debes ingresar primero a este directorio:
> ```bash
> cd galenos_backend
> ```

### 1. Configuración de Variables de Entorno

Antes de iniciar la API, asegúrate de contar con el archivo `.env`:

- **En PowerShell:**
  ```powershell
  if (-not (Test-Path ".env")) { Copy-Item .env.example .env }
  ```
- **En Git Bash:**
  ```bash
  [ ! -f .env ] && cp .env.example .env
  ```

Edita el archivo `.env` y verifica que las variables `SQLSERVER_DSN`, `JWT_SECRET`, `API_USERNAME` y `API_PASSWORD` tengan valores válidos.

---

### 2. Descargar Dependencias

```bash
go mod tidy
```

### En PowerShell / CMD:
```powershell
go run ./cmd/api
# O bien compilar y ejecutar directamente el binario:

```

### 3. Ejecutar y Compilar

#### En PowerShell / Windows CMD:

- **Compilar el ejecutable:**
  ```powershell
  go build -o api.exe ./cmd/api
  ```
- **Ejecutar el binario compilado:**
  ```powershell
  .\api.exe
  ```
- **Ejecutar directamente en desarrollo:**
  ```powershell
  go run ./cmd/api
  ```

#### En Git Bash / Linux:

- **Compilar el ejecutable:**
  ```bash
  go build -o api.exe ./cmd/api
  ```
- **Ejecutar el binario compilado (barra diagonal hacia adelante):**
  ```bash
  ./api.exe
  ```
- **Ejecutar directamente en desarrollo:**
  ```bash
  go run ./cmd/api
  ```

> [!NOTE]
> Al arrancar, la API realiza un `PING` a SQL Server; si el host no es alcanzable (por red o VPN inactiva) o la contraseña es incorrecta, la aplicación se detendrá inmediatamente reportando el error en el registro de salida.

---

### Solución de Problemas Frecuentes

#### 1. Error `no Go files in ...\galenos_backend`
* **Causa:** Ejecutar `go run .` o `go build .`. El archivo de entrada principal está en el subdirectorio `cmd/api`.
* **Solución:** Utilizar siempre `./cmd/api`:
  ```bash
  go run ./cmd/api
  # o para compilar:
  go build -o api.exe ./cmd/api
  ```

#### 2. Error `fork/exec ...\AppData\Local\go-build\...\api.exe: Acceso denegado` al usar `go run`
* **Causa:** Antivirus o EDR corporativos (como WatchGuard Endpoint Security 360 o Windows Defender) bloquean por directiva de seguridad la ejecución de binarios recién generados en directorios temporales de usuario (`%LocalAppData%\go-build`).
* **Solución recomendada:** Compilar el binario localmente en la carpeta del proyecto y ejecutarlo directo:
  ```bash
  go build -o api.exe ./cmd/api
  ./api.exe    # Git Bash
  # .\api.exe  # PowerShell
  ```
* **Limpiar caché de Go si quedó corrupta o bloqueada:**
  ```bash
  go clean -cache
  ```

#### 3. Error `bash: .api.exe: command not found`
* **Causa:** En Git Bash, la barra invertida `\` actúa como carácter de escape, transformando `.\api.exe` en `.api.exe`.
* **Solución:** En Git Bash debes utilizar siempre la barra diagonal hacia adelante:
  ```bash
  ./api.exe
  ```

#### 4. Error `fatal error: SQLSERVER_DSN environment variable is required`
* **Causa:** Falta el archivo `.env` en la carpeta `galenos_backend` o la variable de conexión no fue declarada.
* **Solución:** Crear el archivo `.env` copiando `.env.example` y configurando las credenciales requeridas.


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
