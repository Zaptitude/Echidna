package vulnerabilities

import (
	"fmt"
	"html"
	"regexp"

	"github.com/gookit/color"
)

// XSSscan is the XSS module that will scan a php file for XSS vulns
func XSSscan(content []byte) (VulnResults, error) {

	var vulnResults VulnResults

	signatures := []string{
		`echo +.*\$(_GET|_POST|_COOKIE).*`,
		`print +.*\$(_GET|_POST|_COOKIE).*`,
	}

	filter := "esc_|sanitize|isset|int|htmlentities|htmlspecial|intval|wp_strip"
	reFilter, err := regexp.Compile(filter)
	if err != nil {
		return vulnResults, fmt.Errorf("error compiling XSS filter in XSSScan() with error\n%s", err)
	}

	for _, signature := range signatures {
		re, err := regexp.Compile(signature)
		if err != nil {
			return vulnResults, fmt.Errorf("error compiling signature %s in XSSscan() with error\n%s", signature, err)
		}
		matches := re.FindAllString(string(content), -1)
		if matches != nil {
			for _, match := range matches {
				filteredMatches := reFilter.FindAllString(match, 1)
				if len(filteredMatches) != 0 {
					continue
				} else {
					match := html.UnescapeString(match)
					color.Magenta.Println(match)
					vulnResults.Matches = append(vulnResults.Matches, match)
				}
			}

		}
	}

	return vulnResults, nil

}
