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

// TibiaDataDatetimeV3 func
func TibiaDataDatetimeV3(date string) string {
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

// TibiaDataHTMLRemoveLinebreaksV3 func
func TibiaDataHTMLRemoveLinebreaksV3(data string) string {
	return strings.ReplaceAll(data, "\n", "")
}

var removeUrlRegex = regexp.MustCompile(`<a.*>(.*)<\/a>`)

// TibiaDataRemoveURLsV3 func
func TibiaDataRemoveURLsV3(data string) string {
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

// TibiaDataStringWorldFormatToTitleV3 func
func TibiaDataStringWorldFormatToTitleV3(world string) string {
	return strings.Title(strings.ToLower(world))
}

// TibiaDataQueryEscapeStringV3 func - encode string to be correct formatted
func TibiaDataQueryEscapeStringV3(data string) string {
	// switching "+" to " "
	data = strings.ReplaceAll(data, "+", " ")

	// encoding string to latin-1
	data, _ = TibiaDataConvertEncodingtoISO88591(data)

	// returning with QueryEscape function
	return url.QueryEscape(data)
}

// TibiaDataDateV3 func
func TibiaDataDateV3(date string) string {
	// removing weird spacing and comma
	date = TibiaDataSanitizeStrings(strings.ReplaceAll(date, ",", ""))

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
		tmpDate, _ = time.Parse("January 2 2006", date)
	default:
		log.Printf("Weird format detected: %s", date)
	}

	return tmpDate.UTC().Format("2006-01-02")
}

// TibiaDataStringToIntegerV3 converts a string to an int
func TibiaDataStringToIntegerV3(data string) int {
	str := strings.ReplaceAll(data, ",", "")

	returnData, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("[warning] TibiaDataStringToIntegerV3: couldn't convert string into int. error: %s", err)
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

// TibiaDataSanitizeEscapedString func - run unescape string on string
func TibiaDataSanitizeEscapedString(data string) string {
	return html.UnescapeString(data)
}

// TibiaDataSanitizeDoubleQuoteString func - replaces double quotes to single quotes in strings
func TibiaDataSanitizeDoubleQuoteString(data string) string {
	return strings.ReplaceAll(data, "\"", "'")
}

// TibiaDataSanitizeStrings func - replacing various encoded strings to pure html
func TibiaDataSanitizeStrings(data string) string {
	// replaces weird \u00A0 string to real space
	data = strings.ReplaceAll(data, "\u00A0", " ")
	// replaces weird \u0026 string to amp (&)
	data = strings.ReplaceAll(data, "\u0026", "&")
	// returning string unescaped
	return TibiaDataSanitizeEscapedString(data)
}

// TibiaDataSanitize0026String replaces \u0026#39; with '
func TibiaDataSanitize0026String(data string) string {
	return strings.ReplaceAll(data, "\u0026#39;", "'")
}

// isEnvExist func - check if environment var is set
func isEnvExist(key string) (ok bool) {
	_, ok = os.LookupEnv(key)
	return
}

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
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

// TibiaDataConvertValuesWithK func - convert price strings that contain k, kk or more to 3x0
func TibiaDataConvertValuesWithK(data string) int {
	return TibiaDataStringToIntegerV3(strings.ReplaceAll(data, "k", "") + strings.Repeat("000", strings.Count(data, "k")))
}

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

// TibiaDataGetNewsCategory func - extract news category by newsicon
func TibiaDataGetNewsCategory(data string) string {
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

// TibiaDataGetNewsType func - extract news type
func TibiaDataGetNewsType(data string) string {
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
