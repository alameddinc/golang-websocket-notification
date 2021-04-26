package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http server address")

var WS = NewWSServer()

func main() {
	flag.Parse()
	go WS.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(WS, w, r)
	})
	http.HandleFunc("/coklu", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Toplu Gönderim Sağlandı. \n")
		SendMessageAllUsers([]byte("TRT MAM Projesi için online tüm kullanıcılar için websocket notification testi"))
	})

	http.HandleFunc("/tekil/alameddin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Tekil Gönderim Sağlandı. : Alameddin\n")
		SendMessageUser("alameddin", []byte("Alameddin Kullanıcısına Websocket aracılığı ile Notification gönderildi."))
	})
	fmt.Println("Started Server")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
