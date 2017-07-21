package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/zgs225/hexnuts/monitor"
	"github.com/zgs225/hexnuts/server"
)

func serve(args []string) {
	flags := flag.NewFlagSet("hexnuts server", flag.ExitOnError)
	addr := flags.String("addr", ":5678", "HTTP服务地址")
	tls := flags.Bool("tls", false, "是否使用TLS")
	certFile := flags.String("cert", "", "Cert文件路径")
	keyFile := flags.String("key", "", "Key文件路径")
	dumpsFile := flags.String("dumps", "hexnuts.db", "持久化保存文件")
	monitorAddr := flags.String("monitor", ":5679", "TCP监听服务地址")
	flags.Parse(args)

	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	l := logrus.StandardLogger()
	c := &server.Configurer{Items: make(map[string]interface{})}
	pc := loads(*dumpsFile, c, l)
	s := server.Server{Configer: pc}
	h := s.MakeHTTPServer()
	lh := server.LoggerHandlerMiddleware(l)(h)
	ms := &monitor.TCPServer{
		Addr:      *monitorAddr,
		TLS:       *tls,
		Cert:      *certFile,
		Key:       *keyFile,
		Audiences: make(map[string]*monitor.Audience),
		Ch:        make(chan *monitor.Event),
		Logger:    l,
	}
	s.Monitor = ms
	ch := make(chan error)

	go func(ch chan error) {
		if *tls {
			l.Infof("Listening http server on %s with TLS", *addr)
			ch <- http.ListenAndServeTLS(*addr, *certFile, *keyFile, lh)
		} else {
			l.Infof("Listening http server on %s", *addr)
			ch <- http.ListenAndServe(*addr, lh)
		}
	}(ch)

	go func() {
		if *tls {
			l.Infof("Listening monitor server on %s with TLS", *monitorAddr)
		} else {
			l.Infof("Listening monitor server on %s", *monitorAddr)
		}
		ms.ServeLoop()
	}()

	go func() {
		ch := make(chan os.Signal, 3)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-ch
		dumps(*dumpsFile, pc, l)
		l.Warning("Exiting")
		os.Exit(0)
	}()

	go func() {
		tick := time.Tick(1 * time.Minute)
		for range tick {
			dumps(*dumpsFile, pc, l)
		}
	}()

	l.Fatal(<-ch)
}

func dumps(filepath string, c server.PersistentConfiger, logger *logrus.Logger) {
	if c.Dirty() {
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
