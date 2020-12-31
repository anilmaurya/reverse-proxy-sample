package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	var payload requestPayloadStruct

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// For debugging request body
	// respBytes, err := json.Marshal(payload)
	// res.Write(respBytes)

	proxy_url := getProxyUrl(payload.ProxyCondition)

	log.Printf("proxy url is: %s, proxy condition is: %s", proxy_url, payload.ProxyCondition)

	url, _ := url.Parse(proxy_url)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)

	return
}

func getProxyUrl(proxyCondition string) string {
	if proxyCondition == "a" {
		return os.Getenv("A_CONDITION_URL")
	}

	if proxyCondition == "b" {
		return os.Getenv("B_CONDITION_URL")
	}

	return os.Getenv("DEFAULT_CONDITION_URL")
}

func main() {
	log.Printf("Starting server")
	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
