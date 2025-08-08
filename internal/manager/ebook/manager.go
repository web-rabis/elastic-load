package ebook

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers"
	"github.com/web-rabis/elastic-load/internal/manager/block"
	"github.com/web-rabis/elastic-load/internal/manager/catalog"
	"github.com/web-rabis/elastic-load/internal/manager/dictionary"
	"github.com/web-rabis/elastic-load/internal/model"
)

type IManager interface {
	EbookList(ctx context.Context, paging *model.Paging) ([]model.Ebook, error)
	EbookCount(ctx context.Context) (int64, error)
	EbookElk(ctx context.Context, ebook map[string]interface{}) (model.Ebook, error)
}
type Manager struct {
	ebookRepo drivers.EbookRepository
	blockMan  block.IManager
	dictMan   dictionary.IManager
	catMan    catalog.IManager
}

func NewEbookManager(
	ebookRepo drivers.EbookRepository,
	blockMan block.IManager,
	dictMan dictionary.IManager,
	catMan catalog.IManager) *Manager {
	return &Manager{
		ebookRepo: ebookRepo,
		blockMan:  blockMan,
		dictMan:   dictMan,
		catMan:    catMan,
	}
}

func (m *Manager) EbookList(ctx context.Context, paging *model.Paging) ([]model.Ebook, error) {
	return m.ebookRepo.List(ctx, paging)
}
func (m *Manager) EbookCount(ctx context.Context) (int64, error) {
	return m.ebookRepo.Count(ctx)
}

func (m *Manager) EbookElk(ctx context.Context, ebook map[string]interface{}) (model.Ebook, error) {
	book, err := m.ebookRepo.EbookElk(ctx, ebook, m.blockMan.MapBlocks())
	if err != nil {
		return nil, err
	}
	book["catalog"] = m.catMan.CatalogById(int64(book["catalog_id"].(int32)))
	book["b_level"] = m.dictMan.BibliographicLevelById(int64(book["b_level_id"].(int32)))
	book["type_descr"] = m.dictMan.TypeDescriptionById(int64(book["type_descr_id"].(int32)))
	book["state"] = m.dictMan.StateById(book["state_id"].(int64))
	return book, nil
}
