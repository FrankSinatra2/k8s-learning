package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type HomePageData struct {
	Username  string
	LastLogin string
}

type SessionData struct {
	Username  string
	SessionId string
	LastLogin string
}

const SessionIdCookieName = "k3sLearningSessionId"

var ctx = context.Background()
var rdb *redis.Client

var loginTemplate *template.Template
var homeTemplate *template.Template

func getEnv(envvar string, fallback string) string {
	if val, ok := os.LookupEnv(envvar); ok {
		return val
	}
	return fallback
}

func serializeSessionData(data *SessionData) string {
	return fmt.Sprintf("%s$%s$%s", data.Username, data.SessionId, data.LastLogin)
}

func deserializeSessionData(src string, data *SessionData) bool {
	parts := strings.Split(src, "$")

	if len(parts) != 3 {
		return false
	}

	data.Username = parts[0]
	data.SessionId = parts[1]
	data.LastLogin = parts[2]

	return true
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var username string
		if username = r.FormValue("username"); username == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Println("No Username")
			return
		}

		var sessionId string
		if uuid, err := uuid.NewUUID(); err == nil {
			sessionId = uuid.String()
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Failed to Create UUID")
			return
		}

		now := time.Now()

		sessionData := SessionData{
			SessionId: sessionId,
			Username:  username,
			LastLogin: time.Now().String(),
		}

		sessionKey := fmt.Sprintf("session:%s", sessionId)

		if err := rdb.Set(ctx, sessionKey, serializeSessionData(&sessionData), time.Hour).Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Failed to write to Redis")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    SessionIdCookieName,
			Value:   sessionId,
			Expires: now.Add(time.Hour),
			Path:    "/",
		})
		http.Redirect(w, r, "/home/", http.StatusFound)

		return
	} else if r.Method == http.MethodGet {
		loginTemplate.Execute(w, nil)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sessionId, err := r.Cookie(SessionIdCookieName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	sessionKey := fmt.Sprintf("session:%s", sessionId.Value)
	sessionValue, err := rdb.Get(ctx, sessionKey).Result()

	if err != nil {
		http.Redirect(w, r, "/public/error.html", http.StatusFound)
		return
	}

	sessionData := SessionData{}

	if ok := deserializeSessionData(sessionValue, &sessionData); !ok {
		http.Redirect(w, r, "/public/error.html", http.StatusFound)
		return
	}

	homePageData := HomePageData{
		Username:  sessionData.Username,
		LastLogin: sessionData.LastLogin,
	}

	homeTemplate.Execute(w, homePageData)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	// connect to redis server
	addr := getEnv("REDIS_HOST", "localhost:6379")
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	// load templates
	loginTemplate = template.Must(template.ParseFiles("templates/login.html"))
	homeTemplate = template.Must(template.ParseFiles("templates/home.html"))

	// start server
	http.HandleFunc("/login/", handleLogin)
	http.HandleFunc("/home/", handleHome)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// catch all
	http.HandleFunc("/", handleNotFound)

	http.ListenAndServe(":3000", nil)
}
