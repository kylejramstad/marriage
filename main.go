package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
	"bytes"
	"log"
	"os"

	"golang.org/x/net/publicsuffix"
)

const (
	unavailable        = "unavailable"
	event              = ""
	IFTTTkey           = ""
)

func main() {
  var s,e = getMessage()
  var NotificationResponse = ""
  var runTime = time.Now().Format("2006-01-02 15:04:05")
  if s != "Nothing available" {
    fmt.Println("Sending Notification")
    NotificationResponse = sendNotice(s,e)
  } else {
    fmt.Println("No notification Sent")
    NotificationResponse = "Nothing Sent"
  }

  // If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("marriage.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("Time: "+runTime+"\nMessage: "+s+"\n"+NotificationResponse+"\n\n")); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}


}

func sendNotice(s string,e error) (string){
  url := "https://maker.ifttt.com/trigger/"+event+"/with/key/"+IFTTTkey
  fmt.Println("URL:>", url)

  var jsonStr = []byte(`{"value1":"`+s+`"}`)
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Set("X-Custom-Header", "myvalue")
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  fmt.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))

  var NotificationResponse = "response Status:" + resp.Status + "\nresponse Body:" + string(body)

  return NotificationResponse
}

func getMessage() (string, error) {
	start := time.Now().Add(2 * 24 * time.Hour).Format("2006-01-02")
	end := time.Now().Add(150 * 24 * time.Hour).Format("2006-01-02")
	fmt.Printf("Checking from %s to %s\n", start, end)
	urlFormat := "https://calendly.com/api/booking/event_types/BFEN3HYGXBSSUYOC/calendar/range?timezone=America%%2FNew_York&diagnostics=false&range_start=%s&range_end=%s&single_use_link_uuid=&embed_domain=projectcupid.cityofnewyork.us&embed_type=Inline"

	j, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return "", err
	}
	client := http.Client{
		Jar: j,
	}
	u, err := url.Parse(fmt.Sprintf(urlFormat, start, end))
	if err != nil {
		return "", err
	}
	//fmt.Println(u.String())
	headers := http.Header{}
	headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	headers.Set("Host", "calendly.com")
	headers.Set("User-Agent", " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1 Safari/605.1.15")
	headers.Set("Accept-Language", "en-us")
	// Don't do this or it returns brotli lol
	//	headers.Set("Accept-Encoding", "gzip, deflate, br")
	req := &http.Request{
		Method: "GET",
		URL:    u,
		Header: headers,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	//fmt.Println(resp)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	responseData := &response{}
	//fmt.Fprintln(os.Stderr, string(data))
	if err := json.Unmarshal(data, responseData); err != nil {
		return "", err
	}
	message := ""
	for _, day := range responseData.Days {
		if day.Status != unavailable {
			message += fmt.Sprintf("*%s has availability!*\n", day.Date)
			for _, spot := range day.Spots {
				message += fmt.Sprintf("%v\n", spot)
			}
		}
	}
	if message == "" {
		return "Nothing available", nil
	}
	return message, nil
}

type response struct {
	InviteePublisherError bool
	Today                 string
	AvailabilityTimezone  string
	Days                  []struct {
		Date          string
		Status        string
		Spots         []interface{}
		InviteeEvents []interface{}
	}
}