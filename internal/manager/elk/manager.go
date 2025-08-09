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
	StartFullLoad(ctx context.Context)
	StatusFullLoad() *LoadStatus
	StopFullLoad()
}
type Manager struct {
	indexer           esutil.BulkIndexer
	ebookMan          ebook.IManager
	fullLoadStatus    *LoadStatus
	deltaLoadStatus   *LoadStatus
	partialLoadStatus *LoadStatus

	fullCancel context.CancelFunc
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
		indexer:           indexer,
		ebookMan:          ebookMan,
		fullLoadStatus:    &LoadStatus{},
		partialLoadStatus: &LoadStatus{},
		deltaLoadStatus:   &LoadStatus{},
	}, nil
}

func (m *Manager) StartFullLoad(ctx context.Context) {
	if m.fullLoadStatus.Running {
		log.Printf("[ERROR] уже запущено\n")
		return
	}
	log.Printf("[DEBUG] Full load started")
	// создаём контекст с отменой
	cctx, cancel := context.WithCancel(ctx)
	m.fullCancel = cancel // сохраняем для StopFullLoad
	m.fullLoadStatus.Start()
	cnt, err := m.ebookMan.EbookCount(cctx)
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
		if m.fullLoadStatus.Stopping {
			log.Printf("[DEBUG] Full load stopped")
			break
		}
		ebooks, err := m.ebookMan.EbookList(cctx, paging)
		if err != nil {
			log.Printf("[ERROR] error %s\n", err.Error())
			m.fullLoadStatus.Fail(err)
			return
		}
		if len(ebooks) == 0 {
			break
		}
		err = m.load(cctx, ebooks, m.fullLoadStatus)
		if err != nil {
			log.Printf("[ERROR] error %s\n", err.Error())
		}
		paging.NextPage()

	}
	log.Printf("[DEBUG] Full load finished")
	m.fullLoadStatus.Finish()
}
func (m *Manager) StopFullLoad() {
	log.Printf("[DEBUG] Full load will stopped")
	m.fullLoadStatus.Stopping = true
	if m.fullCancel != nil {
		m.fullCancel() // прерываем все операции с контекстом
	}
}
func (m *Manager) StatusFullLoad() *LoadStatus {
	return m.fullLoadStatus
}
func (m *Manager) load(ctx context.Context, ebooks []model.Ebook, loadStatus *LoadStatus) error {
	var ebooksElk []model.Ebook
	for _, book := range ebooks {
		b, err := m.ebookMan.EbookElk(ctx, book)
		if err != nil {
			continue
		}
		ebooksElk = append(ebooksElk, b)
	}
	err := m.loadToIndex(ctx, ebooksElk, loadStatus)
	if err != nil {
		return err
	}
	loadStatus.AddProcessed(uint64(len(ebooksElk)))
	return nil
}
func (m *Manager) loadToIndex(ctx context.Context, ebooks []model.Ebook, loadStatus *LoadStatus) error {
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
