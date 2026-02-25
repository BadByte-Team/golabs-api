# Guía de integración de cliente — golabs-api

Esta guía explica cómo conectar un frontend o cualquier cliente HTTP contra la API, desde el registro hasta el uso de todos los módulos. Los ejemplos están en **JavaScript (navegador)**, **Python** y **Go**.

---

## Índice

1. [Almacenamiento de tokens en el navegador](#1-almacenamiento-de-tokens-en-el-navegador)
2. [Cliente HTTP base](#2-cliente-http-base)
3. [Autenticación: registro y login](#3-autenticación-registro-y-login)
4. [Manejo de expiración: auto-refresh](#4-manejo-de-expiración-auto-refresh)
5. [Logout](#5-logout)
6. [Módulo: Users](#6-módulo-users)
7. [Módulo: Events](#7-módulo-events)
8. [Módulo: Teams](#8-módulo-teams)
9. [Módulo: Challenges](#9-módulo-challenges)
10. [Flujo completo de participante](#10-flujo-completo-de-participante)
11. [Flujo completo de administrador](#11-flujo-completo-de-administrador)

---

## 1. Almacenamiento de tokens en el navegador

El servidor emite dos tokens al hacer login:

| Token | Tipo | Duración | Propósito |
|---|---|---|---|
| `access_token` | JWT | 15 min | Autenticar cada request |
| `refresh_token` | Opaco | Larga duración | Obtener un nuevo par de tokens |

### Opciones de almacenamiento

**`localStorage`** — más fácil, pero vulnerable a XSS:

```javascript
// Guardar tokens
localStorage.setItem('access_token', response.access_token);
localStorage.setItem('refresh_token', response.refresh_token);

// Leer
const token = localStorage.getItem('access_token');

// Borrar (logout)
localStorage.removeItem('access_token');
localStorage.removeItem('refresh_token');
```

**`sessionStorage`** — igual que localStorage pero se borra al cerrar la pestaña:

```javascript
sessionStorage.setItem('access_token', response.access_token);
sessionStorage.setItem('refresh_token', response.refresh_token);
```

**En memoria (más seguro contra XSS, recomendado para acceso a APIs sensibles)**:

```javascript
// tokens.js — módulo singleton
let _accessToken  = null;
let _refreshToken = null;

export const tokens = {
    set(access, refresh) {
        _accessToken  = access;
        _refreshToken = refresh;
    },
    getAccess()  { return _accessToken; },
    getRefresh() { return _refreshToken; },
    clear()      { _accessToken = null; _refreshToken = null; }
};
```

> **Recomendación práctica**: guardar el `refresh_token` en `localStorage` (para sobrevivir recargas de página) y el `access_token` solo en memoria. De esta forma si hay XSS, el atacante solo obtiene el refresh token (que puede revocarse desde el servidor).

---

## 2. Cliente HTTP base

Construir un cliente centralizado simplifica el manejo de headers, errores y refresh automático.

### JavaScript

```javascript
// api.js
const BASE = 'http://localhost:8080/api/v1';

// Almacenamiento en memoria para el access token
let _accessToken  = localStorage.getItem('access_token')  || null;
let _refreshToken = localStorage.getItem('refresh_token') || null;

export function setTokens(access, refresh) {
    _accessToken  = access;
    _refreshToken = refresh;
    localStorage.setItem('access_token',  access);
    localStorage.setItem('refresh_token', refresh);
}

export function clearTokens() {
    _accessToken = _refreshToken = null;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
}

// Ejecuta un fetch; si recibe 401, intenta refresh y reintenta una vez
async function request(method, path, body = null, retry = true) {
    const headers = { 'Content-Type': 'application/json' };
    if (_accessToken) headers['Authorization'] = `Bearer ${_accessToken}`;

    const res = await fetch(BASE + path, {
        method,
        headers,
        body: body ? JSON.stringify(body) : null,
    });

    // Token expirado → intentar refresh automático
    if (res.status === 401 && retry && _refreshToken) {
        const ok = await refreshTokens();
        if (ok) return request(method, path, body, false); // reintento
        clearTokens();
        window.location.href = '/login'; // redirigir si el refresh también falla
        return;
    }

    if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }));
        throw Object.assign(new Error(err.error || 'API error'), { status: res.status, data: err });
    }

    if (res.status === 204) return null;
    return res.json();
}

async function refreshTokens() {
    try {
        const res = await fetch(`${BASE}/auth/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: _refreshToken }),
        });
        if (!res.ok) return false;
        const data = await res.json();
        setTokens(data.access_token, data.refresh_token);
        return true;
    } catch { return false; }
}

// Helpers públicos
export const api = {
    get:    (path)         => request('GET',    path),
    post:   (path, body)   => request('POST',   path, body),
    put:    (path, body)   => request('PUT',    path, body),
    patch:  (path, body)   => request('PATCH',  path, body),
    delete: (path)         => request('DELETE', path),
};
```

---

### Python

```python
# api_client.py
import requests

BASE = "http://localhost:8080/api/v1"

class GoLabsClient:
    def __init__(self):
        self.session = requests.Session()
        self.access_token  = None
        self.refresh_token = None

    def _headers(self):
        h = {"Content-Type": "application/json"}
        if self.access_token:
            h["Authorization"] = f"Bearer {self.access_token}"
        return h

    def _request(self, method, path, json=None, retry=True):
        res = self.session.request(
            method, BASE + path, json=json, headers=self._headers()
        )

        # Auto-refresh si el access token expiró
        if res.status_code == 401 and retry and self.refresh_token:
            if self._do_refresh():
                return self._request(method, path, json, retry=False)
            raise PermissionError("Sesión expirada. Hace login de nuevo.")

        res.raise_for_status()
        if res.status_code == 204:
            return None
        return res.json()

    def _do_refresh(self):
        try:
            res = self.session.post(
                f"{BASE}/auth/refresh",
                json={"refresh_token": self.refresh_token},
                headers={"Content-Type": "application/json"},
            )
            if not res.ok:
                return False
            data = res.json()
            self.access_token  = data["access_token"]
            self.refresh_token = data["refresh_token"]
            return True
        except Exception:
            return False

    def get(self, path):              return self._request("GET",   path)
    def post(self, path, body=None):  return self._request("POST",  path, json=body)
    def put(self, path, body=None):   return self._request("PUT",   path, json=body)
    def patch(self, path, body=None): return self._request("PATCH", path, json=body)

client = GoLabsClient()
```

---

### Go

```go
// client/client.go
package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
)

const base = "http://localhost:8080/api/v1"

type Client struct {
    http         *http.Client
    mu           sync.Mutex
    accessToken  string
    refreshToken string
}

func New() *Client {
    return &Client{http: &http.Client{}}
}

func (c *Client) SetTokens(access, refresh string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.accessToken  = access
    c.refreshToken = refresh
}

func (c *Client) do(method, path string, body any, retry bool) (*http.Response, error) {
    var buf *bytes.Buffer
    if body != nil {
        b, _ := json.Marshal(body)
        buf = bytes.NewBuffer(b)
    } else {
        buf = &bytes.Buffer{}
    }

    req, err := http.NewRequest(method, base+path, buf)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    c.mu.Lock()
    if c.accessToken != "" {
        req.Header.Set("Authorization", "Bearer "+c.accessToken)
    }
    c.mu.Unlock()

    res, err := c.http.Do(req)
    if err != nil {
        return nil, err
    }

    if res.StatusCode == http.StatusUnauthorized && retry {
        if err := c.doRefresh(); err == nil {
            res.Body.Close()
            return c.do(method, path, body, false)
        }
        return nil, fmt.Errorf("sesión expirada")
    }
    return res, nil
}

func (c *Client) doRefresh() error {
    c.mu.Lock()
    rt := c.refreshToken
    c.mu.Unlock()

    body, _ := json.Marshal(map[string]string{"refresh_token": rt})
    res, err := c.http.Post(base+"/auth/refresh", "application/json", bytes.NewBuffer(body))
    if err != nil || res.StatusCode != http.StatusOK {
        return fmt.Errorf("refresh falló")
    }
    defer res.Body.Close()
    var data struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
    }
    if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
        return err
    }
    c.SetTokens(data.AccessToken, data.RefreshToken)
    return nil
}

// Decode unmarshals the response body into dest.
func Decode(res *http.Response, dest any) error {
    defer res.Body.Close()
    return json.NewDecoder(res.Body).Decode(dest)
}
```

---

## 3. Autenticación: registro y login

### JavaScript

```javascript
import { api, setTokens } from './api.js';

// Registro
async function register(username, email, password) {
    const user = await api.post('/auth/register', { username, email, password });
    console.log('Usuario creado:', user.id);
    return user;
}

// Login
async function login(identifier, password) {
    const data = await api.post('/auth/login', { identifier, password });
    setTokens(data.access_token, data.refresh_token);
    console.log(`Bienvenido! Token expira en ${data.expires_in}s`);
    return data;
}

// Uso
await register('hackerman', 'hack@example.com', 'secret123');
await login('hackerman', 'secret123');
```

---

### Python

```python
# Registro
def register(username, email, password):
    return client.post('/auth/register', {
        'username': username,
        'email':    email,
        'password': password,
    })

# Login — guarda los tokens en el cliente automáticamente
def login(identifier, password):
    data = client.post('/auth/login', {
        'identifier': identifier,
        'password':   password,
    })
    client.access_token  = data['access_token']
    client.refresh_token = data['refresh_token']
    return data

# Uso
register('hackerman', 'hack@example.com', 'secret123')
login('hackerman', 'secret123')
print('Access token:', client.access_token[:20], '...')
```

---

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "golabs-client/client"
)

func main() {
    c := client.New()

    // Registro
    res, _ := c.do("POST", "/auth/register", map[string]string{
        "username": "hackerman",
        "email":    "hack@example.com",
        "password": "secret123",
    }, false)
    var user map[string]any
    client.Decode(res, &user)
    fmt.Println("Usuario creado:", user["id"])

    // Login
    res, _ = c.do("POST", "/auth/login", map[string]string{
        "identifier": "hackerman",
        "password":   "secret123",
    }, false)
    var auth struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    int    `json:"expires_in"`
    }
    client.Decode(res, &auth)
    c.SetTokens(auth.AccessToken, auth.RefreshToken)
    fmt.Println("Login exitoso. Token expira en", auth.ExpiresIn, "segundos")
}
```

---

## 4. Manejo de expiración: auto-refresh

El cliente base en la sección 2 ya maneja el refresh automático. Si preferís controlarlo manualmente, aquí el patrón:

### JavaScript — usando un temporizador

```javascript
// refreshScheduler.js
let refreshTimer = null;

// Programa un refresh 60 segundos antes de que el access token expire
export function scheduleRefresh(expiresInSeconds) {
    clearTimeout(refreshTimer);
    const delay = Math.max((expiresInSeconds - 60) * 1000, 5000);
    refreshTimer = setTimeout(async () => {
        try {
            const data = await fetch(BASE + '/auth/refresh', {
                method:  'POST',
                headers: { 'Content-Type': 'application/json' },
                body:    JSON.stringify({ refresh_token: _refreshToken }),
            }).then(r => r.json());
            setTokens(data.access_token, data.refresh_token);
            scheduleRefresh(data.expires_in); // reprogramar para el siguiente
        } catch {
            clearTokens();
        }
    }, delay);
}

// En el login:
const data = await api.post('/auth/login', { identifier, password });
setTokens(data.access_token, data.refresh_token);
scheduleRefresh(data.expires_in); // iniciar el ciclo de auto-refresh
```

---

### Python — para scripts de larga duración

```python
import time
import threading

def start_auto_refresh(expires_in: int):
    """Inicia un hilo que renueva el token 60s antes de que expire."""
    delay = max(expires_in - 60, 5)

    def _refresh():
        time.sleep(delay)
        if client._do_refresh():
            print("[auth] Token renovado correctamente")
            # No reprogramamos aquí; se re-llamará en el próximo login o cuando sea necesario.
        else:
            print("[auth] Refresh falló, el próximo request requerirá login")

    t = threading.Thread(target=_refresh, daemon=True)
    t.start()

# Uso
data = login('hackerman', 'secret123')
start_auto_refresh(data['expires_in'])
```

---

### Go — con context y goroutine

```go
func (c *Client) StartAutoRefresh(expiresIn int) {
    delay := time.Duration(expiresIn-60) * time.Second
    if delay < 5*time.Second {
        delay = 5 * time.Second
    }
    go func() {
        time.Sleep(delay)
        if err := c.doRefresh(); err != nil {
            fmt.Println("[auth] Auto-refresh falló:", err)
        } else {
            fmt.Println("[auth] Token renovado")
        }
    }()
}

// Uso después del login:
c.StartAutoRefresh(auth.ExpiresIn)
```

---

## 5. Logout

### JavaScript

```javascript
async function logout() {
    const refreshToken = localStorage.getItem('refresh_token');
    if (refreshToken) {
        // Revocar el token en el servidor (silenciar errores, siempre redirigir)
        await api.post('/auth/logout', { refresh_token: refreshToken }).catch(() => {});
    }
    clearTokens();
    window.location.href = '/login';
}
```

---

### Python

```python
def logout():
    if client.refresh_token:
        try:
            client.post('/auth/logout', {'refresh_token': client.refresh_token})
        except Exception:
            pass  # el servidor siempre devuelve 204, el error es de red
    client.access_token  = None
    client.refresh_token = None
    print("Sesión cerrada")
```

---

### Go

```go
func (c *Client) Logout() error {
    c.mu.Lock()
    rt := c.refreshToken
    c.mu.Unlock()

    if rt != "" {
        c.do("POST", "/auth/logout", map[string]string{"refresh_token": rt}, false)
    }
    c.SetTokens("", "")
    return nil
}
```

---

## 6. Módulo: Users

### Obtener el propio perfil by ID

```javascript
// JS
const me = await api.get(`/users/${userId}`);
console.log(me.username, me.points);
```

```python
# Python
me = client.get(f'/users/{user_id}')
print(me['username'], me['points'])
```

```go
// Go
res, _ := c.do("GET", "/users/"+userID, nil, true)
var me map[string]any
client.Decode(res, &me)
```

---

### Buscar usuarios por username

```javascript
// JS
const results = await api.get('/users/search?q=hack');
results.forEach(u => console.log(u.username));
```

```python
# Python
results = client.get('/users/search?q=hack')
for u in results:
    print(u['username'])
```

---

### Actualizar perfil

```javascript
// JS
const updated = await api.post(`/users/${userId}/update`, {
    username: 'nuevo_nombre',
    // email es opcional, omitir si no se cambia
});
```

```python
# Python
updated = client.post(f'/users/{user_id}/update', {
    'username': 'nuevo_nombre',
})
```

---

### Cambiar contraseña

```javascript
// JS
await api.post(`/users/${userId}/password`, {
    current_password: 'actual123',
    new_password:     'nueva456',
});
```

```python
# Python
client.post(f'/users/{user_id}/password', {
    'current_password': 'actual123',
    'new_password':     'nueva456',
})
```

---

## 7. Módulo: Events

### Listar eventos

```javascript
// JS — primera página
const { data, meta } = await api.get('/events/?page=1&size=20');
data.forEach(e => console.log(e.name, e.status));
```

```python
# Python
resp = client.get('/events/?page=1&size=20')
for event in resp['data']:
    print(event['name'], event['status'])
```

```go
// Go
res, _ := c.do("GET", "/events/?page=1&size=20", nil, true)
var resp struct {
    Data []map[string]any `json:"data"`
    Meta map[string]any   `json:"meta"`
}
client.Decode(res, &resp)
for _, e := range resp.Data {
    fmt.Println(e["name"], e["status"])
}
```

---

### Obtener detalle de un evento

```javascript
const event = await api.get(`/events/${eventId}`);
console.log(`${event.name} — ${event.status}`);
// Campos: id, name, description, max_team_size, status, starts_at, ends_at
```

---

### Crear evento _(admin)_

```javascript
// JS
const event = await api.post('/events/', {
    name:          'GoLabs CTF 2025',
    description:   'Competencia anual de seguridad.',
    max_team_size: 4,
    starts_at:     '2025-03-01T00:00:00Z',
    ends_at:       '2025-03-03T23:59:00Z',
});
console.log('Evento creado en estado:', event.status); // "draft"
```

```python
# Python
event = client.post('/events/', {
    'name':          'GoLabs CTF 2025',
    'description':   'Competencia anual de seguridad.',
    'max_team_size': 4,
    'starts_at':     '2025-03-01T00:00:00Z',
    'ends_at':       '2025-03-03T23:59:00Z',
})
print('Evento:', event['id'])
```

---

### Ciclo de vida del evento _(admin)_

```javascript
// JS — las 3 transiciones de estado
const eventId = 'uuid-del-evento';

await api.post(`/events/${eventId}/open`);    // draft → open
await api.post(`/events/${eventId}/start`);   // open  → running
await api.post(`/events/${eventId}/finish`);  // running → finished
```

```python
# Python
for action in ['open', 'start', 'finish']:
    client.post(f'/events/{event_id}/{action}')
    print(f'Evento → {action}')
```

---

## 8. Módulo: Teams

### Crear un equipo

```javascript
// JS
const { team, join_secret } = await api.post(`/events/${eventId}/teams`, {
    name: 'Byte Busters',
});
// IMPORTANTE: guardar join_secret para compartir con el equipo
// Solo se muestra esta única vez
console.log('Secreto de invitación:', join_secret);
sessionStorage.setItem('join_secret', join_secret);
```

```python
# Python
result = client.post(f'/events/{event_id}/teams', {'name': 'Byte Busters'})
join_secret = result['join_secret']
print('Guardar este secreto:', join_secret)
```

---

### Unirse a un equipo

```javascript
// JS
const team = await api.post(`/events/${eventId}/teams/join`, {
    join_secret: 'AbCd3FgHiJ',
});
console.log('Equipo unido:', team.name);
```

```python
# Python
team = client.post(f'/events/{event_id}/teams/join', {
    'join_secret': 'AbCd3FgHiJ',
})
print('Unido a:', team['name'])
```

---

### Ver miembros del equipo

```javascript
// JS
const members = await api.get(`/events/${eventId}/teams/${teamId}/members`);
members.forEach(m => console.log(m.username, m.user_id));
```

---

### Listar equipos y leaderboard

```javascript
// JS
const teams       = await api.get(`/events/${eventId}/teams`);
const leaderboard = await api.get(`/events/${eventId}/leaderboard`);

leaderboard.forEach((t, i) => {
    console.log(`#${i + 1} ${t.name} — ${t.score} pts`);
});
```

```python
# Python
leaderboard = client.get(f'/events/{event_id}/leaderboard')
for i, t in enumerate(leaderboard, 1):
    print(f"#{i} {t['name']} — {t['score']} pts")
