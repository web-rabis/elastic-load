package drivers

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/model"
)

type DataStore interface {
	Base
}

type Base interface {
	Name() string
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Connect(ctx context.Context) error

	Block() BlockRepository
	Ebook() EbookRepository
	Dictionary() DictionaryRepository
	Catalog() CatalogRepository
}
type BlockRepository interface {
	BlockList(ctx context.Context, paging *model.Paging) ([]*model.Block, error)
	BlockFieldList(ctx context.Context, paging *model.Paging) ([]*model.BlockField, error)
}
type EbookRepository interface {
	List(ctx context.Context, paging *model.Paging) ([]model.Ebook, error)
	Count(ctx context.Context) (int64, error)
	EbookElk(ctx context.Context, ebookElk model.Ebook, blocks map[*model.Block][]*model.BlockField) (model.Ebook, error)
}
type DictionaryRepository interface {
	BibliographicLevelList(ctx context.Context, paging *model.Paging) ([]*model.BibliographicLevel, error)
	TypeDescriptionList(ctx context.Context, paging *model.Paging) ([]*model.TypeDescription, error)
	StateList(ctx context.Context, paging *model.Paging) ([]*model.DState, error)
}
type CatalogRepository interface {
	CatalogList(ctx context.Context, paging *model.Paging) ([]*model.Catalog, error)
}
