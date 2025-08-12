package v1

import (
	"context"
	"github.com/go-chi/render"
	"github.com/web-rabis/elastic-load/internal/model"
	"net/http"
)

func (res *Resource) full(w http.ResponseWriter, r *http.Request) {
	filter, err := model.EbookFilterParseFromHttp(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.HTML(w, r, err.Error())
		return
	}
	go res.elkMan.StartFullLoad(context.Background(), filter)
}
