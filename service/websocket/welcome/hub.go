package welcome

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			//注册用户
			h.clients[client] = true
		case client := <-h.unregister:
			//删除注册用户
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send) //释放这个客户send
			}
		case message := <-h.broadcast:

			//广播信息给每个用户
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					//如果case没有执行，也就是，广播信息没有不成功，说明用户已经断开连接了，关闭这个用户占用的chan资源，并且从注册map中删除这个用户
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
