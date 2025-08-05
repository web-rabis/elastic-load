package ebook

import (
	"context"
	"elastic-load/internal/adapter/database/drivers"
	"elastic-load/internal/manager/block"
	"elastic-load/internal/model"
)

type IManager interface {
	EbookList(ctx context.Context, paging *model.Paging) ([]model.Ebook, error)
	EbookElk(ctx context.Context, ebook map[string]interface{}) (model.Ebook, error)
}
type Manager struct {
	ebookRepo drivers.EbookRepository
	blockMan  block.IManager
}

func NewEbookManager(ebookRepo drivers.EbookRepository, blockMan block.IManager) *Manager {
	return &Manager{
		ebookRepo: ebookRepo,
		blockMan:  blockMan,
	}
}

func (m *Manager) EbookList(ctx context.Context, paging *model.Paging) ([]model.Ebook, error) {
	return m.ebookRepo.List(ctx, paging)
}

func (m *Manager) EbookElk(ctx context.Context, ebook map[string]interface{}) (model.Ebook, error) {
	book, err := m.ebookRepo.EbookElk(ctx, ebook, m.blockMan.MapBlocks())
	if err != nil {
		return nil, err
	}
	return book, nil
}
