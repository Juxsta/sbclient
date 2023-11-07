package main

import (
	"net/http"
	"time"

	"github.com/Juxsta/sbclient/seedbox-service/database"
	debridlink "github.com/Juxsta/sbclient/seedbox-service/qbittorrent"
	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
	"github.com/Juxsta/sbclient/seedbox-service/session"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (m *MyServer) AppWebapiVersionGet(c echo.Context) error {
	version := "2.9.2"
	return c.String(http.StatusOK, version)
}

func (m *MyServer) AuthLoginPost(c echo.Context, params qbtp.AuthLoginPostParams) error {
	req := new(qbtp.AuthLoginPostFormdataBody)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// find user
	user := &database.User{}
	if err := m.db.Where("username = ?", req.Username).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return echo.ErrUnauthorized // returns a 401 Unauthorized response
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while fetching user")
	}

	// Validate password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Invalid password
		return echo.ErrUnauthorized // returns a 401 Unauthorized response
	}

	// Password validated, generate a new session id
	sessionID, err := session.GenerateSID()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while generating session id")
	}

	// Set session in Redis
	if err := m.store.SetSession(sessionID, user.Username, time.Hour*24); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while saving session")
	}

	// Set session id in cookie
	cookie := new(http.Cookie)
	cookie.Name = "SID"
	cookie.Value = sessionID
	cookie.HttpOnly = true
	cookie.Secure = true
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Ok.")
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

	var category database.Category
	if req.SavePath == "" {
		req.SavePath = req.Category
	}
	m.db.FirstOrCreate(&category, database.Category{Category: qbtp.Category{Category: req.Category, SavePath: req.SavePath}})

	// Return categories as JSON
	return c.JSON(http.StatusOK, getCategories(m))
}
