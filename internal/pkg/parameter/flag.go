package parameter

import (
	"ScanWebPath/config"
	"flag"
)

func banner() {
	banner := `
 _______                  _______  _______  _______  _______  _       
(  ___  )       |\     /|(  ____ \(  ____ \(  ____ \(  ___  )( (    /|
| (   ) |       | )   ( || (    \/| (    \/| (    \/| (   ) ||  \  ( |
| (___) | _____ | | _ | || (__    | (_____ | |      | (___) ||   \ | |
|  ___  |(_____)| |( )| ||  __)   (_____  )| |      |  ___  || (\ \) |
| (   ) |       | || || || (            ) || |      | (   ) || | \   |
| )   ( |       | () () || )      /\____) || (____/\| )   ( || )  \  |
|/     \|       (_______)|/       \_______)(_______/|/     \||/    )_)
Author: AU9U5T    Version: ` + config.Version + `
`
	println(banner)
}

func Flag(config *config.Config) {
	banner()
	flag.StringVar(&config.Host, "h", "", "IP address of the host you want to scan,for example: 192.168.11.11")
	flag.IntVar(&config.Port, "p", 80, "Port of the host you want to scan,for example: 80")
	flag.StringVar(&config.Url, "u", "", "for example: http://www.baidu.com")
	flag.StringVar(&config.Dict, "d", "../dict", "for example: /usr/share/wordlists/dirb/common.txt")
	flag.IntVar(&config.GoroutineNum, "g", 10, "Set goroutine nums")
	flag.IntVar(&config.SleepTime, "s", 1, "Set sleep time")
	flag.Parse()
}

func PrintHelp() {
	flag.Usage()
}
