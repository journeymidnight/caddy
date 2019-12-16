package pipa

import (
	"encoding/json"
	"encoding/xml"
	"github.com/google/uuid"
	"github.com/journeymidnight/yig-front-caddy/caddydb/types"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"net/http"
	"strings"
)

const (
	HTTP   = "http://"
	HEADER = "?x-oss-process="
)

type TaskData struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

type result struct {
	resByte []byte
	resErr  error
}

func processRequest(r *http.Request, key string) (response []byte, err error) {
	fuc := strings.Split(key, "/")
	styleName := fuc[1]
	switch fuc[0] {
	case "image":
		response, err = processImage(r)
		return
	case "style":
		response, err = processStyle(r, styleName)
		return
	case "put_style":
		err = putStyle(r, styleName)
		return
	case "del_style":
		err = delStyle(r, styleName)
		return
	case "get_style":
		response, err = getStyle(r)
		return
	default:
		return nil, ErrNoRouter
	}
}

func processImage(r *http.Request) (response []byte, err error) {
	PIPA.Log.Println(20, "Enter image process!")
	ch := make(chan result)
	taskData := &TaskData{}
	taskData.Uuid = uuid.New().String()
	taskData.Url = HTTP + r.Host + r.URL.String()
	go image(*taskData, ch)
	res := <-ch
	response = res.resByte
	err = res.resErr
	if err != nil {
		return
	}
	PIPA.Log.Println(20, "Image process successful!")
	return
}

func processStyle(r *http.Request, styleName string) (response []byte, err error) {
	PIPA.Log.Println(20, "Enter style process!")
	bucketHost := r.Host
	host := strings.Split(bucketHost, ".")
	bucket := host[0]
	style, err := PIPA.CaddyClient.GetStyle(bucket, styleName)
	if err != nil {
		return
	}
	if style.StyleName == "" || style.StyleCode == "" {
		return nil, ErrNoSuchKey
	}
	urls := strings.Split(r.URL.String(), "?")
	styleUrl := urls[0] + HEADER + style.StyleCode
	ch := make(chan result)
	taskData := &TaskData{}
	taskData.Uuid = uuid.New().String()
	taskData.Url = HTTP + r.Host + styleUrl
	go image(*taskData, ch)
	res := <-ch
	response = res.resByte
	err = res.resErr
	if err != nil {
		return
	}
	PIPA.Log.Println(20, "Image process with style successful!")
	return
}

func putStyle(r *http.Request, styleName string) (err error) {
	if r.Method != "PUT" {
		return ErrInvalidRequestMethod
	}
	claim, err := GetMethodFromJWT(r, PIPA.SecretKey)
	if err != nil {
		return
	}
	PIPA.Log.Println(20, "Enter put image style!")
	styleCode, err := ParseAndValidStyleCodeFromBody(r)
	if err != nil {
		return
	}
	if ok := validStyleName(styleName); !ok {
		return ErrInvalidStyleCode
	}
	style := new(types.ImageStyle)
	style.Bucket = claim.Bucket
	style.StyleName = styleName
	style.StyleCode = styleCode
	data, err := PIPA.CaddyClient.GetStyle(claim.Bucket, styleName)
	if err != nil {
		return
	}
	if data.StyleName == "" {
		uid, err := PIPA.S3Client.GetBucket(claim.Bucket)
		if err != nil {
			return err
		}
		if uid != claim.ProjectId {
			return ErrInvalidBucketPermission
		}
		validLength, err := PIPA.CaddyClient.GetStyles(claim.Bucket)
		if len(validLength) >= 50 {
			return ErrTooManyImageStyle
		}
		err = PIPA.CaddyClient.InsertStyle(*style)
		if err != nil {
			return err
		}
		PIPA.Log.Println(20, "Put new image style successfully!")
		return nil
	}
	err = PIPA.CaddyClient.UpdateStyle(*style)
	if err != nil {
		return
	}
	PIPA.Log.Println(20, "Put update image style successfully!")
	return
}

func delStyle(r *http.Request, styleName string) (err error) {
	if r.Method != "DELETE" {
		return ErrInvalidRequestMethod
	}
	claim, err := GetMethodFromJWT(r, PIPA.SecretKey)
	if err != nil {
		return
	}
	PIPA.Log.Println(20, "Enter delete image style!")
	style, err := PIPA.CaddyClient.GetStyle(claim.Bucket, styleName)
	if err != nil {
		return
	}
	if style.Bucket == "" {
		return ErrNoSuchKey
	}
	err = PIPA.CaddyClient.DelStyle(style)
	if err != nil {
		return
	}
	PIPA.Log.Println(20, "Delete image style successfully!")
	return
}

func getStyle(r *http.Request) ([]byte, error) {
	if r.Method != "GET" {
		return nil, ErrInvalidRequestMethod
	}
	claim, err := GetMethodFromJWT(r, PIPA.SecretKey)
	if err != nil {
		return nil, err
	}
	PIPA.Log.Println(20, "Enter get image style!")
	styles, err := PIPA.CaddyClient.GetStyles(claim.Bucket)
	if err != nil {
		return nil, err
	}
	if len(styles) <= 0 {
		return nil, ErrNoRow
	}
	response, err := xml.Marshal(styles)
	if err != nil {
		return nil, err
	}
	PIPA.Log.Println(20, "Get image styles successfully!")
	return response, nil
}

func image(data TaskData, ch chan result) {
	var err error
	res := new(result)
	res.resByte, err = getImageFromRedis(data.Url, PIPA.redis)
	if len(res.resByte) > 0 {
		if err != nil {
			PIPA.Log.Println(20, data.Url, data.Uuid, err)
			res.resByte = nil
			res.resErr = err
			ch <- *res
			return
		}
		res.resErr = nil
		ch <- *res
		return
	}
	taskdata, err := json.Marshal(data)
	if err != nil {
		PIPA.Log.Println(20, data.Url, data.Uuid, err)
		res.resByte = nil
		res.resErr = err
		ch <- *res
		return
	}
	err = pushRequest(taskdata, PIPA.redis)
	if err != nil {
		PIPA.Log.Println(20, data.Url, data.Uuid, err)
		res.resByte = nil
		res.resErr = err
		ch <- *res
		return
	}
	res.resByte, err = popResponse(data, PIPA.redis)
	if err != nil {
		PIPA.Log.Println(20, "popResponse err:", res.resByte, err)
		res.resErr = err
		ch <- *res
		return
	}
	res.resErr = nil
	ch <- *res
	return
}
