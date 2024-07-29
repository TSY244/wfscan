package webPathScan

import (
	"strconv"
	"time"
)

type WebPathScannerAttrFunc func(scanner *WebPathScanner)
type WebPathScannerAttrFuncs []WebPathScannerAttrFunc

func (fs WebPathScannerAttrFuncs) Apply(scanner *WebPathScanner) {
	for _, f := range fs {
		f(scanner)
	}
}

func WebPathScannerWithDictPath(dicPath string) WebPathScannerAttrFunc {
	return func(scanner *WebPathScanner) {
		scanner.dict = dicPath
	}
}

func WebPathScannerWithUrl(url string) WebPathScannerAttrFunc {
	return func(scanner *WebPathScanner) {
		scanner.url = url
	}
}

func WebPathScannerWithGoroutineNum(num int) WebPathScannerAttrFunc {
	return func(scanner *WebPathScanner) {
		if MAXGOROUTINENUM < num {
			strMaxPOrt := strconv.Itoa(MAXGOROUTINENUM)
			panic("MaxGoroutineNum < num ,MaxGoroutineNum=" + strMaxPOrt)
		}
		scanner.goroutineNum = num
	}
}
func WebPathScannerWithSleepTime(sleepTime time.Duration) WebPathScannerAttrFunc {
	return func(scanner *WebPathScanner) {
		scanner.sleepTime = sleepTime
	}
}
