package model

import (
	"net/http"
	"strconv"
	"strings"
)

type EbookFilter struct {
	IdList []int
}

func (f *EbookFilter) Sql() string {
	if f == nil {
		return ""
	}
	var sql = ""
	if f.IdList != nil && len(f.IdList) > 0 {
		if len(f.IdList) == 1 {
			sql = "id=" + strconv.Itoa(f.IdList[0])
		} else {
			for _, id := range f.IdList {
				sql += "and id=" + strconv.Itoa(id)
			}
			sql = sql[4:]
		}
	}
	return sql
}
func EbookFilterParseFromHttp(r *http.Request) (*EbookFilter, error) {
	filter := EbookFilter{
		IdList: make([]int, 0),
	}

	if idList := r.URL.Query().Get("idList"); idList != "" {
		for _, id := range strings.Split(idList, ",") {
			_id, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			filter.IdList = append(filter.IdList, _id)
		}
	}

	return &filter, nil
}
