package elk

import (
	"context"
	"log"
	"sync"

	"github.com/web-rabis/elastic-load/internal/model"
)

func (m *Manager) StartFullLoad(ctx context.Context, filter *model.EbookFilter, paging *model.Paging) {
	if m.fullLoadStatus.Running {
		log.Printf("[ERROR] уже запущено\n")
		return
	}
	log.Printf("[DEBUG] Full load started")
	// создаём контекст с отменой
	cctx, cancel := context.WithCancel(ctx)
	m.fullCancel = cancel // сохраняем для StopFullLoad
	m.fullLoadStatus.Start()
	cnt, err := m.ebookMan.EbookCount(cctx, filter)
	if err != nil {
		log.Printf("[ERROR] error %s\n", err.Error())
		m.fullLoadStatus.Fail(err)
		return
	}
	m.fullLoadStatus.InitTotal(cnt)
	if paging.Limit == 0 {
		paging.Limit = 1000
	}
	sem := make(chan struct{}, 5) // семафор на 2 слота
	var wg sync.WaitGroup
	for {
		if m.fullLoadStatus.Stopping {
			log.Printf("[DEBUG] Full load stopped")
			break
		}
		sem <- struct{}{} // займём слот
		wg.Add(1)
		go func(p model.Paging, f *model.EbookFilter) {
			defer wg.Done()
			defer func() {
				<-sem
			}() // освободим слот
			_ = m.BulkLoad(cctx, p, f)
		}(*paging, filter)
		paging.NextPage()
	}
	log.Printf("[DEBUG] Full load finished")
	m.fullLoadStatus.Finish()
}
func (m *Manager) BulkLoad(cctx context.Context, paging model.Paging, filter *model.EbookFilter) error {
	log.Printf("[DEBUG] Bulk load started skip=%v", paging.Skip)
	ebooks, err := m.ebookMan.EbookList(cctx, &paging, filter)
	if err != nil {
		log.Printf("[ERROR] error %s\n", err.Error())
		m.fullLoadStatus.Fail(err)
		return err
	}
	if len(ebooks) == 0 {
		return nil
	}
	err = m.load(cctx, ebooks, []int64{}, m.fullLoadStatus)
	if err != nil {
		log.Printf("[ERROR] error %s\n", err.Error())
		return err
	}
	return err
}
func (m *Manager) StopFullLoad() {
	log.Printf("[DEBUG] Full load will stopped")
	m.fullLoadStatus.Stopping = true
	if m.fullCancel != nil {
		m.fullCancel() // прерываем все операции с контекстом
	}
}
func (m *Manager) StatusFullLoad() *LoadStatus {
	m.fullLoadStatus.EstimateETA()
	return m.fullLoadStatus
}