```

---

### Rotar join secret

```javascript
// JS — solo el owner puede ejecutar esto
const { join_secret } = await api.post(
    `/events/${eventId}/teams/${teamId}/rotate-secret`
);
console.log('Nuevo secreto:', join_secret);
```

---

### Abandonar equipo

```javascript
// JS
await api.post(`/events/${eventId}/teams/${teamId}/leave`);
console.log('Saliste del equipo');
```

---

## 9. Módulo: Challenges

### Listar challenges

```javascript
// JS — listar todos los challenges del evento
const challenges = await api.get(`/events/${eventId}/challenges`);

// Filtrar por categoría
const webChallenges = await api.get(`/events/${eventId}/challenges?category=web`);

// Filtrar por dificultad
const easyChallenges = await api.get(`/events/${eventId}/challenges?difficulty=easy`);

challenges.forEach(c => {
    console.log(`[${c.category}] ${c.title} — ${c.points} pts — Resueltos: ${c.solve_count}`);
});
```

```python
# Python
challenges = client.get(f'/events/{event_id}/challenges')
for c in challenges:
    blood = f" 🩸 first blood: {c['first_blood_team_id']}" if c.get('first_blood_team_id') else ''
    print(f"[{c['category']}] {c['title']} — {c['points']} pts{blood}")
