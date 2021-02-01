package azureBilling

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/adeturner/observability"
)

// RestClient is exported
type RestClient struct {
	client     http.Client
	intialised bool
	defaults   restDefaults
}

// Init is exported
func (r *RestClient) Init() {

	r.defaults.Init()

	r.client = http.Client{
		Timeout: time.Second * 10,
	}

	r.intialised = true

}

// GET is exported
func (r *RestClient) GET(url string) (htmlData []byte, err error) {

	if !r.intialised {
		r.Init()
	}

	resp, err := r.client.Get(url)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("client.Get: %s", err.Error()))
		return nil, err
	}

	htmlData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	observability.Logger("Info", resp.Status+" : "+url)
	//fmt.Printf("%v\n", resp.Status)
	// fmt.Printf(string(htmlData))

	return htmlData, nil
}

// POST is exported
func (r *RestClient) POST(url string, jsonStr []byte, zip bool) {

	if !r.intialised {
		r.Init()
	}

	var buf bytes.Buffer

	if zip {
		g := gzip.NewWriter(&buf)
		if _, err := g.Write(jsonStr); err != nil {
			observability.Logger("Error", fmt.Sprintf("gzip write: %s", err.Error()))
			return
		}
		if err := g.Close(); err != nil {
			observability.Logger("Error", fmt.Sprintf("gzip close: %s", err.Error()))
			return
		}

	} else {

		buf.Write(jsonStr)
		observability.Logger("Debug", fmt.Sprintf("plain JSON send of length=%d to %s", buf.Len(), url))
	}

	req, err := http.NewRequest("POST", url, &buf)
	if zip {
		req.Header.Set("Content-Encoding", "gzip")
	}
	req.Header.Set("X-Custom-Header", "steve")
	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Add("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		observability.Logger("Fatal", fmt.Sprintf("ERROR POSTing to %s", url))
	}
	defer resp.Body.Close()

	observability.Logger("Debug", fmt.Sprintf("Status: %s, Headers: %s", resp.Status, resp.Header))

	// additional debug if needed
	// body, _ := ioutil.ReadAll(resp.Body)
	// observability.Logger("Info", fmt.Sprintf("Body: %s", string(body)))
}
