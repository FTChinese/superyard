package controller

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
)

const appRoot = `
<!DOCTYPE html>
<html lang="en">

<head>
    <base href="/">
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <link rel="apple-touch-icon" sizes="180x180" href="{{.IconURL}}/apple-touch-icon-180x180.png">
    <link rel="apple-touch-icon" sizes="152x152" href="{{.IconURL}}/apple-touch-icon-152x152.png">
    <link rel="apple-touch-icon" sizes="120x120" href="{{.IconURL}}/apple-touch-icon-120x120.png">
    <link rel="apple-touch-icon" sizes="76x76" href="{{.IconURL}}/apple-touch-icon-76x76.png">
    <link href="{{.IconURL}}/favicon.ico" type="image/x-icon" rel="shortcut icon" />

    <title>FTC CMS</title>
    {{if .Debug}}
    <link rel="stylesheet" href="/style/bootstrap.css">
    <link rel="stylesheet" href="/style/main.css">
    {{else}}
    <link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.3.1/css/bootstrap.min.css" rel="stylesheet">
    {{end}}
</head>

<body>
    <app-root></app-root>
</body>

</html>`

var tmpl = template.Must(template.New("appRoot").Parse(appRoot))

var HomeData = struct {
	IconURL string
	Debug   bool
}{
	IconURL: "http://interactive.ftchinese.com/favicons",
}

func Home(c echo.Context) error {

	c.Response().WriteHeader(http.StatusOK)
	return tmpl.Execute(c.Response().Writer, HomeData)
}
