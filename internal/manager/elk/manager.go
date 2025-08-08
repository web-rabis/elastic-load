package elk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	FullLoad(ctx context.Context)
	FullLoadInfo() LoadStatus
}
type Manager struct {
	indexer           esutil.BulkIndexer
	ebookMan          ebook.IManager
	fullLoadStatus    LoadStatus
	deltaLoadStatus   LoadStatus
	partialLoadStatus LoadStatus
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
		FlushBytes:    5 * 1024 * 1024, // 5 MB
		FlushInterval: 5 * time.Second, // 5 секунд
	})
	if err != nil {
		return nil, err
	}
	return &Manager{
		indexer:  indexer,
		ebookMan: ebookMan,
	}, nil
}

func (m *Manager) FullLoad(ctx context.Context) {
	if m.fullLoadStatus.Running {
		log.Printf("[ERROR] уже запущено\n")
		return
	}
	m.fullLoadStatus.Start()
	cnt, err := m.ebookMan.EbookCount(ctx)
	if err != nil {
		log.Printf("[ERROR] error %s\n", err.Error())
		m.fullLoadStatus.Fail(err)
		return
	}
	m.fullLoadStatus.InitTotal(cnt)
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
			m.fullLoadStatus.Fail(err)
			return
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
		err = m.loadToIndex(ctx, ebooksElk)
		if err != nil {
			log.Printf("[ERROR] error %s\n", err.Error())
		}
		m.fullLoadStatus.AddCounters(int64(paging.Limit), 0)
		paging.NextPage()

	}
	m.fullLoadStatus.Finish()
}
func (m *Manager) FullLoadInfo() LoadStatus {
	return m.fullLoadStatus
}
func (m *Manager) loadToIndex(ctx context.Context, ebooks []model.Ebook) error {
	for _, book := range ebooks {
		data, err := json.Marshal(book)
		if err != nil {
			log.Printf("Ошибка сериализации: %v", err)
			continue
		}
		err = m.indexer.Add(
			ctx,
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

	stats := m.indexer.Stats()
	fmt.Printf("✅ Загрузка завершена: %d документов (успешно), %d с ошибками\n",
		stats.NumFlushed, stats.NumFailed)
	return nil
}
