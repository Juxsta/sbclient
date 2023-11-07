package debridlink

import (
	"fmt"

	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
)

func ConvertTorrentToQBT(t *Torrent) (*qbtp.TorrentInfo, error) {
	downloadSpeed := int64(t.DownloadSpeed) // Convert int to int64

	qbtT := &qbtp.TorrentInfo{
		Name:       &t.Name,
		Hash:       &t.HashString,
		Size:       &t.TotalSize,
		Dlspeed:    &downloadSpeed,
		Progress:   func(p int) *float32 { v := float32(p) / 100.0; return &v }(t.DownloadPercent),
		TotalSize:  &t.TotalSize,
		AddedOn:    &t.Created,
		Downloaded: func(size int64, percent int) *int64 { downloaded := size * int64(percent) / 100; return &downloaded }(t.TotalSize, t.DownloadPercent),
		State:      func(s int) *qbtp.TorrentInfoState { v := qbtp.TorrentInfoState(fmt.Sprintf("%d", s)); return &v }(t.Status),
	}

	return qbtT, nil
}