```

```go
// Go
res, _ := c.do("GET", "/events/"+eventID+"/challenges", nil, true)
var challenges []map[string]any
client.Decode(res, &challenges)
for _, ch := range challenges {
    fmt.Printf("[%s] %s — %.0f pts\n", ch["category"], ch["title"], ch["points"])
}
```

---

### Ver detalle de un challenge

```javascript
const challenge = await api.get(`/events/${eventId}/challenges/${challengeId}`);
console.log(challenge.description);
// Campos: id, event_id, title, description, category, points, difficulty,
//         visible, solve_count, first_blood_team_id, created_at, updated_at
```

---

### Enviar una flag

```javascript
// JS
async function submitFlag(eventId, challengeId, flag) {
    try {
        const result = await api.post(
            `/events/${eventId}/challenges/${challengeId}/submit`,
            { flag }
        );

        if (result.correct) {
            alert(`¡Correcto! +${result.points} puntos`);
        } else {
            alert('Flag incorrecta, intentá de nuevo');
        }
        return result.correct;
    } catch (err) {
        if (err.status === 409) {
            alert('Tu equipo ya resolvió este challenge');
        } else {
            alert('Error: ' + err.message);
        }
        return false;
    }
}

// Uso
await submitFlag(eventId, challengeId, 'CTF{mi_flag_aqui}');
```

```python
# Python
def submit_flag(event_id, challenge_id, flag):
    try:
        result = client.post(
            f'/events/{event_id}/challenges/{challenge_id}/submit',
            {'flag': flag}
        )
        if result['correct']:
            print(f"¡Correcto! +{result['points']} puntos")
        else:
            print('Flag incorrecta')
        return result['correct']
    except requests.HTTPError as e:
        if e.response.status_code == 409:
            print('Tu equipo ya resolvió este challenge')
        else:
            print('Error:', e.response.json().get('error'))
        return False

