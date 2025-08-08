package model

type Block struct {
	Id            int64  `json:"id,omitempty" bson:"id"`
	Name          string `json:"name,omitempty" bson:"name"`
	IsRepeat      bool   `json:"isRepeat,omitempty" bson:"is_repeat"`
	Priority      int64  `json:"priority,omitempty" bson:"priority"`
	ExternalTable string `json:"externalTable,omitempty" bson:"external_table"`
	KeyValue      int64  `json:"keyValue,omitempty" bson:"key_value"`
}
type BlockField struct {
	Id          int64  `json:"id,omitempty" bson:"id"`
	BlockId     int64  `json:"blockId,omitempty" bson:"block_id"`
	Name        string `json:"name,omitempty" bson:"name"`
	FieldName   string `json:"fieldName,omitempty" bson:"field_name"`
	FieldType   string `json:"fieldType,omitempty" bson:"field_type"`
	FieldLength int    `json:"fieldLength,omitempty" bson:"field_length"`
	Priority    int    `json:"priority,omitempty" bson:"priority"`
	DirId       int    `json:"dirId,omitempty" bson:"dir_id"`
	Required    bool   `json:"required,omitempty" bson:"required"`
	Search      bool   `json:"search,omitempty" bson:"search"`
}
