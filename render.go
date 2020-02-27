package main

import (
	"errors"
	"github.com/FTChinese/go-rest/render"
	rice "github.com/GeertJohan/go.rice"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

// NewRenderer creates a new instance of Renderer based on runtime configuration.
func NewRenderer(config Config) (Renderer, error) {
	// In debug mode, we use pongo's default local file system loader.
	if config.Debug {
		log.Info("Development environment using local file system loader")

		loader := pongo2.MustNewLocalFileSystemLoader("templates")
		set := pongo2.NewSet("local", loader)
		set.Debug = true

		return Renderer{
			config:      config,
			templateSet: set,
		}, nil
	}

	log.Info("Production environment using rice template loader")
	box, err := rice.FindBox("templates")
	if err != nil {
		return Renderer{}, err
	}
	loader := NewRiceTemplateLoader(box)
	set := pongo2.NewSet("rice", loader)
	set.Debug = false

	return Renderer{
		config:      config,
		templateSet: set,
	}, nil
}

func MustNewRenderer(config Config) Renderer {
	r, err := NewRenderer(config)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	return r
}

// Renderer is used to render pong2 templates.
type Renderer struct {
	templateSet *pongo2.TemplateSet // Load templates from filesystem or rice.
	config      Config
}

func (r Renderer) Render(w io.Writer, name string, data interface{}, e echo.Context) error {
	var ctx = pongo2.Context{}

	if data != nil {
		var ok bool
		ctx, ok = data.(pongo2.Context)
		if !ok {
			return errors.New("no pongo2.Context data was passed")
		}
	}

	var t *pongo2.Template
	var err error

	if r.config.Debug {
		// In development the file is loaded from local
		// file system.
		t, err = r.templateSet.FromFile(name)
	} else {
		// In production the file is loaded from rice.
		t, err = r.templateSet.FromCache(name)
	}

	if err != nil {
		return err
	}

	config.Year = time.Now().Year()
	ctx["env"] = r.config

	return t.ExecuteWriter(ctx, w)
}

// RiceTemplateLoader implements pongo2.TemplateLoader to
// loads templates from compiled binary
type RiceTemplateLoader struct {
	box *rice.Box
}

func NewRiceTemplateLoader(box *rice.Box) *RiceTemplateLoader {

	return &RiceTemplateLoader{box: box}
}

func (loader RiceTemplateLoader) Abs(base, name string) string {
	return name
}

func (loader RiceTemplateLoader) Get(path string) (io.Reader, error) {
	return loader.box.Open(path)
}

// RestfulErrorHandler implements echo's HTTPErrorHandler.
func errorHandler(err error, c echo.Context) {
	re, ok := err.(*render.ResponseError)
	if !ok {
		re = render.NewInternalError(err.Error())
	}

	if re.Message == "" {
		re.Message = http.StatusText(re.StatusCode)
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(re.StatusCode)
		} else {
			err = c.JSON(re.StatusCode, re)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
