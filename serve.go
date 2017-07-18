package main

import (
	"flag"
	"net/http"

	"git.youplus.cc/tiny/hexnuts/server"
	"github.com/Sirupsen/logrus"
)

func serve(args []string) {
	flags := flag.NewFlagSet("server", flag.PanicOnError)
	addr := flags.String("addr", ":5678", "服务监听地址")
	tls := flags.Bool("tls", false, "是否使用TLS")
	certFile := flags.String("cert", "", "Cert文件路径")
	keyFile := flags.String("key", "", "Key文件路径")
	flags.Parse(args)

	c := &server.Configurer{Items: make(map[string]interface{})}
	s := server.Server{Configer: c}
	h := s.MakeHTTPServer()
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	l := logrus.StandardLogger()
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

	l.Fatal(<-ch)
}
