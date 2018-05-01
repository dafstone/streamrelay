package main

import (
	"io"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
)

type RelayServer struct {
	address     string
	key         string
	rtmpServers []RTMPServer
}

func NewRelayServer(address string, key string, rtmpServers []RTMPServer) *RelayServer {
	return &RelayServer{
		address:     address,
		key:         key,
		rtmpServers: rtmpServers,
	}
}

func (server *RelayServer) ListenAndServe() error {
	pathPattern, err := regexp.Compile("^/relay/(.*)$")
	if err != nil {
		return err
	}

	s := rtmp.Server{
		Addr: server.address,
		HandlePublish: func(conn *rtmp.Conn) {
			defer conn.Close()

			log.Println(conn.URL)

			// parse URL path
			m := pathPattern.FindStringSubmatch(conn.URL.Path)
			if m == nil {
				return
			}

			key := m[1]
			if key != server.key {
				return
			}

			// relay to all RTMP servers
			wg := sync.WaitGroup{}
			queue := pubsub.NewQueue()
			defer queue.Close()

			for _, rtmpServer := range server.rtmpServers {
				wg.Add(1)
				go relayConnection(&wg, rtmpServer, queue)
			}

			// copy from connected client to queue
			err := avutil.CopyFile(queue, conn)
			if err != nil && err != io.EOF {
				log.Println(err)
			}
		},
	}

	ip := "<this IP>"
	log.Printf("Set up your broadcasting software to publish to:\n")
	log.Printf("  rtmp://%s/relay/%s\n", ip, server.key)

	return s.ListenAndServe()
}

func relayConnection(wg *sync.WaitGroup, rtmpServer RTMPServer, queue *pubsub.Queue) {
	defer wg.Done()

	log.Println("Connecting to", rtmpServer.URL)

	dst, err := rtmp.DialTimeout(rtmpServer.URL, 3*time.Second)
	if err != nil {
		log.Println(err)
	}

	log.Println("Streaming to", rtmpServer.URL)
	cursor := queue.Latest()
	err = avutil.CopyFile(dst, cursor)
	if err != nil && err != io.EOF {
		log.Println(rtmpServer.URL, err)
	}
}
