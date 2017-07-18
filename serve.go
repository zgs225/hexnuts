package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.youplus.cc/tiny/hexnuts/server"
	"github.com/Sirupsen/logrus"
)

func serve(args []string) {
	flags := flag.NewFlagSet("hexnuts server", flag.ExitOnError)
	addr := flags.String("addr", ":5678", "服务监听地址")
	tls := flags.Bool("tls", false, "是否使用TLS")
	certFile := flags.String("cert", "", "Cert文件路径")
	keyFile := flags.String("key", "", "Key文件路径")
	dumpsFile := flags.String("dumps", "dumps.db", "持久化保存文件")
	flags.Parse(args)

	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	l := logrus.StandardLogger()
	c := &server.Configurer{Items: make(map[string]interface{})}
	pc := loads(*dumpsFile, c, l)
	s := server.Server{Configer: pc}
	h := s.MakeHTTPServer()
	lh := server.LoggerHandlerMiddleware(l)(h)
	ch := make(chan error)

	go func(ch chan error) {
		l.Info("Listening on ", *addr)
		if *tls {
			ch <- http.ListenAndServeTLS(*addr, *certFile, *keyFile, lh)
		} else {
			ch <- http.ListenAndServe(*addr, lh)
		}
	}(ch)

	go func() {
		ch := make(chan os.Signal, 3)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-ch
		dumps(*dumpsFile, pc, l)
		l.Warning("Exiting")
		os.Exit(0)
	}()

	go func() {
		tick := time.Tick(30 * time.Second)
		for range tick {
			dumps(*dumpsFile, pc, l)
		}
	}()

	l.Fatal(<-ch)
}

func dumps(filepath string, c server.PersistentConfiger, logger *logrus.Logger) {
	logger.Infof("Dumping to %s...", filepath)
	f, err := os.Create(filepath)
	if err != nil {
		logger.Panicln(err)
	}
	defer f.Close()

	if err := c.Dumps(f); err != nil {
		logger.Panicln(err)
	}
}

func loads(filepath string, c server.PersistentConfiger, logger *logrus.Logger) server.PersistentConfiger {
	f, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return c
		}
		logger.Panicln(err)
	}
	defer f.Close()

	logger.Infof("Loading from %s...", filepath)
	pc, _ := c.Loads(f)
	return pc
}
