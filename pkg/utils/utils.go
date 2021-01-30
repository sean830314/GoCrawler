package utils

import (
	"net/http"
)

func GetResponseWithCookie(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "over18", Value: "1"})
	// req.Header.Set("Cookie", "over18=1")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
