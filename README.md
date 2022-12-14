<!--
SPDX-FileCopyrightText: 2022 Intel

SPDX-License-Identifier: Apache-2.0

-->

# Application Function

This repository implements a 5G Application Function (AF) that interacts with sample face detection application running as client/server mode over 5G networks and dynamically provisioning the underlying 5G network QoS by calling the 5G Network Exposure Function (NEF) AsSessionWithQoS API. The [face detection](https://github.com/joshuazhu78/dlstreamer/tree/devQoS/samples/gstreamer/gst_launch/face_detection_and_classification) application is implemented using Intel® Deep Learning Streamer (Intel® DL Streamer). Below figure illustrates the overall system architecture:

![af-aether](./docs/images/af-aether-ueransim.svg)

## 5G

This AF has to be tested with 5G end to end system with NEF to open network API for the applications. For 5G, [Aether](https://docs.aetherproject.org/master/index.html) can be used, which contains a light weight 5G core network and a ROC based management plane. Real gNB or gNB emulator can be used to create an end to end system for testing the applications. For testing with a gNB and UE emulator, [UERANSIM](https://github.com/aligungr/UERANSIM) can be used. After successfully setup a 5G network with UE attached to it and got an IP address allocated, you are ready to run the face detection application over the 5G network.

## Application

The [face detection](https://github.com/joshuazhu78/dlstreamer/tree/devQoS/samples/gstreamer/gst_launch/face_detection_and_classification) application needs to be run in client/server mode. With client running at the UE and lively streaming a web camera video over the 5G network using real time RTP protocol. Suppose the server is running at the DN on IP 192.168.250.1 and UE gets IP allocated as 172.250.0.1 by the 5G network. Firstly configure the route table at UE to ensure the traffic towards DN going through the 5G network.

Route configuration at the client/UE side:
```
$ sudo ip route replace 192.168.250.0/24 via 172.250.0.1
```

Then running the face detection server at DN:
```
$ ./face_detection_and_classification.sh port=9001 CPU display-and-json fifo nofps
```

Lastly run the client to stream the video from the USB camera to the DN:
```
$ ./face_detection_and_classification.sh /dev/video0 CPU "host=192.168.250.1 port=9001"
```

### Install prerequisite for AF

AF relies on `v4l-utils` and **imagemagic** `magick` to capture image lively. `v4l-utils` can be installed as:

```
sudo apt update && sudo apt install v4l-utils
```

And `magick` can be installed as:

```
$ wget https://imagemagick.org/archive/binaries/magick
$ chmod a+x magick
$ sudo mv magick /usr/local/bin/
```

### Run the AF at the server

After the face detection application is running successfully on the 5G network, AF can be run as:

```
$ cd cmd/af
$ go run af.go
```

By default the AF will read the FIFO at /tmp/output.json which the face detection server publishes its results.

AF function usage is as below:

```
Usage of af:
  -captureDir string
    	directory to save captured picture (default "./capture")
  -capturePeriod uint
    	capture period in seconds (default 5)
  -fifoFile string
    	fifo filename (default "/tmp/output.json")
  -inactiveTime uint
    	inactive length before firing NEF delete (default 30)
  -nefJson string
    	NEF post json for QoS provisioning (default "nef-modify.json")
  -nefSvcEndpoint string
    	NEF service endpoint
  -scsAsId string
    	application ID (default "facedetection")
  -ueIpv4 string
    	UE IPv4 address (default "172.250.0.1")
```

The AF logic is described as two states machine:

![af_logic](./docs/images/af.svg)

## Sample output

Below video shows a demo of dynamic provisioned 5G QoS driven by face detection application.

[![facedetection_af](./docs/images/fd-af.png)](http://weip-bj.bj.intel.com/facedetection-af.mp4)

