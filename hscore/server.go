package hscore

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hexis-revival/hexagon/common"
)

type ScoreServer struct {
	Host   string
	Port   int
	Logger *common.Logger
	State  *common.State
}

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Server   *ScoreServer
}

func (server *ScoreServer) Serve() {
	bind := fmt.Sprintf("%s:%d", server.Host, server.Port)
	server.Logger.Infof("Listening on %s", bind)

	r := mux.NewRouter()
	r.HandleFunc("/score/submit", server.contextMiddleware(ScoreSubmissionHandler)).Methods("POST")

	loggedMux := server.loggingMiddleware(r)
	http.ListenAndServe(bind, loggedMux)
}

func NewServer(host string, port int, logger *common.Logger, state *common.State) *ScoreServer {
	return &ScoreServer{
		Host:   host,
		Port:   port,
		Logger: logger,
		State:  state,
	}
}

func (server *ScoreServer) contextMiddleware(handler func(*Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := &Context{
			Response: w,
			Request:  r,
			Server:   server,
		}

		handler(context)
	}
}

func (server *ScoreServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		time := time.Since(start)

		server.Logger.Infof(
			"%s - %s %s (%v)",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			time,
		)
	})
}
