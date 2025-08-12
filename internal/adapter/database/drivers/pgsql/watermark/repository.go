package watermark

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/web-rabis/elastic-load/internal/adapter/database/orm"
	"github.com/web-rabis/elastic-load/internal/model"
	"gorm.io/gorm"
	"strings"
	"time"
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

func (r *Repository) ByJob(ctx context.Context, jobName string) (*model.Watermark, error) {
	var watermark = &model.Watermark{
		Job:           jobName,
		LastTimestamp: time.Now(),
		LastId:        0,
		UpdatedAt:     time.Now(),
	}
	var f = strings.Join(orm.Fields(model.Watermark{}).SqlFields("elastic_sync"), ",")

	var sql = "select " + f + " from elastic_sync where job='" + jobName + "'"

	result, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	if result.Next() {
		watermark = orm.NewObjectFromResult(&model.Watermark{}, result, "", model.MappingObjects).(*model.Watermark)
	}
	return watermark, nil
}

func (r *Repository) Update(ctx context.Context, jobName string, lastId int64, lastTimestamp time.Time) error {
	const sql = `
		INSERT INTO elastic_sync (job, watermark_id, watermark_ts, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (job) DO UPDATE
		SET watermark_id = EXCLUDED.watermark_id,
		    watermark_ts = EXCLUDED.watermark_ts,
		    updated_at   = EXCLUDED.updated_at
	`
	_, err := r.pool.Exec(ctx, sql, jobName, lastId, lastTimestamp.UTC(), time.Now().UTC())
	return err
}
