package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"flag"
)

var(
	api_path = flag.String("api_path","http://localhost:1080","HTTP URL where the SmartPi API lives")
)

type SmartPiData struct {
	Datasets []struct {
		Phases []struct {
			Name   string `json:"name"`
			Phase  int    `json:"phase"`
			Values []struct {
				Data  float64 `json:"data"`
				Info  string  `json:"info"`
				Type  string  `json:"type"`
				Unity string  `json:"unity"`
			} `json:"values"`
		} `json:"phases"`
		Time string `json:"time"`
	} `json:"datasets"`
	Ipaddress       string  `json:"ipaddress"`
	Lat             float64 `json:"lat"`
	Lng             float64 `json:"lng"`
	Name            string  `json:"name"`
	Serial          string  `json:"serial"`
	Softwareversion string  `json:"softwareversion"`
	Time            string  `json:"time"`
}

func main() {

	res, err := http.Get(*api_path)
	if err != nil{
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var dat SmartPiData
	if err := json.Unmarshal(body, &dat); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(dat.Datasets[0].Phases) ; i++ {

		for j := 0; j < len(dat.Datasets[0].Phases[i].Values) ; j++ {
			log.Println(dat.Datasets[0].Phases[i].Values[j].Type)
			log.Println(dat.Datasets[0].Phases[i].Values[j].Data)
		}
	}
}
