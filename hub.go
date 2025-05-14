package main

import (
	"encoding/json"
	"fmt"
)

type hub struct {
	cons   map[*connection]bool
	b      chan []byte
	newcon chan *connection
	r      chan *connection
	u      chan *connection
}

var user_list = []string{}
var h = hub{
	cons:   make(map[*connection]bool),
	b:      make(chan []byte),
	newcon: make(chan *connection),
	r:      make(chan *connection),
	u:      make(chan *connection),
}

func slice_del(users []string, user string) []string {
	n := len(users)
	if n == 0 {
		return users
	}
	if n == 1 && users[0] == user {
		return []string{}
	}
	new_users := []string{}
	for _, user_i := range users {
		if user_i != user {
			new_users = append(new_users, user_i)
		}
	}
	return new_users

}
func (h *hub) run() {
	for {
		select {
		case c := <-h.newcon:
			fmt.Println("data_register")
			h.cons[c] = true
			fmt.Println("user_list", user_list)
			c.data.Type = "con_resp"
			c.data.Ip = c.con.RemoteAddr().String()
			c.data.User = "server"
			c.data.Content = "connection compelte,update_ip"
			data_send, _ := json.Marshal(c.data)
			c.writeChan <- data_send
		case c := <-h.r:
			fmt.Println("data_register")
			user_list = append(user_list, c.data.Ip+"-"+c.data.User)
			fmt.Println("user_list", user_list)
			c.data.Type = "login_resp"
			c.data.User = "server"
			c.data.Content = "login compelte,update_user_list"
			c.data.UserList = user_list
			data_send, _ := json.Marshal(c.data)
			// c.writeChan <- data_send
			for c := range h.cons {
				select {
				case c.writeChan <- data_send:
				default:
					delete(h.cons, c)
					close(c.writeChan)
				}
			}
		case c := <-h.u:
			fmt.Println("data_unregister")
			user_list = slice_del(user_list, c.data.Ip+"-"+c.data.User)
			fmt.Println("user_list", user_list)
			delete(h.cons, c)
			c.data.Type = "user_update"
			c.data.Ip = "server_IP"
			c.data.User = "server"
			c.data.Content = "someone offline,update user_list"
			// c.data.UserList = user_list
			c.data.UserList = user_list
			data_send, _ := json.Marshal(c.data)
			for c := range h.cons {
				select {
				case c.writeChan <- data_send:
				default:
					delete(h.cons, c)
					close(c.writeChan)
				}
			}
			// h.b <- data_send
		case data := <-h.b:
			fmt.Println("data_broadcast", data)
			for c := range h.cons {
				select {
				case c.writeChan <- data:
				default:
					delete(h.cons, c)
					close(c.writeChan)
				}
			}

		}

	}
}
