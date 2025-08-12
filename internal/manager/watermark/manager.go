package watermark

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers"
	"github.com/web-rabis/elastic-load/internal/model"
	"time"
)

type IManager interface {
	ByJob(ctx context.Context, jobName string) (*model.Watermark, error)
	Update(ctx context.Context, jobName string, lastId int64, lastTimestamp time.Time) error
}

type Manager struct {
	wmRepo drivers.WatermarkRepository
}

func NewWatermarkManager(ds drivers.DataStore) *Manager {
	return &Manager{wmRepo: ds.Watermark()}
}
func (m *Manager) ByJob(ctx context.Context, jobName string) (*model.Watermark, error) {
	return m.wmRepo.ByJob(ctx, jobName)
}
func (m *Manager) Update(ctx context.Context, jobName string, lastId int64, lastTimestamp time.Time) error {
	return m.wmRepo.Update(ctx, jobName, lastId, lastTimestamp)
}
