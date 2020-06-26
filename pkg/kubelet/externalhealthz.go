package kubelet

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type externalHealthz struct {
	interval   time.Duration
	port       int32
	lastStatus bool
	mux        sync.Mutex
}

func (ez *externalHealthz) getLastStatus() (bool, error) {
	return ez.lastStatus, nil
}

func (ez *externalHealthz) periodicallyCheckExternalHealthz() {
	ticker := time.NewTicker(ez.interval)
	go func() {
		for range ticker.C {
			result := ez.checkExternalHealthz()
			ez.mux.Lock()
			ez.lastStatus = result
			ez.mux.Unlock()
		}
	}()
}

func (ez *externalHealthz) checkExternalHealthz() bool {
	if ez.port == 0 {
		return true
	}
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz", ez.port))
	if err != nil {
		return false
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}
	return false
}
