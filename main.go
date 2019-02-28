package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/dags-/ping/status"
)

func main() {
	port := flag.Int("port", 8085, "The http server port")
	flag.Parse()

	router := routing.New()
	router.Get("/<server>", handler)
	router.Get("/<server>/<port>", handler)
	server := fasthttp.Server{
		Handler:            router.HandleRequest,
		GetOnly:            true,
		DisableKeepalive:   true,
		ReadBufferSize:     0,
		WriteBufferSize:    0,
		ReadTimeout:        time.Duration(time.Second * 2),
		WriteTimeout:       time.Duration(time.Second * 2),
		MaxRequestBodySize: 0,
	}

	panic(server.ListenAndServe(fmt.Sprintf(":%v", *port)))
}

func handler(c *routing.Context) error {
	c.SetContentType("application/json")
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")
	server := c.Param("server")
	port := parsePort(c.Param("port"))
	stats := status.GetStatus(server, port)
	enc := json.NewEncoder(c)
	enc.SetIndent("", "  ")
	return enc.Encode(stats)
}

func parsePort(port string) int {
	if port == "" {
		return 25565
	}

	i, e := strconv.ParseInt(port, 10, 32)
	if e != nil {
		return 25565
	}

	return int(i)
}
