package probe

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/quanxiang-cloud/cabin/logger"
)

const (
	readinessPending int32 = iota
	readinessTrue
	readinessFalse
)

// Probe probe
type Probe struct {
	readiness int32

	log logger.AdaptedLogger
}

// New return *Probe
func New() *Probe {
	return &Probe{
		readiness: readinessPending,
		log:       logger.Logger.WithName("probe"),
	}
}

func (p *Probe) setTrue() {
	atomic.StoreInt32(&p.readiness, readinessTrue)
}

func (p *Probe) setFalse() {
	atomic.StoreInt32(&p.readiness, readinessFalse)
}

func (p *Probe) getReadiness() int32 {
	return atomic.LoadInt32(&p.readiness)
}

// SetRunning set running
func (p *Probe) SetRunning() {
	p.log.Info("probe ready")
	p.setTrue()
}

// LivenessProbe liveness probe
func (p *Probe) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	if p.getReadiness() != readinessFalse {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (p *Probe) isSafe(r *http.Request) bool {
	if strings.HasPrefix(r.Host, "127.0.0.1") ||
		strings.HasPrefix(r.Host, "localhost") {
		return true
	}
	return false
}

// ReadinessProbe readiness probe
func (p *Probe) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-readiness-shutdown") != "" {
		if !p.isSafe(r) {
			p.log.Info("try to shutdown,but is not safe. refuse!", "host", r.Host)
			return
		}
		p.log.Info("readiness shutdown")
		p.setFalse()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if p.getReadiness() == readinessTrue {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
