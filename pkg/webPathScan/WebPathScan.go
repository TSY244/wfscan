package webPathScan

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/panjf2000/ants"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// MAXGOROUTINENUM Represents the maximum allowed number of concurrent coroutines.
	MAXGOROUTINENUM = 20

	// DEFAULTGOROUTINENUM Is the default number of coroutines to use when there is no specific configuration.
	DEFAULTGOROUTINENUM = 10

	// DEFAULTDICT default dict is a directory
	DEFAULTDICT = "../dict"

	// DEFAULTSLEEPTIME
	DEFAULTSLEEPTIME = 0
)

type WebPathScanner struct {
	dict         string
	url          string
	goroutineNum int
	sleepTime    time.Duration
}

func NewWebPathScanner(fs ...WebPathScannerAttrFunc) *WebPathScanner {
	scanner := &WebPathScanner{
		sleepTime:    DEFAULTSLEEPTIME,
		goroutineNum: DEFAULTGOROUTINENUM,
		dict:         DEFAULTDICT,
	}
	WebPathScannerAttrFuncs(fs).Apply(scanner)
	return scanner
}

func (scanner *WebPathScanner) SetDict(dict string) {
	scanner.dict = dict
}

func (scanner *WebPathScanner) SetUrl(url string) {
	scanner.url = url
}

func (scanner *WebPathScanner) SetGoroutineNum(num int) {
	scanner.goroutineNum = num
}

func (scanner *WebPathScanner) SetSleepTime(sleepTime time.Duration) {
	scanner.sleepTime = sleepTime
}

func colorPrint(url, webPath string, length int64, status string) {

	switch status {
	case "200 OK":
		fmt.Printf("\033[32m%s length:%d status:%s \n\033[0m", url+webPath, length, status)
	case "301 Moved Permanently", "302 Found":
		fmt.Printf("\033[34m%s length:%d status:%s \n\033[0m", url+webPath, length, status)
	case "403 Forbidden":
		fmt.Printf("\033[31m%s length:%d status:%s \n\033[0m", url+webPath, length, status)
	default:
		fmt.Printf("\033[37m%s length:%d status:%s \n\033[0m", url+webPath, length, status)
	}
}

func sendRequest(method, url, webPath string) error {
	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" {
		return errors.New("method not support")
	}

	req, err := http.NewRequest(method, url+webPath, nil)
	if err != nil {
		return err
	}

	clinet := &http.Client{}
	resp, err := clinet.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	length := resp.ContentLength
	status := resp.Status
	colorPrint(url, webPath, length, status)
	return nil

}
func logError(url, webPath, err string) {
	logFile, _ := os.OpenFile("../log/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer logFile.Close()
	logFile.WriteString(fmt.Sprintf("%s get error: %s\n", url+webPath, err))

}

func worker(url, webPath string, count *int64, lock *sync.Mutex) {
	webPath = strings.Replace(webPath, "\r", "", -1)
	if !strings.HasPrefix(webPath, "/") {
		webPath = "/" + webPath
	}
	// http
	err := sendRequest("GET", url, webPath)
	if err != nil {
		logError(url, webPath, err.Error())
	}

	// post
	err = sendRequest("POST", url, webPath)
	if err != nil {
		logError(url, webPath, err.Error())
	}

	if lock.TryLock() {
		*count++
		lock.Unlock()
	}

}

func (scanner *WebPathScanner) getFileList() ([]string, error) {
	s, err := os.Stat(scanner.dict)
	if os.IsNotExist(err) {
		return nil, errors.New("dict not exit")
	}

	fileList := make([]string, 0)
	if s.IsDir() {
		fmt.Println(scanner.dict, " is a dir")
		fmt.Println("traversal dir...")
		fileInfo, err := os.ReadDir(scanner.dict)
		if err != nil {
			return nil, errors.New("read dict error")
		}
		for _, file := range fileInfo {
			fileName := file.Name()
			filePath := scanner.dict + "/" + fileName
			f, err := os.Stat(filePath)
			if err != nil {
				return nil, err
			}
			if !f.IsDir() {
				fileList = append(fileList, filePath)
			} else {
				fmt.Println(filePath, " is a directory, Skip!!!")
			}
		}
	} else {
		fmt.Println(scanner.dict, " is a file")
		fileList = append(fileList, scanner.dict)
	}
	fmt.Println("find file list: ", fileList)
	return fileList, nil
}

func (scanner *WebPathScanner) getLineData(file string) ([]string, error) {
	filePoint, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer filePoint.Close()

	lineData := make([]string, 0)
	fileReader := bufio.NewReader(filePoint)
	for {
		if Data, err := fileReader.ReadString('\n'); err != nil {
			if err == io.EOF { // The last line doesn't \n
				if len(Data) > 0 {
					lineData = append(lineData, Data)
				}
				break
			}
		} else {
			lineData = append(lineData, strings.TrimRight(Data, "\n"))
		}
	}
	return lineData, nil
}

func (scanner *WebPathScanner) scanFile() {
	lineData, err := scanner.getLineData(scanner.dict)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var count int64
	pool, _ := ants.NewPool(scanner.goroutineNum)
	defer pool.Release()

	for _, data := range lineData {
		wg.Add(1)
		err = pool.Submit(func() {
			worker(scanner.url, data, &count, &sync.Mutex{})
			wg.Done()
		})
		if err != nil {
			panic(err)
		}
	}
	wg.Wait()
}

func (scanner *WebPathScanner) scanDict() {
	fileList, err := scanner.getFileList()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	var count int64
	pool, _ := ants.NewPool(scanner.goroutineNum)
	defer pool.Release()
	for _, fileName := range fileList {
		lineData, err := scanner.getLineData(fileName)
		if err != nil {
			panic(err)
		}
		for _, data := range lineData {
			wg.Add(1)
			err = pool.Submit(func() {
				worker(scanner.url, data, &count, &sync.Mutex{})
				wg.Done()
			})
			if err != nil {
				panic(err)
			}
		}
	}
	wg.Wait()
	fmt.Println("count: ", count)
}

func (scanner *WebPathScanner) Run() (bool, error) {
	if scanner.dict == "" {
		return false, errors.New("dict is empty")
	}
	if scanner.url == "" {
		return false, errors.New("url is empty")
	}

	file, err := os.Stat(scanner.dict)
	if os.IsNotExist(err) {
		return false, errors.New("dict not exist")
	}
	if file.IsDir() {
		scanner.scanDict()
	} else {
		scanner.scanFile()
	}
	return true, nil
}
