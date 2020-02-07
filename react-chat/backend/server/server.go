package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/theoremoon/kaidsuka/react-chat/backend/model"
	"github.com/theoremoon/kaidsuka/react-chat/backend/resolver"
	"github.com/theoremoon/kaidsuka/react-chat/backend/service"
)

type Server interface {
	Start(addr string) error
}

type server struct {
	e       *echo.Echo
	service service.Service
}

func New(s service.Service) Server {
	e := echo.New()
	return &server{
		e:       e,
		service: s,
	}
}

func (s *server) playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/query")
	return echo.WrapHandler(h)
}

func (s *server) Start(addr string) error {
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://front.web.localhost:1234"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
	}))

	s.e.POST("/login", s.loginHandler())
	s.e.POST("/logout", s.logoutHandler())
	s.e.POST("/register", s.registerHandler())

	s.e.GET("/play", s.playgroundHandler())

	queryHandler := s.queryHandler()
	s.e.POST("/query", queryHandler, s.loginUserMiddleware())
	s.e.GET("/query", queryHandler, s.loginUserMiddleware())
	return s.e.Start(addr)
}

const sessionKey = "token"

func (s *server) loginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		params := new(struct {
			Username string `from:"username"`
		})
		if err := c.Bind(params); err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid request",
			})
		}
		user, err := s.service.LoginUser(params.Username)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "failed to login",
			})
		}
		c.SetCookie(&http.Cookie{
			Name:     sessionKey,
			Value:    strconv.FormatUint(uint64(user.ID), 10),
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Domain:   ".web.localhost",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})
		return c.NoContent(http.StatusOK)
	}
}
func (s *server) logoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     sessionKey,
			Value:    "",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Time{},
		})
		return c.NoContent(http.StatusOK)
	}
}
func (s *server) registerHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		params := new(struct {
			Username string `from:"username"`
		})
		if err := c.Bind(params); err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid request",
			})
		}
		_, err := s.service.RegisterUser(params.Username)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "failed to register",
			})
		}
		return c.NoContent(http.StatusOK)
	}
}

func (s *server) queryHandler() echo.HandlerFunc {
	h := handler.New(resolver.NewExecutableSchema(resolver.Config{
		Resolvers: resolver.New(s.service),
	}))
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	})
	return echo.WrapHandler(h)
}

func (s *server) getLoginUser(c echo.Context) *model.User {
	cookie, err := c.Cookie(sessionKey)
	if err != nil || cookie.Value == "" {
		return nil
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return nil
	}

	user, err := s.service.GetUser(uint32(id))
	if err != nil {
		return nil
	}
	return user
}

func (s *server) loginUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			user := s.getLoginUser(c)
			newCtx := resolver.AttachUser(r.Context(), user)
			c.SetRequest(r.WithContext(newCtx))
			return next(c)
		}
	}
}