submit_flag(event_id, challenge_id, 'CTF{mi_flag_aqui}')
```

```go
// Go
res, _ := c.do("POST", "/events/"+eventID+"/challenges/"+chalID+"/submit",
    map[string]string{"flag": "CTF{mi_flag_aqui}"}, true)

var result struct {
    Correct bool `json:"correct"`
    Points  int  `json:"points"`
}
client.Decode(res, &result)
if result.Correct {
    fmt.Println("¡Correcto! Puntos:", result.Points)
} else {
    fmt.Println("Flag incorrecta")
}
```

---

### Crear y configurar un challenge _(admin)_

```javascript
// 1. Crear el challenge (queda oculto por defecto)
const challenge = await api.post(`/events/${eventId}/challenges`, {
    title:       'SQL Injection 101',
    description: 'Encontrá la flag mediante SQL injection en el formulario de login.',
    category:    'web',
    points:      200,
    difficulty:  'easy',
});
console.log('Challenge creado:', challenge.id);

// 2. Configurar la flag (sha-256 del texto plano se almacena internamente)
await api.post(`/events/${eventId}/challenges/${challenge.id}/flag`, {
    flag: 'CTF{sql_1nj3ct10n_ftw}',
});
console.log('Flag configurada');

// 3. Publicar para que los participantes lo vean
const published = await api.post(`/events/${eventId}/challenges/${challenge.id}/publish`);
console.log('Challenge visible:', published.visible); // true
```

```python
# Python — flujo completo de admin para crear challenge
challenge = client.post(f'/events/{event_id}/challenges', {
    'title':       'SQL Injection 101',
    'description': 'Encontrá la flag mediante SQL injection.',
    'category':    'web',
    'points':      200,
    'difficulty':  'easy',
})
challenge_id = challenge['id']

