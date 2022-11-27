package evthandler

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/joshuazhu78/af/pkg/util"
)

type CapEvtHandler struct {
	captureDir    string
	capturePeriod uint
	done          chan bool
}

func NewCapEventHandler(captureDir string, capturePeriod uint) EvtHandler {
	return &CapEvtHandler{
		captureDir:    captureDir,
		capturePeriod: capturePeriod,
		done:          make(chan bool),
	}
}

func (e CapEvtHandler) OnFaceDetected(metaChan chan util.GvaMeta) error {
	t := time.Now()
	t = t.Add(time.Duration(-e.capturePeriod) * time.Second)
	var err error
	for {
		select {
		case <-e.done:
			return nil
		case <-metaChan:
			tNow := time.Now()
			d := tNow.Sub(t)
			if d.Seconds() >= float64(e.capturePeriod) {
				t, err = e.captureOne()
				if err != nil {
					return err
				}
			}
		}
	}
}

func (e CapEvtHandler) OnDeactivated() error {
	e.done <- true
	return nil
}

func (e CapEvtHandler) captureOne() (time.Time, error) {
	t := time.Now()
	filename := fmt.Sprintf("%s/%s.png", e.captureDir, t.Format("2006-01-02-15-04-05"))
	cmd := exec.Command("magick", "import", "-window", "gst-launch-1.0", filename)
	err := cmd.Run()

	if err != nil {
		log.Println(err)
		return t, err
	}
	log.Printf("%s saved", filename)
	return t, nil
}
