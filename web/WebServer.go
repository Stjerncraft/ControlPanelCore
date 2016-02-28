package web

import (
	"github.com/stjerncraft/controlpanelcore/core"
	"net/http"
"log"
	"github.com/gorilla/mux"
"github.com/gorilla/websocket"
	"fmt"
)

type WebServer struct {
	coreInst *core.Core

	templates *TemplateManager
	modules *ModuleManager
}

func NewWebServer(core *core.Core) (*WebServer, error) {
	newServer := new(WebServer)

	newServer.coreInst = core

	newServer.templates = NewTemplateManager()
	err := newServer.templates.ReloadTemplates()
	if err != nil {
		return nil, err
	}

	newServer.modules = NewModuleManager()
	_, err = newServer.modules.NewModule("core", AccessTypePublic)
	if err != nil {
		return nil, err
	}

	return newServer, nil
}

func (server *WebServer) StartServer() {
	router := mux.NewRouter()

	router.StrictSlash(true)

	router.HandleFunc("/api", apiSocketHandler)
	router.PathPrefix("/public").Handler(http.FileServer(http.Dir("public")))
	router.HandleFunc("/", server.mainHandler)
	router.HandleFunc("/login", server.loginHandler)
	router.HandleFunc("/module/{name}/{file:.*}", server.moduleHandler)
	//router.HandleFunc("/template/{file:.*}", server.templateHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func apiSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error WebSocket:", err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("SocketReadErr:",err)
			return
		}

		fmt.Println("Websocket got:", string(p))

		err = conn.WriteMessage(messageType, p)
		if err != nil {
			fmt.Println("SocketWriteErr:",err)
			return
		}
	}
}

func (server *WebServer) mainHandler(w http.ResponseWriter, r *http.Request) {
	//Check AuthCode from Cookie, redirect to login if missing or wrong.

	tmpl := server.templates.GetCurrentTemplate()
	if tmpl == nil {
		fmt.Println("Error: Missing current template! Reloading templates.")
		server.templates.ReloadTemplates()
		tmpl = server.templates.GetCurrentTemplate()

		if tmpl == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("No valid template found!"))
			return
		}
	}

	tmpl.WriteMain(w)
}

func (server *WebServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	//Receives: Username and password
	//Sends: AuthCode which is used when establishing WebSocket connection
	//AuthCode is stored in cookie if remember-me enabled.

	//If no username and password is sent in, return the login page.
}

func (server *WebServer) moduleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleName := vars["name"]
	if moduleName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing Module name"))
		return
	}

	module := server.modules.GetModule(moduleName)
	if module == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Module not found: " + moduleName))
		return
	}

	if(module.access == AccessTypePrivate) {
		//TODO: Check AuthCode from cookie unless trying to access 'moduleName/public/'
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Module is private, requires user to be logged in"))
		return
	}

	fs :=  http.FileServer(http.Dir("./modules/" + moduleName + "/"))
	http.StripPrefix("/module/" + moduleName, fs).ServeHTTP(w, r)
}