client.post(f'/events/{event_id}/challenges/{challenge_id}/flag', {
    'flag': 'CTF{sql_1nj3ct10n_ftw}'
})

published = client.post(f'/events/{event_id}/challenges/{challenge_id}/publish')
print('Visible:', published['visible'])  # True
```

---

## 10. Flujo completo de participante

End-to-end desde cero hasta resolver un challenge:

### JavaScript

```javascript
import { api, setTokens, clearTokens } from './api.js';

async function flujoParticipante() {
    // 1. Registrarse
    await api.post('/auth/register', {
        username: 'hackerman',
        email:    'hack@example.com',
        password: 'secret123',
    });

    // 2. Login
    const auth = await api.post('/auth/login', {
        identifier: 'hackerman',
        password:   'secret123',
    });
    setTokens(auth.access_token, auth.refresh_token);

    // 3. Ver eventos disponibles
    const { data: events } = await api.get('/events/');
    const event = events.find(e => e.status === 'open');
    if (!event) return console.log('No hay eventos abiertos');

    // 4. Crear equipo
    const { team, join_secret } = await api.post(`/events/${event.id}/teams`, {
        name: 'Byte Busters',
    });
    console.log('Equipo creado. Invitar con secret:', join_secret);

    // 5. Esperar a que el evento empiece (status = 'running')...

    // 6. Ver challenges
    const challenges = await api.get(`/events/${event.id}/challenges`);
    const target = challenges.find(c => c.difficulty === 'easy');

    // 7. Enviar flag
    const result = await api.post(
        `/events/${event.id}/challenges/${target.id}/submit`,
        { flag: 'CTF{la_flag_encontrada}' }
    );
    console.log(result.correct ? `¡+${result.points} puntos!` : 'Incorrecta');

    // 8. Ver leaderboard
    const board = await api.get(`/events/${event.id}/leaderboard`);
    board.slice(0, 3).forEach((t, i) => console.log(`#${i+1} ${t.name} ${t.score}pts`));

    // 9. Logout
    await api.post('/auth/logout', { refresh_token: localStorage.getItem('refresh_token') });
    clearTokens();
}

