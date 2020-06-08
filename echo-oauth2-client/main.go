package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"flag"

	"golang.org/x/oauth2"

	echo "github.com/labstack/echo/v4"
)

type server struct {
	OAuthConf oauth2.Config
	ServerID  string
}

func (s *server) Start(addr string) error {
	e := echo.New()

	e.GET("/callback", s.callbackHandler())
	return e.Start(addr)
}

func (s *server) callbackHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")

		token, err := s.OAuthConf.Exchange(context.Background(), code)
		if err != nil {
			log.Println(err)
			return err
		}
		client := s.OAuthConf.Client(context.Background(), token)
		r, err := client.Get("https://discord.com/api/users/@me")
		if err != nil {
			log.Println(err)
			return err
		}
		defer r.Body.Close()
		jbody, _ := ioutil.ReadAll(r.Body)

		return c.String(200, string(jbody))
	}
}

func run() error {
	clientID := flag.String("client-id", "", "OAuth2 Client ID")
	clientSecret := flag.String("client-secret", "", "OAuth2 Client Secret")
	authURL := flag.String("authurl", "", "OAuth2 Authentication URL")
	tokenURL := flag.String("tokenurl", "", "OAuth2 Authentication get token URL")
	serverID := flag.String("server-id", "", "Discord Server ID to invite users")
	scope := flag.String("scope", "identify", "")

	flag.Usage = func() {
		fmt.Printf("Usage:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if clientID == nil || *clientID == "" {
		flag.Usage()
		return nil
	}
	if clientSecret == nil || *clientSecret == "" {
		flag.Usage()
		return nil
	}
	if authURL == nil || *authURL == "" {
		flag.Usage()
		return nil
	}
	if tokenURL == nil || *tokenURL == "" {
		flag.Usage()
		return nil
	}
	if serverID == nil || *serverID == "" {
		flag.Usage()
		return nil
	}

	srv := server{
		OAuthConf: oauth2.Config{
			ClientID:     *clientID,
			ClientSecret: *clientSecret,
			Scopes:       strings.Split(*scope, " "),
			Endpoint: oauth2.Endpoint{
				AuthURL:  *authURL,
				TokenURL: *tokenURL,
			},
		},
		ServerID: *serverID,
	}
	url := srv.OAuthConf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	log.Println(url)
	return srv.Start(":5000")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
