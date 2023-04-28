package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiKey string
	port   string
	link   string
)

func init() {
	flag.StringVar(&apiKey, "api", "", "API key")
	flag.StringVar(&port, "port", "8080", "HTTP port number")
	flag.StringVar(&link, "link", "http://localhost/blacklistmonitor/api.php", "API link")
	flag.Parse()
}

func main() {
	// Register Prometheus gauge
	blacklistStatus := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "blacklist_status",
			Help: "Blacklist status of hosts",
		},
		[]string{"host", "blacklist_result"},
	)
	prometheus.MustRegister(blacklistStatus)

	// Define request body
	requestBody := url.Values{}
	requestBody.Set("apiKey", apiKey)
	requestBody.Set("type", "blacklistStatus")
	requestBody.Set("data", "all")


	if apiKey == "" {
		fmt.Println("Error: API key is required.")
		fmt.Println("Usage: go run main.go -api=<API key> [-port=<port number | default :8080>] [-link=<API link | default http://localhost/blacklistmonitor/api.php >]")
		return
	}

	// Make HTTP request
	resp, err := http.Post(link, "application/x-www-form-urlencoded", strings.NewReader(requestBody.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse response JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	// Get blacklist results
	results := response["result"].([]interface{})

	// Loop through results and update Prometheus gauge
	for _, result := range results {
		result := result.(map[string]interface{})
		host := result["host"].(string)

		switch isBlocked := result["isBlocked"].(type) {
		case float64:
			if isBlocked == 0 {
				blacklistStatus.WithLabelValues(host, "").Set(0)
			} else {
				status := result["status"].([]interface{})
				for _, s := range status {
					s := s.([]interface{})
					if s[1] == nil || s[1] == "" {
						continue
					}
					blacklistStatus.WithLabelValues(host, fmt.Sprintf("%v - %v", s[0], s[1])).Set(1)
				}
			}
		case string:
			if isBlocked == "0" {
				blacklistStatus.WithLabelValues(host, "").Set(0)
			} else {
				status := result["status"].([]interface{})
				for _, s := range status {
					s := s.([]interface{})
					if s[1] == nil || s[1] == "" {
						continue
					}
					blacklistStatus.WithLabelValues(host, fmt.Sprintf("%v - %v", s[0], s[1])).Set(1)
				}
			}
		default:
			panic(fmt.Sprintf("Unknown type for isBlocked: %T", isBlocked))
		}
	}

	// Start Prometheus HTTP server
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Starting server on :%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	fmt.Printf("")
}