flujoParticipante();
```

---

### Python

```python
from api_client import client

# 1. Registro
client.post('/auth/register', {
    'username': 'hackerman',
    'email':    'hack@example.com',
    'password': 'secret123',
})

# 2. Login
data = client.post('/auth/login', {'identifier': 'hackerman', 'password': 'secret123'})
client.access_token  = data['access_token']
client.refresh_token = data['refresh_token']

# 3. Buscar evento abierto
events = client.get('/events/')['data']
event = next(e for e in events if e['status'] == 'open')

# 4. Crear equipo
result   = client.post(f"/events/{event['id']}/teams", {'name': 'Byte Busters'})
team     = result['team']
secret   = result['join_secret']
print(f"Equipo '{team['name']}' creado. Secret: {secret}")

# 5. (El evento arranca... status cambia a 'running')

# 6. Ver challenges
challenges = client.get(f"/events/{event['id']}/challenges")
target = next(c for c in challenges if c['difficulty'] == 'easy')
print(f"Resolviendo: {target['title']} ({target['points']} pts)")

# 7. Enviar flag
res = client.post(
    f"/events/{event['id']}/challenges/{target['id']}/submit",
    {'flag': 'CTF{la_flag_encontrada}'}
)
print('Correcto:', res['correct'], '— Puntos:', res.get('points', 0))

