package test

import (
	"encoding/base64"
	. "github.com/journeymidnight/yig-front-caddy/test/lib"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	HTTP     = "http://"
	REQUEST  = "?x-oss-process="
	PUTSTYLE = "put_style/"
	GETSTYLE = "get_style/"
	DELSTYLE = "del_style/"
)

type PipaTestUnit struct {
	Cases    []Case
	TestName string
	Fn       func(t *testing.T, input []string, directly bool) (code int, output string)
}

type Case struct {
	Input              []string
	ExpectedStatusCode int
	ExpectedContent    string
	IsDirectly         bool
}

func doPutStyle(t *testing.T, input []string, directly bool) (code int, output string) {
	header := make(map[string]string)
	url := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + PUTSTYLE + input[0]
	auth, err := GetJwtForPipa()
	if err != nil {
		t.Error(err)
	}
	header["Authorization"] = auth
	styleCode := input[1]
	payload := strings.NewReader(base64.StdEncoding.EncodeToString([]byte(styleCode)))
	code, output = HttpPut(url, header, payload)
	urlDel := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + DELSTYLE + input[0]
	status, _ := HttpDelete(urlDel, header)
	if status != 200 {
		t.Error("doPutStyle have something wrong")
	}
	return
}

func doDelStyle(t *testing.T, input []string, directly bool) (code int, output string) {
	header := make(map[string]string)
	url := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + DELSTYLE + input[0]
	auth, err := GetJwtForPipa()
	if err != nil {
		t.Error(err)
	}
	header["Authorization"] = auth
	urlPut := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + PUTSTYLE + input[0]
	styleCode := input[1]
	payload := strings.NewReader(base64.StdEncoding.EncodeToString([]byte(styleCode)))
	status, _ := HttpPut(urlPut, header, payload)
	if status != 200 {
		t.Error("doDelStyle have something wrong")
	}
	code, output = HttpDelete(url, header)
	return
}

func doGetStyle(t *testing.T, input []string, directly bool) (code int, output string) {
	header := make(map[string]string)
	url := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + GETSTYLE
	auth, err := GetJwtForPipa()
	if err != nil {
		t.Error(err)
	}
	header["Authorization"] = auth
	styleCode := input[1]
	payload := strings.NewReader(base64.StdEncoding.EncodeToString([]byte(styleCode)))
	urlPut := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + PUTSTYLE + input[0]
	status, _ := HttpPut(urlPut, header, payload)
	if status != 200 {
		t.Error("doGetStyle have something wrong")
	}
	code, output = HttpGet(url, header)
	urlDel := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + DELSTYLE + input[0]
	status, _ = HttpDelete(urlDel, header)
	if status != 200 {
		t.Error("doGetStyle have something wrong")
	}
	return
}

func doProcessImage(t *testing.T, input []string, directly bool) (code int, output string) {
	var i int
	if directly {
		i = 2
	} else {
		i = 1
	}
	defer RedisFlushAll()
	for j := 0; j < i; j++ {
		url := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + input[0]
		go RedisPipa()
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			t.Error(err)
		}
		code = resp.StatusCode
		bodyByte, _ := ioutil.ReadAll(resp.Body)
		output = string(bodyByte)
	}
	return
}

func doProcessStyle(t *testing.T, input []string, directly bool) (code int, output string) {
	var i int
	if directly {
		i = 2
	} else {
		i = 1
	}
	defer RedisFlushAll()
	for j := 0; j < i; j++ {
		headerPut := make(map[string]string)
		urlPut := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + PUTSTYLE + input[1]
		auth, err := GetJwtForPipa()
		if err != nil {
			t.Error(err)
		}
		headerPut["Authorization"] = auth
		var styleCode string
		styleCode = input[2]
		payload := strings.NewReader(base64.StdEncoding.EncodeToString([]byte(styleCode)))
		status, _ := HttpPut(urlPut, headerPut, payload)
		if status != 200 {
			t.Error("doProcessStyle have something wrong")
		}
		url := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + input[0]
		go RedisPipa()
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			t.Error(err)
		}
		code = resp.StatusCode
		bodyByte, _ := ioutil.ReadAll(resp.Body)
		output = string(bodyByte)
		urlDel := HTTP + TEST_BUCKET + "." + End_Point + REQUEST + DELSTYLE + input[1]
		status, _ = HttpDelete(urlDel, headerPut)
		if status != 200 {
			t.Error("doProcessStyle have something wrong")
		}
	}
	return
}

