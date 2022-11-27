// SPDX-FileCopyrightText: 2020-present Intel
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joshuazhu78/af/pkg/evthandler"
	"github.com/joshuazhu78/af/pkg/util"
)

type STATE int

const (
	NORMAL STATE = iota
	CRITICAL
)

func producer(fifoFile string, ch chan []byte) error {

	// Open pipe for read only
	log.Printf("Starting read from fifo %s", fifoFile)
	pipe, err := os.OpenFile(fifoFile, os.O_RDONLY, 0640)
	if err != nil {
		return fmt.Errorf("couldn't open pipe with error: %+v", err)
	}
	defer pipe.Close()

	// Read the content of named pipe
	reader := bufio.NewReader(pipe)
	log.Println("READER >> created")

	// Infinite loop
	for {
		line, err := reader.ReadBytes('\n')
		// Close the pipe once EOF is reached
		if err != nil {
			log.Printf("Reading FIFO err %+v", err)
			return err
		}

		ch <- line
	}
}

func consumer(ch chan []byte, inactiveTime uint, evtHandlers []evthandler.EvtHandler) {
	state := NORMAL
	timer := time.NewTimer(time.Duration(inactiveTime) * time.Second)
	metaChans := make([]chan util.GvaMeta, len(evtHandlers))
	for {
		select {
		case line := <-ch:
			meta := util.GvaMeta{}
			json.Unmarshal(line, &meta)
			if state == NORMAL {
				log.Println("Object detected")
				for i, evtHandler := range evtHandlers {
					metaChans[i] = make(chan util.GvaMeta)
					go evtHandler.OnFaceDetected(metaChans[i])
				}
				state = CRITICAL
			}
			for _, metaChan := range metaChans {
				metaChan <- meta
			}
			timer = time.NewTimer(time.Duration(inactiveTime) * time.Second)
		case <-timer.C:
			if state == CRITICAL {
				log.Printf("No object detected for %d secs", inactiveTime)
				for _, evtHandler := range evtHandlers {
					err := evtHandler.OnDeactivated()
					if err != nil {
						log.Printf("%+v", err)
					}
				}
				state = NORMAL
			}
		}
	}
}

func main() {
	fifoFile := flag.String("fifoFile", "/tmp/output.json", "fifo filename")
	inactiveTime := flag.Uint("inactiveTime", 30, "inactive length before firing NEF delete")
	//nefSvcEndpoint := flag.String("nefSvcEndpoint", "http://172.16.113.107:29512", "NEF service endpoint")
	nefSvcEndpoint := flag.String("nefSvcEndpoint", "", "NEF service endpoint")
	nefJson := flag.String("nefJson", "nef-modify.json", "NEF post json for QoS provisioning")
	ueIpv4 := flag.String("ueIpv4", "172.250.0.1", "UE IPv4 address")
	scsAsId := flag.String("scsAsId", "facedetection", "application ID")
	captureDir := flag.String("captureDir", "./capture", "directory to save captured picture")
	capturePeriod := flag.Uint("capturePeriod", 5, "capture period in seconds")

	flag.Parse()
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	evtHandlers := make([]evthandler.EvtHandler, 0, 2)
	if *nefSvcEndpoint != "" {
		evtHandlers = append(evtHandlers, evthandler.NewAfEventHandler(*nefSvcEndpoint, *nefJson, *ueIpv4, *scsAsId))
	}
	if *captureDir != "" {
		evtHandlers = append(evtHandlers, evthandler.NewCapEventHandler(*captureDir, *capturePeriod))
	}

	ch := make(chan []byte)
	go consumer(ch, *inactiveTime, evtHandlers)
	err := producer(*fifoFile, ch)
	if err != nil {
		fmt.Printf("error: %+v", err)
	}
}
