# tools/api-tester

Herramienta web de pruebas manuales para **golabs-api**.

Interfaz Flask que corre localmente y permite probar todos los endpoints de la API desde el navegador, con manejo automático de tokens JWT.

---

## Requisitos

- Python 3.10+
- pip

---

## Instalación y ejecución

```bash
cd tools/api-tester

# Crear entorno virtual e instalar dependencias
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Ejecutar
python app.py
# → http://localhost:5000
```

Por defecto apunta a `http://localhost:8080`. Para apuntar a otra URL:

```bash
API_BASE_URL=http://otra-url:8080 python app.py
```

---

## Uso

1. Abrir `http://localhost:5000` en el navegador.
2. Ingresar con un usuario válido de golabs-api (email + contraseña).
3. El token JWT se guarda en sesión automáticamente y se envía en cada request.
4. Navegar por las secciones para probar los distintos módulos.

---

## Variables de entorno

| Variable | Default | Descripción |
|---|---|---|
| `API_BASE_URL` | `http://localhost:8080` | URL base de la API a testear |

---

> [!NOTE]
> Esta herramienta es solo para desarrollo local. No exponer en entornos accesibles públicamente, ya que no tiene autenticación propia más allá de la sesión Flask.
