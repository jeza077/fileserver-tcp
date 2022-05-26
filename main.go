package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer listener.Close()
	log.Printf("Servidor iniciado en puerto: 5000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("No se puede conectar: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
