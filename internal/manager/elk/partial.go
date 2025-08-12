package elk

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/model"
	"log"
	"time"
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
		Limit:   5000,
		SortKey: "id",
		SortVal: 1,
	}
	for {
		if m.partialLoadStatus.Stopping {
			log.Printf("[DEBUG] Partial load stopped")
			break
		}
		ebooks, err := m.ebookMan.EbookList(cctx, paging, filter)
		if err != nil {
			log.Printf("[ERROR] Partial error %s\n", err.Error())
			m.partialLoadStatus.Fail(err)
			return
		}
		if len(ebooks) == 0 {
			break
		}
		err = m.load(cctx, ebooks, updateFields, m.partialLoadStatus)
		if err != nil {
			log.Printf("[ERROR] Partial error %s\n", err.Error())
		}
		paging.NextPage()

	}
	time.Sleep(time.Second * 30)
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
	return m.partialLoadStatus
}
