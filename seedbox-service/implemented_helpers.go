package main

import (
	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
)

func getCategories(m *MyServer) map[string]qbtp.Category {
	var categories []qbtp.Category
	m.db.Find(&categories)

	catMap := make(map[string]qbtp.Category)
	for _, category := range categories {
		catMap[category.Category] = category
	}
	return catMap
}
