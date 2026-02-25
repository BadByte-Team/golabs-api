# Seguridad — golabs-api

Este documento describe las decisiones de seguridad implementadas en la API, el razonamiento detrás de cada una y las consideraciones operativas importantes.

---

## Autenticación: JWT + Refresh Token

### Access Token (JWT)

- **Algoritmo:** HS256 (HMAC-SHA256)
- **Duración:** 15 minutos (configurable con `JWT_EXP_MINUTES`)
- **Payload:**

  ```json
  {
    "sub":    "<user_id>",
    "role":   "user|admin",
    "iss":    "golabs-api",
    "exp":    1234567890
  }
  ```

- **Validación:** firma + expiración en cada request protected.
- **Clave:** `JWT_SECRET` (variable de entorno requerida, mínimo 32 caracteres recomendado).

> [!CAUTION]
> `JWT_SECRET` debe ser una cadena aleatoria y larga. Nunca usar valores por defecto en producción.

### Refresh Token

- **Tipo:** token opaco (string aleatorio criptográfico)
- **Almacenamiento:** solo se persiste el hash SHA-256 — el token en texto plano nunca toca la BD
- **Rotación:** cada uso de un refresh token lo revoca e inmediatamente emite uno nuevo
- **Revocación:** el campo `revoked_at` en la tabla `refresh_tokens` marca el token como inválido

#### ¿Por qué token rotation?

Si un atacante roba un refresh token pero el usuario legítimo hace una petición antes, el token rotado invalida el robado. El atacante no puede reutilizarlo.

---

## Contraseñas: bcrypt

- **Algoritmo:** bcrypt con costo **12**
- **Almacenamiento:** solo el hash bcrypt en `users.password_hash`
- **Verificación:** `bcrypt.CompareHashAndPassword()` en tiempo constante

> El costo 12 implica ~250ms por hash en hardware moderno — suficiente para proteger contra ataques de fuerza bruta offline.

---

## Flags: SHA-256

- Las flags de los retos **nunca se almacenan en texto plano**
- Al configurar una flag: `SHA-256(texto_plano)` → almacenado en `flags.hash`
- Al verificar un intento: `SHA-256(intento)` se compara con el hash almacenado
- La respuesta al intento incorrecto es **idéntica** al caso sin flag configurada, para no revelar el estado del reto

---

## Join Secrets: SHA-256

- Al crear un equipo: se genera un secreto aleatorio criptográfico de 10 caracteres
- En la BD se almacena `SHA-256(secreto_crudo)` en `event_teams.join_secret_hash`
- El secreto crudo se retorna **solo una vez** al creador (o al rotar)
- Al unirse: `SHA-256(intento)` se compara con el hash almacenado

---

## Rate Limiting

### Login (`LoginRateLimit`)

Aplicado a `POST /auth/login`, `POST /auth/register`, `POST /auth/refresh`.

- Limita intentos por IP para mitigar ataques de fuerza bruta
- Implementado en memoria (sin Redis): reinicia con el servidor

### Por usuario (`UserRateLimit`)

Aplicado a todos los endpoints protegidos y especialmente a `POST /challenges/{id}/submit`.

- Limita requests por `user_id` identificado en el JWT
- Previene automatización masiva de envíos de flags

> [!NOTE]
> El rate limiting en memoria no persiste entre reinicios del servidor y no es efectivo en despliegues multi-instancia. Para producción a escala, reemplazar con Redis.

---

## Control de acceso

### Middleware en cadena

```
JWTAuth          → valida firma y expiración del token
  └→ LoadUser    → recarga rol y estado banned desde BD (evita usar datos stale del token)
       └→ RequireNotBanned → retorna 403 si banned=true
            └→ UserRateLimit
                 └→ RequireRole      → verifica rol mínimo
                      └→ RequireSelfOrAdmin → verifica identidad o rol admin
```

### Separación `JWTAuth` / `LoadUser`

- Las rutas de **admin** solo usan `JWTAuth` (el rol está en el token), evitando una consulta extra a BD por petición.
- Las rutas de **participantes** usan `LoadUser` para recargar el estado de baneo en tiempo real — un usuario baneado no puede esperar a que expire su JWT.

---

## Protecciones HTTP globales

| Middleware | Propósito |
|---|---|
| `RequestID` | Traza única por petición para correlacionar logs |
| `RealIP` | Respeta headers `X-Forwarded-For` / `X-Real-IP` de proxies |
| `Recoverer` | Captura panics y retorna 500 sin crashear el servidor |
| `MaxBodySize` | Limita el tamaño del body a ~1MB para prevenir ataques de memoria |
| `CORS` | Permite solo orígenes definidos en `ALLOWED_ORIGINS` |

---

## Manejo de errores seguro

- Los errores de dominio se mapean a códigos HTTP en `apperrors` sin exponer detalles internos
- Los errores 5xx se loguean con `slog`; los errores 4xx no se loguean para evitar spam
- Los mensajes de error al cliente son intencionalmente vagos donde la precisión podría dar información a un atacante (ej. "flag incorrecta" vs "flag no configurada")

---

## Checklist de producción

- [ ] `JWT_SECRET` generado con `openssl rand -hex 32` o similar
- [ ] `DB_PASSWORD` y `MARIADB_ROOT_PASSWORD` únicos y complejos
- [ ] `ALLOWED_ORIGINS` configurado con los orígenes exactos (no `*`)
- [ ] HTTPS en el reverse proxy (nginx, Caddy, etc.)
- [ ] Rate limiting respaldado por Redis en despliegues multi-instancia
- [ ] Backups automáticos de BD habilitados (`make backup-auto`)
- [ ] Logs enviados a un sistema central (Loki, CloudWatch, etc.)
- [ ] Healthchecks configurados en el orquestador (`/healthz/ready`)
