package xkcdcom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/toadharvard/goxkcd/internal/config"
)

type XKCDClient struct {
	URL      string
	Language config.ISOCode639_1
	Timeout  time.Duration
	client   *http.Client
}

func New(URL string, language config.ISOCode639_1, timeout time.Duration) *XKCDClient {
	client := &XKCDClient{
		URL:      URL,
		Language: language,
		client:   &http.Client{Timeout: timeout},
	}
	return client
}

func (c *XKCDClient) GetByID(ctx context.Context, id int) (comixInfo *ComixInfo, err error) {
	comixInfo = NewComixInfo(c.Language)
	strID := strconv.Itoa(id)
	urlPath, _ := url.JoinPath(c.URL, strID, "info.0.json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
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
	urlPath, err := url.JoinPath(c.URL, "rss.xml")
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
