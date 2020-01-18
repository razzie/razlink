package razlink

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// Location ...
type Location struct {
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionCode  string `json:"region_code"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
	ZipCode     string `json:"zip_code"`
	TimeZone    string `json:"time_zone"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	MetroCode   string `json:"metro_code"`
}

func (loc *Location) String() string {
	return fmt.Sprintf("%s/%s/%s", loc.CountryName, loc.RegionName, loc.City)
}

// GetLocation returns the geolocation of an IP address
func GetLocation(ip string) (*Location, error) {
	if IsPrivateIP(net.ParseIP(ip)) {
		ip = ""
	}

	url := "https://freegeoip.app/json/" + ip
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var loc Location
	return &loc, json.NewDecoder(resp.Body).Decode(&loc)
}
