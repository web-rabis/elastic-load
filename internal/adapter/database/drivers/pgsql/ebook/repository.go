package ebook

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/web-rabis/elastic-load/internal/model"
	"gorm.io/gorm"
	"log"
	"strconv"
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
func (r *Repository) List(ctx context.Context, paging *model.Paging) ([]model.Ebook, error) {
	var books []model.Ebook
	fields := r.ebookListFields()
	var sql = "select " + strings.Join(fields, ",") + " from ebook "
	if paging != nil {
		sql = sql + paging.Sql()
	}
	result, err := r.pool.Query(ctx, sql)
	if err != nil {
		log.Printf("[ERROR] error query %s\n", err.Error())
		println(err.Error())
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		values, err := result.Values()
		if err != nil {
			log.Printf("[ERROR] error values %s\n", err.Error())
			println(err.Error())
			continue
		}
		books = append(books, r.newEbook(fields, values))
	}
	return books, nil
}
func (r *Repository) EbookElk(ctx context.Context, ebookElk model.Ebook, blocks map[*model.Block][]*model.BlockField) (model.Ebook, error) {
	ebookId := int(ebookElk["id"].(int32))
	for block, fields := range blocks {
		tableName := "ebook"
		where := "id="
		if block.ExternalTable != "" {
			tableName = tableName + "_" + block.ExternalTable
			if block.KeyValue != 0 {
				where = "key_value=" + strconv.FormatInt(block.KeyValue, 10) + " and ebook_id="
			} else {
				where = "ebook_id="
			}
		}
		where = where + strconv.Itoa(ebookId)
		var f []string
		for _, blockField := range fields {
			if blockField.Search {
				f = append(f, strings.ToLower(blockField.FieldName))
			}
		}
		if len(f) == 0 {
			continue
		}
		sql := "select " + strings.Join(f, ",") + " from " + tableName + " where " + where
		result, err := r.pool.Query(ctx, sql)
		if err != nil {
			return nil, err
		}
		if block.IsRepeat {
			ebookElk[tableName] = []model.Ebook{}
		}
		for result.Next() {
			values, err := result.Values()
			if err != nil {
				continue
			}
			blockValue := r.newEbook(f, values)
			if !block.IsRepeat {
				if block.ExternalTable != "" {
					ebookElk[tableName] = blockValue
				} else {
					for key, value := range blockValue {
						ebookElk[key] = value
					}
				}
			} else {
				ebookElk[tableName] = append(ebookElk[tableName].([]model.Ebook), blockValue)
			}

		}
		result.Close()
	}

	return ebookElk, nil
}
func (r *Repository) Count(ctx context.Context) (int64, error) {
	sql := "select count(id) from ebook"

	result, err := r.pool.Query(ctx, sql)
	if err != nil {
		return 0, err
	}
	defer result.Close()

	count := int64(0)

	for result.Next() {
		err = result.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}
func (r *Repository) ebookListFields() []string {
	return []string{
		"id",
		"catalog_id",
		"create_date",
		"edit_date",
		"state_id",
		"krv",
		"digest",
		"b_level_id",
		"type_descr_id",
		"volume_number",
		"author",
		"title",
	}
}
func (r *Repository) newEbook(fields []string, values []any) model.Ebook {
	book := model.Ebook{}
	for idx, field := range fields {
		if values[idx] == nil {
			continue
		}
		book[field] = values[idx]
	}
	return book
}
