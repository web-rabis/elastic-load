package model

type TypeDescription struct {
	Id     int64  `json:"id" bson:"id"`
	Code   string `json:"code" bson:"code"`
	Name   string `json:"name" bson:"name"`
	NameKz string `json:"nameKz" bson:"name_kz"`
}
