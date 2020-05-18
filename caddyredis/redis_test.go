package caddyredis_test

import (
	"encoding/json"
	. "github.com/journeymidnight/yig-front-caddy/caddyredis"
	"strconv"
	"testing"
	"time"
)

type TaskData struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

func BenchmarkRedisConnTest(b *testing.B) {
	conf := &Config{
		Address: []string{"10.0.42.4:7000", "10.0.42.3:7000", "10.0.47.215:7000", "10.0.42.4:7001", "10.0.42.3:7001", "10.0.47.215:7001"},
	}
	var configs []*Config
	configs = append(configs, conf)
	redis := MakeRedisConfig(configs)
	b.Run("Test whether 100 consecutive LPUSH requests have failed to connect", func(b *testing.B) {
		for i := 0; i < 100; i++ {
			time.Sleep(1 * time.Second)
			data := TaskData{
				Url: "http://uuo.bj.unicloud.com/%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20190716150310.jpg?x-oss-process=image/rotate,52&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20200514T021503Z&X-Amz-SignedHeaders=host&X-Amz-Expires=3600&X-Amz-Credential=l5ZPKvgf2N95pLGp%2F20200514%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=68441e48113a593e35f529ffcc7f73acf0da5f9172b26724428660c6dbce7253",
			}
			data.Uuid = "726eed7b-3ac7-45e5-bf6a-ef5f6013507b" + "===" + strconv.Itoa(i)
			taskdata, err := json.Marshal(data)
			if err != nil {
				b.Error("err marshal is :", err)
			}
			err = redis.PushRequest(taskdata)
			if err != nil {
				b.Error("err is :", err)
			}
		}
	})
}
