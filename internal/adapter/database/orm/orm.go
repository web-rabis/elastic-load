package orm

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FieldsMap map[string]Field
type Field struct {
	Info   reflect.StructField
	Fields FieldsMap
}
type MappingObjectsFn func(t reflect.Type, v reflect.Value, result pgx.Rows, fdm map[string]int, bson string, isPtr bool)

func Fields(o interface{}) FieldsMap {
	rf := reflect.ValueOf(o)
	return fields(rf.Type())
}

func fields(rftt reflect.Type) FieldsMap {
	var fields1 = FieldsMap{}
	for i := 0; i < rftt.NumField(); i++ {
		field := rftt.Field(i)
		bson := field.Tag.Get("bson")
		if bson != "" {
			f := Field{
				Info: field,
			}
			if field.Type.Kind().String() == "ptr" {
				f.Fields = fields(field.Type.Elem())
			}
			if field.Type.Kind().String() == "struct" {
				f.Fields = fields(field.Type)
			}
			fields1[bson] = f
		}
	}
	return fields1
}

func (fm FieldsMap) SqlFields(prefix string) []string {
	return fm.sqlFields("", prefix)
}
func (fm FieldsMap) sqlFields(parent, prefix string) []string {
	var sqlFields []string
	for bson, f := range fm {
		if len(f.Fields) > 0 {
			if strings.HasPrefix(bson, "_") {
				sqlFields = append(sqlFields, f.Fields.sqlFields(parent, "")...)
			} else {
				sqlFields = append(sqlFields, bson)
				sqlFields = append(sqlFields, f.Fields.sqlFields(bson, "")...)
			}
		} else {
			if parent != "" {
				bson = parent + "." + bson + " " + parent + "_" + bson
			} else if prefix != "" {
				bson = prefix + "." + bson + " " + bson
			}
			sqlFields = append(sqlFields, bson)
		}
	}
	return sqlFields
}

func NewObjectFromResult(r interface{}, result pgx.Rows, prefix string, mappingFunc MappingObjectsFn) interface{} {
	newObjectFromResult(reflect.ValueOf(r).Elem(), result, prefix, mappingFunc)
	return r
}
func newObjectFromResult(rf reflect.Value, result pgx.Rows, prefix string, mappingFunc MappingObjectsFn) {

	var fdm = map[string]int{}
	for i, ff := range result.FieldDescriptions() {
		fdm[ff.Name] = i
	}
	for bson, field := range fields(rf.Type()) {
		if strings.HasPrefix(bson, "_") {
			setValue(field.Info.Type, rf.FieldByName(field.Info.Name), result, fdm, prefix[0:len(prefix)-1], false, mappingFunc)
		} else {
			setValue(field.Info.Type, rf.FieldByName(field.Info.Name), result, fdm, prefix+bson, false, mappingFunc)
		}

	}
}
func setValue(t reflect.Type, v reflect.Value, result pgx.Rows, fdm map[string]int, bson string, isPtr bool, mappingFunc MappingObjectsFn) {
	//println(t.Name())
	//println(bson)
	if t.Name() == "string" {
		v.Set(reflect.ValueOf(string(result.RawValues()[fdm[bson]])))
	} else if t.Name() == "int64" {
		v1, _ := result.Values()
		//fmt.Printf("bson %s\n", bson)
		//fmt.Printf("fdm %v\n", result.RawValues()[fdm[bson]])
		//fmt.Printf("fdm %v\n", v1[fdm[bson]])
		//fmt.Printf("int64 %s\n", string(result.RawValues()[fdm[bson]]))
		i, err := strconv.ParseInt(fmt.Sprintf("%v", v1[fdm[bson]]), 10, 64)
		if err == nil {
			v.Set(reflect.ValueOf(i))
		}
	} else if t.Name() == "bool" {
		v1, _ := result.Values()
		i, err := strconv.ParseBool(fmt.Sprintf("%v", v1[fdm[bson]]))
		if err == nil {
			v.Set(reflect.ValueOf(i))
		}
	} else if t.Name() == "Time" {
		time1, err := time.Parse("2006-01-02", string(result.RawValues()[fdm[bson]]))
		if err == nil {
			if isPtr {
				v.Set(reflect.ValueOf(&time1))
			} else {
				v.Set(reflect.ValueOf(time1))
			}
		}
	} else if t.Kind().String() == "ptr" {
		if string(result.RawValues()[fdm[bson]]) != "" {
			setValue(t.Elem(), v, result, fdm, bson, true, mappingFunc)
		}
	} else if t.Kind().String() == "slice" {
		v.Set(reflect.ValueOf(result.RawValues()[fdm[bson]]))
	}
	mappingFunc(t, v, result, fdm, bson, isPtr)

}

func (fm FieldsMap) FieldsValues(rf reflect.Value) ([]string, []any) {
	var fieldList []string
	var valueList []any
	for bson, field := range fm {
		if rf.FieldByName(field.Info.Name).Type().Kind() == reflect.String {
			fieldList = append(fieldList, bson)
			valueList = append(valueList, rf.FieldByName(field.Info.Name).String())
		}
		if rf.FieldByName(field.Info.Name).Type().Kind() == reflect.Int64 {
			if rf.FieldByName(field.Info.Name).String() != "" {
				fieldList = append(fieldList, bson)
				valueList = append(valueList, strconv.Itoa(int(rf.FieldByName(field.Info.Name).Int())))
			}
		}
		if rf.FieldByName(field.Info.Name).Type().Kind() == reflect.Bool {
			if rf.FieldByName(field.Info.Name).String() != "" {
				fieldList = append(fieldList, bson)
				valueList = append(valueList, strconv.FormatBool(rf.FieldByName(field.Info.Name).Bool()))
			}
		}
		if rf.FieldByName(field.Info.Name).Type().Kind() == reflect.Pointer {
			if !rf.FieldByName(field.Info.Name).IsZero() {
				if rf.FieldByName(field.Info.Name).Type().Elem().Kind() == reflect.Struct {
					if rf.FieldByName(field.Info.Name).Type().Elem().Name() == "Time" {
						fieldList = append(fieldList, bson)
						valueList = append(valueList, rf.FieldByName(field.Info.Name).Elem().Interface().(time.Time).Format("2006-01-02 15:04:05"))
					} else {
						fieldList = append(fieldList, bson)
						valueList = append(valueList, strconv.Itoa(int(rf.FieldByName(field.Info.Name).Elem().FieldByName("Id").Int())))
					}
				}
			}
		}
	}
	return fieldList, valueList
}
