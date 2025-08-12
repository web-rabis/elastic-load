package v1

import (
	"github.com/go-chi/render"
	"net/http"
)

func (res *Resource) partialStatus(w http.ResponseWriter, r *http.Request) {
	partialStatus := res.elkMan.StatusPartialLoad()
	render.JSON(w, r, partialStatus)
}
