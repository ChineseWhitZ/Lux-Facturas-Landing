package httpserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"lux-facturas/backend/internal/leads"
)

//go:embed docs/openapi.yaml
var openAPISpec []byte

var dashboardTemplate = template.Must(template.New("dashboard").Funcs(template.FuncMap{
	"planLabel":         planLabel,
	"customerTypeLabel": customerTypeLabel,
	"statusLabel":       statusLabel,
	"whatsappHref":      whatsappHref,
}).Parse(dashboardHTML))

var loginTemplate = template.Must(template.New("login").Parse(loginHTML))

const adminSessionCookie = "lux_admin_session"

type dashboardView struct {
	Leads         []leads.Lead
	Total         int
	New           int
	Contacted     int
	Quoted        int
	Closed        int
	Discarded     int
	StatusOptions []leads.Status
}

func planLabel(plan leads.Plan) string {
	switch plan {
	case leads.PlanBasic:
		return "Kit Basico"
	case leads.PlanStandard:
		return "Kit Estandar"
	case leads.PlanComplete:
		return "Kit Completo"
	default:
		return string(plan)
	}
}

func customerTypeLabel(customerType leads.CustomerType) string {
	switch customerType {
	case leads.CustomerHasMac:
		return "Ya tiene Mac"
	case leads.CustomerNeedsKit:
		return "Necesita kit completo"
	case leads.CustomerWantsQuote:
		return "Quiere cotizar"
	default:
		return string(customerType)
	}
}

func statusLabel(status leads.Status) string {
	switch status {
	case leads.StatusNew:
		return "Nuevo"
	case leads.StatusContacted:
		return "Contactado"
	case leads.StatusQuoted:
		return "Cotizado"
	case leads.StatusClosed:
		return "Cerrado"
	case leads.StatusDiscarded:
		return "Descartado"
	default:
		return string(status)
	}
}

func whatsappHref(phone string) string {
	clean := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "", "+", "").Replace(phone)
	if clean == "" {
		return "#"
	}
	if strings.HasPrefix(clean, "51") {
		return "https://wa.me/" + clean
	}
	return "https://wa.me/51" + clean
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if s.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	s.renderLogin(w, "")
}

func (s *Server) loginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderLogin(w, "No pudimos leer el formulario.")
		return
	}

	user := r.FormValue("user")
	password := r.FormValue("password")
	if !s.validAdminCredentials(user, password) {
		s.renderLogin(w, "Usuario o clave incorrectos.")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     adminSessionCookie,
		Value:    s.signedSessionValue(user),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
		MaxAge:   int((12 * time.Hour).Seconds()),
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     adminSessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}

	items := s.leadRepository.List()
	view := dashboardView{
		Leads: items,
		Total: len(items),
		StatusOptions: []leads.Status{
			leads.StatusNew,
			leads.StatusContacted,
			leads.StatusQuoted,
			leads.StatusClosed,
			leads.StatusDiscarded,
		},
	}

	for _, item := range items {
		switch item.Status {
		case leads.StatusNew:
			view.New++
		case leads.StatusContacted:
			view.Contacted++
		case leads.StatusQuoted:
			view.Quoted++
		case leads.StatusClosed:
			view.Closed++
		case leads.StatusDiscarded:
			view.Discarded++
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := dashboardTemplate.Execute(w, view); err != nil {
		s.logger.Error("could not render dashboard", "error", err)
	}
}

func (s *Server) updateLeadStatusFromDashboard(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	status := leads.Status(r.FormValue("status"))
	if _, err := s.leadRepository.UpdateStatus(r.PathValue("id"), status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) renderLogin(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := loginTemplate.Execute(w, map[string]string{"Message": message}); err != nil {
		s.logger.Error("could not render login", "error", err)
	}
}

func (s *Server) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	if s.isAuthenticated(r) {
		return true
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return false
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(adminSessionCookie)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(s.signedSessionValue(s.config.AdminUser))) == 1
}

func (s *Server) validAdminCredentials(user string, password string) bool {
	userOK := subtle.ConstantTimeCompare([]byte(user), []byte(s.config.AdminUser)) == 1
	passOK := subtle.ConstantTimeCompare([]byte(password), []byte(s.config.AdminPass)) == 1
	return userOK && passOK
}

func (s *Server) signedSessionValue(user string) string {
	mac := hmac.New(sha256.New, []byte(s.config.SessionKey))
	mac.Write([]byte(user))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return user + "." + signature
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "lux-facturas-api",
		"status":  "ready",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) createLead(w http.ResponseWriter, r *http.Request) {
	var input leads.CreateLeadInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	lead, err := s.leadRepository.Create(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, lead)
}

func (s *Server) listLeads(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": s.leadRepository.List(),
	})
}

