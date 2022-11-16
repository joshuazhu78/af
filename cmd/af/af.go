// SPDX-FileCopyrightText: 2020-present Intel
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/free5gc/openapi/models"
	"github.com/joshuazhu78/af/pkg/util"
)

type STATE int

const (
	NORMAL STATE = iota
	CRITICAL
)

func producer(fifoFile string, ch chan []byte) error {

	// Open pipe for read only
	fmt.Printf("Starting read from fifo %s\n", fifoFile)
	pipe, err := os.OpenFile(fifoFile, os.O_RDONLY, 0640)
	if err != nil {
		return fmt.Errorf("couldn't open pipe with error: %+v", err)
	}
	defer pipe.Close()

	// Read the content of named pipe
	reader := bufio.NewReader(pipe)
	fmt.Println("READER >> created")

	// Infinite loop
	for {
		line, err := reader.ReadBytes('\n')
		// Close the pipe once EOF is reached
		if err != nil {
			fmt.Println("FINISHED!")
			return err
		}

		ch <- line
	}
}

func consumer(ch chan []byte, inactiveTime uint, nefSvcEndpoint string, nefJson string, ueIpv4 string, scsAsId string) {
	client := util.GetAsSessionWithQoSClient(nefSvcEndpoint)
	req := models.AsSessionWithQoSSubscription{}
	nefJsonObj, _ := ioutil.ReadFile(nefJson)
	json.Unmarshal(nefJsonObj, &req)
	req.UeIpv4Addr = &ueIpv4
	state := NORMAL
	var subId string
	timer := time.NewTimer(time.Duration(inactiveTime) * time.Second)
	for {
		select {
		case line := <-ch:
			meta := util.GvaMeta{}
			json.Unmarshal(line, &meta)
			//fmt.Printf("%+v\n", meta)
			if state == NORMAL {
				t := time.Now()
				fmt.Printf("%s: Object detected=>Fire NEF Post\n", t.Format("2006-01-02 15:04:05"))
				postRequest := client.AsSessionWithQoSAPISubscriptionLevelPOSTOperationApi.ScsAsIdSubscriptionsPost(context.Background(), scsAsId)
				postRequest = postRequest.AsSessionWithQoSSubscription(req)
				_, http_response, err := postRequest.Execute()
				if err != nil {
					fmt.Printf("http_response: %+v, err: %+v\n", http_response, err)
				} else {
					loc := http_response.Header["Location"][0]
					fmt.Printf("    %+v created\n", loc)
					ls := strings.Split(loc, "/")
					subId = ls[len(ls)-1]
				}
				state = CRITICAL
			}
			timer = time.NewTimer(time.Duration(inactiveTime) * time.Second)
		case <-timer.C:
			if state == CRITICAL {
				t := time.Now()
				fmt.Printf("%s: No object detected for %d secs=>Fire NEF Del\n", t.Format("2006-01-02 15:04:05"), inactiveTime)
				deleteRequest := client.AsSessionWithQoSAPISubscriptionLevelDELETEOperationApi.ScsAsIdSubscriptionsSubscriptionIdDelete(context.Background(), scsAsId, subId)
				_, http_response, err := deleteRequest.Execute()
				if err != nil {
					fmt.Printf("http_response: %+v, err: %+v\n", http_response, err)
				}
				state = NORMAL
			}
		}
	}
}

func main() {
	fifoFile := flag.String("fifoFile", "/tmp/output.json", "fifo filename")
	inactiveTime := flag.Uint("inactiveTime", 30, "Inactive length before firing NEF delete")
	nefSvcEndpoint := flag.String("nefSvcEndpoint", "http://172.16.113.107:29512", "NEF service endpoint")
	nefJson := flag.String("nefJson", "nef-modify.json", "NEF post json for QoS provisioning")
	ueIpv4 := flag.String("ueIpv4", "172.250.0.1", "UE IPv4 address")
	scsAsId := flag.String("scsAsId", "facedetection", "Application ID")

	flag.Parse()
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	ch := make(chan []byte)
	go consumer(ch, *inactiveTime, *nefSvcEndpoint, *nefJson, *ueIpv4, *scsAsId)
	err := producer(*fifoFile, ch)
	if err != nil {
		fmt.Printf("error: %+v", err)
	}
}
