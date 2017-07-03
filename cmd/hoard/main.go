package main

import (
	"fmt"
	"os"

	"os/signal"
	"syscall"

	"code.monax.io/platform/hoard/config"
	"code.monax.io/platform/hoard/core/storage"
	"code.monax.io/platform/hoard/server"
	"github.com/go-kit/kit/log"
	"github.com/jawher/mow.cli"
)

func main() {
	hoardApp := cli.App("hoard",
		"A content-addressed deterministically encrypted blob storage system")
	listenURL := hoardApp.StringOpt("a address", config.DefaultListenAddress,
		"local address for hoard to listen on encoded as a URL with the "+
			"network protocol as the scheme, for example 'tcp://localhost:54192' "+
			"or 'unix:///tmp/hoard.sock'")

	logging := hoardApp.BoolOpt("l logging", false,
		"Whether to emit any operational logging")
	// This string spec is parsed by mow.cli and has actual semantic significance
	// around optionality and ordering of options and arguments
	hoardApp.Spec = "[--address=<address to listen on>] [--logging]"

	hoardApp.Action = func() {
		var logger log.Logger
		if *logging {
			logger = log.NewLogfmtLogger(os.Stderr)
		}
		serv := server.New(*listenURL, storage.NewMemoryStore(), logger)
		// Catch interrupt etc
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)
		go func(c chan os.Signal) {
			sig := <-c
			printf("\nCaught %s signal: shutting down...", sig)
			// Make sure we clean up
			serv.Stop()
			os.Exit(0)
		}(signalCh)

		printf("Starting hoard daemon on %s...", *listenURL)
		err := serv.Serve()
		if err != nil {
			fatalf("Could not start hoard server: %s", err)
		}
	}

	hoardApp.Run(os.Args)
}

func printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
