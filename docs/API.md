# API Reference — golabs-api

Base URL: `http://localhost:8080`  
Versión actual: `/api/v1`

---

## Tabla de contenidos

1. [Convenciones generales](#1-convenciones-generales)
2. [Autenticación](#2-autenticación)
3. [Errores](#3-errores)
4. [Módulo: Auth](#4-módulo-auth)
5. [Módulo: Users](#5-módulo-users)
6. [Módulo: Events](#6-módulo-events)
7. [Módulo: Event Teams](#7-módulo-event-teams)
8. [Módulo: Challenges](#8-módulo-challenges)
9. [Módulo: Health](#9-módulo-health)

---

## 1. Convenciones generales

- Todas las rutas del módulo están prefijadas con `/api/v1/` salvo `/healthz/*`.
- El cuerpo de las peticiones debe enviarse como `application/json`.
- Todas las respuestas son `application/json`.
- Los UUIDs se representan como strings en formato `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.
- Los timestamps siguen el formato **RFC 3339** en UTC (ej. `"2025-01-15T14:30:00Z"`).
- La paginación se controla con `?page=1&size=20` (defaults: página 1, 20 ítems).

### Respuesta paginada

Todos los endpoints que retornan listas paginadas responden con:

```json
{
  "data": [ ... ],
  "meta": {
    "page":  1,
    "size":  20,
    "total": 100
  }
}
```

---

## 2. Autenticación

Los endpoints protegidos requieren el header:

```
Authorization: Bearer <access_token>
```

El **access token** es un JWT firmado con HS256, con una duración de 15 minutos.  
El **refresh token** es un token opaco de larga duración que se usa para renovar el access token.

### Roles

| Rol | Acceso |
|---|---|
| `user` | Participante en CTFs (rol por defecto) |
| `admin` | Acceso completo a endpoints de gestión |

---

## 3. Errores

Todos los errores siguen el formato:

```json
{
  "error": "mensaje descriptivo del error"
}
```

### Códigos HTTP utilizados

| Código | Descripción |
|---|---|
| `400 Bad Request` | Body inválido, parámetros faltantes o con formato incorrecto |
| `401 Unauthorized` | Token ausente, inválido, expirado o credenciales incorrectas |
| `403 Forbidden` | Autenticado pero sin permisos suficientes (rol incorrecto o usuario baneado) |
| `404 Not Found` | Recurso no encontrado |
| `409 Conflict` | Conflicto de unicidad (email/username duplicado, solve duplicado, etc.) |
| `422 Unprocessable Entity` | Validación de campos fallida |
| `500 Internal Server Error` | Error interno del servidor |

---

## 4. Módulo: Auth

Endpoints públicos de autenticación (sin token requerido). Todos aplican **rate limiting** por IP.

---

### `POST /api/v1/auth/register`

Crea una nueva cuenta de usuario.

**Body**

```json
{
  "username": "string (3–30 chars)",
  "email":    "string (email válido)",
  "password": "string (mín. 6 chars)"
}
```

**Respuesta exitosa — `201 Created`**

```json
{
  "id":         "uuid",
  "username":   "hackerman",
  "email":      "hacker@example.com",
  "role":       "user",
  "points":     0,
  "banned":     false,
  "banned_at":  null,
  "created_at": "2025-01-15T14:30:00Z",
  "updated_at": "2025-01-15T14:30:00Z"
}
```

**Errores**

| Código | Causa |
|---|---|
| `400` | Body inválido o campos faltantes |
| `409` | Email o username ya registrado |

---

### `POST /api/v1/auth/login`

Autentica al usuario y retorna un par de tokens.  
El campo `identifier` acepta **email o username**.

**Body**

```json
{
  "identifier": "string (email o username)",
  "password":   "string"
}
```

**Respuesta exitosa — `200 OK`**

```json
{
  "access_token":  "eyJhbGci...",
  "refresh_token": "a1b2c3d4...",
  "expires_in":    900
}
```

> `expires_in` está en segundos. El access token dura 15 minutos (900s).

**Errores**

| Código | Causa |
|---|---|
| `400` | Body inválido |
| `401` | Credenciales incorrectas |
| `403` | Usuario baneado |

---

### `POST /api/v1/auth/refresh`

Rota el refresh token y emite un nuevo par de tokens.  
El token anterior queda **revocado** inmediatamente (token rotation).

**Body**

```json
{
  "refresh_token": "string"
}
```

**Respuesta exitosa — `200 OK`**

```json
{
  "access_token":  "eyJhbGci...",
  "refresh_token": "e5f6g7h8...",
  "expires_in":    900
}
```

**Errores**

| Código | Causa |
|---|---|
| `400` | Body inválido |
| `401` | Token inválido, expirado o ya revocado |

---

### `POST /api/v1/auth/logout`

Revoca el refresh token para invalidar la sesión actual.  
**Siempre responde `204`** independientemente del estado del token (idempotente).

**Body**

```json
{
  "refresh_token": "string"
}
```

**Respuesta exitosa — `204 No Content`**

---

## 5. Módulo: Users

Gestión de perfiles de usuario. La mayoría de endpoints requieren autenticación.

---

### `GET /api/v1/users/` _(admin)_

Lista todos los usuarios de forma paginada.

**Auth:** JWT + rol `admin`  
**Query params:** `page`, `size`

**Respuesta exitosa — `200 OK`**

```json
{
  "data": [
    {
      "id":         "uuid",
      "username":   "hackerman",
      "email":      "hacker@example.com",
      "role":       "user",
      "points":     150,
      "banned":     false,
      "banned_at":  null,
      "created_at": "2025-01-15T14:30:00Z",
      "updated_at": "2025-01-15T14:30:00Z"
    }
  ],
  "meta": { "page": 1, "size": 20, "total": 45 }
}
```

**Errores:** `401`, `403`

---

### `POST /api/v1/users/` _(admin)_

Crea un usuario directamente (equivalente a `/auth/register` pero para admins).

**Auth:** JWT + rol `admin`  
**Body:** igual que `POST /auth/register`  
**Respuesta exitosa — `201 Created`:** igual que el registro

**Errores:** `400`, `401`, `403`, `409`

---

### `GET /api/v1/users/{id}`

Retorna el perfil de un usuario por su UUID.

**Auth:** JWT (cualquier rol)

**Respuesta exitosa — `200 OK`**

```json
{
  "id":         "uuid",
  "username":   "hackerman",
  "email":      "hacker@example.com",
  "role":       "user",
  "points":     150,
  "banned":     false,
  "banned_at":  null,
  "created_at": "2025-01-15T14:30:00Z",
  "updated_at": "2025-01-15T14:30:00Z"
}
```

**Errores:** `401`, `404`

---

### `GET /api/v1/users/by-username/{username}`

Retorna el perfil de un usuario por su username exacto.

**Auth:** JWT (cualquier rol)  
**Respuesta:** igual que `GET /users/{id}`  
**Errores:** `401`, `404`

---

### `GET /api/v1/users/search?q={term}`

Busca usuarios cuyo username contenga el término (búsqueda parcial, máx. 20 resultados).

**Auth:** JWT (cualquier rol)  
**Query params:** `q` — término de búsqueda (requerido)

**Respuesta exitosa — `200 OK`**

```json
[
  { "id": "uuid", "username": "hackerman", "email": "...", "role": "user", "points": 0, "banned": false, "banned_at": null, "created_at": "...", "updated_at": "..." }
]
```

**Errores:** `400` (q vacío), `401`

---

### `POST /api/v1/users/{id}/update`

Actualiza username y/o email del usuario (semántica PATCH).  
Solo el propio usuario o un admin pueden modificar el perfil.

**Auth:** JWT + ser el propio usuario o `admin`

**Body**

```json
{
  "username": "nuevo_nombre",
  "email":    "nuevo@email.com"
}
```

> Ambos campos son opcionales. Los campos omitidos no se modifican.

**Respuesta exitosa — `200 OK`:** perfil actualizado (igual que `GET /users/{id}`)

**Errores:** `400`, `401`, `403`, `409` (username/email duplicado)

---

### `POST /api/v1/users/{id}/password`

Cambia la contraseña del usuario. Requiere confirmar la contraseña actual.

**Auth:** JWT + ser el propio usuario o `admin`

**Body**

```json
{
  "current_password": "string",
  "new_password":     "string (mín. 6 chars)"
}
```

**Respuesta exitosa — `204 No Content`**

**Errores:** `400`, `401`, `403` (contraseña actual incorrecta)

---

### `POST /api/v1/admin/users/{id}/role` _(admin)_

Cambia el rol de un usuario.

**Auth:** JWT + rol `admin`

**Body**

```json
{
  "role": "admin" | "user"
}
```

**Respuesta exitosa — `204 No Content`**

**Errores:** `400`, `401`, `403`, `404`

---

### `POST /api/v1/admin/users/{id}/points` _(admin)_

Establece los puntos de un usuario de forma absoluta (no incremental).

**Auth:** JWT + rol `admin`

**Body**

```json
{
  "points": 500
}
```

**Respuesta exitosa — `204 No Content`**

**Errores:** `400`, `401`, `403`, `404`

---

### `POST /api/v1/admin/users/{id}/ban` _(admin)_

Suspende el acceso del usuario.

**Auth:** JWT + rol `admin`

**Respuesta exitosa — `200 OK`**

```json
{
  "banned":    true,
  "banned_at": "2025-01-15T15:00:00Z"
}
```

**Errores:** `401`, `403`, `404`

---

### `POST /api/v1/admin/users/{id}/unban` _(admin)_

Reactiva el acceso de un usuario baneado.

**Auth:** JWT + rol `admin`

**Respuesta exitosa — `200 OK`**

```json
{
  "banned":    false,
  "banned_at": null
}
```

**Errores:** `401`, `403`, `404`

---

## 6. Módulo: Events

Gestión del ciclo de vida de eventos CTF.

### Estados de un evento

```
draft ──→ open ──→ running ──→ finished
```

| Estado | Descripción |
|---|---|
| `draft` | Creado, no visible para participantes |
| `open` | Aceptando inscripciones de equipos |
| `running` | Evento en curso, flag submissions activas |
| `finished` | Evento terminado, solo lectura |

---

### `GET /api/v1/events/`

Lista todos los eventos (paginado).

**Auth:** ninguna  
**Query params:** `page`, `size`

**Respuesta exitosa — `200 OK`**

```json
{
  "data": [
    {
      "id":           "uuid",
      "name":         "GoLabs CTF 2025",
      "description":  "Competencia anual de CTF.",
      "max_team_size": 4,
      "status":       "open",
      "starts_at":    "2025-03-01T00:00:00Z",
      "ends_at":      "2025-03-03T00:00:00Z",
      "created_at":   "2025-01-15T14:30:00Z",
      "updated_at":   "2025-01-15T14:30:00Z"
    }
  ],
  "meta": { "page": 1, "size": 20, "total": 3 }
}
```

---

### `GET /api/v1/events/{event_id}`

Retorna el detalle de un evento.

**Auth:** ninguna  
**Respuesta exitosa — `200 OK`:** igual que un ítem de la lista  
**Errores:** `404`

---

### `POST /api/v1/events/` _(admin)_

Crea un nuevo evento en estado `draft`.

**Auth:** JWT + rol `admin`

**Body**

```json
{
  "name":         "GoLabs CTF 2025",
  "description":  "Descripción del evento (máx. 1000 chars)",
  "max_team_size": 4,
  "starts_at":    "2025-03-01T00:00:00Z",
  "ends_at":      "2025-03-03T00:00:00Z"
}
```

**Respuesta exitosa — `201 Created`:** objeto `EventResponse`

**Errores:** `400`, `401`, `403`, `422`

---

### `POST /api/v1/events/{event_id}/open` _(admin)_

Transiciona el evento de `draft` → `open`.

**Auth:** JWT + rol `admin`  
**Body:** vacío  
**Respuesta exitosa — `204 No Content`**

**Errores:** `400` (estado incorrecto), `401`, `403`, `404`

---

### `POST /api/v1/events/{event_id}/start` _(admin)_

Transiciona el evento de `open` → `running`.

**Auth:** JWT + rol `admin`  
**Body:** vacío  
**Respuesta exitosa — `204 No Content`**

**Errores:** `400` (estado incorrecto), `401`, `403`, `404`

---

### `POST /api/v1/events/{event_id}/finish` _(admin)_

Transiciona el evento de `running` → `finished`. No se aceptan más flag submissions.

**Auth:** JWT + rol `admin`  
**Body:** vacío  
**Respuesta exitosa — `204 No Content`**

**Errores:** `400` (estado incorrecto), `401`, `403`, `404`

---

## 7. Módulo: Event Teams

Gestión de equipos dentro de un evento. El evento debe estar en estado `open` o superior para la mayoría de operaciones.

---

### `GET /api/v1/events/{event_id}/teams`

Lista todos los equipos de un evento.

**Auth:** JWT (cualquier rol no baneado)

**Respuesta exitosa — `200 OK`**

```json
[
  {
    "id":           "uuid",
    "event_id":     "uuid",
    "name":         "Byte Busters",
    "owner_id":     "uuid",
    "score":        300,
    "created_at":   "2025-01-15T14:30:00Z",
    "updated_at":   "2025-01-15T14:30:00Z"
  }
]
```

**Errores:** `401`, `403` (baneado), `404` (evento no existe)

---

### `POST /api/v1/events/{event_id}/teams`

Crea un nuevo equipo en el evento. El creador queda como owner y primer miembro.

**Auth:** JWT (cualquier rol no baneado)

**Body**

```json
{
  "name": "Byte Busters"
}
```

**Respuesta exitosa — `201 Created`**

```json
{
  "team": {
    "id":         "uuid",
    "event_id":   "uuid",
    "name":       "Byte Busters",
    "owner_id":   "uuid",
    "score":      0,
    "created_at": "2025-01-15T14:30:00Z",
    "updated_at": "2025-01-15T14:30:00Z"
  },
  "join_secret": "AbCd3FgHiJ"
}
```

> `join_secret` es el código para invitar miembros. **Solo se muestra en este momento**; si se pierde, usar `/rotate-secret`.

**Errores**

| Código | Causa |
|---|---|
| `400` | Body inválido o evento no en estado correcto |
| `401` | Sin autenticación |
| `403` | Usuario baneado |
| `409` | El usuario ya está en un equipo en este evento, o el nombre de equipo ya existe |

---

### `POST /api/v1/events/{event_id}/teams/join`

Uniirse a un equipo existente usando el join secret.

**Auth:** JWT (cualquier rol no baneado)

**Body**

```json
{
  "join_secret": "AbCd3FgHiJ"
}
```

**Respuesta exitosa — `200 OK`:** objeto del equipo (`EventTeamResponse`)

**Errores**

| Código | Causa |
|---|---|
| `400` | Join secret incorrecto, equipo lleno, o evento no acepta inscripciones |
| `401` | Sin autenticación |
| `403` | Usuario baneado |
| `409` | El usuario ya pertenece a un equipo en este evento |

---

### `GET /api/v1/events/{event_id}/teams/{team_id}/members`

Lista los miembros de un equipo junto con sus usernames.

**Auth:** JWT (cualquier rol no baneado)

**Respuesta exitosa — `200 OK`**

```json
[
  {
    "user_id":  "uuid",
    "username": "hackerman",
    "team_id":  "uuid"
  }
]
```

**Errores:** `401`, `403`, `404`

---

### `POST /api/v1/events/{event_id}/teams/{team_id}/leave`

Abandona el equipo actual. El owner no puede abandonar si quedan miembros.

**Auth:** JWT (cualquier rol no baneado)  
**Body:** vacío  
**Respuesta exitosa — `204 No Content`**

**Errores**

| Código | Causa |
|---|---|
| `400` | El owner no puede abandonar con miembros restantes |
| `401` | Sin autenticación |
| `403` | Usuario baneado o no pertenece al equipo |
| `404` | Equipo no encontrado |

---

### `POST /api/v1/events/{event_id}/teams/{team_id}/rotate-secret`

Genera un nuevo join secret e invalida el anterior. Solo el owner puede ejecutarlo.

**Auth:** JWT (owner del equipo, no baneado)  
**Body:** vacío

**Respuesta exitosa — `200 OK`**

```json
{
  "join_secret": "NuEv0S3cr3t"
}
```

**Errores**

| Código | Causa |
|---|---|
| `401` | Sin autenticación |
| `403` | Usuario baneado o no es el owner |
| `404` | Equipo no encontrado |

---

### `GET /api/v1/events/{event_id}/leaderboard`

Retorna la tabla de posiciones del evento ordenada por score descendente.

**Auth:** JWT (cualquier rol no baneado)

**Respuesta exitosa — `200 OK`**

```json
[
  {
    "id":       "uuid",
    "event_id": "uuid",
    "name":     "Byte Busters",
    "owner_id": "uuid",
    "score":    750,
    "created_at": "...",
    "updated_at": "..."
  }
]
```

**Errores:** `401`, `403`

---

## 8. Módulo: Challenges

Gestión de retos CTF dentro de un evento.

### Categorías válidas

`web` · `pwn` · `rev` · `crypto` · `forensics` · `misc`

### Dificultades válidas

`easy` · `medium` · `hard`

---

### `GET /api/v1/events/{event_id}/challenges`

Lista los challenges del evento.

- Los **admins** ven todos los challenges (incluyendo ocultos).
- Los **participantes** solo ven los challenges con `visible=true`.

**Auth:** JWT (cualquier rol no baneado)  
**Query params (opcionales):**

| Param | Descripción |
|---|---|
| `category` | Filtrar por categoría (ej. `web`, `pwn`) |
| `difficulty` | Filtrar por dificultad (`easy`, `medium`, `hard`) |

**Respuesta exitosa — `200 OK`**

```json
[
  {
    "id":                  "uuid",
    "event_id":            "uuid",
    "title":               "SQL Injection 101",
    "description":         "Encuentra la flag inyectando SQL.",
    "category":            "web",
    "points":              200,
    "difficulty":          "easy",
    "visible":             true,
    "solve_count":         12,
    "first_blood_team_id": "uuid",
    "created_at":          "2025-01-15T14:30:00Z",
    "updated_at":          "2025-01-15T14:30:00Z"
  }
]
```

> `first_blood_team_id` es `null` si nadie ha resuelto el challenge aún.

**Errores:** `401`, `403`

---

### `GET /api/v1/events/{event_id}/challenges/{challenge_id}`

Retorna el detalle de un challenge.

- Los challenges **ocultos** retornan `404` para no-admins (para no revelar su existencia).

**Auth:** JWT  
**Respuesta exitosa — `200 OK`:** objeto `ChallengeResponse` (igual que un ítem de la lista)

**Errores:** `401`, `404`

---

### `POST /api/v1/events/{event_id}/challenges` _(admin)_

Crea un challenge en estado **oculto** (`visible=false`).

> El challenge no será visible para participantes hasta llamar a `/publish` y tener una flag configurada via `/flag`.

**Auth:** JWT + rol `admin`

**Body**

```json
{
  "title":       "SQL Injection 101",
  "description": "Encuentra la flag inyectando SQL en el formulario de login.",
  "category":    "web",
  "points":      200,
  "difficulty":  "easy"
}
```

**Respuesta exitosa — `201 Created`:** objeto `ChallengeResponse`

**Errores:** `400`, `401`, `403`, `404` (evento no existe), `422`

---

### `PUT /api/v1/events/{event_id}/challenges/{challenge_id}` _(admin)_

Actualiza los metadatos de un challenge existente. La flag se gestiona por separado.

**Auth:** JWT + rol `admin`  
**Body:** igual que `POST /challenges`  
**Respuesta exitosa — `200 OK`:** objeto `ChallengeResponse` actualizado

**Errores:** `400`, `401`, `403`, `404`, `422`

---

### `POST /api/v1/events/{event_id}/challenges/{challenge_id}/publish` _(admin)_

Hace visible el challenge para los participantes (`visible=true`).

**Auth:** JWT + rol `admin`  
**Body:** vacío  
**Respuesta exitosa — `200 OK`:** objeto `ChallengeResponse` con `visible: true`

**Errores:** `401`, `403`, `404`

---

### `POST /api/v1/events/{event_id}/challenges/{challenge_id}/unpublish` _(admin)_

Oculta el challenge para los participantes (`visible=false`).

**Auth:** JWT + rol `admin`  
**Body:** vacío  
**Respuesta exitosa — `200 OK`:** objeto `ChallengeResponse` con `visible: false`

**Errores:** `401`, `403`, `404`

---

### `POST /api/v1/events/{event_id}/challenges/{challenge_id}/flag` _(admin)_

Establece o reemplaza la flag del challenge.

> El texto plano **nunca se almacena**. El servidor almacena exclusivamente el hash SHA-256.

**Auth:** JWT + rol `admin`

**Body**

```json
{
  "flag": "CTF{s3cr3t_fl4g_here}"
}
```

**Respuesta exitosa — `204 No Content`**

**Errores:** `400`, `401`, `403`, `404`, `422`

---

### `POST /api/v1/events/{event_id}/challenges/{challenge_id}/submit`

Envía una flag para un challenge. Si es correcta, registra el solve y suma los puntos al equipo.

**Auth:** JWT (participante no baneado, con rate limiting por usuario)

**Body**

```json
{
  "flag": "CTF{mi_intento}"
}
```

**Respuesta exitosa — `200 OK`**

```json
{
  "correct": true,
  "points":  200
}
```

Si la flag es incorrecta:

```json
{
  "correct": false
}
```

> La respuesta es intencionalmente **vaga**: "flag sin configurar" y "flag incorrecta" retornan el mismo resultado para no revelar el estado interno del challenge.

**Errores**

| Código | Causa |
|---|---|
| `400` | Body inválido, evento no en curso, challenge no visible, sin equipo |
| `401` | Sin autenticación |
| `403` | Usuario baneado |
| `409` | El equipo ya resolvió este challenge |

---

## 9. Módulo: Health

Endpoints de verificación de salud del servidor, destinados a orquestadores (Kubernetes, Docker Compose, load balancers).

---

### `GET /healthz/live`

Liveness probe — indica que el servidor está corriendo.

**Auth:** ninguna  
**Respuesta exitosa — `200 OK`**

```json
{ "status": "ok" }
```

---

### `GET /healthz/ready`

Readiness probe — indica que el servidor está listo para recibir tráfico (incluye verificación de conexión a BD).

**Auth:** ninguna  
**Respuesta exitosa — `200 OK`**

```json
{ "status": "ok" }
```

**Error — `503 Service Unavailable`**

```json
{ "status": "unavailable" }
```

---

## Apéndice: Flujo típico de participante

```
1. POST /auth/register              → crear cuenta
2. POST /auth/login                 → obtener access_token + refresh_token
3. GET  /events/                    → explorar eventos disponibles
4. POST /events/{id}/teams          → crear equipo (obtienes join_secret)
   — o —
   POST /events/{id}/teams/join     → unirse con join_secret de otro equipo
5. GET  /events/{id}/challenges     → ver retos del evento
6. POST /events/{id}/challenges/{cid}/submit  → enviar flag
7. GET  /events/{id}/leaderboard    → ver posiciones
8. POST /auth/refresh               → renovar access_token antes de que expire
9. POST /auth/logout                → cerrar sesión
```

## Apéndice: Flujo típico de admin

```
1. POST /auth/login                        → login como admin
2. POST /events/                           → crear evento (estado: draft)
3. POST /events/{id}/open                  → abrir inscripciones (estado: open)
4. POST /events/{id}/challenges            → crear challenge (oculto)
5. POST /events/{id}/challenges/{cid}/flag → configurar flag
6. POST /events/{id}/challenges/{cid}/publish  → publicar challenge
7. POST /events/{id}/start                 → iniciar evento (estado: running)
   ... el CTF corre ...
8. POST /events/{id}/finish                → finalizar evento (estado: finished)
```
