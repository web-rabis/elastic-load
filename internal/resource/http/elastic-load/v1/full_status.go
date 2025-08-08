package v1

import (
	"github.com/go-chi/render"
	"net/http"
)

func (res *Resource) fullStatus(w http.ResponseWriter, r *http.Request) {
	fullStatus := res.elkMan.FullLoadInfo()
	render.JSON(w, r, fullStatus)
}
