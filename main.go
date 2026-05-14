package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"
)

const pageHTML = `<!DOCTYPE html>
<html lang="th">
<head>
<meta charset="UTF-8">
<title>BlockGo Test 🛡️</title>
<style>
  * { box-sizing: border-box; }
  body {
    margin: 0; min-height: 100vh;
    font-family: 'Segoe UI', system-ui, sans-serif;
    background: linear-gradient(135deg, #1e3a8a 0%%, #7c3aed 100%%);
    display: flex; align-items: center; justify-content: center;
    color: #fff;
  }
  .card {
    background: rgba(255,255,255,0.1);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255,255,255,0.2);
    border-radius: 20px;
    padding: 40px;
    width: 90%%; max-width: 500px;
    box-shadow: 0 20px 60px rgba(0,0,0,0.3);
  }
  h1 { margin: 0 0 8px; font-size: 28px; }
  .sub { opacity: 0.8; margin-bottom: 24px; font-size: 14px; }
  .ip-badge {
    background: rgba(0,0,0,0.3);
    padding: 12px 16px;
    border-radius: 10px;
    margin-bottom: 20px;
    font-family: monospace;
    font-size: 14px;
  }
  .ip-badge b { color: #fbbf24; }
  input, textarea {
    width: 100%%;
    padding: 14px;
    border-radius: 10px;
    border: 1px solid rgba(255,255,255,0.3);
    background: rgba(0,0,0,0.2);
    color: #fff;
    font-size: 15px;
    margin-bottom: 14px;
    font-family: inherit;
  }
  input::placeholder, textarea::placeholder { color: rgba(255,255,255,0.5); }
  button {
    width: 100%%;
    padding: 14px;
    border: none;
    border-radius: 10px;
    background: #fbbf24;
    color: #1e3a8a;
    font-weight: 700;
    font-size: 16px;
    cursor: pointer;
    transition: transform 0.2s;
  }
  button:hover { transform: translateY(-2px); }
  .result {
    margin-top: 20px;
    padding: 16px;
    background: rgba(251,191,36,0.15);
    border-left: 4px solid #fbbf24;
    border-radius: 8px;
    word-break: break-word;
  }
  .footer {
    margin-top: 24px; text-align: center;
    font-size: 12px; opacity: 0.6;
  }
</style>
</head>
<body>
  <div class="card">
    <h1>🛡️ BlockGo Test</h1>
    <div class="sub">Web server พร้อมระบบ block IP สำหรับทดสอบ</div>
    <div class="ip-badge">Your IP: <b>%s</b></div>
    <form method="POST" action="/echo">
      <input name="name" placeholder="ชื่อของคุณ" required>
      <textarea name="message" rows="3" placeholder="พิมพ์ข้อความอะไรก็ได้..." required></textarea>
      <button type="submit">ส่งข้อความ ✨</button>
    </form>
    %s
    <div class="footer">เวลาเซิร์ฟเวอร์: %s</div>
  </div>
</body>
</html>`

func home(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, pageHTML, html.EscapeString(ip), "", now)
}

func echo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad form", 400)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	msg := strings.TrimSpace(r.FormValue("message"))
	ip := getClientIP(r)
	now := time.Now().Format("2006-01-02 15:04:05")

	// log ข้อความที่ user ส่ง
	log.Printf("📩 MESSAGE from %s — %s: %q", ip, name, msg)

	result := fmt.Sprintf(
		`<div class="result">👋 สวัสดี <b>%s</b>!<br>ข้อความของคุณ: <i>%s</i></div>`,
		html.EscapeString(name),
		html.EscapeString(msg),
	)
	fmt.Fprintf(w, pageHTML, html.EscapeString(ip), result, now)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/echo", echo)

	addr := ":8080"
	log.Printf("🚀 BlockGo server running on %s", addr)
	log.Printf("📋 Blocked IPs: %v", blockedIPs)

	if err := http.ListenAndServe(addr, Logger(mux)); err != nil {
		log.Fatal(err)
	}
}
