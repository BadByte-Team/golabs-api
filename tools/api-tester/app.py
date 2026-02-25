from flask import Flask, request, jsonify, render_template, session, redirect, url_for
import requests as req
import os

app = Flask(__name__)
app.secret_key = "golabs-api-tester-dev-key-2025"

BASE_URL = os.environ.get("API_BASE_URL", "http://localhost:8080")
V1 = "/api/v1"


# ── helpers ──────────────────────────────────────────────────────────────────

def _headers():
    h = {"Content-Type": "application/json"}
    if session.get("token"):
        h["Authorization"] = f"Bearer {session['token']}"
    return h


def _call(method, path, body=None, params=None):
    url = BASE_URL + path
    try:
        r = req.request(method, url, json=body, params=params, headers=_headers(), timeout=10)
        try:
            data = r.json()
        except Exception:
            data = {"raw": r.text or "(empty)"}
        return {"status": r.status_code, "body": data}
    except req.exceptions.ConnectionError:
        return {"status": 0, "body": {"error": f"No se pudo conectar a {BASE_URL}. ¿Está corriendo el servidor?"}}
    except Exception as e:
        return {"status": 0, "body": {"error": str(e)}}


# ── pages ─────────────────────────────────────────────────────────────────────

@app.route("/login")
def login_page():
    if session.get("token"):
        return redirect(url_for("index"))
    return render_template("login.html", base_url=BASE_URL)


@app.route("/logout")
def logout():
    session.clear()
    return redirect(url_for("login_page"))


@app.route("/")
def index():
    if not session.get("token"):
        return redirect(url_for("login_page"))
    return render_template("index.html",
                           base_url=BASE_URL,
                           token=session.get("token", ""),
                           username=session.get("username", ""))


# ── auth API ──────────────────────────────────────────────────────────────────

@app.route("/api/auth/login", methods=["POST"])
def api_login():
    body = request.json
    identifier = body.get("identifier") or body.get("email", "")
    result = _call("POST", f"{V1}/auth/login", {"identifier": identifier, "password": body.get("password", "")})
    if result["status"] == 200:
        body_data = result.get("body") or {}
        access_token = body_data.get("access_token") or body_data.get("token", "")
        if access_token:
            session["token"] = access_token
            session["username"] = identifier
        if body_data.get("refresh_token"):
            session["refresh_token"] = body_data["refresh_token"]
    return jsonify(result)


@app.route("/api/auth/refresh", methods=["POST"])
def api_refresh():
    raw_rt = (request.json or {}).get("refresh_token") or session.get("refresh_token", "")
    result = _call("POST", f"{V1}/auth/refresh", {"refresh_token": raw_rt})
    if result["status"] == 200:
        body_data = result.get("body") or {}
        if body_data.get("access_token"):
            session["token"] = body_data["access_token"]
        if body_data.get("refresh_token"):
            session["refresh_token"] = body_data["refresh_token"]
    return jsonify(result)


@app.route("/api/auth/logout", methods=["POST"])
def api_logout():
    raw_rt = (request.json or {}).get("refresh_token") or session.pop("refresh_token", "")
    result = _call("POST", f"{V1}/auth/logout", {"refresh_token": raw_rt})
    session.pop("token", None)
    return jsonify(result)


@app.route("/api/auth/register", methods=["POST"])
def api_register():
    return jsonify(_call("POST", f"{V1}/auth/register", request.json))


# ── health ────────────────────────────────────────────────────────────────────

@app.route("/api/health")
def api_health():
    return jsonify(_call("GET", "/health"))


@app.route("/api/health/live")
def api_health_live():
    return jsonify(_call("GET", "/healthz/live"))


@app.route("/api/health/ready")
def api_health_ready():
    return jsonify(_call("GET", "/healthz/ready"))


# ── users ─────────────────────────────────────────────────────────────────────

@app.route("/api/users", methods=["GET"])
def api_user_list():
    page = request.args.get("page", "1")
    size = request.args.get("size", "20")
    return jsonify(_call("GET", f"{V1}/users/", params={"page": page, "size": size}))


@app.route("/api/users", methods=["POST"])
def api_user_create():
    return jsonify(_call("POST", f"{V1}/users/", request.json))


@app.route("/api/users/search", methods=["GET"])
def api_user_search():
    q = request.args.get("q", "")
    return jsonify(_call("GET", f"{V1}/users/search", params={"q": q}))


@app.route("/api/users/by-username/<username>", methods=["GET"])
def api_user_by_username(username):
    return jsonify(_call("GET", f"{V1}/users/by-username/{username}"))


@app.route("/api/users/<uid>", methods=["GET"])
def api_user_get(uid):
    return jsonify(_call("GET", f"{V1}/users/{uid}"))


@app.route("/api/users/<uid>/update", methods=["POST"])
def api_user_update(uid):
    return jsonify(_call("POST", f"{V1}/users/{uid}/update", request.json))


