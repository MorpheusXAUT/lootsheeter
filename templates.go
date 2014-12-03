// templates
package main

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

var (
	templates         = template.Must(template.New("").Funcs(templateFunctions).ParseGlob("web/template/*"))
	templateFunctions = template.FuncMap{
		"FormatFloat":     FormatFloat,
		"IsPositiveFloat": IsPositiveFloat,
	}
)

func FormatFloat(f float64) string {
	fString := humanize.Ftoa(f)

	var formattedFloat string

	if strings.Contains(fString, ".") {
		fInt, err := strconv.ParseInt(fString[:strings.Index(fString, ".")], 10, 64)
		if err != nil {
			return fString
		}

		digitsAfterPoint := len(fString) - strings.Index(fString, ".")
		if digitsAfterPoint > 3 {
			digitsAfterPoint = 3
		}

		formattedFloat = fmt.Sprintf("%s%s", humanize.Comma(fInt), fString[strings.Index(fString, "."):strings.Index(fString, ".")+digitsAfterPoint])
	} else {
		fInt, err := strconv.ParseInt(fString, 10, 64)
		if err != nil {
			return fString
		}

		formattedFloat = fmt.Sprintf("%s.00", humanize.Comma(fInt))
	}

	return formattedFloat
}

func IsPositiveFloat(f float64) bool {
	return f > 0
}
