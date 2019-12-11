package pipa

import (
	"encoding/base64"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"
)

var fuc = []string{"resize", "watermark"}

func ParseAndValidStyleCodeFromBody(r *http.Request) (code string, err error) {
	PIPA.Log.Println(20, "Enter get style code from body")
	codeBuffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", ErrImageStyleParsing
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(string(codeBuffer))
	if err != nil {
		return
	}
	code = string(decodeBytes)
	if ok := valid(code); !ok {
		return "", ErrInvalidStyleCode
	}
	return
}

func valid(code string) bool {
	styles := strings.Split(code, "/")
	for _, style := range styles {
		if style == "image" {
			continue
		}
		keys := strings.Split(style, ",")
		isTrue := false
		for _, f := range fuc {
			if f == keys[0] {
				isTrue = true
			}
		}
		if !isTrue {
			return false
		}
		for n, key := range keys {
			if n > 0 && len(strings.Split(key, "_")) < 2 {
				return false
			}
		}
	}
	return true
}

func validStyleName(styleName string) bool {
	if len(styleName) > 64 || len(styleName) < 1 {
		return false
	}
	if !utf8.ValidString(styleName) {
		return false
	}
	for _, n := range styleName {
		if (n >= 0 && n <= 31) || (n >= 127 && n <= 255) {
			return false
		}
		c := string(n)
		if strings.ContainsAny(c, "&$=;:+ ,?\\^`><{}][#%\"'~|") {
			return false
		}
	}
	return true
}
