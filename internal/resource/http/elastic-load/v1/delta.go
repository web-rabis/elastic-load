package v1

import (
	"context"
	"net/http"
)

func (res *Resource) delta(w http.ResponseWriter, r *http.Request) {
	go res.elkMan.StartDeltaLoad(context.Background())
}
