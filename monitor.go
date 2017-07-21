package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	stdpath "path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/zgs225/hexnuts/client"
	stdmonitor "github.com/zgs225/hexnuts/monitor"
	stdsync "github.com/zgs225/hexnuts/sync"
)

func monitoring(args []string) {
	flags := flag.NewFlagSet("hexnuts monitor", flag.ExitOnError)
	server := flags.String("server", "", "HTTP服务")
	monitorServer := flags.String("monitor.server", "", "TCP监听服务")
	tls := flags.Bool("tls", false, "是否使用TLS")
	in := flags.String("in", "", "监听的目录")
	out := flags.String("out", "", "监听的文件")
	flags.Parse(args)

	if len(*server) == 0 || len(*monitorServer) == 0 || len(*in) == 0 || len(*out) == 0 {
		flags.PrintDefaults()
		os.Exit(1)
	}

	adr, err := net.ResolveTCPAddr("tcp", *monitorServer)
	if err != nil {
		log.Fatalln(err)
	}
	cli := &client.HTTPClient{Addr: *server, TLS: *tls}
	syn := &stdsync.HTTPSyncer{Client: cli, Symbols: make(map[string]string)}
	ctx := context.Background()
	mon := &stdmonitor.Client{Ctx: ctx, Name: getName(), RemoteAddr: adr, TLS: *tls, Syncer: syn, Pairs: getFilePairs(*in, *out)}
	che := make(chan error)

	if err := mon.Dial(); err != nil {
		log.Fatal(err)
	}

	if err := mon.Register(); err != nil {
		log.Fatal(err)
	}

	go handleSignals(mon)

	go func() {
		tick := time.Tick(time.Second)
		for range tick {
			if err := mon.Live(); err != nil {
				che <- err
			}
		}
	}()

	go func() {
		for {
			if err := mon.ReadEvent(); err != nil {
				che <- err
			}
		}
	}()

	go func() {
		if err := mon.SyncPairs(); err != nil {
			che <- err
		}
	}()

	for err := range che {
		if err == io.EOF {
			log.Println("Connection closed")
			os.Exit(0)
		}
		log.Println(err)
	}
}

func handleSignals(mon *stdmonitor.Client) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	<-c
	log.Println("Exiting...")
	mon.Deregister()
	os.Exit(0)
}

func getName() string {
	n, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	t := time.Now().Unix()
	return fmt.Sprintf("%s.%d", n, t)
}

func walkFunc(root, out string, pairs map[string]*stdsync.Pair) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && stdpath.Ext(path) == ".hexnuts" {
			log.Printf("minitoring %s\n", path)
			pair := &stdsync.Pair{Src: path, Dst: getOutputFile(root, path, out)}
			pairs[path] = pair
		}
		return nil
	}
}

func getOutputFile(root, file, outdir string) string {
	ext := stdpath.Ext(file)
	s := len(root) - 1
	e := len(file) - len(ext)
	if ext == ".hexnuts" {
		return stdpath.Join(outdir, file[s:e])
	} else {
		return stdpath.Join(outdir, file[s:])
	}
}

func getFilePairs(in, out string) map[string]*stdsync.Pair {
	v := make(map[string]*stdsync.Pair)
	f := walkFunc(in, out, v)
	if err := filepath.Walk(in, f); err != nil {
		log.Fatalln(err)
	}
	return v
}
