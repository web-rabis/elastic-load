package v1

import (
	"net/http"
)

func (res *Resource) deltaStop(w http.ResponseWriter, r *http.Request) {
	res.elkMan.StopDeltaLoad()
}
