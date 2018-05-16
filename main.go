package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/dags-/ping/status"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
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
		MaxConnsPerIP:      3,
		MaxRequestsPerConn: 2,
		MaxRequestBodySize: 0,
	}

	go handleStop()

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

func handleStop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "stop" {
			fmt.Println("Stopping...")
			os.Exit(0)
			break
		}
	}
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
