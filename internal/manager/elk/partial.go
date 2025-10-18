package elk

import (
	"context"
	"log"
	"sync"

	"github.com/web-rabis/elastic-load/internal/model"
)

func (m *Manager) StartPartialLoad(ctx context.Context, filter *model.EbookFilter, updateFields []int64) {
	if m.partialLoadStatus.Running {
		log.Printf("[ERROR] Partial load уже запущено\n")
		return
	}
	log.Printf("[DEBUG] Partial load started")
	// создаём контекст с отменой
	cctx, cancel := context.WithCancel(ctx)
	m.partialCancel = cancel // сохраняем для StopFullLoad
	m.partialLoadStatus.Start()
	cnt, err := m.ebookMan.EbookCount(cctx, filter)
	if err != nil {
		log.Printf("[ERROR] Partial error %s\n", err.Error())
		m.partialLoadStatus.Fail(err)
		return
	}
	m.partialLoadStatus.InitTotal(cnt)
	paging := &model.Paging{
		Skip:    0,
		Limit:   1000,
		SortKey: "id",
		SortVal: 1,
	}
	sem := make(chan struct{}, 5) // семафор на 2 слота
	var wg sync.WaitGroup
	for {
		if int64(paging.Skip) > cnt {
			break
		}
		if m.partialLoadStatus.Stopping {
			log.Printf("[DEBUG] Partial load stopped")
			break
		}
		sem <- struct{}{} // займём слот
		wg.Add(1)
		go func(updateFields []int64, p model.Paging, f *model.EbookFilter) {
			defer wg.Done()
			defer func() {
				<-sem
			}() // освободим слот
			_ = m.BulkLoad(cctx, updateFields, m.partialLoadStatus, p, f)
		}(updateFields, *paging, filter)
		paging.NextPage()

	}
	log.Printf("[DEBUG] Partial load finished")
	m.partialLoadStatus.Finish()

}

func (m *Manager) StopPartialLoad() {
	log.Printf("[DEBUG] Partial load will stopped")
	m.partialLoadStatus.Stopping = true
	if m.partialCancel != nil {
		m.partialCancel() // прерываем все операции с контекстом
	}
}
func (m *Manager) StatusPartialLoad() *LoadStatus {
	m.partialLoadStatus.EstimateETA()
	return m.partialLoadStatus
}
