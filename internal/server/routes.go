package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tredstart/scrolly/internal/models"
	"github.com/tredstart/scrolly/internal/views"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("819423a698f9ea9ba3577f20993cb0da98a79ea22ce5d6550b65b69fb36fd438"))))
	e.Static("/static", "static")

	e.GET("/", IndexPage)
	e.GET("/text/:id", GetText)
	e.PUT("/text/:id", PutText)
	e.GET("/new", NewText)

	e.GET("/text/update/:id", UpdateArea)

	return e
}

func NewText(c echo.Context) error {
	id := uuid.NewString()
	text := models.Text{Id: id, Text: "No content for now. Click on me to add."}
	sesh, err := session.Get("session", c)
	if err != nil {
		return err
	}

	if id := sesh.Values["id"]; id == nil {

		sesh.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}

		text.Session = uuid.NewString()
		err = models.CreateSession(text.Session)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		sesh.Values["id"] = text.Session
		if err = sesh.Save(c.Request(), c.Response()); err != nil {
			return err
		}
	} else {
		text.Session = id.(string)
	}
	if err := models.CreateText(text); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	texts, err := models.FetchTexts(text.Session)
	if err != nil {
		return err
	}
	return views.IndexPage(texts, text).Render(c.Request().Context(), c.Response())

}

func PutText(c echo.Context) error {
	id := c.Param("id")
	text := c.Request().FormValue("text")
	err := models.UpdateText(id, text)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	t := models.Text{Text: text, Id: id}

	return views.Text(t).Render(c.Request().Context(), c.Response())
}

func UpdateArea(c echo.Context) error {
	id := c.Param("id")
	text, err := models.FetchTextById(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return views.Textarea(text).Render(c.Request().Context(), c.Response())
}

func GetText(c echo.Context) error {
	id := c.Param("id")
	text, err := models.FetchTextById(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return views.Text(text).Render(c.Request().Context(), c.Response())
}

func IndexPage(c echo.Context) error {
	sesh, err := session.Get("session", c)
	if err != nil {
		return err
	}

	if id := sesh.Values["id"]; id == nil {

		sesh.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}

		session_id := uuid.NewString()
		err = models.CreateSession(session_id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		sesh.Values["id"] = session_id

		id := uuid.NewString()
		text := models.Text{Id: id, Session: session_id, Text: "No content for now. Click on me to add."}
		if err := models.CreateText(text); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if err = sesh.Save(c.Request(), c.Response()); err != nil {
			return err
		}

		return views.IndexPage([]models.Text{}, text).Render(c.Request().Context(), c.Response())
	} else {
		texts, err := models.FetchTexts(id.(string))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return views.IndexPage(texts, texts[len(texts)-1]).Render(c.Request().Context(), c.Response())

	}

}