func (s *Server) getLead(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}

	lead, err := s.leadRepository.FindByID(r.PathValue("id"))
	if errors.Is(err, leads.ErrNotFound) {
		writeError(w, http.StatusNotFound, "lead not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not get lead")
		return
	}

	writeJSON(w, http.StatusOK, lead)
}

func (s *Server) updateLeadStatus(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}

	var input leads.UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	lead, err := s.leadRepository.UpdateStatus(r.PathValue("id"), input.Status)
	if errors.Is(err, leads.ErrNotFound) {
		writeError(w, http.StatusNotFound, "lead not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, lead)
}

func (s *Server) openapi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(openAPISpec)
}

func (s *Server) swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(swaggerHTML))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

const swaggerHTML = `<!doctype html>
<html lang="es">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Lux Facturas API - Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: "#swagger-ui"
      });
    </script>
  </body>
</html>`

const loginHTML = `<!doctype html>
<html lang="es">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Lux Facturas - Admin</title>
    <style>
      :root {
        color-scheme: dark;
        --bg: #070910;
        --surface: #10131d;
        --text: #f4f7fb;
        --muted: #9aa4b6;
        --border: rgba(255, 255, 255, .12);
        --cyan: #19c7f3;
        --red: #ff6b57;
      }

      * {
        box-sizing: border-box;
      }

      body {
        align-items: center;
        background:
          linear-gradient(140deg, rgba(25, 199, 243, .12), transparent 34%),
          radial-gradient(circle at 82% 18%, rgba(61, 220, 132, .10), transparent 26%),
          var(--bg);
        color: var(--text);
        display: grid;
        font-family: "Avenir Next", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        margin: 0;
        min-height: 100vh;
        padding: 24px;
      }

      .shell {
        margin: 0 auto;
        max-width: 430px;
        width: 100%;
      }

      .brand {
        align-items: center;
        display: flex;
        gap: 12px;
        margin-bottom: 24px;
      }

      .mark {
        background: var(--cyan);
        border-radius: 8px;
        color: #071018;
        display: grid;
        font-weight: 900;
        height: 42px;
        place-items: center;
        width: 42px;
      }

      h1 {
        font-size: 22px;
        line-height: 1.2;
        margin: 0;
      }

      p {
        color: var(--muted);
        line-height: 1.7;
        margin: 8px 0 0;
      }

      form {
        background: rgba(16, 19, 29, .9);
        border: 1px solid var(--border);
        border-radius: 8px;
        display: grid;
        gap: 16px;
        padding: 22px;
      }

      label {
        display: grid;
        gap: 8px;
      }

      span {
        color: rgba(244, 247, 251, .76);
        font-size: 13px;
        font-weight: 800;
      }

      input,
      button {
        border-radius: 8px;
        font: inherit;
        min-height: 44px;
        padding: 10px 12px;
        width: 100%;
      }

      input {
        background: #070910;
        border: 1px solid var(--border);
        color: var(--text);
        outline: none;
      }

      input:focus {
        border-color: rgba(25, 199, 243, .72);
      }

      button {
        background: var(--cyan);
        border: 1px solid var(--cyan);
        color: #071018;
        cursor: pointer;
        font-weight: 900;
        transition: transform 160ms ease-out;
      }

      button:active {
        transform: scale(.97);
      }

      .error {
        background: rgba(255, 107, 87, .12);
        border: 1px solid rgba(255, 107, 87, .26);
        border-radius: 8px;
        color: #ffd0ca;
        font-size: 14px;
        font-weight: 800;
        padding: 11px 12px;
      }
    </style>
  </head>
  <body>
    <main class="shell">
      <div class="brand">
        <span class="mark">L</span>
        <div>
          <h1>Lux Facturas Admin</h1>
          <p>Ingresa para revisar leads, compras y seguimiento comercial.</p>
        </div>
      </div>

      <form method="post" action="/login">
        {{if .Message}}<div class="error">{{.Message}}</div>{{end}}
        <label>
          <span>Usuario</span>
          <input name="user" autocomplete="username" required>
        </label>
        <label>
          <span>Clave</span>
          <input name="password" type="password" autocomplete="current-password" required>
        </label>
        <button type="submit">Entrar al dashboard</button>
      </form>
    </main>
  </body>
</html>`

