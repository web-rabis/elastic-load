package v1

import (
	"context"
	"github.com/go-chi/render"
	"github.com/web-rabis/elastic-load/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (res *Resource) partial(w http.ResponseWriter, r *http.Request) {
	filter, err := model.EbookFilterParseFromHttp(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.HTML(w, r, err.Error())
		return
	}
	var updateFields []int64
	uf := strings.Split(r.URL.Query().Get("updateFields"), ",")
	for _, f := range uf {
		field, err := strconv.ParseInt(f, 10, 64)
		if err == nil {
			updateFields = append(updateFields, field)
		}
	}
	go res.elkMan.StartPartialLoad(context.Background(), filter, updateFields)
}
