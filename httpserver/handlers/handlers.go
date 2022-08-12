package handlers

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Handlers struct{}

func New() *Handlers {
	return &Handlers{}
}

type JSONResponse struct {
	Hostname   string              `json:"hostname"`
	IP         []string            `json:"ip"`
	Host       string              `json:"host"`
	URL        string              `json:"url"`
	Params     map[string][]string `json:"params,omitempty"`
	Method     string              `json:"method"`
	Proto      string              `json:"proto"`
	Headers    http.Header         `json:"headers,omitempty"`
	UserAgent  string              `json:"user_agent"`
	RemoteAddr string              `json:"remote_addr"`
	RequestID  string              `json:"request_id"`
	// IP        []net.IP            `json:"ips"`
}

func getRequestInfo(c echo.Context) (*JSONResponse, error) {
	req := c.Request()
	resp := c.Response()

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	var localIPs []string
	// var localIPs []net.IP

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
			localIPs = append(localIPs, v4.String())
		}
	}

	remoteAddr := req.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = req.RemoteAddr
	}

	params := make(map[string][]string)
	for k, v := range req.URL.Query() {
		params[k] = v
	}

	return &JSONResponse{
		Hostname:   hostname,
		IP:         localIPs,
		Host:       req.Host,
		URL:        req.RequestURI,
		Params:     params,
		Method:     req.Method,
		Proto:      req.Proto,
		Headers:    req.Header.Clone(),
		UserAgent:  req.UserAgent(),
		RemoteAddr: remoteAddr,
		RequestID:  resp.Header().Get(echo.HeaderXRequestID),
	}, nil
}

// Check query params on custom status value
func getStatus(c echo.Context) int {
	statusQuery := c.Request().URL.Query().Get("status")
	if statusQuery == "" {
		return 200
	}
	status, err := strconv.Atoi(statusQuery)
	if err != nil || status < 200 || status > 599 {
		return 400
	}
	return status
}

func (h *Handlers) WhoamiJSON(c echo.Context) error {
	result, err := getRequestInfo(c)
	if err != nil {
		return err
	}
	if len(c.Request().URL.Query().Get("pretty")) > 0 {
		return c.JSONPretty(getStatus(c), result, "  ")
	}
	return c.JSON(getStatus(c), result)
}

func (h *Handlers) WhoamiPlain(c echo.Context) error {
	result, err := getRequestInfo(c)
	if err != nil {
		return err
	}
	return c.String(getStatus(c), strings.Join([]string{
		fmt.Sprintf("Hostname: %s", result.Hostname),
		fmt.Sprintf("IP: %s", strings.Join(result.IP, ",")),
		fmt.Sprintf("Host: %s", result.Host),
		fmt.Sprintf("URL: %s", result.URL),
		fmt.Sprintf("Method: %s", result.Method),
		fmt.Sprintf("Proto: %s", result.Proto),
		fmt.Sprintf("UserAgent: %s", result.UserAgent),
		fmt.Sprintf("RemoteAddr: %s", result.RemoteAddr),
		fmt.Sprintf("RequestID: %s", result.RequestID),
	}, "\n"))
}
