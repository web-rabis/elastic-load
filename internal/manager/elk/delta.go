package elk

import (
	"context"
	"github.com/web-rabis/elastic-load/internal/model"
	"log"
	"time"
)

func (m *Manager) StartDeltaLoad(ctx context.Context) {
	if m.deltaLoadStatus.Running {
		log.Printf("[ERROR] Delta load уже запущено\n")
		return
	}
	log.Printf("[DEBUG] Delta load started")
	// создаём контекст с отменой
	cctx, cancel := context.WithCancel(ctx)
	m.deltaCancel = cancel // сохраняем для StopFullLoad
	m.deltaLoadStatus.Start()
	wm, err := m.wmMan.ByJob(cctx, "books_index")
	if err != nil {
		log.Printf("[ERROR] Delta error %s\n", err.Error())
		m.deltaLoadStatus.Fail(err)
		return
	}
	filter := &model.EbookFilter{
		LastId: &wm.LastId,
	}
	cnt, err := m.ebookMan.EbookCount(cctx, filter)
	if err != nil {
		log.Printf("[ERROR] Delta error %s\n", err.Error())
		m.deltaLoadStatus.Fail(err)
		return
	}
	m.deltaLoadStatus.InitTotal(cnt)
	paging := &model.Paging{
		Skip:    0,
		Limit:   5000,
		SortKey: "id",
		SortVal: 1,
	}
	for {
		if m.deltaLoadStatus.Stopping {
			log.Printf("[DEBUG] Delta load stopped")
			break
		}
		ebooks, err := m.ebookMan.EbookList(cctx, paging, filter)
		if err != nil {
			log.Printf("[ERROR] Delta error %s\n", err.Error())
			m.deltaLoadStatus.Fail(err)
			return
		}
		if len(ebooks) == 0 {
			break
		}
		err = m.load(cctx, ebooks, []int64{}, m.deltaLoadStatus)
		if err != nil {
			log.Printf("[ERROR] Delta error %s\n", err.Error())
		}
		paging.NextPage()
		last := ebooks[len(ebooks)-1]
		wm.LastId = int64(last["id"].(int32))
		if last["edit_date"] != nil {
			wm.LastTimestamp = last["edit_date"].(time.Time)
		} else if last["create_date"] != nil {
			wm.LastTimestamp = last["create_date"].(time.Time)
		}
	}
	log.Printf("[DEBUG] Delta load finished")
	m.deltaLoadStatus.Finish()
	err = m.wmMan.Update(cctx, wm.Job, wm.LastId, wm.LastTimestamp)
	if err != nil {
		log.Printf("[ERROR] Delta update watermark error %s \n", err.Error())
	}

}
func (m *Manager) StopDeltaLoad() {
	log.Printf("[DEBUG] Delta load will stopped")
	m.deltaLoadStatus.Stopping = true
	if m.deltaCancel != nil {
		m.deltaCancel() // прерываем все операции с контекстом
	}
}
func (m *Manager) StatusDeltaLoad() *LoadStatus {
	m.deltaLoadStatus.EstimateETA()
	return m.deltaLoadStatus
}
