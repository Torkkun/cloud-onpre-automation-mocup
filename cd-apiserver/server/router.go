package server

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ServerE struct {
	Engine *echo.Echo
}

func NewServer() *ServerE {
	server := &ServerE{
		Engine: echo.New(),
	}
	server.Routes()

	server.testRoutes()
	return server
}

func (s *ServerE) Routes() {

}

func (s *ServerE) testRoutes() {
	testg := s.Engine.Group("/test")
	testg.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "test ok")
	})
}

func (se *ServerE) Run(address string) {
	err := se.Engine.Start(address)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
