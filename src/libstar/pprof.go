package libstar

import (
	"net/http"
	_ "net/http/pprof"
)

type PProf struct {
	Listen string
}

func (p *PProf) Start() {
	if p.Listen == "" {
		p.Listen = "localhost:6062"
	}
	go func() {
		Info("PProf.Start %s", p.Listen)
		if err := http.ListenAndServe(p.Listen, nil); err != nil {
			Error("PProf.Start %v", err)
		}
	}()
}
