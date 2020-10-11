package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/render"
)

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, `{"ping":"pong"}`)
	}
}

func envHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp string

		if key := r.URL.Query().Get("key"); key != "" {
			value := os.Getenv(key)
			resp = fmt.Sprintf(`{"%s":"%s"}`, key, value)
		}

		render.PlainText(w, r, resp)
	}
}

func whoamiHandler() http.HandlerFunc {
	type Response struct {
		Hostname      string      `json:"hostname"`
		IP            []net.IP    `json:"ip"`
		Host          string      `json:"host"`
		Headers       http.Header `json:"headers"`
		RemoteAddr    string      `json:"remote_addr"`
		UserAgent     string      `json:"user_agent"`
		ContentType   string      `json:"content_type"`
		ContentLength int64       `json:"content_length"`
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
			Hostname:      hostname,
			IP:            localIPs,
			Host:          r.Host,
			Headers:       r.Header.Clone(),
			RemoteAddr:    remoteAddr,
			UserAgent:     r.UserAgent(),
			ContentType:   "application/json",
			ContentLength: r.ContentLength,
		}

		render.JSON(w, r, resp)
	}
}
