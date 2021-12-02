package web

import (
	"embed"
	"errors"
	"github.com/FTChinese/go-rest/render"
	"github.com/flosch/pongo2/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
)

// The embed files' path starts with `template/xxx`

//go:embed template
var templates embed.FS

type Config struct {
	Debug   bool
	Version string
	BuiltAt string
}

// NewRenderer creates a new instance of Renderer based on runtime configuration.
func NewRenderer(conf Config) (Renderer, error) {
	// In debug mode, we use pongo's default local file system loader.
	if conf.Debug {
		log.Info("Development environment using local file system loader")

		// Relative to current working directory.
		loader := pongo2.MustNewLocalFileSystemLoader("web/template")
		set := pongo2.NewSet("local", loader)
		set.Debug = true

		return Renderer{
			config:      conf,
			templateSet: set,
		}, nil
	}

	log.Info("Production environment embed.FS template loader")
	loader := NewEmbedFSTemplateLoader(templates)
	set := pongo2.NewSet("rice", loader)
	set.Debug = false

	return Renderer{
		config:      conf,
		templateSet: set,
	}, nil
}

func MustNewRenderer(config Config) Renderer {
	r, err := NewRenderer(config)
	if err != nil {
		log.Fatal(err)
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

	return t.ExecuteWriter(ctx, w)
}

type EmbedFSTemplateLoader struct {
	f embed.FS
}

func NewEmbedFSTemplateLoader(f embed.FS) EmbedFSTemplateLoader {
	return EmbedFSTemplateLoader{
		f: f,
	}
}

func (loader EmbedFSTemplateLoader) Abs(base, name string) string {
	return name
}

// Get a template file. The path is relative to `template`, such as `b2b/home.html`.
// Since go embed.FS starts from the `template` level,
// you have to prefix the path with `template/`
func (loader EmbedFSTemplateLoader) Get(path string) (io.Reader, error) {
	return loader.f.Open("template/" + path)
}

// ErrorHandler implements echo's HTTPErrorHandler.
func ErrorHandler(err error, c echo.Context) {
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
