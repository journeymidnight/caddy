package pipa

import (
	"encoding/base64"
	"encoding/json"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"io/ioutil"
	"net/http"
	"strings"
)

var fuc = []string{"resize", "watermark"}

func ParseAndValidStyleCodeFromBody(r *http.Request) (code string, err error) {
	styleCode := make(map[string]string)
	PIPA.Log.Println(20, "Enter get style code from body")
	codeBuffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", ErrImageStyleParsing
	}
	err = json.Unmarshal(codeBuffer, &styleCode)
	if err != nil {
		return "", ErrImageStyleParsing
	}
	code = styleCode["style_code"]
	if ok := valid(code); !ok {
		return "", ErrInvalidStyleCode
	}
	return
}

func valid(code string) bool {
	decodeBytes, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return false
	}
	decodeStyle := string(decodeBytes)
	styles := strings.Split(decodeStyle, "/")
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
