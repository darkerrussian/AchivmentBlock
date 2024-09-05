package main

import (
	"AchivmentBlock/AchivmentNames"
	"AchivmentBlock/Server"
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

func main() {

	names := AchivmentNames.InitNames()
	db, err := Server.InitDB()
	if err == nil {
		fmt.Println(err)
	}
	Server.UpdateTable(db, names)
	r := mux.NewRouter()
	Server.RegisterRoutes(r)
	//new
	log.Println("Server is listening on port 8090")
	log.Fatal(http.ListenAndServe(":8090", r))

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Received ", message)

}
