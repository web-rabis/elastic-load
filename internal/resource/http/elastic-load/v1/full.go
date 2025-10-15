package v1

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/web-rabis/elastic-load/internal/model"
)

func (res *Resource) full(w http.ResponseWriter, r *http.Request) {
	filter, err := model.EbookFilterParseFromHttp(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.HTML(w, r, err.Error())
		return
	}
	paging, err := model.PagingParseFromHttp(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.HTML(w, r, err.Error())
		return
	}
	go res.elkMan.StartFullLoad(context.Background(), filter, paging)
}
