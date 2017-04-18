package main

import (
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"
	"log"
	"flag"
	"strconv"
	"fmt"
	"bufio"
	"bytes"
)


var(
	api_path = flag.String("api_path","http://localhost:1080/api/all/all/now","HTTP URL where the SmartPi API lives")
	listenAddress = flag.String("web.listen-address", ":9122", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	site   = flag.String("smartpi.site", "yourSite", "Your site's name in the SmartPi")
	smartPiName   = flag.String("smartpi.name", "yourSmartPi", "Your SmartPi's name")
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


func handleMetricsRequest(w io.Writer, r *http.Request) error {
        getMetrics(w)
	return nil
}


func getMetrics(w io.Writer) {

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
			fmt.Fprintf(w,"smartpi_exporter" + "{site=\""+ *site +"\",sensor=\","+ *smartPiName + ",\"phase=\""+strconv.Itoa(i)+"\",type=\""+ dat.Datasets[0].Phases[i].Values[j].Type + "\"}" + strconv.FormatFloat(dat.Datasets[0].Phases[i].Values[j].Data,'f', -1, 64 )+ "\n" )
		}
	}
}


func errorHandler(f func(io.Writer, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)
		err := f(wr, r)
		wr.Flush()

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = w.Write(buf.Bytes())

		if err != nil {
			log.Println(err)
		}
	}
}


func startServer() {
	fmt.Printf("Starting exporter")
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}


func main() {
	startServer()
}

