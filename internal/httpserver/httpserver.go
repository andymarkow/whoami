package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andymarkow/whoami/internal/config"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"github.com/urfave/negroni"
)

// Data size units.
const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
)

var healthStatus = 200

type HTTPServer struct {
	server *http.Server
}

type jsonResponse struct {
	RequestID   string              `json:"request_id"`
	Hostname    string              `json:"hostname"`
	IP          []string            `json:"ip"`
	Host        string              `json:"host"`
	URL         string              `json:"url"`
	Params      map[string][]string `json:"params,omitempty"`
	Method      string              `json:"method"`
	Proto       string              `json:"proto"`
	Headers     http.Header         `json:"headers,omitempty"`
	UserAgent   string              `json:"user_agent"`
	RemoteAddr  string              `json:"remote_addr"`
	Environment map[string]string   `json:"environment,omitempty"`
}

// NewHTTPServer creates a new HTTP server with the given configuration.
//
// It takes a pointer to a config.Config struct as a parameter.
// It returns a pointer to an HTTPServer struct.
func NewHTTPServer(cfg *config.Config) *HTTPServer {
	mux := http.NewServeMux()

	mux.Handle("/metrics", useMiddleware(promhttp.Handler(), cfg.AccessLogEnabled, cfg.AccessLogSkipPaths))
	mux.Handle("/health", useMiddleware(healthHandler(), cfg.AccessLogEnabled, cfg.AccessLogSkipPaths))
	mux.Handle("/upload", useMiddleware(uploadHandler(), cfg.AccessLogEnabled, cfg.AccessLogSkipPaths))
	mux.Handle("/data", useMiddleware(dataHandler(), cfg.AccessLogEnabled, cfg.AccessLogSkipPaths))
	mux.Handle("/api/", useMiddleware(apiHandler(), cfg.AccessLogEnabled, cfg.AccessLogSkipPaths))
	mux.Handle("/api", useMiddleware(apiHandler(), cfg.AccessLogEnabled, cfg.AccessLogSkipPaths))
	mux.Handle("/", useMiddleware(whoamiHandler(), cfg.AccessLogEnabled, cfg.AccessLogSkipPaths))

	metricsMW := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			HandlerIDLabel: "path",
			DurationBuckets: []float64{
				0.05, // 50ms
				0.1,  // 100ms
				0.5,  // 500ms
				1,    // 1s
				2.5,  // 2.5s
				5,    // 5s
				10,   // 10s
			},
		}),
	})

	h := std.Handler("", metricsMW, mux)

	srv := &http.Server{
		Addr:              cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:           h,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	return &HTTPServer{
		server: srv,
	}
}

// Start starts the HTTP server.
//
// It does not take any parameters.
// It returns an error.
func (s *HTTPServer) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}

// Shutdown shuts down the HTTP server.
//
// It uses a context with a timeout of 5 seconds to gracefully shutdown the server.
// It returns an error if the server fails to shutdown.
func (s *HTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server.Shutdown: %w", err)
	}

	return nil
}

// skipURLPath checks if the given path should be skipped based on a list of excluded paths.
//
// Parameters:
// - path: The path to be checked.
// - pathExcludes: The list of excluded paths.
//
// Returns:
// - bool: True if the path should be skipped, false otherwise.
func skipURLPath(path string, pathExcludes []string) bool {
	for _, exclude := range pathExcludes {
		if strings.HasPrefix(path, exclude) {
			return true
		}
	}

	return false
}

