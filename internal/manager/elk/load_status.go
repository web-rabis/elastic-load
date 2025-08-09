package elk

import (
	"sync/atomic"
	"time"
)

type LoadStatus struct {
	Started        *time.Time
	Finished       *time.Time
	Stopping       bool
	Running        bool
	Error          error
	SuccessCount   uint64
	ErrorCount     uint64
	ProcessedCount uint64
	TotalCount     int64
}

func (r *LoadStatus) Start() {
	now := time.Now()
	r.Started = &now
	r.Finished = nil
	r.Running = true
	r.Stopping = false
	r.Error = nil
	atomic.StoreUint64(&r.ProcessedCount, 0)
	atomic.StoreUint64(&r.SuccessCount, 0)
	atomic.StoreUint64(&r.ErrorCount, 0)
	r.TotalCount = 0
}
func (r *LoadStatus) Finish() {
	now := time.Now()
	r.Finished = &now
	r.Running = false
	r.Stopping = false
}
func (r *LoadStatus) Fail(err error) {
	now := time.Now()
	r.Finished = &now
	r.Running = false
	r.Stopping = false
	r.Error = err
}
func (r *LoadStatus) Stop() {
	r.Stopping = true
}
func (r *LoadStatus) InitTotal(count int64) {
	r.TotalCount = count
}
func (r *LoadStatus) AddCounters(success, fail uint64) {
	atomic.AddUint64(&r.SuccessCount, success)
	atomic.AddUint64(&r.ErrorCount, fail)
}
func (r *LoadStatus) AddProcessed(delta uint64) {
	atomic.AddUint64(&r.ProcessedCount, delta)
}
