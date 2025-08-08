package model

type Catalog struct {
	Id            int64  `json:"id" bson:"id"`
	Code          string `json:"code" bson:"code"`
	Name          string `json:"name" bson:"name"`
	Visible       int64  `json:"visible" bson:"visible"`
	Priority      int64  `json:"priority" bson:"priority"`
	TypeEbooks    string `json:"typeEbooks" bson:"type_ebooks"`
	ReaderCatalog string `json:"readerCatalog" bson:"readercatalog"`
}
