package main

import (
	"bufio"
	"fmt"
	"golang-/logcatchsys/logconfig"
	"os"
	"path"
	"sync"
	"time"
)

func writeLog(datapath string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	filew, err := os.OpenFile(datapath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("open file error ", err.Error())
		return
	}

	w := bufio.NewWriter(filew)
	for i := 0; i < 20; i++ {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintln(w, "Hello current time is "+timeStr)
		time.Sleep(time.Millisecond * 100)
		w.Flush()
	}
	logBak := time.Now().Format("20060102150405") + ".txt"
	logBak = path.Join(path.Dir(datapath), logBak)
	filew.Close()
	err = os.Rename(datapath, logBak)
	if err != nil {
		fmt.Println("Rename error ", err.Error())
		return
	}
}

func main() {
	v := logconfig.InitVipper()
	configPaths, confres := logconfig.ReadConfig(v, "collectlogs")
	if !confres {
		fmt.Println("config read failed")
		return
	}
	wg := &sync.WaitGroup{}
	for _, configData := range configPaths.([]interface{}) {
		for ckey, cval := range configData.(map[interface{}]interface{}) {
			if ckey == "logpath" && cval != "" {
				wg.Add(1)
				go writeLog(cval.(string), wg)
				continue
			}
		}
	}
	wg.Wait()
}
