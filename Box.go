package box

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var accessToken string
var configFile string

func BoxRequest(method string, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	} else {
		if len(headers) == 0 {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	if accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+accessToken)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	respBytes, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(" >> URL    :", url)
		log.Println(" >> Status :", response.Status)
	}
	return respBytes, nil
}
