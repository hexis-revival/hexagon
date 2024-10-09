package hscore

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AvatarHandler(ctx *Context) {
	vars := mux.Vars(ctx.Request)

	id, ok := vars["id"]
	if !ok {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	avatar, err := ctx.Server.State.Storage.GetAvatar(userId)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Response.Header().Set("Content-Type", "image/png")
	ctx.Response.Write(avatar)
}
