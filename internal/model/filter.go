package model

import (
	"net/http"
	"strconv"
	"strings"
)

type EbookFilter struct {
	IdList []int64
	LastId *int64
}

func (f *EbookFilter) Sql() string {
	if f == nil {
		return ""
	}
	var sql = ""
	if f.IdList != nil && len(f.IdList) > 0 {
		if len(f.IdList) == 1 {
			sql = "id=" + strconv.FormatInt(f.IdList[0], 10)
		} else {
			for _, id := range f.IdList {
				sql += "and id=" + strconv.FormatInt(id, 10)
			}
			sql = sql[4:]
		}
	}
	if f.LastId != nil {
		if sql != "" {
			sql += " and "
		}
		sql += " id>" + strconv.FormatInt(*f.LastId, 10)
	}
	return sql
}
func EbookFilterParseFromHttp(r *http.Request) (*EbookFilter, error) {
	filter := EbookFilter{
		IdList: make([]int64, 0),
	}

	if idList := r.URL.Query().Get("idList"); idList != "" {
		for _, id := range strings.Split(idList, ",") {
			_id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, err
			}
			filter.IdList = append(filter.IdList, _id)
		}
	}

	return &filter, nil
}