@app.route("/api/users/<uid>/password", methods=["POST"])
def api_user_password(uid):
    return jsonify(_call("POST", f"{V1}/users/{uid}/password", request.json))


@app.route("/api/admin/users/<uid>/role", methods=["POST"])
def api_user_role(uid):
    return jsonify(_call("POST", f"{V1}/admin/users/{uid}/role", request.json))


@app.route("/api/admin/users/<uid>/points", methods=["POST"])
def api_user_points(uid):
    return jsonify(_call("POST", f"{V1}/admin/users/{uid}/points", request.json))


@app.route("/api/admin/users/<uid>/ban", methods=["POST"])
def api_user_ban(uid):
    return jsonify(_call("POST", f"{V1}/admin/users/{uid}/ban"))


@app.route("/api/admin/users/<uid>/unban", methods=["POST"])
def api_user_unban(uid):
    return jsonify(_call("POST", f"{V1}/admin/users/{uid}/unban"))


# ── events ────────────────────────────────────────────────────────────────────

@app.route("/api/events", methods=["GET"])
def api_events_list():
    page = request.args.get("page", "1")
    size = request.args.get("size", "20")
    return jsonify(_call("GET", f"{V1}/events", params={"page": page, "size": size}))


@app.route("/api/events", methods=["POST"])
def api_event_create():
    return jsonify(_call("POST", f"{V1}/events", request.json))


@app.route("/api/events/<eid>", methods=["GET"])
def api_event_get(eid):
    return jsonify(_call("GET", f"{V1}/events/{eid}"))


@app.route("/api/events/<eid>/open", methods=["POST"])
def api_event_open(eid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/open"))


@app.route("/api/events/<eid>/start", methods=["POST"])
def api_event_start(eid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/start"))


@app.route("/api/events/<eid>/finish", methods=["POST"])
def api_event_finish(eid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/finish"))


# ── teams ─────────────────────────────────────────────────────────────────────

@app.route("/api/events/<eid>/teams", methods=["GET"])
def api_team_list(eid):
    return jsonify(_call("GET", f"{V1}/events/{eid}/teams/"))


@app.route("/api/events/<eid>/teams", methods=["POST"])
def api_team_create(eid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/teams/", request.json))


@app.route("/api/events/<eid>/teams/join", methods=["POST"])
def api_team_join(eid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/teams/join", request.json))


@app.route("/api/events/<eid>/teams/<tid>/members", methods=["GET"])
def api_team_members(eid, tid):
    return jsonify(_call("GET", f"{V1}/events/{eid}/teams/{tid}/members"))


@app.route("/api/events/<eid>/teams/<tid>/leave", methods=["POST"])
def api_team_leave(eid, tid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/teams/{tid}/leave"))


@app.route("/api/events/<eid>/teams/<tid>/rotate-secret", methods=["POST"])
def api_team_rotate(eid, tid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/teams/{tid}/rotate-secret"))


# ── leaderboard ───────────────────────────────────────────────────────────────

@app.route("/api/events/<eid>/leaderboard", methods=["GET"])
def api_leaderboard(eid):
    return jsonify(_call("GET", f"{V1}/events/{eid}/leaderboard"))


# ── challenges ────────────────────────────────────────────────────────────────

@app.route("/api/events/<eid>/challenges", methods=["GET"])
def api_challenges_list(eid):
    params = {}
    if request.args.get("category"):
        params["category"] = request.args.get("category")
    if request.args.get("difficulty"):
        params["difficulty"] = request.args.get("difficulty")
    return jsonify(_call("GET", f"{V1}/events/{eid}/challenges", params=params or None))


@app.route("/api/events/<eid>/challenges", methods=["POST"])
def api_challenge_create(eid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/challenges", request.json))


@app.route("/api/events/<eid>/challenges/<cid>", methods=["GET"])
def api_challenge_get(eid, cid):
    return jsonify(_call("GET", f"{V1}/events/{eid}/challenges/{cid}"))


@app.route("/api/events/<eid>/challenges/<cid>", methods=["PUT"])
def api_challenge_update(eid, cid):
    return jsonify(_call("PUT", f"{V1}/events/{eid}/challenges/{cid}", request.json))


@app.route("/api/events/<eid>/challenges/<cid>/publish", methods=["POST"])
def api_challenge_publish(eid, cid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/challenges/{cid}/publish"))


@app.route("/api/events/<eid>/challenges/<cid>/unpublish", methods=["POST"])
def api_challenge_unpublish(eid, cid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/challenges/{cid}/unpublish"))


@app.route("/api/events/<eid>/challenges/<cid>/flag", methods=["POST"])
def api_challenge_flag(eid, cid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/challenges/{cid}/flag", request.json))


@app.route("/api/events/<eid>/challenges/<cid>/submit", methods=["POST"])
def api_challenge_submit(eid, cid):
    return jsonify(_call("POST", f"{V1}/events/{eid}/challenges/{cid}/submit", request.json))


if __name__ == "__main__":
    app.run(debug=True, port=5000)