const dashboardHTML = `<!doctype html>
<html lang="es">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Lux Facturas - Leads</title>
    <style>
      :root {
        color-scheme: dark;
        --bg: #070910;
        --surface: #10131d;
        --surface-2: #161a26;
        --surface-3: #1d2331;
        --text: #f4f7fb;
        --muted: #9aa4b6;
        --border: rgba(255, 255, 255, .11);
        --cyan: #19c7f3;
        --mint: #3ddc84;
        --gold: #f3b248;
        --red: #ff6b57;
        --shadow: 0 24px 70px rgba(0, 0, 0, .28);
      }

      * {
        box-sizing: border-box;
      }

      html {
        background: var(--bg);
      }

      body {
        margin: 0;
        min-height: 100vh;
        background:
          linear-gradient(140deg, rgba(25, 199, 243, .12), transparent 34%),
          radial-gradient(circle at 88% 12%, rgba(61, 220, 132, .10), transparent 28%),
          var(--bg);
        color: var(--text);
        font-family: "Avenir Next", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      }

      a {
        color: inherit;
      }

      .topbar,
      main {
        width: min(1240px, calc(100% - 32px));
        margin: 0 auto;
      }

      header {
        border-bottom: 1px solid var(--border);
        background: rgba(7, 9, 16, .88);
        backdrop-filter: blur(18px);
        position: sticky;
        top: 0;
        z-index: 10;
      }

      .topbar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
        padding: 18px 0;
      }

      .brand {
        display: flex;
        align-items: center;
        gap: 12px;
      }

      .mark {
        display: grid;
        width: 40px;
        height: 40px;
        place-items: center;
        border-radius: 8px;
        background: var(--cyan);
        color: #071018;
        font-weight: 900;
      }

      h1 {
        margin: 0;
        font-size: 18px;
        line-height: 1.2;
      }

      .subtitle {
        margin: 3px 0 0;
        color: var(--muted);
        font-size: 14px;
      }

      .nav {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
      }

      .nav a,
      .nav button {
        border: 1px solid var(--border);
        border-radius: 8px;
        color: var(--text);
        font-size: 14px;
        font-weight: 700;
        padding: 10px 12px;
        text-decoration: none;
        transition: border-color 160ms ease-out, transform 160ms ease-out;
      }

      .nav form {
        display: inline-flex;
      }

      .nav button {
        background: transparent;
        cursor: pointer;
        font-family: inherit;
      }

      .nav a:active,
      button:active {
        transform: scale(.97);
      }

      .nav a:hover,
      .nav button:hover {
        border-color: rgba(25, 199, 243, .55);
      }

      main {
        padding: 34px 0 48px;
      }

      .hero {
        display: grid;
        gap: 28px;
        grid-template-columns: 1fr auto;
        margin-bottom: 26px;
      }

      .eyebrow {
        border-top: 1px solid var(--border);
        border-bottom: 1px solid var(--border);
        color: rgba(244, 247, 251, .72);
        display: inline-flex;
        font-size: 13px;
        font-weight: 800;
        line-height: 1.5;
        margin: 0 0 20px;
        padding: 9px 0;
      }

      .hero h2 {
        font-size: clamp(34px, 5vw, 64px);
        line-height: .98;
        margin: 0;
        max-width: 760px;
      }

      .hero p {
        color: var(--muted);
        font-size: 17px;
        line-height: 1.7;
        margin: 18px 0 0;
        max-width: 690px;
      }

      .ops-card {
        align-self: end;
        background: var(--surface);
        border: 1px solid var(--border);
        border-radius: 8px;
        min-width: 260px;
        padding: 18px;
        box-shadow: var(--shadow);
      }

      .ops-card span {
        color: var(--muted);
        display: block;
        font-size: 13px;
        font-weight: 800;
      }

      .ops-card strong {
        display: block;
        font-size: 44px;
        line-height: 1;
        margin-top: 10px;
      }

      .ops-card p {
        font-size: 14px;
        margin: 12px 0 0;
      }

      .stats {
        display: grid;
        grid-template-columns: repeat(6, minmax(130px, 1fr));
        gap: 12px;
        margin-bottom: 24px;
      }

      .stat {
        background: var(--surface);
        border: 1px solid var(--border);
        border-radius: 8px;
        padding: 16px;
        min-height: 94px;
      }

      .stat span {
        color: var(--muted);
        display: block;
        font-size: 12px;
        margin-bottom: 8px;
      }

      .stat strong {
        display: block;
        font-size: 26px;
        line-height: 1;
      }

      .stat.new strong { color: var(--cyan); }
      .stat.contacted strong { color: #75d7ff; }
      .stat.quoted strong { color: var(--gold); }
      .stat.closed strong { color: var(--mint); }
      .stat.discarded strong { color: var(--red); }

      .panel {
        background: rgba(16, 19, 29, .86);
        border: 1px solid var(--border);
        border-radius: 8px;
        box-shadow: var(--shadow);
      }

      .panel-head {
        align-items: center;
        border-bottom: 1px solid var(--border);
        display: flex;
        justify-content: space-between;
        gap: 12px;
        padding: 16px;
      }

      .panel-head h2 {
        font-size: 22px;
        margin: 0;
      }

      .panel-head p {
        color: var(--muted);
        font-size: 14px;
        margin: 4px 0 0;
      }

      .lead-list {
        display: grid;
        gap: 14px;
        padding: 16px;
      }

      .lead-card {
        background: var(--surface-2);
        border: 1px solid var(--border);
        border-radius: 8px;
        display: grid;
        gap: 18px;
        grid-template-columns: minmax(260px, 1.05fr) minmax(320px, 1.25fr) minmax(250px, .8fr);
        padding: 18px;
      }

      .lead-main {
        display: flex;
        flex-direction: column;
        min-width: 0;
      }

      .lead-main h3 {
        font-size: 20px;
        line-height: 1.25;
        margin: 0;
      }

      .company {
        color: var(--muted);
        font-size: 14px;
        line-height: 1.5;
        margin: 8px 0 0;
      }

      .lead-meta {
        display: grid;
        gap: 10px;
        grid-template-columns: repeat(2, minmax(0, 1fr));
      }

      .field span {
        color: var(--muted);
        display: block;
        font-size: 11px;
        font-weight: 700;
        text-transform: uppercase;
      }

      .field strong,
      .field a {
        color: var(--text);
        display: block;
        font-size: 14px;
        line-height: 1.5;
        margin-top: 3px;
        overflow-wrap: anywhere;
        text-decoration: none;
      }

      .field a:hover {
        color: var(--cyan);
      }

      .message {
        border-top: 1px solid var(--border);
        color: var(--muted);
        font-size: 14px;
        line-height: 1.6;
        margin-top: auto;
        padding-top: 14px;
      }

      .pill {
        border-radius: 999px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        font-weight: 700;
        line-height: 1;
        padding: 7px 10px;
      }

      .pill.nuevo { background: rgba(25, 199, 243, .14); color: var(--cyan); }
      .pill.contactado { background: rgba(117, 215, 255, .14); color: #75d7ff; }
      .pill.cotizado { background: rgba(243, 178, 72, .16); color: var(--gold); }
      .pill.cerrado { background: rgba(61, 220, 132, .14); color: var(--mint); }
      .pill.descartado { background: rgba(255, 107, 87, .14); color: var(--red); }

      .lead-actions {
        align-content: start;
        display: grid;
        gap: 12px;
      }

      form {
        display: grid;
        gap: 8px;
      }

      select,
      button {
        border: 1px solid var(--border);
        border-radius: 8px;
        font: inherit;
        min-height: 40px;
        padding: 9px 11px;
        width: 100%;
      }

      select {
        background: var(--surface-3);
        color: var(--text);
      }

      button {
        background: var(--cyan);
        border-color: var(--cyan);
        color: #071018;
        cursor: pointer;
        font-weight: 700;
        transition: transform 160ms ease-out, opacity 160ms ease-out;
      }

      .quick {
        border: 1px solid var(--border);
        border-radius: 8px;
        color: var(--text);
        display: inline-flex;
        font-size: 14px;
        font-weight: 800;
        justify-content: center;
        min-height: 40px;
        padding: 10px 12px;
        text-decoration: none;
      }

      .quick:hover {
        border-color: rgba(25, 199, 243, .55);
      }

      .empty {
        color: var(--muted);
        font-size: 15px;
        line-height: 1.7;
        padding: 28px;
        text-align: center;
      }

      @media (max-width: 1080px) {
        .hero,
        .lead-card {
          grid-template-columns: 1fr;
        }

        .ops-card {
          min-width: 0;
        }
      }

      @media (max-width: 760px) {
        .topbar {
          align-items: flex-start;
          flex-direction: column;
        }

        .stats {
          grid-template-columns: repeat(2, minmax(0, 1fr));
        }

        .lead-meta {
          grid-template-columns: 1fr;
        }

        .panel-head {
          align-items: flex-start;
          flex-direction: column;
        }
      }
    </style>
  </head>
  <body>
    <header>
      <div class="topbar">
        <div class="brand">
          <span class="mark">L</span>
          <div>
            <h1>Lux Facturas</h1>
            <p class="subtitle">Dashboard interno de leads y compras</p>
          </div>
        </div>
        <nav class="nav" aria-label="Accesos">
          <a href="/swagger">Swagger</a>
          <a href="/openapi.yaml">OpenAPI</a>
          <a href="/api/v1/leads">JSON</a>
          <form method="post" action="/logout">
            <button type="submit">Salir</button>
          </form>
        </nav>
      </div>
    </header>

    <main>
      <section class="hero">
        <div>
          <p class="eyebrow">Operacion comercial</p>
          <h2>Prospectos listos para cotizar, sin perder el contexto de compra.</h2>
          <p>Revisa quien llego desde la landing, que paquete pidio, si necesita hardware Apple y en que estado esta la conversacion.</p>
        </div>
        <aside class="ops-card" aria-label="Leads totales">
          <span>Total registrado</span>
          <strong>{{.Total}}</strong>
          <p>Usa el estado para llevar cada lead de nuevo a cerrado.</p>
        </aside>
      </section>

      <section class="stats" aria-label="Resumen">
        <div class="stat"><span>Total</span><strong>{{.Total}}</strong></div>
        <div class="stat new"><span>Nuevos</span><strong>{{.New}}</strong></div>
        <div class="stat contacted"><span>Contactados</span><strong>{{.Contacted}}</strong></div>
        <div class="stat quoted"><span>Cotizados</span><strong>{{.Quoted}}</strong></div>
        <div class="stat closed"><span>Cerrados</span><strong>{{.Closed}}</strong></div>
        <div class="stat discarded"><span>Descartados</span><strong>{{.Discarded}}</strong></div>
      </section>

      <section class="panel">
        <div class="panel-head">
          <div>
            <h2>Leads registrados</h2>
            <p>Datos capturados desde la pagina y seguimiento comercial interno.</p>
          </div>
          <span class="pill nuevo">{{.Total}} registros</span>
        </div>

        {{if .Leads}}
        <div class="lead-list">
          {{range .Leads}}
          <article class="lead-card">
            <div class="lead-main">
              <span class="pill {{.Status}}">{{statusLabel .Status}}</span>
              <h3>{{.CustomerName}}</h3>
              <p class="company">{{if .CompanyName}}{{.CompanyName}}{{else}}Empresa sin registrar{{end}} / {{.DocumentNumber}}</p>
              <p class="message">{{if .Message}}{{.Message}}{{else}}Sin mensaje adicional.{{end}}</p>
            </div>

            <div class="lead-meta">
              <div class="field">
                <span>Plan elegido</span>
                <strong>{{planLabel .SelectedPlan}}</strong>
              </div>
              <div class="field">
                <span>Tipo de cliente</span>
                <strong>{{customerTypeLabel .CustomerType}}</strong>
              </div>
              <div class="field">
                <span>Rubro</span>
                <strong>{{if .BusinessCategory}}{{.BusinessCategory}}{{else}}No indicado{{end}}</strong>
              </div>
              <div class="field">
                <span>Productos</span>
                <strong>{{.ApproxProductQuantity}} aprox.</strong>
              </div>
              <div class="field">
                <span>Ciudad</span>
                <strong>{{if .City}}{{.City}}{{else}}No indicada{{end}}</strong>
              </div>
              <div class="field">
                <span>Registro</span>
                <strong>{{.RegisteredAt.Format "2006-01-02 15:04"}}</strong>
              </div>
              <div class="field">
                <span>Telefono</span>
                <a href="{{whatsappHref .Phone}}" target="_blank" rel="noreferrer">{{.Phone}}</a>
              </div>
              <div class="field">
                <span>Email</span>
                <a href="mailto:{{.Email}}">{{.Email}}</a>
              </div>
            </div>

            <div class="lead-actions">
              <form method="post" action="/dashboard/leads/{{.ID}}/status">
                <select name="status" aria-label="Estado">
                  <option value="nuevo" {{if eq .Status "nuevo"}}selected{{end}}>Nuevo</option>
                  <option value="contactado" {{if eq .Status "contactado"}}selected{{end}}>Contactado</option>
                  <option value="cotizado" {{if eq .Status "cotizado"}}selected{{end}}>Cotizado</option>
                  <option value="cerrado" {{if eq .Status "cerrado"}}selected{{end}}>Cerrado</option>
                  <option value="descartado" {{if eq .Status "descartado"}}selected{{end}}>Descartado</option>
                </select>
                <button type="submit">Guardar estado</button>
              </form>
              <a class="quick" href="{{whatsappHref .Phone}}" target="_blank" rel="noreferrer">Abrir WhatsApp</a>
            </div>
          </article>
          {{end}}
        </div>
        {{else}}
        <p class="empty">Todavia no hay leads registrados. Cuando alguien complete el diagnostico de compra en la landing, aparecera aqui con plan, rubro, equipo Apple y datos de contacto.</p>
        {{end}}
      </section>
    </main>
  </body>
</html>`
