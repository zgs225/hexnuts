package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/zgs225/hexnuts/client"
	stdsync "github.com/zgs225/hexnuts/sync"
)

func sync(args []string) {
	flags := flag.NewFlagSet("hexnuts sync", flag.ExitOnError)
	server := flags.String("server", "", "服务地址")
	tls := flags.Bool("tls", false, "使用TLS")
	in := flags.String("in", "", "输入文件")
	out := flags.String("out", "", "输出文件")
	flags.Parse(args)

	if len(*server) == 0 {
		flags.PrintDefaults()
		os.Exit(0)
	}

	cli := &client.HTTPClient{Addr: *server, TLS: *tls}
	syn := &stdsync.HTTPSyncer{Client: cli, Symbols: make(map[string]string)}
	ctx := context.Background()

	var (
		r io.Reader
		w io.Writer
	)

	if len(*in) > 0 {
		f, err := os.Open(*in)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		r = f
	} else {
		r = os.Stdin
	}

	if len(*out) > 0 {
		f, err := os.Create(*out)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		w = f
	} else {
		w = os.Stdin
	}

	if err := syn.Sync(ctx, r, w); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
