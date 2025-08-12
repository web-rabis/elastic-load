package drivers

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/model"
	"time"
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

	Watermark() WatermarkRepository
}
type BlockRepository interface {
	BlockList(ctx context.Context, paging *model.Paging) ([]*model.Block, error)
	BlockFieldList(ctx context.Context, paging *model.Paging) ([]*model.BlockField, error)
}
type EbookRepository interface {
	List(ctx context.Context, paging *model.Paging, filter *model.EbookFilter) ([]model.Ebook, error)
	Count(ctx context.Context, filter *model.EbookFilter) (int64, error)
	EbookElk(ctx context.Context, ebookElk model.Ebook, blocks map[*model.Block][]*model.BlockField) (model.Ebook, error)
	EbookYearEditions(ctx context.Context, ebookId int64) ([]int64, error)
}
type DictionaryRepository interface {
	BibliographicLevelList(ctx context.Context, paging *model.Paging) ([]*model.BibliographicLevel, error)
	TypeDescriptionList(ctx context.Context, paging *model.Paging) ([]*model.TypeDescription, error)
	StateList(ctx context.Context, paging *model.Paging) ([]*model.DState, error)
}
type CatalogRepository interface {
	CatalogList(ctx context.Context, paging *model.Paging) ([]*model.Catalog, error)
}
type WatermarkRepository interface {
	ByJob(ctx context.Context, jobName string) (*model.Watermark, error)
	Update(ctx context.Context, jobName string, lastId int64, lastTimestamp time.Time) error
}
