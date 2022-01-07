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
)

// TibiadataHTMLRemoveLinebreaksV3 func
func TibiadataHTMLRemoveLinebreaksV3(data string) string {
	return strings.ReplaceAll(data, "\n", "")
}

// TibiadataRemoveURLsV3 func
func TibiadataRemoveURLsV3(data string) string {
	// prepare return value
	var returnData string

	// Regex to remove URLs
	regex := regexp.MustCompile(`<a.*>(.*)<\/a>`)
	result := regex.FindAllStringSubmatch(data, -1)

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

// TibiadataDatetimeV3 func
func TibiadataDatetimeV3(date string) string {
	var returnDate string

	// If statement to determine if date string is filled or empty
	if date == "" {
		// The string that should be returned is the current timestamp
		returnDate = time.Now().Format(time.RFC3339)
	} else {
		// Converting: Jan 02 2007, 19:20:30 CET -> RFC1123 -> RFC3339

		// regex to exact values..
		regex1 := regexp.MustCompile(`(.*).([0-9][0-9]).([0-9][0-9][0-9][0-9]),.([0-9][0-9]:[0-9][0-9]:[0-9][0-9]).(.*)`)
		subma1 := regex1.FindAllStringSubmatch(date, -1)

		if len(subma1) > 0 {
			// Adding fake-Sun for valid RFC1123 convertion..
			dateDate, err := time.Parse(time.RFC1123, "Sun, "+subma1[0][2]+" "+subma1[0][1]+" "+subma1[0][3]+" "+subma1[0][4]+" "+subma1[0][5])
			if err != nil {
				// log.Fatal(err)
				log.Println(err)
			}

			// Set data to return
			returnDate = dateDate.Format(time.RFC3339)

		} else {
			// Format not defined yet..
			log.Println("Incoming date: " + date)
			log.Println("UNKNOWN FORMAT YET!")

			// Parse the given string to be formatted correct later
			dateDate, err := time.Parse(time.RFC3339, string(date))
			if err != nil {
				log.Fatal(err)
			}

			// Set data to return
			returnDate = dateDate.Format(time.RFC3339)

		}
	}

	// Return of formatted date and time string to functions..
	return returnDate
}

// TibiadataDateV3 func
func TibiadataDateV3(date string) string {
	// use regex to skip weird formatting on "spaces"
	regex1 := regexp.MustCompile(`([a-zA-Z]{3}).*([0-9]{2}).*([0-9]{4})`)
	subma1 := regex1.FindAllStringSubmatch(date, -1)
	date = (subma1[0][1] + " " + subma1[0][2] + " " + subma1[0][3])

	// parsing and setting format of return
	tmpDate, _ := time.Parse("Jan 02 2006", date)
	date = tmpDate.Format("2006-01-02")

	return date
}

// TibiadataStringToIntegerV3 func
func TibiadataStringToIntegerV3(data string) int {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(data, "")
	returnData, _ := strconv.Atoi(processedString)

	// Return of formatted date and time string to functions..
	return returnData
}

// match html tag and replace it with ""
func RemoveHtmlTag(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
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
	data, err := charmap.ISO8859_1.NewEncoder().String(data)
	return data, err
}

// TibiaDataConvertEncodingtoUTF8 func - convert string from latin1 (ISO 8859-1) to UTF-8
func TibiaDataConvertEncodingtoUTF8(data io.Reader) io.Reader {
	return charmap.ISO8859_1.NewDecoder().Reader(data)
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
