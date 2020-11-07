package videorecorder

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DefaultImageXsz = 640
const DefaultImageYsz = 480
const DefaultFps = 25.0
const DefaultFragmentLength = 10.0
const DefaultRepository = "recorder_video"
const DefaultMaxFragmentsCnt = 10
const DefaultLogPath = "videorecorder.log"

type IVideoRecorder interface {
	Init(SyncChan chan gocv.Mat)
	SetRepo(RepoPath string)
	SetFps(Fps float32)
	SetRecordResolution(xsz int, ysz int)
	SetMaxFragmentLength(length float32)
	Start()
	Stop()
	IsWorking() bool
}

type VideoRecorder struct {
	syncChan chan gocv.Mat
	writer *gocv.VideoWriter
	xsz int
	ysz int
	fps float32
	reporitoryPath string
	recordsList []string
	isWorking bool
	recordingInProcess bool
	fragmentMaxLenSec float32
	curFrameCnt int
	maxFrameCnt int
	currentFragmentId int
	maxFragmentsCnt int
	toSave string
	logpath string
	logfile *os.File
}

func (self *VideoRecorder) Init(SyncChan chan gocv.Mat) {
	self.syncChan = SyncChan
	self.reporitoryPath = DefaultRepository
	self.xsz = DefaultImageXsz
	self.ysz = DefaultImageYsz
	self.fps = DefaultFps
	self.maxFragmentsCnt = DefaultMaxFragmentsCnt
	self.SetMaxFragmentLength(DefaultFragmentLength)
	self.isWorking = false
	self.recordingInProcess = false
	self.recordsList = make([]string, 0, 0)
	self.SetRepo(DefaultRepository)

	self.logpath = DefaultLogPath
	var err error
	self.logfile, err = os.OpenFile(self.logpath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Print("Can't create log file")
	}
	log.SetOutput(self.logfile)
}

func (self *VideoRecorder) SetRepo(RepoPath string) {
	self.reporitoryPath = RepoPath
	if _, err := os.Stat(self.reporitoryPath); os.IsNotExist(err) {
		os.Mkdir(self.reporitoryPath, 0777)
	}
	self.checkFileRepo()
}

func (self *VideoRecorder) checkFileRepo() {
	self.currentFragmentId = 0
	self.recordsList = GetVideoFileList(self.reporitoryPath)
	sort.Strings(self.recordsList)
	overflow := len(self.recordsList) - self.maxFragmentsCnt

	for i := 0; i < overflow; i++ {
		path := self.reporitoryPath + "/" + self.recordsList[i]
		if PathIsCorrect(path) {
			os.Remove(path)
			log.Print("Record deleted: " + self.recordsList[i])
		}
		self.currentFragmentId = GetVideoIdFromPath(self.recordsList[i])
	}

	if self.currentFragmentId > self.maxFragmentsCnt - 1 {
		self.currentFragmentId = self.maxFragmentsCnt - 1
	}
}

func (self *VideoRecorder) SetFps(Fps float32) {
	if self.recordingInProcess {
		self.Stop()
		self.fps = Fps
		self.Start()
	} else {
		self.fps = Fps
	}
}

func (self *VideoRecorder) SetRecordResolution(xsz int, ysz int) {
	if self.recordingInProcess {
		self.Stop()
	}
	self.xsz = xsz
	self.ysz = ysz
	if self.recordingInProcess {
		self.Start()
	}
}

func (self *VideoRecorder) SetMaxFragmentLength(length float32) {
	self.fragmentMaxLenSec = length
	self.maxFrameCnt = int(self.fragmentMaxLenSec * 25)
	self.curFrameCnt = 0
}

func (self *VideoRecorder) Start() {
	self.isWorking = true
	if !self.recordingInProcess {
		self.toSave, _ = self.getNextFileName("")
		self.createVideoFile()
		go self.writingRoutine()
		log.Print("Recording started")
	} else {
		log.Print("Recording already in progress")
	}
}

func (self *VideoRecorder) Stop() {
	if self.recordingInProcess {
		self.closeVideoFile()
		log.Print("Recording is stopped")
	}
	self.isWorking = false
}

func (self *VideoRecorder) IsWorking() bool {
	return self.isWorking
}

func (self *VideoRecorder) createVideoFile() {
	if !self.recordingInProcess {
		var err error
		self.writer, _ = gocv.VideoWriterFile(self.toSave, "MJPG", float64(self.fps), self.xsz, self.ysz, true)
		if err != nil {
			log.Print("error opening video writer device: %v\n", self.toSave)
			return
		} else {
			self.curFrameCnt = 0
		}
		self.recordingInProcess = true
	} else {
		log.Print("Error: Cant create new file while current is not closed")
	}
}

func (self *VideoRecorder) closeVideoFile() {
	if self.isWorking {
		self.writer.Close()
		self.recordingInProcess = false
	} else {
		log.Print("Error: Cant close file, no one is opened")
	}
}

func (self *VideoRecorder) writingRoutine() {
	for self.isWorking {
		frame := <- self.syncChan
		if !frame.Empty() {
			self.addFrameToVideoFile(frame)
			if self.curFrameCnt > self.maxFrameCnt {
				self.curFrameCnt = 0
				self.saveVideofile()
			}
		}
	}
}

func (self *VideoRecorder) addFrameToVideoFile(frame gocv.Mat) {
	if self.recordingInProcess {
		self.writer.Write(frame)
		self.curFrameCnt++
	} else {
		log.Print("Error: Cant record frame, no opened files")
	}
}

func (self *VideoRecorder) saveVideofile() {
	self.closeVideoFile()
	log.Print("Record saved: " + self.toSave)
	self.toSave, _ = self.getNextFileName("")
	self.createVideoFile()
	self.checkFileRepo()
}

func (self *VideoRecorder) getNextFileName(label string) (string, string) {

	newname := ""
	delname := ""
	currentTime := time.Now()
	timestamp := currentTime.Format("2006-01-02 15:04:05")
	timestamp = strings.Replace(timestamp, " ", "-", -1)
	timestamp = strings.Replace(timestamp, ":", "-", -1)

	probe := "-" + label + "_" + strconv.Itoa(self.currentFragmentId) +".avi"
	newname += timestamp + probe

	self.recordsList = GetVideoFileList(self.reporitoryPath)
	for i:=0; i<len(self.recordsList); i++ {
		if strings.Contains(self.recordsList[i], probe) {
			delname = self.recordsList[i]
		}
	}

	self.currentFragmentId++
	if self.currentFragmentId > self.maxFragmentsCnt - 1 {
		self.currentFragmentId = 0
	}
	return self.reporitoryPath + "/" + newname, self.reporitoryPath + "/" + delname
}
