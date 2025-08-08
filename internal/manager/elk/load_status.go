package elk

import "time"

type LoadStatus struct {
	Started      *time.Time
	Finished     *time.Time
	Running      bool
	Error        error
	SuccessCount int64
	ErrorCount   int64
	TotalCount   int64
}

func (r *LoadStatus) Start() {
	now := time.Now()
	r.Started = &now
	r.Finished = nil
	r.Running = true
	r.Error = nil
	r.SuccessCount = 0
	r.ErrorCount = 0
	r.TotalCount = 0
}
func (r *LoadStatus) Finish() {
	now := time.Now()
	r.Finished = &now
	r.Running = false
}
func (r *LoadStatus) Fail(err error) {
	now := time.Now()
	r.Finished = &now
	r.Running = false
	r.Error = err
}
func (r *LoadStatus) InitTotal(count int64) {
	r.TotalCount = count
}
func (r *LoadStatus) AddCounters(success, fail int64) {
	r.SuccessCount = r.SuccessCount + success
	r.ErrorCount = r.ErrorCount + fail
}