# 8. Leaderboard
board = client.get(f"/events/{event['id']}/leaderboard")
for i, t in enumerate(board[:3], 1):
    print(f"#{i} {t['name']} — {t['score']} pts")

# 9. Logout
client.post('/auth/logout', {'refresh_token': client.refresh_token})
client.access_token = client.refresh_token = None
```

---

## 11. Flujo completo de administrador

### JavaScript

```javascript
import { api, setTokens } from './api.js';

async function flujoAdmin() {
    // 1. Login como admin
    const auth = await api.post('/auth/login', {
        identifier: 'admin',
        password:   'adminpassword',
    });
    setTokens(auth.access_token, auth.refresh_token);

    // 2. Crear evento en draft
    const event = await api.post('/events/', {
        name:          'GoLabs CTF 2025',
        description:   'Competencia anual de CTF.',
        max_team_size: 4,
        starts_at:     '2025-03-01T00:00:00Z',
        ends_at:       '2025-03-03T23:59:00Z',
    });
    const eid = event.id;
    console.log('Evento creado:', eid, '— Estado:', event.status); // draft

    // 3. Crear challenges
    const challengeData = [
        { title: 'Warmup',      category: 'misc',    points: 50,  difficulty: 'easy',   flag: 'CTF{warmup}' },
        { title: 'SQLi 101',    category: 'web',     points: 200, difficulty: 'easy',   flag: 'CTF{sql_ftw}' },
        { title: 'Buffer Vuln', category: 'pwn',     points: 500, difficulty: 'hard',   flag: 'CTF{pwn_it}' },
    ];

    for (const ch of challengeData) {
        const created = await api.post(`/events/${eid}/challenges`, {
            title:      ch.title,
            category:   ch.category,
            points:     ch.points,
            difficulty: ch.difficulty,
            description: `Descripción de ${ch.title}`,
        });
        await api.post(`/events/${eid}/challenges/${created.id}/flag`, { flag: ch.flag });
        await api.post(`/events/${eid}/challenges/${created.id}/publish`);
        console.log(`Challenge listo: ${ch.title}`);
    }

    // 4. Abrir inscripciones
    await api.post(`/events/${eid}/open`);
    console.log('Evento abierto para inscripciones');

    // ... equipos se registran ...

    // 5. Iniciar competencia
    await api.post(`/events/${eid}/start`);
    console.log('Competencia en curso');

    // ... el CTF corre ...

    // 6. Finalizar
    await api.post(`/events/${eid}/finish`);
    console.log('Competencia finalizada');

    // 7. Ver resultados finales
    const leaderboard = await api.get(`/events/${eid}/leaderboard`);
    console.log('\n=== Resultados finales ===');
    leaderboard.slice(0, 5).forEach((t, i) => {
        console.log(`#${i+1} ${t.name.padEnd(20)} ${t.score} pts`);
    });
}

