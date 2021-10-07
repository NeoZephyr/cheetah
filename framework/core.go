package framework

import (
	"log"
	"net/http"
)

type Core struct {
	router map[string]Controller
}

func NewCore() *Core {
	return &Core{router: map[string]Controller{}}
}

func (c *Core) Get(url string, controller Controller) {
	c.router[url] = controller
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("begin serve http")
	ctx := NewContext(r, w)

	router := c.router["foo"]

	if router == nil {
		return
	}

	router(ctx)

	log.Println("end serve http")
}
