# golabs-api

REST API backend para la plataforma **GoLabs** — sistema de Capture The Flag (CTF) con gestión de eventos, equipos, retos y tabla de posiciones. Construido en Go con arquitectura hexagonal.

---

## Tabla de contenidos

- [Descripción](#descripción)
- [Requisitos](#requisitos)
- [Inicio rápido](#inicio-rápido)
- [Variables de entorno](#variables-de-entorno)
- [Configuración YAML](#configuración-yaml)
- [Comandos Make](#comandos-make)
- [Estructura del proyecto](#estructura-del-proyecto)
- [Módulos de la API](#módulos-de-la-api)
- [Endpoints principales](#endpoints-principales)
- [Autenticación](#autenticación)
- [Decisiones de diseño](#decisiones-de-diseño)

---

## Descripción

`golabs-api` expone una API HTTP REST versionada bajo `/api/v1`. Soporta:

- Registro, login y gestión de usuarios con roles (`admin` / `user`)
- Creación y ciclo de vida de eventos CTF (`draft → open → running → finished`)
- Inscripción de equipos, invitaciones con join secret y tabla de posiciones
- Retos con soporte de flags hasheadas, publicación controlada y envío de soluciones
- Refresh tokens con rotación automática
- Health checks para orquestadores (`/healthz/live`, `/healthz/ready`)

---

## Requisitos

| Herramienta | Versión mínima |
|---|---|
| Go | 1.22 |
| Docker | 24 |
| Docker Compose | v2 |
| MariaDB | 11 (via Docker) |

---

## Inicio rápido

### Con Docker (recomendado)

```bash
# 1. Copiar y configurar variables de entorno
cp .env.example .env
# Editar .env: JWT_SECRET, DB_PASSWORD, etc.

# 2. Levantar base de datos + API
make docker-up

# La API estará disponible en http://localhost:8080
```

### Local (sin Docker)

```bash
# 1. Levantar solo la base de datos
cd deployments/database && make up

# 2. Volver a la raíz y ejecutar la API
cd ../..
cp .env.example .env   # configurar credenciales apuntando a localhost
make run
```

### Hot-reload para desarrollo

```bash
# Requiere: go install github.com/air-verse/air@latest
make dev
```

---

## Variables de entorno

| Variable | Requerida | Descripción | Default |
|---|---|---|---|
| `JWT_SECRET` | ✅ | Clave para firmar tokens JWT (HS256) | — |
| `JWT_ISSUER` | ❌ | Issuer del token | `golabs-api` |
| `JWT_EXP_MINUTES` | ❌ | Duración del access token en minutos | `15` |
| `DB_HOST` | ✅ | Host de la base de datos | — |
| `DB_PORT` | ✅ | Puerto de la base de datos | `3306` |
| `DB_NAME` | ✅ | Nombre de la base de datos | — |
| `DB_USER` | ✅ | Usuario de la base de datos | — |
| `DB_PASSWORD` | ✅ | Contraseña de la base de datos | — |
| `DB_MAX_OPEN_CONNS` | ❌ | Conexiones abiertas máximas | `25` |
| `DB_MAX_IDLE_CONNS` | ❌ | Conexiones inactivas máximas | `25` |
| `DB_CONN_MAX_LIFETIME` | ❌ | Tiempo máximo de vida de conexión | `5m` |
| `ALLOWED_ORIGINS` | ❌ | Orígenes CORS permitidos (separados por coma) | `*` |

---

## Configuración YAML

El archivo `configs/config.yaml` controla ajustes adicionales del servidor:

```yaml
server:
  port: 8080

app:
  name: golabs-api
  env: development   # development | production
```

Para entornos locales, crear `configs/config.local.yaml` con las sobreescrituras necesarias (ya ignorado por `.gitignore`).

---

## Comandos Make

| Target | Descripción |
|---|---|
| `make build` | Compila el binario en `./bin/api` |
| `make run` | Compila y ejecuta localmente |
| `make dev` | Hot-reload con `air` |
| `make test` | Ejecuta todos los tests |
| `make lint` | `go vet` + `staticcheck` |
| `make tidy` | `go mod tidy` + `go mod verify` |
| `make docker-up` | Build y levanta todos los contenedores |
| `make docker-down` | Detiene los contenedores |
| `make migrate` | Lista las migraciones pendientes de `deployments/database/init/` |

---

## Estructura del proyecto

```
golabs-api/
├── cmd/api/               # Punto de entrada: main.go
├── configs/               # Configuración YAML del servidor
├── deployments/
│   └── database/          # MariaDB + Flyway (ver deployments/database/README.md)
├── internal/
│   ├── apperrors/         # Errores centinela y respuestas HTTP estandarizadas
│   ├── challenges/        # Módulo de retos CTF
│   ├── event/             # Módulo de eventos
│   ├── eventteam/         # Módulo de equipos por evento
│   ├── health/            # Health checks
│   ├── infrastructure/    # Config, DB, logging, seguridad
│   ├── interfaces/http/   # Middleware compartido (auth, access, ratelimit, pagination)
│   ├── refreshtoken/      # Módulo de refresh tokens
│   └── user/              # Módulo de usuarios
├── tools/
│   └── api-tester/        # Herramienta de pruebas manual (ver tools/api-tester/README.md)
├── Dockerfile             # Build multi-stage para producción
├── docker-compose.yml     # Orquestación de la API + DB para desarrollo
└── Makefile               # Comandos de desarrollo
```

Cada módulo en `internal/` sigue la arquitectura hexagonal:

```
<modulo>/
├── domain/          # Entidades, interfaces de repositorio, reglas de negocio
├── application/     # Casos de uso
├── interfaces/      # Handlers HTTP, DTOs, registro de rutas
└── infrastructure/  # Implementación MySQL del repositorio
```

---

## Módulos de la API

| Módulo | Descripción |
|---|---|
| `user` | Registro, login, perfil, roles, puntos, ban/unban |
| `refreshtoken` | Emisión, rotación y revocación de tokens |
| `event` | Ciclo de vida de eventos CTF |
| `eventteam` | Equipos, join secrets, leaderboard |
| `challenges` | Retos, flags, envío de soluciones |
| `health` | Liveness y readiness probes |

---

## Endpoints principales

### Autenticación

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
```

### Usuarios

```
GET    /api/v1/users                    (admin)
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
POST   /api/v1/users/:id/change-password
PUT    /api/v1/users/:id/role           (admin)
PUT    /api/v1/users/:id/points         (admin)
POST   /api/v1/users/:id/ban            (admin)
POST   /api/v1/users/:id/unban          (admin)
```

### Eventos

```
GET    /api/v1/events
GET    /api/v1/events/:id
POST   /api/v1/events                   (admin)
POST   /api/v1/events/:id/open          (admin)
POST   /api/v1/events/:id/start         (admin)
POST   /api/v1/events/:id/finish        (admin)
```

### Equipos

```
GET    /api/v1/events/:id/teams
POST   /api/v1/events/:id/teams
POST   /api/v1/events/:id/teams/join
GET    /api/v1/events/:id/teams/:team_id/members
POST   /api/v1/events/:id/teams/:team_id/leave
POST   /api/v1/events/:id/teams/:team_id/rotate-secret
GET    /api/v1/events/:id/leaderboard
```

### Retos

```
GET    /api/v1/events/:id/challenges
GET    /api/v1/events/:id/challenges/:cid
POST   /api/v1/events/:id/challenges              (admin)
PUT    /api/v1/events/:id/challenges/:cid         (admin)
POST   /api/v1/events/:id/challenges/:cid/publish (admin)
POST   /api/v1/events/:id/challenges/:cid/flag    (admin)
POST   /api/v1/events/:id/challenges/:cid/submit
```

### Health

```
GET    /healthz/live
GET    /healthz/ready
```

---

## Autenticación

La API usa JWT firmados con HS256. El acceso se gestiona con dos tokens:

| Token | Duración | Descripción |
|---|---|---|
| **Access token** | 15 min (configurable) | Enviado en el header `Authorization: Bearer <token>` |
| **Refresh token** | Larga vida | Usado para renovar el access token. Se rota en cada uso. |

El payload del access token incluye `user_id`, `role` e `issuer`.

---

## Decisiones de diseño

- **Arquitectura hexagonal**: el dominio no depende de la infraestructura. Los repositorios se definen como interfaces en el dominio e implementados en `infrastructure/`.
- **Errores centinela**: `internal/apperrors` mapea errores de negocio a códigos HTTP en un único lugar.
- **Flags hasheadas**: las flags se almacenan como SHA-256. El texto plano nunca toca la base de datos.
- **Join secrets**: se generan aleatoriamente y se almacenan hasheados. Solo se exponen en texto plano al momento de creación o rotación.
- **Token rotation**: cada uso de un refresh token lo revoca y genera uno nuevo.
- **Logging estructurado**: `log/slog` con salida JSON. Errores de servidor (5xx) se loguean; errores de cliente (4xx) no, para evitar spam.
- **Rate limiting**: aplicado en login y en envío de flags para mitigar ataques de fuerza bruta.
