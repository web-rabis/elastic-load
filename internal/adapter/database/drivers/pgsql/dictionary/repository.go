package dictionary

import (
	"context"
	"elastic-load/internal/adapter/database/orm"
	"elastic-load/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
	"log"
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

func (r *Repository) BibliographicLevelList(ctx context.Context, paging *model.Paging) ([]*model.BibliographicLevel, error) {
	var levels []*model.BibliographicLevel
	fields := []string{
		"id",
		"code",
		"name",
		"type_ebooks",
	}
	var sql = "select " + strings.Join(fields, ",") + " from bibliographic_level "
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
		levels = append(levels, orm.NewObjectFromResult(&model.BibliographicLevel{}, result, "", model.MappingObjects).(*model.BibliographicLevel))
	}
	return levels, nil
}

func (r *Repository) TypeDescriptionList(ctx context.Context, paging *model.Paging) ([]*model.TypeDescription, error) {
	var descriptions []*model.TypeDescription
	fields := []string{
		"id",
		"code",
		"name",
		"type_ebooks",
	}
	var sql = "select " + strings.Join(fields, ",") + " from bibliographic_level "
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
		descriptions = append(descriptions, orm.NewObjectFromResult(&model.TypeDescription{}, result, "", model.MappingObjects).(*model.TypeDescription))
	}
	return descriptions, nil
}
