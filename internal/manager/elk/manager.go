package elk

import (
	"bytes"
	"context"
	"elastic-load/internal/manager/ebook"
	"elastic-load/internal/model"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
	"strconv"
)

type Manager struct {
	indexer  esutil.BulkIndexer
	ebookMan ebook.IManager
}

func NewElkManager(ebookMan ebook.IManager) (*Manager, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         "books_index", // имя индекса по умолчанию
		Client:        es,
		FlushBytes:    5 * 1024 * 1024, // 5 MB
		FlushInterval: 5_000_000_000,   // 5 секунд
	})
	if err != nil {
	}
	return &Manager{
		indexer:  indexer,
		ebookMan: ebookMan,
	}, nil
}

func (m *Manager) Load(ctx context.Context) error {
	paging := &model.Paging{
		Skip:    0,
		Limit:   5000,
		SortKey: "id",
		SortVal: 1,
	}
	for {
		ebooks, err := m.ebookMan.EbookList(ctx, paging)
		if err != nil {
			log.Printf("[ERROR] error %s\n", err.Error())
			return err
		}
		if len(ebooks) == 0 {
			log.Printf("[ERROR] o count \n")
			break
		}
		var ebooksElk []model.Ebook
		for _, book := range ebooks {
			b, err := m.ebookMan.EbookElk(ctx, book)
			if err != nil {
				continue
			}
			ebooksElk = append(ebooksElk, b)
		}
		err = m.LoadToIndex(ctx, ebooksElk)
		if err != nil {
			log.Printf("[ERROR] error %s\n", err.Error())
		}
		paging.NextPage()
	}
	return nil
}
func (m *Manager) LoadToIndex(ctx context.Context, ebooks []model.Ebook) error {
	for _, book := range ebooks {
		data, err := json.Marshal(book)
		if err != nil {
			log.Printf("Ошибка сериализации: %v", err)
			continue
		}
		err = m.indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(int(book["id"].(int32))),
				Body:       bytes.NewReader(data),
				OnFailure: func(
					ctx context.Context,
					item esutil.BulkIndexerItem,
					resp esutil.BulkIndexerResponseItem,
					err error,
				) {
					if err != nil {
						log.Printf("Ошибка индексации: %v", err)
					} else {
						log.Printf("Ошибка от ES [%s]: %s", resp.Status, resp.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Printf("Ошибка добавления документа в буфер: %v", err)
			continue
		}
	}

	stats := m.indexer.Stats()
	fmt.Printf("✅ Загрузка завершена: %d документов (успешно), %d с ошибками\n",
		stats.NumFlushed, stats.NumFailed)
	return nil
}
