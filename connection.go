package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var wu = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type connection struct {
	con       *websocket.Conn
	writeChan chan []byte
	data      *Data
}

func myws(w http.ResponseWriter, r *http.Request) {
	websockt_con, err := wu.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{con: websockt_con, writeChan: make(chan []byte, 256), data: &Data{}}
	h.newcon <- c
	go c.writer()
	c.reader()

}
func (c *connection) reader() {
	for {
		_, msg, err := c.con.ReadMessage()
		if err != nil {

			break
		}
		json.Unmarshal(msg, &c.data)
		fmt.Println("reader", c.data)
		switch c.data.Type {
		case "login_req":
			h.r <- c
		case "user_msg":
			date_send, _ := json.Marshal(c.data)
			h.b <- date_send
		case "logout_msg":
			h.u <- c
		default:
			fmt.Print("========default================")
		}
	}
	defer func() {
		c.data.Type = "logout"
		h.u <- c
	}()

}
func (c *connection) writer() {
	for msg_bytes := range c.writeChan {
		json.Unmarshal(msg_bytes, c.data)
		fmt.Println("writer", c.data.Content)
		c.con.WriteMessage(websocket.TextMessage, msg_bytes)
	}

}
