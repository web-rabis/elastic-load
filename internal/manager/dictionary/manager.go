package dictionary

import (
	"context"
	"elastic-load/internal/adapter/database/drivers"
	"elastic-load/internal/model"
)

type IManager interface {
	BibliographicLevelById(id int64) *model.BibliographicLevel
	TypeDescriptionById(id int64) *model.TypeDescription
}
type Manager struct {
	biblLevels       map[int64]*model.BibliographicLevel
	typeDescriptions map[int64]*model.TypeDescription
}

func NewDictionaryManager(ctx context.Context, dictRepo drivers.DictionaryRepository) (*Manager, error) {
	biblLevels := map[int64]*model.BibliographicLevel{}
	typeDescriptions := map[int64]*model.TypeDescription{}
	td, err := dictRepo.TypeDescriptionList(ctx, nil)
	if err != nil {
		return nil, err
	}
	bl, err := dictRepo.BibliographicLevelList(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, t := range td {
		typeDescriptions[t.Id] = t
	}
	for _, b := range bl {
		biblLevels[b.Id] = b
	}
	return &Manager{
		biblLevels:       biblLevels,
		typeDescriptions: typeDescriptions,
	}, nil
}
func (m *Manager) BibliographicLevelById(id int64) *model.BibliographicLevel {
	return m.biblLevels[id]
}
func (m *Manager) TypeDescriptionById(id int64) *model.TypeDescription {
	return m.typeDescriptions[id]
}
