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

func (e CapEvtHandler) OnFaceDetected(meta util.GvaMeta) error {
	err := e.captureOne()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Duration(e.capturePeriod) * time.Second)
	for {
		select {
		case <-e.done:
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := e.captureOne()
			if err != nil {
				return err
			}
		}
	}
}

func (e CapEvtHandler) OnDeactivated(inactiveTime uint) error {
	e.done <- true
	return nil
}

func (e CapEvtHandler) captureOne() error {
	t := time.Now()
	filename := fmt.Sprintf("%s/%s.png", e.captureDir, t.Format("2006-01-02-15-04-05"))
	cmd := exec.Command("magick", "import", "-window", "gst-launch-1.0", filename)
	err := cmd.Run()

	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("%s saved", filename)
	return nil
}
