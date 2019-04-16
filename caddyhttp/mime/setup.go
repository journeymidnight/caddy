// Copyright 2015 Light Code Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mime

import (
	"fmt"
	"strings"

	"github.com/journeymidnight/yig-front-caddy"
	"github.com/journeymidnight/yig-front-caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("mime", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new mime middleware instance.
func setup(c *caddy.Controller) error {
	configs, setquerys, setheaders, err := mimeParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Mime{Next: next, Configs: configs, SetQuerys:setquerys, SetHeaders: setheaders}
	})
	return nil
}

func mimeParse(c *caddy.Controller) (Config, SetQuery, SetHeader, error) {
	configs := Config{}
	setquerys := SetQuery{}
	setheaders := SetHeader{}
	for c.Next() {
		// At least one extension is required

		args := c.RemainingArgs()
		switch len(args) {
		case 2:
			if err := validateExtOutside(configs, args[0]); err != nil {
				return configs, setquerys, setheaders, err
			}
			configs[args[0]] = args[1]
		case 1:
			return configs, setquerys, setheaders, c.ArgErr()
		case 0:
			// At least one extension is required
			//Change: this func must be like this
			//  mime {
			//      Query Key [Req.value] value
			//      Header key [Req.value] value
			//      .file-type value
			//      ...
			//  }
			for c.NextBlock() {
				ext := c.Val()
				//check the correctness of key-data
				if err := valiKeydateExtInside(ext); err != nil {
					return configs, setquerys, setheaders, err
				}
				switch ext {
				case "Query": //when the key is header, get the line
					c.NextArg()
					key := c.Val()
					if !c.NextArg() {
						return configs, setquerys, setheaders, c.ArgErr()
					}
					value := c.Val()
					if value =="}" {
						return configs, setquerys, setheaders, c.ArgErr()
					}
					setquerys[key] = new(ParameterValue)
					setquerys[key].SettingContentType = value
					setquerys[key].Parameter = ""
					if c.NextArg() {
						setquerys[key].Parameter = value
						setquerys[key].SettingContentType = c.Val()
					}
					break
				case "Header": //when the key is header, get the line
					c.NextArg()
					key := c.Val()
					if !c.NextArg() {
						return configs, setquerys, setheaders, c.ArgErr()
					}
					value := c.Val()
					if value =="}" {
						return configs, setquerys, setheaders, c.ArgErr()
					}
					setheaders[key] = new(ParameterValue)
					setheaders[key].SettingContentType = value
					setheaders[key].Parameter = ""
					if c.NextArg() {
						setheaders[key].Parameter = value
						setheaders[key].SettingContentType = c.Val()
					}
					break
				default: //set the content-type directly
					if _, err := validataExt(c, configs, ext); err != nil {
						return configs, setquerys, setheaders, err
					}
					configs[ext] = c.Val()
				}
			}
		}

	}

	return configs, setquerys, setheaders, nil
}

// validateExt checks for valid file name extension.
func validateExtOutside(configs Config, ext string) error {
	if !strings.HasPrefix(ext, ".") {
		return fmt.Errorf(`mime: invalid extension "%v" (must start with dot)`, ext)
	}
	if _, ok := configs[ext]; ok {
		return fmt.Errorf(`mime: duplicate extension "%v" found`, ext)
	}
	return nil
}

// check the correctness of key-data
func valiKeydateExtInside(ext string) error {
	if strings.HasPrefix(ext, ".") || ext == "Query" || ext == "Header" {
		return nil
	}
	return fmt.Errorf(`mime: invalid extension "%v" (must start with dot or "Query" or "Header")`, ext)
}

// When we need set the file's response directly, we should check it at first.
func validataExt(c *caddy.Controller,configs Config,ext string) (Config,error) {
	if _, ok := configs[ext]; ok {
		return configs,fmt.Errorf(`mime: duplicate extension "%v" found`, ext)
	}
	if !c.NextArg() {
		return configs, c.ArgErr()
	}
	return configs,nil
}