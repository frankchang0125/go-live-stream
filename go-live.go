package main

import (
	"flag"
	"net"
	"os"
	
	"github.com/frankchang0125/go-live-stream/rtmp"
	log "github.com/sirupsen/logrus"
)

var (
	rtmpAddr = flag.String("rtmp-addr", ":1935", "RTMP server address:port")
)

func init() {
	switch os.Getenv("GO_ENV") {
	case "dev":
		log.SetLevel(log.DebugLevel)
	}

	flag.Parse()
}

func main() {
	log.Info("Starting RTMP server...")
	
	listener, err := net.Listen("tcp", *rtmpAddr)
	if err != nil {
		log.WithField("err", err).Fatal("Cannot start server.")
		os.Exit(1)
	}
	defer listener.Close()
	
	rtmpServer := rtmp.NewRTMPServer()
	log.Info("RTMP server started, waiting for connections.")

	// Montior incoming new streamers and viewers
	go rtmpServer.Monitor()
	
	for {
		conn, err := listener.Accept()
		
		if err != nil {
			log.WithField("err", err).Error("Cannot accept connection.")
			conn.Close()
			continue
		}
		
		log.Info("Connection accepted.")
		go rtmpServer.HandleRTMPRequest(conn)	
	}
}
