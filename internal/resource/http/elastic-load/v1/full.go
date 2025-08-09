package v1

import (
	"context"
	"net/http"
)

func (res *Resource) full(w http.ResponseWriter, r *http.Request) {
	go res.elkMan.StartFullLoad(context.Background())
}
