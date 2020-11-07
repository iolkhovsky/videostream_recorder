package videorecorder

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GetVideoFileList(folder string) []string {
	out := make([]string, 0, 0)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		print("<GetVideoFilesCnt>: Invalid path to directory.")
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".avi") {
			out = append(out, f.Name())
		}
	}
	return out
}

func PathIsCorrect(path string) bool {
	res := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		res = false
	}
	return res
}

func GetVideoIdFromPath(name string) int {
	substr := strings.Split(name,"_")
	id, _ := strconv.Atoi(substr[1])
	return id
}