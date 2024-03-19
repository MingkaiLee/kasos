package main

import (
	"flag"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	url          string
	qpsFile      string
	logOutput    string
	logger       *log.Logger
	client       *http.Client
	secondTicker *time.Ticker
)

const timeStep = 15

func init() {
	flag.StringVar(&url, "url", "http://localhost:10188/stress-test", "url of the service")
	flag.StringVar(&qpsFile, "qps", "", "qps file")
	flag.StringVar(&logOutput, "log", "", "log file")
	flag.Parse()
	logFile, err := os.OpenFile(logOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("panic: %v\n", err)
	}
	logger = log.New(logFile, "[INFO]", log.LstdFlags|log.Lshortfile)
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	secondTicker = time.NewTicker(time.Second)
}

func parallelGet(parallel int) {
	for i := 0; i < parallel; i++ {
		go func() {
			for {
				_, err := client.Get(url)
				if err != nil {
					logger.Println(err)
				}
			}
		}()
	}
}

func qpsControl(qpsList []int) {
	for {
		t := <-secondTicker.C
		diff := t.Sub(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
		qps := qpsList[diff/(timeStep*time.Second)]
		go parallelGet(qps)
	}
}

func genQpsList(qpsFile string) []int {
	r := make([]int, 0)
	file, err := os.Open(qpsFile)
	if err != nil {
		log.Panicf("panic: %v\n", err)
	}
	defer file.Close()
	contentBytes, err := io.ReadAll(file)
	if err != nil {
		log.Panicf("panic: %v\n", err)
	}
	content := string(contentBytes)
	cols := strings.Split(strings.Trim(content, "\n"), "\n")
	for idx := range cols {
		qpsStr := strings.Split(cols[idx], "\t")[1]
		qpsVal, _ := strconv.ParseFloat(qpsStr, 64)
		r = append(r, int(math.Round(qpsVal)))
	}

	return r
}

func main() {
	qpsList := genQpsList(qpsFile)
	go qpsControl(qpsList)
}
