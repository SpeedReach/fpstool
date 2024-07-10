package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
)

type YoloV5FigureDetection struct {
	FigureDetector
}

func NewYoloV5Detection() YoloV5FigureDetection {
	return YoloV5FigureDetection{}
}

func (fd YoloV5FigureDetection) Detect(image image.Image) []Detected {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	formFile, err := writer.CreateFormFile("file", "image.png")
	if err != nil {
		panic(err)
	}

	if err = png.Encode(formFile, image); err != nil {
		panic(err)
	}
	writer.Close()
	// Send the POST request to the Flask server
	resp, err := http.Post("http://127.0.0.1:5000/upload", writer.FormDataContentType(), &b)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var detects [][6]float32

	if err = json.Unmarshal(body, &detects); err != nil {
		panic(errors.Join(err, errors.New(string(body))))
	}

	var result = make([]Detected, len(detects))

	for i, detect := range detects {
		var dType DetectedType
		var confidence float32
		if detect[4] > detect[5] {
			dType = DetectedTypeHead
			confidence = detect[4]
		} else {
			dType = DetectedTypeBody
			confidence = detect[5]
		}
		result[i] = Detected{
			Type:       dType,
			X:          int((detect[0] + detect[2]) / 2),
			Y:          int((detect[1] + detect[3]) / 2),
			Confidence: confidence,
		}
	}
	return result
}
