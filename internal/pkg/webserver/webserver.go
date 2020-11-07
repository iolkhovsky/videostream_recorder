package webserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gocv.io/x/gocv"
	"io/ioutil"
	"log"
	"net/http"
)

const PutFrameUrl = "/frame"
const NotValidUrlMessage = "404 not found."
const PutMethodName = "PUT"

type IRecorderHttpServer interface {
	Init(SyncChan chan gocv.Mat)
	Start(Port int)
	Stop()
}

type RequestData struct {
	RequestId int
	EncodedImg string
}

type RecorderHttpServer struct {
	port int
	frameChan chan gocv.Mat
	ok bool
}

func DecodeImageFromString(s string) (gocv.Mat, error) {
	var outFrame gocv.Mat
	sDec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return outFrame, err
	}
	outFrame, err = gocv.IMDecode(sDec, gocv.IMReadColor)
	return outFrame, err
}

func (self *RecorderHttpServer) RequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != PutFrameUrl {
		http.Error(w, NotValidUrlMessage, http.StatusNotFound)
		return
	}
	switch r.Method {
	case PutMethodName:
		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		rawBodyText := string(responseData)
		var requestData RequestData
		json.Unmarshal([]byte(rawBodyText), &requestData)

		img, err := DecodeImageFromString(requestData.EncodedImg)

		select {
			case self.frameChan <- img: {
				self.ok = true
				log.Print("Received frame was sent through channel")
			}
			default: {
				self.ok = false
				log.Print("Frame was not sent, channel is busy")
			}
		}
	default:
		fmt.Fprintf(w, "Only PUT request is implemented. Send jpg-encoded image.")
	}
}

func (self *RecorderHttpServer) Init(SyncChan chan gocv.Mat) {
	self.frameChan = SyncChan
	http.HandleFunc("/", self.RequestHandler)
}

func (self *RecorderHttpServer) Start(Port int) {
	self.port = Port
	go http.ListenAndServe(":8000", nil)
}

func (self *RecorderHttpServer) Stop() {

}
