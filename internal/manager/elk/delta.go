package elk

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/web-rabis/elastic-load/internal/model"
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
		LastTimestamp: &wm.LastTimestamp,
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
		SortKey: "edit_date",
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
	if !m.deltaLoadStatus.Stopping {
		err = m.wmMan.Update(cctx, wm.Job, wm.LastId, wm.LastTimestamp)
	}
	m.deltaLoadStatus.Finish()

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
func (m *Manager) StartDeltaScheduler(ctx context.Context) error {
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return err
	}
	s, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
		gocron.WithGlobalJobOptions(gocron.WithSingletonMode(gocron.LimitModeReschedule)),
	)
	if err != nil {
		return err
	}
	_, err = s.NewJob(
		gocron.CronJob("*/5 * * * *", false),
		gocron.NewTask(func(jobCtx context.Context) {
			m.StartDeltaLoad(jobCtx)
		}),
		gocron.WithContext(ctx),
		gocron.WithName("delta-books_index"),
	)
	if err != nil {
		return err
	}
	s.Start()
	m.deltaSched = s
	return nil
}

// метод остановки
func (m *Manager) StopDeltaScheduler() error {
	if m.deltaSched != nil {
		return m.deltaSched.Shutdown()
	}
	return nil
}
