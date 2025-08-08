package block

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/web-rabis/ela
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type Repository struct {
	pool *pgxpool.Pool
	db   *gorm.DB
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) BlockList(ctx context.Context, paging *model.Paging) ([]*model.Block, error) {
	var blocks []*model.Block
	var f = strings.Join(orm.Fields(model.Block{}).SqlFields("block"), ",")

	var sql = "select " + f + " from block"
	if paging != nil {
		sql = sql + paging.Sql()
	}

	result, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		block := orm.NewObjectFromResult(&model.Block{}, result, "", MappingObjects).(*model.Block)
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (r *Repository) BlockFieldList(ctx context.Context, paging *model.Paging) ([]*model.BlockField, error) {
	var blockFields []*model.BlockField
	var f = strings.Join(orm.Fields(model.BlockField{}).SqlFields("block_fields"), ",")
	var sql = "select " + f + " from block_fields"
	if paging != nil {
		sql = sql + paging.Sql()
	}
	result, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		blockField := orm.NewObjectFromResult(&model.BlockField{}, result, "", MappingObjects).(*model.BlockField)
		blockFields = append(blockFields, blockField)
	}
	return blockFields, nil
}
func MappingObjects(t reflect.Type, v reflect.Value, result pgx.Rows, fdm map[string]int, bson string, isPtr bool) {

}
