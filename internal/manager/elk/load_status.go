package elk

import (
	"sync/atomic"
	"time"
)

type LoadStatus struct {
	Started          *time.Time
	Finished         *time.Time
	EstimateFinished *time.Time
	Stopping         bool
	Running          bool
	Error            error
	SuccessCount     uint64
	ErrorCount       uint64
	ProcessedCount   uint64
	TotalCount       int64
	Rpm              uint64
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

func (r *LoadStatus) EstimateETA() {
	if r.Started == nil {
		return
	}
	elapsedMin := time.Now().Sub(*r.Started).Minutes()
	if elapsedMin <= 0 {
		return
	}
	r.Rpm = uint64(float64(r.ProcessedCount) / elapsedMin)
	if r.Rpm <= 0 {
		return
	}
	var remaining uint64 = 0
	if uint64(r.TotalCount) > r.ProcessedCount {
		remaining = uint64(r.TotalCount) - r.ProcessedCount
	}
	etaMin := float64(remaining) / float64(r.Rpm)
	eta := time.Duration(etaMin * float64(time.Minute))
	finishAt := time.Now().Add(eta)
	r.EstimateFinished = &finishAt
	return
}
