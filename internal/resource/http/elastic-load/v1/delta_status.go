package v1

import (
	"github.com/go-chi/render"
	"net/http"
)

func (res *Resource) deltaStatus(w http.ResponseWriter, r *http.Request) {
	deltaStatus := res.elkMan.StatusDeltaLoad()
	render.JSON(w, r, deltaStatus)
}
