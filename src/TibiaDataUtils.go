package main

import (
	"html"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/unicode/norm"
)

// TibiadataDatetimeV3 func
func TibiadataDatetimeV3(date string) string {
	//TODO: Normalization needs to happen above this layer
	date = norm.NFKC.String(date)

	var (
		returnDate time.Time
		err        error
	)

	// If statement to determine if date string is filled or empty
	if date == "" {
		// The string that should be returned is the current timestamp
		returnDate = time.Now()
	} else {
		// timezone use in html: CET/CEST
		loc, _ := time.LoadLocation("Europe/Berlin")

		// format used in datetime on html: Jan 02 2007, 19:20:30 CET
		formatting := "Jan 02 2006, 15:04:05 MST"

		// parsing html in time with location set in loc
		returnDate, err = time.ParseInLocation(formatting, date, loc)

		// parsing html in tiem without loc
		//returnDate, err = time.Parse("Jan 02 2006, 15:04:05 MST", date)

		if err != nil {
			log.Println(err)
		}
	}

	// Return of formatted date and time string to functions..
	return returnDate.UTC().Format(time.RFC3339)
}

// TibiadataHTMLRemoveLinebreaksV3 func
func TibiadataHTMLRemoveLinebreaksV3(data string) string {
	return strings.ReplaceAll(data, "\n", "")
}

var removeUrlRegex = regexp.MustCompile(`<a.*>(.*)<\/a>`)

// TibiadataRemoveURLsV3 func
func TibiadataRemoveURLsV3(data string) string {
	// prepare return value
	var returnData string

	result := removeUrlRegex.FindAllStringSubmatch(data, -1)

	if len(result) > 0 {
		returnData = result[0][1]
	} else {
		returnData = ""
	}

	return returnData
}

// TibiadataStringWorldFormatToTitleV3 func
func TibiadataStringWorldFormatToTitleV3(world string) string {
	return strings.Title(strings.ToLower(world))
}

// TibiadataQueryEscapeStringV3 func - encode string to be correct formatted
func TibiadataQueryEscapeStringV3(data string) string {
	// switching "+" to " "
	data = strings.ReplaceAll(data, "+", " ")

	// encoding string to latin-1
	data, _ = TibiaDataConvertEncodingtoISO88591(data)

	// returning with QueryEscape function
	return url.QueryEscape(data)
}

// TibiadataDateV3 func
func TibiadataDateV3(date string) string {
	// removing weird spacing and comma
	date = TibiaDataSanitizeNbspSpaceString(strings.ReplaceAll(date, ",", ""))

	// var time parser
	var tmpDate time.Time

	// parsing and setting format of return
	switch dateLength := len(date); {
	case dateLength == 5:
		// date that contains special formatting only used in date a world was created
		tmpDate, _ = time.Parse("01/06", date)
		// we need to return earlier as well, since we don't have the day
		return tmpDate.UTC().Format("2006-01")
	case dateLength == 11:
		// dates that contain first 3 letters in month
		tmpDate, _ = time.Parse("Jan 02 2006", date)
	case dateLength > 11:
		// dates that contain month fully written
		tmpDate, _ = time.Parse("January 02 2006", date)
	default:
		log.Printf("Weird format detected: %s", date)
	}

	return tmpDate.UTC().Format("2006-01-02")
}

// TibiadataStringToIntegerV3 func
func TibiadataStringToIntegerV3(data string) int {
	returnData, err := strconv.Atoi(strings.ReplaceAll(data, ",", ""))
	if err != nil {
		log.Printf("[warning] TibiadataStringToIntegerV3: couldn't convert string into int. error: %s", err)
	}

	return returnData
}

var removeHtmlTagRegex = regexp.MustCompile(`(<\/?[a-zA-A]+?[^>]*\/?>)*`)

// match html tag and replace it with ""
func RemoveHtmlTag(in string) string {
	groups := removeHtmlTagRegex.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})

	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}

	return in
}

// TibiaDataConvertEncodingtoISO88591 func - convert string from UTF-8 to latin1 (ISO 8859-1)
func TibiaDataConvertEncodingtoISO88591(data string) (string, error) {
	return charmap.ISO8859_1.NewEncoder().String(data)
}

// TibiaDataConvertEncodingtoUTF8 func - convert string from latin1 (ISO 8859-1) to UTF-8
func TibiaDataConvertEncodingtoUTF8(data io.Reader) io.Reader {
	return norm.NFKC.Reader(charmap.ISO8859_1.NewDecoder().Reader(data))
}

// isEnvExist func - check if environment var is set
func isEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}

	return false
}

// TibiaDataSanitizeEscapedString func - run unescape string on string
func TibiaDataSanitizeEscapedString(data string) string {
	return html.UnescapeString(data)
}

// TibiaDataSanitizeDoubleQuoteString func - replaces double quotes to single quotes in strings
func TibiaDataSanitizeDoubleQuoteString(data string) string {
	return strings.ReplaceAll(data, "\"", "'")
}

// TibiaDataSanitizeNbspSpaceString func - replaces weird \u00A0 string to real space
func TibiaDataSanitizeNbspSpaceString(data string) string {
	return strings.ReplaceAll(data, "\u00A0", " ")
}

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// getEnvAsBool func - read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

/*
// getEnvAsFloat func - read an environment variable into a float64 or return default value
func getEnvAsFloat(name string, defaultVal float64) float64 {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return defaultVal
}

// getEnvAsInt func - read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
*/

// TibiaDataVocationValidator func - return valid vocation string and vocation id
func TibiaDataVocationValidator(vocation string) (string, string) {
	// defining return vars
	var vocationid string

	switch strings.ToLower(vocation) {
	case "none":
		vocationid = "1"
		vocation = "none"
	case "knight", "knights":
		vocationid = "2"
		vocation = "knights"
	case "paladin", "paladins":
		vocationid = "3"
		vocation = "paladins"
	case "sorcerer", "sorcerers":
		vocationid = "4"
		vocation = "sorcerers"
	case "druid", "druids":
		vocationid = "5"
		vocation = "druids"
	default:
		vocationid = "0"
		vocation = "all"
	}

	// returning vars
	return vocation, vocationid
}

// TibiadataGetNewsCategory func - extract news category by newsicon
func TibiadataGetNewsCategory(data string) string {
	switch {
	case strings.Contains(data, "newsicon_cipsoft"):
		return "cipsoft"
	case strings.Contains(data, "newsicon_community"):
		return "community"
	case strings.Contains(data, "newsicon_development"):
		return "development"
	case strings.Contains(data, "newsicon_support"):
		return "support"
	case strings.Contains(data, "newsicon_technical"):
		return "technical"
	default:
		return "unknown"
	}
}

// TibiadataGetNewsType func - extract news type
func TibiadataGetNewsType(data string) string {
	switch data {
	case "News Ticker":
		return "ticker"
	case "Featured Article":
		return "article"
	case "News":
		return "news"
	default:
		return "unknown"
	}
}
