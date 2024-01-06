package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/konstfish/angler/geoip-api/models"
)

var url string

func init() {
	url = "https://reallyfreegeoip.org/json/"
}

func isValidAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func GetIpInfo(address string) (models.GeoIP, error) {
	if !isValidAddress(address) {
		return models.GeoIP{}, errors.New("Invalid IP address")
	}

	var geoip models.GeoIP

	geoip.Address = address
	geoip.AddressAge = time.Now()

	requestAddress := url + address

	response, err := http.Get(requestAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &geoip); err != nil {
		log.Fatal(err)
	}

	return geoip, err
}
