package dictionary

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers"
	"github.com/web-rabis/elastic-load/internal/model"
)

type IManager interface {
	BibliographicLevelById(id int64) *model.BibliographicLevel
	TypeDescriptionById(id int64) *model.TypeDescription
	StateById(id int64) *model.DState
}
type Manager struct {
	biblLevels       map[int64]*model.BibliographicLevel
	typeDescriptions map[int64]*model.TypeDescription
	states           map[int64]*model.DState
}

func NewDictionaryManager(ctx context.Context, dictRepo drivers.DictionaryRepository) (*Manager, error) {
	biblLevels := map[int64]*model.BibliographicLevel{}
	typeDescriptions := map[int64]*model.TypeDescription{}
	states := map[int64]*model.DState{}
	td, err := dictRepo.TypeDescriptionList(ctx, nil)
	if err != nil {
		return nil, err
	}
	bl, err := dictRepo.BibliographicLevelList(ctx, nil)
	if err != nil {
		return nil, err
	}
	st, err := dictRepo.StateList(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, t := range td {
		typeDescriptions[t.Id] = t
	}
	for _, b := range bl {
		biblLevels[b.Id] = b
	}
	for _, s := range st {
		states[s.Id] = s
	}
	return &Manager{
		biblLevels:       biblLevels,
		typeDescriptions: typeDescriptions,
		states:           states,
	}, nil
}
func (m *Manager) BibliographicLevelById(id int64) *model.BibliographicLevel {
	return m.biblLevels[id]
}
func (m *Manager) TypeDescriptionById(id int64) *model.TypeDescription {
	return m.typeDescriptions[id]
}
func (m *Manager) StateById(id int64) *model.DState {
	return m.states[id]
}
