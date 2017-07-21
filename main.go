package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "server":
		serve(os.Args[2:])
	case "sync":
		sync(os.Args[2:])
	case "monitor":
		monitoring(os.Args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:
	hexnuts command [arguments]

Commands:

	server	启动配置服务
	sync	同步配置文件
	monitor	启动监听程序
	`)
	os.Exit(0)
}
