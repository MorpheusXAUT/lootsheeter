// templates
package main

import (
	"html/template"
)

var (
	templates = template.Must(template.New("").ParseGlob("web/template/*"))
)
