package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/mcuadros/go-syslog.v2"
)

// RequestKeyMessage used to extract the message from the metadata.
const RequestKeyMessage = "message"

var (
	cliPort  = kingpin.Flag("port", "Port which will receive requests").Short('p').Default(":514").String()
	cliDebug = kingpin.Flag("debug", "Turns on debugging").Short('d').Bool()
)

func main() {
	kingpin.Parse()

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)

	err := server.ListenUDP(*cliPort)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server")

	err = server.Boot()
	if err != nil {
		panic(err)
	}

	go func(channel syslog.LogPartsChannel) {
		for req := range channel {
			if *cliDebug {
				fmt.Println(req)
				continue
			}

			if message, ok := req[RequestKeyMessage]; ok {
				fmt.Println(message)
			}
		}
	}(channel)

	// Handle common process-killing signals so we can gracefully shut down.
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL.
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)

		err = server.Kill()
		if err != nil {
			panic(err)
		}

		// And we're done!
		os.Exit(0)
	}(shutdown)

	server.Wait()
}
