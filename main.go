package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/dags-/ping/status"
)

func main() {
	port := flag.Int("port", 8085, "The http server port")
	dial := flag.Int("dial", 250, "Dial timeout")
	con := flag.Int("con", 250, "Connection timeout")
	flag.Parse()

	status.DialTimeout = time.Duration(*dial) * time.Millisecond
	status.ConTimeout = time.Duration(*con) * time.Millisecond

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

	log.Printf("running on: http://127.0.0.1:%v\n", *port)
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
