package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"videostream_recorder/internal/pkg/webserver"
)

func main() {
	var WebServer webserver.RecorderHttpServer
	frameChan := make(chan gocv.Mat, 10)
	WebServer.Init(frameChan)
	WebServer.Start(8000)

	ok := true
	for {
		select {
			case img:= <- frameChan:
				fmt.Print("Got new image")
				gocv.IMWrite("gocv.jpg", img)
				ok = true
			default:
				ok = false
		}
	}
	fmt.Print("Ok: ", ok)
}