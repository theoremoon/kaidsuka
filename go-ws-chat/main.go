package main

import (
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

type Server struct {
	Msg     chan []byte
	Clients map[*Client]struct{}
	Add     chan *Client
	Remove  chan *Client
}

func newServer() *Server {
	return &Server{
		Msg:     make(chan []byte),
		Clients: make(map[*Client]struct{}),
		Add:     make(chan *Client),
		Remove:  make(chan *Client),
	}
}
func (s *Server) Run() {
	for {
		select {
		case m := <-s.Msg:
			for c, _ := range s.Clients {
				c.Msg <- m
			}
		case c := <-s.Add:
			s.Clients[c] = struct{}{}
		case c := <-s.Remove:
			delete(s.Clients, c)
		}
	}
}

type Client struct {
	Msg    chan []byte
	Socket *websocket.Conn
	Srv    *Server
}

func newClient(s *Server, ws *websocket.Conn) *Client {
	c := &Client{
		Msg:    make(chan []byte),
		Socket: ws,
		Srv:    s,
	}
	s.Add <- c
	return c
}
func (c *Client) Close() {
	c.Srv.Remove <- c
	c.Socket.Close()
}
func (c *Client) SendHandler() {
	defer c.Close()

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := c.Socket.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		case m := <-c.Msg:
			c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Socket.WriteMessage(websocket.TextMessage, m); err != nil {
				return
			}
		}
	}
}
func (c *Client) ReceiveHandler() {
	defer c.Close()
	for {
		_, m, err := c.Socket.ReadMessage()
		if err != nil {
			break
		}
		c.Srv.Msg <- m
	}
}

func wsConnect(s *Server, c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := newClient(s, ws)
	go client.ReceiveHandler()
	go client.SendHandler()
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	srv := newServer()
	go srv.Run()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "./public")
	e.GET("/ws", func(c echo.Context) error {
		return wsConnect(srv, c)
	})
	e.Logger.Fatal(e.Start(":" + port))
}
