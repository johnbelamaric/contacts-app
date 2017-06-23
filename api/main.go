package main

import (
	"flag"
	"log"
)

func main() {
	var dsn, addr string
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&dsn, "dsn", "", "Mysql database DSN")
	flag.StringVar(&addr, "s", ":80", "Start server on specified address")
	flag.Parse()

	s, err := NewContactServer(verbose, dsn)
	if err != nil {
		log.Println(err)
		panic("failed to create contacts server")
	}

	log.Println("Starting server on %s", addr)
	s.Serve(addr)
}
