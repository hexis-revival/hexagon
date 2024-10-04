package hscore

import (
	"fmt"
	"net/http"

	"github.com/lekuruu/hexagon/common"
)

type ScoreServer struct {
	Host   string
	Port   int
	Logger *common.Logger
}

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Server   *ScoreServer
}

func (server *ScoreServer) Serve() {
	bind := fmt.Sprintf("%s:%d", server.Host, server.Port)
	server.Logger.Infof("Listening on %s", bind)

	http.HandleFunc("/score/submit", withContext(ScoreSubmissionHandler))
	http.ListenAndServe(bind, nil)
}

func withContext(handler func(*Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := &Context{
			Response: w,
			Request:  r,
			Server:   nil,
		}

		handler(context)
	}
}
