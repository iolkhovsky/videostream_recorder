package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"time"
	"videostream_recorder/internal/pkg/videorecorder"
	"videostream_recorder/internal/pkg/webserver"
)

func main() {
	var WebServer webserver.RecorderHttpServer
	var Recorder videorecorder.VideoRecorder
	frameChan := make(chan gocv.Mat, 10)
	WebServer.Init(frameChan)
	WebServer.Start(8000)
	Recorder.Init(frameChan)
	Recorder.Start()

	for true {
		time.Sleep(5 * time.Second)
		fmt.Print("Recording server is working...\n")
	}
}