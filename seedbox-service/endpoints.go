package main

import (
	"net/http"

	debridlink "github.com/Juxsta/sbclient/seedbox-service/qbittorrent"
	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
	"github.com/labstack/echo/v4"
)

func (m *MyServer) AppWebapiVersionGet(c echo.Context) error {
	version := "2.9.2"
	return c.String(http.StatusOK, version)
}

func (m *MyServer) AuthLoginPost(c echo.Context, params qbtp.AuthLoginPostParams) error {
	// Here you would verify the params against your user database.
	// As it's not implemented, let's use hardcoded username and password.
	username := "admin"
	password := "adminadmin"

	req := new(qbtp.AuthLoginPostFormdataBody)

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
	preferences := qbtp.Preferences{
		AddTrackers:        nil,
		AddTrackersEnabled: &addTrackersEnabled,
		AltDlLimit:         &altDlLimit,
	}
	return c.JSON(http.StatusOK, preferences)
}

func (m *MyServer) TorrentsInfoPost(c echo.Context) error {
	torrents := []qbtp.TorrentInfo{
		{},
	}
	return c.JSON(http.StatusOK, torrents)
}

func (m *MyServer) TorrentsInfoGet(ctx echo.Context, params qbtp.TorrentsInfoGetParams) error {

	// Call ListTorrents method
	torrents, err := m.client.ListTorrents(&debridlink.ListTorrentsParams{
		Page:    params.Offset,
		PerPage: params.Limit,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Unable to retrieve torrents from DebridLink")
	}
	return ctx.JSON(http.StatusOK, torrents)
}

func (m *MyServer) TorrentsAddPost(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented yet")
}

func (m *MyServer) TorrentsDeletePost(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented yet")
}

// Override due to improper generated code
type TorrentsCreateCategoryPostFormdataBody struct {
	Category string `json:"category" form:"category"`
	SavePath string `json:"savePath" form:"savePath"`
}

func (m *MyServer) TorrentsCategoriesGet(c echo.Context) error {
	return c.JSON(http.StatusOK, getCategories(m))
}

func (m *MyServer) TorrentsCreateCategoryPost(c echo.Context) error {
	req := new(TorrentsCreateCategoryPostFormdataBody)
	if err := c.Bind(req); err != nil {
		return err
	}

	var category Category
	if req.SavePath == "" {
		req.SavePath = req.Category
	}
	m.db.FirstOrCreate(&category, Category{Category: qbtp.Category{Category: req.Category, SavePath: req.SavePath}})

	// Return categories as JSON
	return c.JSON(http.StatusOK, getCategories(m))
}
