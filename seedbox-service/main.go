package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Juxsta/sbclient/seedbox-service/database"
	qbt "github.com/Juxsta/sbclient/seedbox-service/qbittorrent"
	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
	"github.com/Juxsta/sbclient/seedbox-service/session"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type MyServer struct {
	db     *gorm.DB
	client qbt.Client
	store  *session.Store
}

func main() {
	e := echo.New()

	db := database.InitDB()
	defer db.Close()

	client := qbt.NewClient("")
	store := session.NewSessionStore("localhost:6379", "", 0)

	server := &MyServer{
		db:     db,
		client: *client,
		store:  store,
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			path := c.Request().URL.EscapedPath()

			log.Printf("Received %s request for %s at %s", req.Method, req.URL, path)
			return next(c)
		}
	})

	qbtp.RegisterHandlersWithBaseURL(e, server, "/api/v2")

	doc, err := qbtp.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		log.Fatal(err)
	}

	// Each request is validated against OpenAPI spec before processing
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			httpReq := c.Request()
			route, pathParams, err := router.FindRoute(httpReq)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    httpReq,
				PathParams: pathParams,
				Route:      route,
				Options: &openapi3filter.Options{
					AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
						cookie, err := input.RequestValidationInput.Request.Cookie("SID")
						sid, err := store.GetSession(cookie)
						if sid == "" || err != nil {
							return &ErrAuthenticationFailed{
								errors.New("unauthorized"),
							}
						}
						return nil
					},
				},
			}

			err = openapi3filter.ValidateRequest(c.Request().Context(), requestValidationInput)
			if err != nil {
				if secErr, ok := err.(*openapi3filter.SecurityRequirementsError); ok {
					for _, subErr := range secErr.Errors {
						if _, ok := subErr.(*ErrAuthenticationFailed); ok {
							return c.String(http.StatusForbidden, "Forbidden")
						}
					}
				}

				return c.String(http.StatusBadRequest, err.Error())
			}

			return next(c)
		}
	})

	// Add Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Start Server
	e.Logger.Fatal(e.Start(":8080"))
}

type ErrAuthenticationFailed struct {
	error
}
