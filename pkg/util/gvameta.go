// SPDX-FileCopyrightText: 2020-present Intel
//
// SPDX-License-Identifier: Apache-2.0
//

package util

type Model struct {
	Name string `json:"name"`
}

type Age struct {
	Label uint  `json:"label"`
	Model Model `json:"model"`
}

type BoundingBox struct {
	XMax float64 `json:"x_max"`
	XMin float64 `json:"x_min"`
	YMax float64 `json:"y_max"`
	YMin float64 `json:"y_min"`
}

type Detection struct {
	BoundingBox BoundingBox `json:"bounding_box"`
	Confidence  float64     `json:"confidence"`
	LabelID     int         `json:"label_id"`
}

type Emotion struct {
	Confidence float64 `json:"confidence"`
	Label      string  `json:"label"`
	LabelID    int     `json:"label_id"`
	Model      Model   `json:"model"`
}

type Gender struct {
	Confidence float64 `json:"confidence"`
	Label      string  `json:"label"`
	LabelID    int     `json:"label_id"`
	Model      Model   `json:"model"`
}

type Object struct {
	Age       Age       `json:"age"`
	Detection Detection `json:"detection"`
	Emotion   Emotion   `json:"emotion"`
	Gender    Gender    `json:"gender"`
	H         int       `json:"h"`
	RegionID  int       `json:"region_id"`
	W         int       `json:"w"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
}

type Resolution struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type GvaMeta struct {
	Objects    []Object   `json:"objects"`
	Resolution Resolution `json:"resolution"`
	TimeStamp  uint64     `json:"timestamp"`
}
