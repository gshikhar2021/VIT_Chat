package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

// Configuration stores config info of server
type Configuration struct {
	Address      string
	RedisURL     string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

// Config captures parsed input from config.json
var Config Configuration

// Mux contains all the HTTP handlers
var (
	Mux *mux.Router
)

// registerHandlers will register all HTTP handlers
func registerHandlers() *mux.Router {
	api := mux.NewRouter()
	// index
	api.HandleFunc("/", logConsole(index))
	// handle static assets by routing requests from /static/ => "public" directory
	staticDir := "/static/"
	api.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir(Config.Static))))

	//REST-API for chat room [JSON]
	api.Handle("/chats", errHandler(handlePost)).Methods(http.MethodPost)
	api.Handle("/chats/{titleOrID}", errHandler(authorize(handleRoom))).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	// List all rooms in [HTML]
	api.HandleFunc("/chats", logConsole(listChats)).Methods(http.MethodGet)

	// Entrance [HTML]
	api.Handle("/chats/{titleOrID}/entrance", errHandler(joinRoom)).Methods(http.MethodGet)

	// Check password matches room
	api.Handle("/chats/{titleOrID}/token", errHandler(login)).Methods(http.MethodPost)

	// Check password matches room
	api.Handle("/chats/{titleOrID}/token/renew", errHandler(renewToken)).Methods(http.MethodGet)

	// Load chat box [HTML]
	api.HandleFunc("/chats/{titleOrID}/chatbox", logConsole(chatbox)).Methods(http.MethodGet)

	// Chat Sessions (WebSocket)
	// Do not authorize since you can't add headers to WebSockets. We will do authorization when actually receiving chat messages
	api.Handle("/chats/{titleOrID}/ws", errHandler(authorize(webSocketHandler))).Methods(http.MethodGet)

	api.HandleFunc("/favicon.ico", faviconHandler)

	// Error page
	api.HandleFunc("/err", logConsole(handleError)).Methods(http.MethodGet)
	return api
}

func init() {
	loadConfig()
	loadEnvs()
	loadLog()
	// initialize chat server
	data.CS.Init()
	Mux = registerHandlers()
}

func loadLog() {
	file, err := os.OpenFile("vitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		file, err = os.OpenFile("../vitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", err)
		}

	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		file, err = os.Open("../config.json")
		if err != nil {
			log.Fatalln("Cannot open config file", err)
		}
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

func loadEnvs() {
	if key, ok := os.LookupEnv("SECRET_KEY"); ok {
		secretKey = key
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "image/x-icon")
	//w.Header().Set("Cache-Control", "public, max-age=7776000")
	http.ServeFile(w, r, "public/img/favicon.ico")
}
