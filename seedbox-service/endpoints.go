package main

import (
	"net/http"

	qbt "github.com/Juxsta/sbclient/seedbox-service/qbittorrent"
	"github.com/labstack/echo/v4"
)

func (m *MyServer) AppWebapiVersionGet(c echo.Context) error {
	version := "2.9.2"
	return c.String(http.StatusOK, version)
}

func (m *MyServer) AuthLoginPost(c echo.Context, params qbt.AuthLoginPostParams) error {
	// Here you would verify the params against your user database.
	// As it's not implemented, let's use hardcoded username and password.
	username := "admin"
	password := "adminadmin"

	req := new(qbt.AuthLoginPostFormdataBody)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Username == username && req.Password == password {
		cookie := new(http.Cookie)
		cookie.Name = "SID"
		cookie.Value = "tjZcBehXWvgh4Eb6ilHmlhkFEsd2nGfu" // you should generate a secure session id here
		cookie.HttpOnly = true
		cookie.Secure = true
		c.SetCookie(cookie)
		return c.String(http.StatusOK, "Ok.")
	} else {
		return echo.ErrUnauthorized // returns a 401 Unauthorized response
	}
}

func (m *MyServer) AppPreferencesGet(c echo.Context) error {
	addTrackersEnabled := false // Create a bool variable
	altDlLimit := int64(10240)  // Create an int64 variable
	preferences := qbt.Preferences{
		AddTrackers:        nil,
		AddTrackersEnabled: &addTrackersEnabled,
		AltDlLimit:         &altDlLimit,
	}
	return c.JSON(http.StatusOK, preferences)
}

func (m *MyServer) TorrentsCategoriesGet(c echo.Context) error {
	categories := map[string]qbt.Category{
		"tv": {
			Category: "tv",
			SavePath: "tv",
		},
	}
	return c.JSON(http.StatusOK, categories)
}

func (m *MyServer) TorrentsInfoPost(c echo.Context) error {
	torrents := []qbt.TorrentInfo{
		{},
	}
	return c.JSON(http.StatusOK, torrents)
}

func (m *MyServer) TorrentsInfoGet(ctx echo.Context, params qbt.TorrentsInfoGetParams) error {
	torrents := []qbt.TorrentInfo{
		{},
	}
	return ctx.JSON(http.StatusOK, torrents)
}
