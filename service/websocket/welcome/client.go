package welcome

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// 允许向对等方写入消息的时间。
	writeWait = 10 * time.Second

	// 读取下一个乒乓消息的时间。
	pongWait = 60 * time.Second

	// ping必须小于pong
	pingPeriod = (pongWait * 9) / 10

	// peer允许的最大消息大小。
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize * 2, //读缓存设置为消息的2倍
	WriteBufferSize: 1024,
}

// 客户端是websocket连接和集线器之间的中间人。
type Client struct {
	hub *Hub

	// websocket连接.
	conn *websocket.Conn

	// 缓冲出站消息的通道。
	send chan []byte
}

//读处理
func (c *Client) readPump() {
	defer func() {
		//连接发生了中断
		//取消注册表
		c.hub.unregister <- c
		//关闭连接资源
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)                                                                        //设置读取的消息最大数
	c.conn.SetReadDeadline(time.Now().Add(pongWait))                                                           //设置读取消息时间
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) // 这只pong响应
	for {
		_, message, err := c.conn.ReadMessage() //读取一个512消息,如果读取超时或者什么的，会err
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		//消息处理
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//广播消息
		c.hub.broadcast <- message
	}
}

//推送消息给用户处理
func (c *Client) writePump() {
	//定时ping，做心跳检测
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		//用户异常的时候释放一些资源
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //设置写信息时间
			if !ok {
				// 房间关闭，了所有人都要走
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			//写消息
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			//定时发送ping给客户
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// websocket处理
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	//注册用户
	client.hub.register <- client

	// 开始用户的读写
	go client.writePump()
	go client.readPump()
}
