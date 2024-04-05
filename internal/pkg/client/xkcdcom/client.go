package xkcdcom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/toadharvard/goxkcd/internal/config"
)

type XKSDClient struct {
	cfg    config.XkcdCom
	client *http.Client
}

func New(cfg config.XkcdCom) XKSDClient {
	client := XKSDClient{
		cfg:    cfg,
		client: &http.Client{},
	}
	return client
}

func (c XKSDClient) GetById(id int) (comixInfo ComixInfo, err error) {
	comixInfo = NewComixInfo(id, c.cfg.Language)
	strId := strconv.Itoa(id)
	urlPath, _ := url.JoinPath(c.cfg.URL, strId, "info.0.json")

	resp, err := c.client.Get(urlPath)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code is not 200, recieved: %d", resp.StatusCode)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&comixInfo)
	if err != nil {
		return
	}
	return
}

func (c XKSDClient) GetLastComixNum() (num int, err error) {
	urlPath, _ := url.JoinPath(c.cfg.URL, "info.0.json")
	resp, err := c.client.Get(urlPath)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code is not 200, recieved: %d", resp.StatusCode)
		return
	}

	var info ComixInfo
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return
	}

	num = info.Num
	return
}