func useMiddleware(next http.Handler, withAccessLog bool, pathExcludes []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set("X-Request-ID", requestID)
		}

		rw := negroni.NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		if withAccessLog && !skipURLPath(r.URL.Path, pathExcludes) {
			fmt.Printf(
				`{"time":"%s","request_id":"%s","remote_ip":"%s",`+
					`"host":"%s","method":"%s","uri":"%s","status":%d,`+
					`"proto":"%s","user_agent":"%s","duration":"%s","bytes_in":%d,"bytes_out":%d}`+"\n",
				time.Now().Format("2006-01-02T15:04:05.000Z"),
				requestID,
				r.RemoteAddr,
				r.Host,
				r.Method,
				r.RequestURI,
				rw.Status(),
				r.Proto,
				r.UserAgent(),
				time.Since(startTime).String(),
				r.ContentLength,
				rw.Size(),
			)
		}
	})
}

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
			defer r.Body.Close()

			if len(body) == 0 {
				http.Error(w, "post request payload required", http.StatusBadRequest)

				return
			}

			status, err := strconv.Atoi(string(body))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			} else if status < 100 || status > 599 {
				http.Error(w, fmt.Sprintf("invalid status code: %d", status), http.StatusBadRequest)

				return
			}

			healthStatus = status
			w.WriteHeader(http.StatusAccepted)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(healthStatus)
		fmt.Fprintf(w, `{"status":%d}`, healthStatus)
	})
}

func uploadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func dataHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()

		var dataSize int64 = 1

		if queryParams.Has("size") {
			var err error
			dataSize, err = strconv.ParseInt(queryParams.Get("size"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			if dataSize < 0 {
				dataSize *= -1
			} else if dataSize == 0 {
				http.Error(w, errors.New("size cannot be 0").Error(), http.StatusBadRequest)

				return
			}
		}

		if queryParams.Has("unit") {
			switch strings.ToLower(queryParams.Get("unit")) {
			case "kb":
				dataSize *= KB
			case "mb":
				dataSize *= MB
			case "gb":
				dataSize *= GB
			case "tb":
				dataSize *= TB
			}
		}

		content := &contentReader{size: dataSize}

		if queryParams.Has("attachment") {
			w.Header().Set("Content-Disposition", "Attachment")
			http.ServeContent(w, r, "data.txt", time.Now(), content)

			return
		}

		if _, err := io.Copy(w, content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	})
}

func apiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("delay") {
			duration, err := time.ParseDuration(r.URL.Query().Get("delay"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}
			time.Sleep(duration)
		}

		data, err := getWhoamiData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	})
}

func whoamiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("delay") {
			duration, err := time.ParseDuration(r.URL.Query().Get("delay"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}
			time.Sleep(duration)
		}

		data, err := getWhoamiData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		resp := fmt.Sprintf(`RequestID: %s
Hostname: %s
IP: %s
Host: %s
URL: %s
Method: %s
Proto: %s
Params: %v
Headers: %v
UserAgent: %s
RemoteAddr: %s
Environment: %s
`,
			data.RequestID,
			data.Hostname,
			data.IP,
			data.Host,
			data.URL,
			data.Method,
			data.Proto,
			data.Params,
			data.Headers,
			data.UserAgent,
			data.RemoteAddr,
			data.Environment,
		)

		fmt.Fprintln(w, resp)
	})
}

func getWhoamiData(r *http.Request) (*jsonResponse, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("os.Hostname: %w", err)
	}

	var localIPs []string

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

	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	params := make(map[string][]string)
	for k, v := range r.URL.Query() {
		params[k] = v
	}

	requestID := r.Header.Get("X-Request-Id")

	environment := make(map[string]string)

	if r.URL.Query().Has("env") {
		for _, env := range os.Environ() {
			split := strings.SplitN(env, "=", 2)
			environment[split[0]] = split[1]
		}
	}

	return &jsonResponse{
		RequestID:   requestID,
		Hostname:    hostname,
		IP:          localIPs,
		Host:        r.Host,
		URL:         r.RequestURI,
		Params:      params,
		Method:      r.Method,
		Proto:       r.Proto,
		Headers:     r.Header.Clone(),
		UserAgent:   r.UserAgent(),
		RemoteAddr:  remoteAddr,
		Environment: environment,
	}, nil
}
