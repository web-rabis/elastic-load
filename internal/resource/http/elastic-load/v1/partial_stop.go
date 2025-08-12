package v1

import (
	"net/http"
)

func (res *Resource) partialStop(w http.ResponseWriter, r *http.Request) {
	res.elkMan.StopFullLoad()
}
