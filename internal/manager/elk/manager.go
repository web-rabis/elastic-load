package elk

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/web-rabis/elastic-load/internal/config"
	"github.com/web-rabis/elastic-load/internal/manager/ebook"
	"github.com/web-rabis/elastic-load/internal/model"
	"log"
	"strconv"
	"time"
)

type IManager interface {
	StartFullLoad(ctx context.Context, filter *model.EbookFilter)
	StatusFullLoad() *LoadStatus
	StopFullLoad()
	StartPartialLoad(ctx context.Context, filter *model.EbookFilter, updateFields []int64)
	StatusPartialLoad() *LoadStatus
	StopPartialLoad()
}
type Manager struct {
	indexer           esutil.BulkIndexer
	ebookMan          ebook.IManager
	fullLoadStatus    *LoadStatus
	deltaLoadStatus   *LoadStatus
	partialLoadStatus *LoadStatus

	fullCancel    context.CancelFunc
	partialCancel context.CancelFunc
}

func NewElkManager(opts *config.APIServer, ebookMan ebook.IManager) (*Manager, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:         []string{opts.ESURL},
		EnableDebugLogger: true,
	})
	if err != nil {
		return nil, err
	}
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         opts.ESINDEX, // имя индекса по умолчанию
		Client:        es,
		NumWorkers:    3,
		FlushBytes:    5 * 1024 * 1024, // 5 MB
		FlushInterval: 5 * time.Second, // 5 секунд
	})
	if err != nil {
		return nil, err
	}
	return &Manager{
		indexer:           indexer,
		ebookMan:          ebookMan,
		fullLoadStatus:    &LoadStatus{},
		partialLoadStatus: &LoadStatus{},
		deltaLoadStatus:   &LoadStatus{},
	}, nil
}

func (m *Manager) load(ctx context.Context, ebooks []model.Ebook, updateFields []int64, loadStatus *LoadStatus) error {
	var ebooksElk []model.Ebook
	var action = "index"
	if len(updateFields) > 0 {
		action = "update"
	}
	for _, book := range ebooks {
		b, err := m.ebookMan.EbookElk(ctx, book, updateFields)
		if err != nil {
			continue
		}
		ebooksElk = append(ebooksElk, b)
	}
	err := m.loadToIndex(ctx, action, ebooksElk, loadStatus)
	if err != nil {
		return err
	}
	loadStatus.AddProcessed(uint64(len(ebooksElk)))
	return nil
}
func (m *Manager) loadToIndex(ctx context.Context, action string, ebooks []model.Ebook, loadStatus *LoadStatus) error {
	retryOnConflict := 3
	for _, book := range ebooks {
		var body map[string]any = book
		if action == "update" {
			body = map[string]any{
				"doc": book,
			}
		}
		data, err := json.Marshal(body)
		if err != nil {
			log.Printf("Ошибка сериализации: %v", err)
			continue
		}
		err = m.indexer.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action:          action,
				DocumentID:      strconv.Itoa(int(book["id"].(int32))),
				Body:            bytes.NewReader(data),
				RetryOnConflict: &retryOnConflict,
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem) {
					loadStatus.AddCounters(1, 0)
				},
				OnFailure: func(
					ctx context.Context,
					item esutil.BulkIndexerItem,
					resp esutil.BulkIndexerResponseItem,
					err error,
				) {
					loadStatus.AddCounters(0, 1)
					if err != nil {
						log.Printf("[ERROR] Ошибка индексации: %v", err)
					} else {
						log.Printf("[ERROR] Ошибка от ES [%s]: %s", resp.Status, resp.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			println(err.Error())
			log.Printf("[ERROR] Ошибка добавления документа в буфер: %v", err)
			continue
		}
	}
	return nil
}
