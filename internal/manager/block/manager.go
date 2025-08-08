package block

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers"
	"github.com/web-rabis/elastic-load/internal/model"
)

type IManager interface {
	Blocks() []*model.Block
	BlockFields() []*model.BlockField
	MapBlocks() map[*model.Block][]*model.BlockField
	MapBlockFields() map[int64]*model.BlockField
}
type Manager struct {
	blocks          []*model.Block
	blockFields     []*model.BlockField
	mapBlocks       map[*model.Block][]*model.BlockField
	mapBlockFields  map[int64]*model.BlockField
	mapBlockFields1 map[int64][]*model.BlockField
	blockRepo       drivers.BlockRepository
}

func NewBlockManager(ctx context.Context, blockRepo drivers.BlockRepository) (*Manager, error) {
	blocks, err := blockRepo.BlockList(ctx, nil)
	if err != nil {
		return nil, err
	}
	blockFields, err := blockRepo.BlockFieldList(ctx, nil)
	if err != nil {
		return nil, err
	}
	mapBlockFields := map[int64]*model.BlockField{}
	mapBlockFields1 := map[int64][]*model.BlockField{}
	for _, blockField := range blockFields {
		mapBlockFields[blockField.Id] = blockField
		mapBlockFields1[blockField.BlockId] = append(mapBlockFields1[blockField.BlockId], blockField)
	}
	mapBlocks := map[*model.Block][]*model.BlockField{}
	for _, block := range blocks {
		mapBlocks[block] = mapBlockFields1[block.Id]
	}
	return &Manager{
		blocks:         blocks,
		blockFields:    blockFields,
		mapBlocks:      mapBlocks,
		mapBlockFields: mapBlockFields,
		blockRepo:      blockRepo,
	}, nil
}

func (m *Manager) Blocks() []*model.Block {
	return m.blocks
}

func (m *Manager) BlockFields() []*model.BlockField {
	return m.blockFields
}

func (m *Manager) MapBlocks() map[*model.Block][]*model.BlockField {
	return m.mapBlocks
}

func (m *Manager) MapBlockFields() map[int64]*model.BlockField {
	return m.mapBlockFields
}
