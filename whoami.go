package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func (s *server) whoamiHandler() http.HandlerFunc {

	type Response struct {
		Hostname   string      `json:"hostname"`
		IP         []net.IP    `json:"ip"`
		Host       string      `json:"host"`
		URL        string      `json:"url"`
		Method     string      `json:"method"`
		Headers    http.Header `json:"headers"`
		UserAgent  string      `json:"user_agent"`
		RemoteAddr string      `json:"remote_addr"`
		RequestID  string      `json:"request_id"`
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err)
	}

	var localIPs []net.IP

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				continue
			}
			localIPs = append(localIPs, v4)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}

		resp := &Response{
			Hostname:   hostname,
			IP:         localIPs,
			Host:       r.Host,
			URL:        r.RequestURI,
			Method:     r.Method,
			Headers:    r.Header.Clone(),
			UserAgent:  r.UserAgent(),
			RemoteAddr: remoteAddr,
			RequestID:  middleware.GetReqID(r.Context()),
		}

		render.Status(r, 200)
		render.JSON(w, r, resp)
	}
}
