package util

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func GetPublicIP() string {
	resp, err := http.Get("http://ipinfo.io/ip")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	ip, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(ip))
}
