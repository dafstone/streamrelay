package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type RTMPServer struct{ URL string }

func main() {
	// parse arguments
	f := flag.NewFlagSet("", flag.ContinueOnError)

	var address, key string
	f.StringVar(&address, "bind", ":1935", "Server bind address and port")
	f.StringVar(&key, "key", "default", "Key for RTMP relay server")
	err := f.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	var filename string
	if f.NArg() >= 1 {
		filename = f.Arg(0)
	} else {
		filename, err = getDefaultFilename()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// read RTMP servers from file
	rtmpServers, err := readServerList(filename)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Ready to stream to:")
	if len(rtmpServers) == 0 {
		log.Println(" (no servers)")
	} else {
		for _, rtmpServer := range rtmpServers {
			log.Println(" -", rtmpServer)
		}
	}
	log.Println()

	// start relay server
	server := NewRelayServer(address, key, rtmpServers)
	server.ListenAndServe()
}

func getDefaultFilename() (string, error) {
	exec, err := os.Executable()
	if err != nil {
		return "", err
	}

	return path.Join(filepath.Dir(exec), "rtmp-servers.txt"), nil
}

func readServerList(filename string) ([]RTMPServer, error) {
	list := []RTMPServer{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			list = append(list, RTMPServer{URL: line})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
