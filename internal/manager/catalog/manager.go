package catalog

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers"
	"github.com/web-rabis/elastic-load/internal/model"
)

type IManager interface {
	CatalogById(id int64) *model.Catalog
}
type Manager struct {
	mapCatalogs map[int64]*model.Catalog
	catalogRepo drivers.CatalogRepository
}

func NewCatalogManager(ctx context.Context, catalogRepo drivers.CatalogRepository) (*Manager, error) {
	cats, err := catalogRepo.CatalogList(ctx, nil)
	if err != nil {
		return nil, err
	}
	mapCatalogs := map[int64]*model.Catalog{}
	for _, c := range cats {
		mapCatalogs[c.Id] = c
	}
	return &Manager{
		mapCatalogs: mapCatalogs,
	}, nil
}
func (m *Manager) CatalogById(id int64) *model.Catalog {
	return m.mapCatalogs[id]
}
