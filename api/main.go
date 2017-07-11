package main

import (
	"flag"
	"log"
)

func main() {
	var dsn, path, addr, certPath, keyPath, caPath string
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&dsn, "dsn", "", "Mysql database DSN")
	flag.StringVar(&path, "path", "/contacts", "Path for the RESTful API")
	flag.StringVar(&addr, "s", ":80", "Start server on specified address")
	flag.StringVar(&certPath, "cert", "", "Path to server certificate")
	flag.StringVar(&keyPath, "key", "", "Path to server private key")
	flag.StringVar(&caPath, "ca", "", "Path to CA certificate")
	flag.Parse()

	s, err := NewContactServer(verbose, dsn, path)
	if err != nil {
		log.Println(err)
		panic("failed to create contacts server")
	}

	s.Serve(addr, certPath, keyPath, caPath)
}
