package xkcdcom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/mmcdole/gofeed"
	"github.com/toadharvard/goxkcd/internal/config"
)

type XKCDClient struct {
	cfg    config.XkcdCom
	client *http.Client
}

func New(cfg config.XkcdCom) *XKCDClient {
	client := &XKCDClient{
		cfg:    cfg,
		client: &http.Client{},
	}
	return client
}

func (c *XKCDClient) GetById(id int) (comixInfo *ComixInfo, err error) {
	comixInfo = NewComixInfo(c.cfg.Language)
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
	return
}

func (c *XKCDClient) GetLastComixNum() (guid int, err error) {
	urlPath, err := url.JoinPath(c.cfg.URL, "rss.xml")
	if err != nil {
		return
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(urlPath)
	if err != nil {
		return
	}
	link := feed.Items[0].Link
	guid, err = strconv.Atoi(path.Base(link))
	return
}
