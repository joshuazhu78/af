package evthandler

import "github.com/joshuazhu78/af/pkg/util"

type EvtHandler interface {
	OnFaceDetected(metaChan chan util.GvaMeta) error

	OnDeactivated() error
}