flujoAdmin();
```

---

### Python

```python
from api_client import client

# 1. Login admin
data = client.post('/auth/login', {'identifier': 'admin', 'password': 'adminpassword'})
client.access_token  = data['access_token']
client.refresh_token = data['refresh_token']

# 2. Crear evento
event = client.post('/events/', {
    'name':          'GoLabs CTF 2025',
    'description':   'Competencia anual de CTF.',
    'max_team_size': 4,
    'starts_at':     '2025-03-01T00:00:00Z',
    'ends_at':       '2025-03-03T23:59:00Z',
})
eid = event['id']
print(f"Evento creado: {eid}")

# 3. Crear, configurar y publicar challenges
challenges_config = [
    ('Warmup',      'misc', 50,  'easy', 'CTF{warmup}'),
    ('SQLi 101',    'web',  200, 'easy', 'CTF{sql_ftw}'),
    ('Buffer Vuln', 'pwn',  500, 'hard', 'CTF{pwn_it}'),
]

for title, category, points, difficulty, flag in challenges_config:
    ch = client.post(f'/events/{eid}/challenges', {
        'title':       title,
        'category':    category,
        'points':      points,
        'difficulty':  difficulty,
        'description': f'Descripción de {title}',
    })
    client.post(f"/events/{eid}/challenges/{ch['id']}/flag", {'flag': flag})
    client.post(f"/events/{eid}/challenges/{ch['id']}/publish")
    print(f"Challenge listo: {title}")

# 4-6. Ciclo de vida del evento
for action, label in [('open', 'abierto'), ('start', 'iniciado'), ('finish', 'finalizado')]:
    input(f"\nPresioná Enter para marcar el evento como {label}...")
    client.post(f'/events/{eid}/{action}')
    print(f'Evento {label}')

# 7. Resultados finales
board = client.get(f'/events/{eid}/leaderboard')
print('\n=== Resultados finales ===')
for i, t in enumerate(board[:5], 1):
    print(f"#{i} {t['name']:<20} {t['score']} pts")
```
