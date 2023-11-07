package debridlink

type File struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	DownloadUrl     string `json:"downloadUrl"`
	Size            int64  `json:"size"`
	DownloadPercent int    `json:"downloadPercent"`
}

type Tracker struct {
	Announce string `json:"announce"`
}

type Torrent struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	HashString      string    `json:"hashString"`
	UploadRatio     float64   `json:"uploadRatio"`
	ServerID        string    `json:"serverId"`
	Wait            bool      `json:"wait"`
	PeersConnected  int       `json:"peersConnected"`
	Status          int       `json:"status"`
	TotalSize       int64     `json:"totalSize"`
	Files           []File    `json:"files"`
	Trackers        []Tracker `json:"trackers"`
	Created         int64     `json:"created"`
	DownloadPercent int       `json:"downloadPercent"`
	DownloadSpeed   int       `json:"downloadSpeed"`
	UploadSpeed     int       `json:"uploadSpeed"`
}

type Pagination struct {
	Page     int `json:"page"`
	Pages    int `json:"pages"`
	Next     int `json:"next"`
	Previous int `json:"previous"`
}

type TorrentListResponse struct {
	Success    bool       `json:"success"`
	Value      []Torrent  `json:"value"`
	Pagination Pagination `json:"pagination"`
}
