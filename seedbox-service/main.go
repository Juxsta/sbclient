package main

import (
	"context"
	"log"
	"net/http"

	qbt "github.com/Juxsta/sbclient/seedbox-service/qbittorrent"
	qbtp "github.com/Juxsta/sbclient/seedbox-service/qbittorrentproxy"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type MyServer struct {
	db     *gorm.DB
	client qbt.Client
}

func main() {
	e := echo.New()

	db := initDB()
	defer db.Close()

	client := qbt.NewClient("xn-1cePhXJ2Jrf4ijMY8eC_bhZcN1ZxeDRX5IWx_plsnZTFqrvKtmRb8Y5jKNl-c")
	server := &MyServer{
		db:     db,
		client: *client,
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

			// For valid requests, add path params and then continue.
			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    httpReq,
				PathParams: pathParams,
				Route:      route,
				Options: &openapi3filter.Options{
					AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
						// Your authentication logic here
						// For example:
						// cookie, err := input.RequestValidationInput.Request.Cookie("SID")
						// if err != nil || !isValidSessionID(cookie.Value) {
						// 	return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
						// }

						return nil
					},
				},
			}

			if err := openapi3filter.ValidateRequest(c.Request().Context(), requestValidationInput); err != nil {
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

func isValidSessionID(sessionID string) bool {
	return true
}
