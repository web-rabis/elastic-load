package model

import (
	"github.com/jackc/pgx/v5"
	"reflect"
)

type BibliographicLevel struct {
	Id         int64  `json:"id" bson:"id"`
	Code       string `json:"code" bson:"code"`
	Name       string `json:"name" bson:"name"`
	TypeEbooks string `json:"typeEbooks" bson:"type_ebooks"`
}

func MappingObjects(t reflect.Type, v reflect.Value, result pgx.Rows, fdm map[string]int, bson string, isPtr bool) {

}
