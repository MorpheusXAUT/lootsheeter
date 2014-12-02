// utilities
package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

func FormatFloat(f float64) string {
	fString := humanize.Ftoa(f)

	var formattedFloat string

	if strings.Contains(fString, ",") {
		fInt, err := strconv.ParseInt(fString[:strings.Index(fString, ",")], 10, 64)
		if err != nil {
			return fString
		}

		formattedFloat = fmt.Sprintf("%s%s", humanize.Comma(fInt), fString[strings.Index(fString, ","):])
	} else {
		fInt, err := strconv.ParseInt(fString, 10, 64)
		if err != nil {
			return fString
		}

		formattedFloat = humanize.Comma(fInt)
	}

	return formattedFloat
}
