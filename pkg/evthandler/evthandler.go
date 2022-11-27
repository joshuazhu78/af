package evthandler

import "github.com/joshuazhu78/af/pkg/util"

type EvtHandler interface {
	OnFaceDetected(meta util.GvaMeta) error

	OnDeactivated(inactiveTime uint) error
}
