package debridlink

import (
	"fmt"

	qbt "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	httpClient *resty.Client
}

func NewClient(token string) *Client {
	client := resty.New()
	client.SetHeader("Authorization", "Bearer "+token)

	client.BaseURL = "https://debrid-link.com/api/v2"
	return &Client{
		httpClient: client,
	}
}

type ListTorrentsParams struct {
	Page    *int64
	PerPage *int64
}

func (c *Client) ListTorrents(params *ListTorrentsParams) ([]qbt.TorrentInfo, error) {
	request := c.httpClient.R().SetResult(&TorrentListResponse{})

	if params != nil {
		// if params.Ids != "" {
		// 	request.SetQueryParam("ids", params.Ids)
		// }
		if params.Page != nil {
			request.SetQueryParam("page", fmt.Sprintf("%d", params.Page))
		} else {
			request.SetQueryParam("page", "0")
		}
		if params.PerPage != nil {
			request.SetQueryParam("perPage", fmt.Sprintf("%d", params.PerPage))
		} else {
			request.SetQueryParam("perPage", "20")
		}
	} else {
		request.SetQueryParam("page", "0")
		request.SetQueryParam("perPage", "20")
	}

	response, err := request.Get("https://debrid-link.com/api/v2/seedbox/list")

	if err != nil {
		return nil, err
	}

	result := response.Result().(*TorrentListResponse)
	// convert result to qbt.TorrentInfo
	torrents := make([]qbt.TorrentInfo, 0)
	for _, torrent := range result.Value {
		qbtTorrent, err := ConvertTorrentToQBT(&torrent)
		if err != nil {
			return nil, err
		}
		torrents = append(torrents, *qbtTorrent)
	}

	return torrents, nil
}
