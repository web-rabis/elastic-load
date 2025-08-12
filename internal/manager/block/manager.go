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
	PartialBlocks(updateFields []int64) map[*model.Block][]*model.BlockField
	MapBlockFields() map[int64]*model.BlockField
}
type Manager struct {
	blocks         []*model.Block
	blockFields    []*model.BlockField
	mapBlocks      map[*model.Block][]*model.BlockField
	mapBlocks1     map[int64]*model.Block
	mapBlockFields map[int64]*model.BlockField
	blockRepo      drivers.BlockRepository
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
	mapBlocks1 := map[int64]*model.Block{}
	for _, blockField := range blockFields {
		mapBlockFields[blockField.Id] = blockField
		mapBlockFields1[blockField.BlockId] = append(mapBlockFields1[blockField.BlockId], blockField)
	}
	mapBlocks := map[*model.Block][]*model.BlockField{}
	for _, block := range blocks {
		mapBlocks[block] = mapBlockFields1[block.Id]
		mapBlocks1[block.Id] = block
	}
	return &Manager{
		blocks:         blocks,
		blockFields:    blockFields,
		mapBlocks:      mapBlocks,
		mapBlockFields: mapBlockFields,
		blockRepo:      blockRepo,
		mapBlocks1:     mapBlocks1,
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

func (m *Manager) PartialBlocks(updateFields []int64) map[*model.Block][]*model.BlockField {
	var updateBlocks = map[*model.Block][]*model.BlockField{}
	blocks := m.mapBlocks1
	blockFields := m.mapBlockFields

	for _, updateField := range updateFields {
		bf, bfExists := blockFields[updateField]
		if bfExists {
			b := blocks[bf.BlockId]
			_, bExists := updateBlocks[b]
			if !bExists {
				updateBlocks[b] = []*model.BlockField{}
			}
			_bf := *bf
			_bf.Search = true
			updateBlocks[b] = append(updateBlocks[b], &_bf)
		}
	}

	return updateBlocks
}

func (m *Manager) MapBlockFields() map[int64]*model.BlockField {
	return m.mapBlockFields
}
