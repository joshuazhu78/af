package evthandler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/free5gc/openapi/AsSessionWithQoS"
	"github.com/free5gc/openapi/models"
	"github.com/joshuazhu78/af/pkg/util"
)

type AfEvtHandler struct {
	client  *AsSessionWithQoS.APIClient
	req     models.AsSessionWithQoSSubscription
	subId   string
	scsAsId string
}

func NewAfEventHandler(nefSvcEndpoint string, nefJson string, ueIpv4 string, scsAsId string) EvtHandler {
	e := &AfEvtHandler{
		client:  util.GetAsSessionWithQoSClient(nefSvcEndpoint),
		req:     models.AsSessionWithQoSSubscription{},
		scsAsId: scsAsId,
	}
	nefJsonObj, _ := ioutil.ReadFile(nefJson)
	json.Unmarshal(nefJsonObj, &e.req)
	e.req.UeIpv4Addr = &ueIpv4
	return e
}

func (e AfEvtHandler) OnFaceDetected(meta util.GvaMeta) error {
	postRequest := e.client.AsSessionWithQoSAPISubscriptionLevelPOSTOperationApi.ScsAsIdSubscriptionsPost(context.Background(), e.scsAsId)
	postRequest = postRequest.AsSessionWithQoSSubscription(e.req)
	_, http_response, err := postRequest.Execute()
	if err != nil {
		log.Printf("http_response: %+v, err: %+v", http_response, err)
		return err
	}
	loc := http_response.Header["Location"][0]
	log.Printf("    %+v created", loc)
	ls := strings.Split(loc, "/")
	e.subId = ls[len(ls)-1]
	return nil
}

func (e AfEvtHandler) OnDeactivated(inactiveTime uint) error {
	log.Printf("No object detected for %d secs", inactiveTime)
	if e.client != nil {
		log.Println("Fire NEF Del")
		deleteRequest := e.client.AsSessionWithQoSAPISubscriptionLevelDELETEOperationApi.ScsAsIdSubscriptionsSubscriptionIdDelete(context.Background(), e.scsAsId, e.subId)
		_, http_response, err := deleteRequest.Execute()
		if err != nil {
			log.Printf("http_response: %+v, err: %+v", http_response, err)
			return err
		}
	}
	return nil
}
