package vkclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	VERSION       string = "5.103"
	METHOD_SEARCH string = "https://api.vk.com/method/newsfeed.search"
)

type Client struct {
	httpClient *http.Client
	vkApiKey   string
}

func (c *Client) GetData(search string, offset string) (VKModel, error) {
	fmt.Println("send request")
	var result VKModel
	req, err := http.NewRequest("GET", METHOD_SEARCH, nil)
	if err != nil {
		log.Print(err)
		return result, err
	}

	q := req.URL.Query()
	q.Add("access_token", c.vkApiKey)
	q.Add("v", VERSION)
	q.Add("q", search)
	q.Add("count", "200")
	if offset != "" {
		q.Add("start_from", offset)
	}
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
	response, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error")
		return result, err
	}
	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func NewClient(httpClient *http.Client, vkApiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		vkApiKey:   vkApiKey,
	}
}
