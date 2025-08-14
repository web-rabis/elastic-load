package elk

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/model"
	"log"
	"time"
)

func (m *Manager) StartFullLoad(ctx context.Context, filter *model.EbookFilter) {
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
	paging := &model.Paging{
		Skip:    0,
		Limit:   5000,
		SortKey: "id",
		SortVal: 1,
	}
	var lastId int64 = 0
	var lastTimestamp time.Time
	for {
		if m.fullLoadStatus.Stopping {
			log.Printf("[DEBUG] Full load stopped")
			break
		}
		ebooks, err := m.ebookMan.EbookList(cctx, paging, filter)
		if err != nil {
			log.Printf("[ERROR] error %s\n", err.Error())
			m.fullLoadStatus.Fail(err)
			return
		}
		if len(ebooks) == 0 {
			break
		}
		err = m.load(cctx, ebooks, []int64{}, m.fullLoadStatus)
		if err != nil {
			log.Printf("[ERROR] error %s\n", err.Error())
		}
		paging.NextPage()
		last := ebooks[len(ebooks)-1]
		lastId = int64(last["id"].(int32))
		if last["edit_date"] != nil {
			lastTimestamp = last["edit_date"].(time.Time)
		} else if last["create_date"] != nil {
			lastTimestamp = last["create_date"].(time.Time)
		}
	}
	log.Printf("[DEBUG] Full load finished")
	if !m.fullLoadStatus.Stopping {
		err = m.wmMan.Update(cctx, "books_index", lastId, lastTimestamp)
	}
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
	m.fullLoadStatus.EstimateETA()
	return m.fullLoadStatus
}
