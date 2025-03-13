package req

import (
	"net"
	"net/http"
)

func GetIPAddress(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(req.RemoteAddr)
	}
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}
