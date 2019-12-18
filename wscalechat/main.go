package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	upgrader = websocket.Upgrader{}
	name     = randomdata.SillyName()
)

type server struct {
	r       *redis.Client
	msg     chan []byte
	clients map[*client]struct{}
	add     chan *client
	remove  chan *client
}

func newServer(r *redis.Client) *server {
	return &server{
		r:       r,
		msg:     make(chan []byte),
		clients: make(map[*client]struct{}),
		add:     make(chan *client),
		remove:  make(chan *client),
	}
}
func (s *server) Run() {
	pubsub := s.r.Subscribe("chat")
	sub := pubsub.Channel()
	for {
		select {
		case m := <-s.msg:
			s.r.Publish("chat", string(m))
			// do not send directly, when redis pub/sub system is available
			// for c, _ := range s.Clients {
			// 	c.Msg <- m
			// }
		case m := <-sub:
			for c, _ := range s.clients {
				c.Msg([]byte(m.Payload))
			}
		case c := <-s.add:
			s.clients[c] = struct{}{}
		case c := <-s.remove:
			delete(s.clients, c)
		}
	}
}
func (s *server) Add(c *client) {
	s.add <- c
}
func (s *server) Remove(c *client) {
	s.remove <- c
}

type client struct {
	msg chan []byte
	ws  *websocket.Conn
	s   *server
}

func newClient(s *server, ws *websocket.Conn) *client {
	c := &client{
		msg: make(chan []byte),
		ws:  ws,
		s:   s,
	}
	s.Add(c)
	return c
}
func (c *client) Close() {
	c.s.Remove(c)
	c.ws.Close()
}
func (c *client) Msg(m []byte) {
	c.msg <- m
}
func (c *client) SendHandler() {
	defer c.Close()

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := c.ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		case m := <-c.msg:
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.ws.WriteMessage(websocket.TextMessage, m); err != nil {
				return
			}
		}
	}
}
func (c *client) ReceiveHandler() {
	defer c.Close()
	for {
		_, m, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.s.msg <- m
	}
}

func wsConnect(s *server, c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := newClient(s, ws)
	go client.ReceiveHandler()
	go client.SendHandler()
	return nil
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	redis_a := os.Getenv("REDIS")
	if redis_a == "" {
		redis_a = ":6379"
	}
	r := redis.NewClient(&redis.Options{
		Addr:       redis_a,
		MaxRetries: 8,
	})
	_, err := r.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	srv := newServer(r)
	go srv.Run()

	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}

	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", map[string]interface{}{
			"Name": name,
		})
	})
	e.GET("/ws", func(c echo.Context) error {
		return wsConnect(srv, c)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
