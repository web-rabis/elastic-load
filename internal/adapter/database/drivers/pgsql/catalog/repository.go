package catalog

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/web-rabis/elastic-load/internal/adapter/database/orm"
	"github.com/web-rabis/elastic-load/internal/model"
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

func (r *Repository) CatalogList(ctx context.Context, paging *model.Paging) ([]*model.Catalog, error) {
	var catalogs []*model.Catalog
	fields := orm.Fields(model.Catalog{})
	var sql = "select " + strings.Join(fields.SqlFields(""), ",") + " from catalog "
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
		catalogs = append(catalogs, orm.NewObjectFromResult(&model.Catalog{}, result, "", model.MappingObjects).(*model.Catalog))
	}
	return catalogs, nil
}
