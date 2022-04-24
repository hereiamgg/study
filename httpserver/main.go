package main

import (
	"fmt"
	"log"
	"math/rand"
	"metrics/metrics"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics.Register()
	// NewServerMux
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/images", images)
	mux.HandleFunc("/healthz", healthz)

	// start
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("start http server failed,error:%s\n", err.Error())
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	//接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Printf("Header key: %s, value: %s \n", k, v)
			w.Header().Set(k, vv)
		}
	}

	//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	os.Setenv("VERSION", "v1.0")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	fmt.Printf("os version: %s \n", version)

	//Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	clientip := ClientIP(r)
	log.Printf("client ip : %s", clientip)
	log.Printf("response:%d", 200)
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// other
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok!")
}

// 随机延时函数
func images(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	randInt := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(randInt))
	w.Write([]byte(fmt.Sprintf("<h1>%d<h1>", randInt)))
}
