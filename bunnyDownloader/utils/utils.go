package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// PrintStructFields print structure's all exported fields.
func PrintStructFields(s interface{}) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if v.Field(i).CanInterface() {
			value := v.Field(i).Interface()
			fmt.Printf("%s: %v\n", field.Name, value)
		}
	}
}

// MatchOneOf match one of the patterns
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

// FileSize return the file size of the specified path file
func FileSize(filePath string) (int64, bool, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return file.Size(), true, nil
}

// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := MatchOneOf(url, domainPattern)
	if domain != nil {
		return domain[1]
	}
	return ""
}

// LimitLength Handle overly long strings
func LimitLength(s string, length int) string {
	// 0 means unlimited
	if length == 0 {
		return s
	}

	const ELLIPSES = "..."
	str := []rune(s)
	if len(str) > length {
		return string(str[:length-len(ELLIPSES)]) + ELLIPSES
	}
	return s
}

// FileName Converts a string to a valid filename
func FileName(name, ext string, length int) string {
	rep := strings.NewReplacer("\n", " ", "/", " ", "|", "-", ": ", "：", ":", "：", "'", "’")
	name = rep.Replace(name)
	if runtime.GOOS == "windows" {
		rep = strings.NewReplacer("\"", " ", "?", " ", "*", " ", "\\", " ", "<", " ", ">", " ")
		name = rep.Replace(name)
	}
	limitedName := LimitLength(name, length)
	if ext == "" {
		return limitedName
	}
	return fmt.Sprintf("%s.%s", limitedName, ext)
}

// FilePath gen valid file path
func FilePath(name, ext string, length int, outputPath string, escape bool) (string, error) {
	if outputPath != "" {
		if _, err := os.Stat(outputPath); err != nil {
			return "", err
		}
	}
	var fileName string
	if escape {
		fileName = FileName(name, ext, length)
	} else {
		fileName = fmt.Sprintf("%s.%s", name, ext)
	}
	return filepath.Join(outputPath, fileName), nil
}

// Range generate a sequence of numbers by range
func Range(min, max int) []int {
	items := make([]int, max-min+1)
	for index := range items {
		items[index] = min + index
	}
	return items
}
