package main

import (
	config2 "ScanWebPath/config"
	"ScanWebPath/internal/pkg/parameter"
	"ScanWebPath/pkg/webPathScan"
	"fmt"
	"time"
)

func main() {
	// get flag
	config := &config2.Config{}
	parameter.Flag(config)
	if config.Url == "" && config.Host == "" {
		fmt.Println("Please input url or host")
		parameter.PrintHelp()
		return
	}
	var url string
	if config.Url != "" {
		url = config.Url
	} else {
		url = "http://" + config.Host
	}
	dict := config.Dict
	goroutineNum := config.GoroutineNum
	sleepTime := time.Duration(config.SleepTime) * time.Second

	// begin time
	beginTime := time.Now()
	fmt.Println("begin time: ", beginTime)
	webPathScanner := webPathScan.NewWebPathScanner(
		webPathScan.WebPathScannerWithUrl(url),
		webPathScan.WebPathScannerWithDictPath(dict),
		webPathScan.WebPathScannerWithGoroutineNum(goroutineNum),
		webPathScan.WebPathScannerWithSleepTime(sleepTime),
	)
	_, err := webPathScanner.Run()
	if err != nil {
		panic(err)
	}
	// end time
	fmt.Println("end time: ", time.Now())
	// time consuming
	fmt.Println("time consuming: ", time.Since(beginTime))
}