var testUnits = []PipaTestUnit{
	{
		Cases: []Case{
			{
				Input:              []string{"test", "image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10"},
				ExpectedStatusCode: 200,
				ExpectedContent:    "",
				IsDirectly:         false,
			},
		},
		TestName: "Put style",
		Fn:       doPutStyle,
	},
	{
		Cases: []Case{
			{
				Input:              []string{"test", "image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10"},
				ExpectedStatusCode: 200,
				ExpectedContent:    "",
				IsDirectly:         false,
			},
		},
		TestName: "Delete style",
		Fn:       doDelStyle,
	},
	{
		Cases: []Case{
			{
				Input:              []string{"test", "image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10"},
				ExpectedStatusCode: 200,
				ExpectedContent:    "<ImageStyle><Bucket>mybucket</Bucket><StyleName>test</StyleName><StyleCode>image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10</StyleCode></ImageStyle>",
				IsDirectly:         false,
			},
		},
		TestName: "Get style",
		Fn:       doGetStyle,
	},
	{
		Cases: []Case{
			{
				Input:              []string{"style/test", "test", "image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10"},
				ExpectedStatusCode: 200,
				ExpectedContent:    "hehehehehehe",
				IsDirectly:         false,
			},
		},
		TestName: "Process style",
		Fn:       doProcessStyle,
	},
	{
		Cases: []Case{
			{
				Input:              []string{"image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10"},
				ExpectedStatusCode: 200,
				ExpectedContent:    "hehehehehehe",
				IsDirectly:         false,
			},
		},
		TestName: "Process image",
		Fn:       doProcessImage,
	},
	{
		Cases: []Case{
			{
				Input:              []string{"style/test", "test", "image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10"},
				ExpectedStatusCode: 200,
				ExpectedContent:    "hehehehehehe",
				IsDirectly:         true,
			},
		},
		TestName: "Process style from redis directly",
		Fn:       doProcessStyle,
	},
	{
		Cases: []Case{
			{
				Input:              []string{"image/watermark,image_cGljcy9KZWxseWZpc2guanBnP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMjA,g_nw,x_10,y_10"},
				ExpectedStatusCode: 200,
				ExpectedContent:    "hehehehehehe",
				IsDirectly:         true,
			},
		},
		TestName: "Process image from redis directly",
		Fn:       doProcessImage,
	},
}

func Test_Pipa(t *testing.T) {
	sc := NewS3()
	err := sc.MakeBucket(TEST_BUCKET)
	if err != nil {
		t.Fatal("MakeBucket err:", err)
	}
	defer sc.DeleteBucket(TEST_BUCKET)
	for _, unit := range testUnits {
		for _, c := range unit.Cases {
			status, output := unit.Fn(t, c.Input, c.IsDirectly)
			if status != c.ExpectedStatusCode {
				t.Fatal(unit.TestName, "test case failed.", "Input:", c.Input, "Code:", status, "ExpectedStatusCode:", c.ExpectedStatusCode, "Directly:", c.IsDirectly, "Content:", output)
			}
			if output != c.ExpectedContent {
				t.Fatal(unit.TestName, "Test case failed.", "Input:", c.Input, "Output:", output, "Directly:", c.IsDirectly, "ExpectedContent:", c.ExpectedContent)
			}
		}
		t.Log(unit.TestName, "test successful")
	}
}
