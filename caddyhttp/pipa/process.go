package pipa

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/journeymidnight/yig-front-caddy/caddydb/types"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
)

const (
	HTTP     = "http://"
	IMAGEKEY = "x-oss-process"
)

type TaskData struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

type result struct {
	resByte []byte
	resErr  error
}

func imageFunc(r *http.Request, key string) (response []byte, err error) {
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
		PIPA.Log.Error("Image process err:", err)
		if err == ErrInternalServer {
			return writeErrorResponseWithPipa(string(response))
		}
		return writeErrorResponse(ErrPipaProcesss)
	}
	PIPA.Log.Info("Image process successful!")
	return
}

func processStyle(r *http.Request, styleName string) (response []byte, err error) {
	PIPA.Log.Info("Enter style process!")
	bucketHost := r.Host
	host := strings.Split(bucketHost, ".")
	bucket := host[0]
	style, err := PIPA.CaddyClient.GetStyle(bucket, styleName)
	if err != nil {
		PIPA.Log.Error("Image style process get style err:", err)
		piErrorCode, _ := err.(HandleError)
		return writeErrorResponse(piErrorCode)
	}
	if style.StyleName == "" || style.StyleCode == "" {
		PIPA.Log.Error("Image style process no such style")
		return writeErrorResponse(ErrNoSuchStyle)
	}
	unprocessedURL := r.URL.String()
	key := r.URL.Query().Get(IMAGEKEY)
	ch := make(chan result)
	taskData := &TaskData{}
	taskData.Uuid = uuid.New().String()
	taskData.Url = HTTP + r.Host + strings.Replace(unprocessedURL, key, style.StyleCode, 1)
	go image(*taskData, ch)
	res := <-ch
	response = res.resByte
	err = res.resErr
	if err != nil {
		PIPA.Log.Error("Image style process err:", err)
		if err == ErrInternalServer {
			return writeErrorResponseWithPipa(string(response))
		}
		piErrorCode, _ := err.(HandleError)
		return writeErrorResponse(piErrorCode)
	}
	PIPA.Log.Info("Image process with style successful!")
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
	PIPA.Log.Info("Enter put image style!")
	styleCode, err := ParseAndValidStyleCodeFromBody(r)
	if err != nil {
		PIPA.Log.Error("Put image style with parse and valid body err:", err)
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
		PIPA.Log.Error("Insert image style with get style err:", err)
		return
	}
	if data.StyleName == "" {
		uid, err := PIPA.S3Client.GetBucket(claim.Bucket)
		if err != nil {
			PIPA.Log.Error("Insert image style with get bucket err:", err)
			return err
		}
		if uid != claim.ProjectId {
			return ErrInvalidBucketPermission
		}
		validLength, err := PIPA.CaddyClient.GetStyles(claim.Bucket)
		if len(validLength.ImageStyle) >= 50 {
			PIPA.Log.Error("Insert image style with get style err:", err)
			return ErrTooManyImageStyle
		}
		err = PIPA.CaddyClient.InsertStyle(*style)
		if err != nil {
			PIPA.Log.Error("Insert new image style err:", err)
			return err
		}
		PIPA.Log.Info("Put new image style successfully!")
		return nil
	}
	err = PIPA.CaddyClient.UpdateStyle(*style)
	if err != nil {
		PIPA.Log.Error("Update image style err:", err)
		return
	}
	PIPA.Log.Info("Put update image style successfully!")
	return
}

func delStyle(r *http.Request, styleName string) (err error) {
	if r.Method != "DELETE" {
		return ErrInvalidRequestMethod
	}
	claim, err := GetMethodFromJWT(r, PIPA.SecretKey)
	if err != nil {
		PIPA.Log.Error("Delete image style with get params err:", err)
		return
	}
	PIPA.Log.Info("Enter delete image style!")
	style, err := PIPA.CaddyClient.GetStyle(claim.Bucket, styleName)
	if err != nil {
		PIPA.Log.Error("Delete image style with get style err:", err)
		return
	}
	if style.Bucket == "" {
		PIPA.Log.Error("Delete image style with no such key")
		return ErrNoSuchStyle
	}
	err = PIPA.CaddyClient.DelStyle(style)
	if err != nil {
		PIPA.Log.Error("Delete image style with sql err:", err)
		return
	}
	PIPA.Log.Info("Delete image style successfully!")
	return
}

func getStyle(r *http.Request) ([]byte, error) {
	if r.Method != "GET" {
		return writeErrorResponse(ErrInvalidRequestMethod)
	}
	claim, err := GetMethodFromJWT(r, PIPA.SecretKey)
	if err != nil {
		PIPA.Log.Error("Get image style with JWT err", err)
		piErrorCode, _ := err.(HandleError)
		return writeErrorResponse(piErrorCode)
	}
	PIPA.Log.Info("Enter get image style!")
	styles, err := PIPA.CaddyClient.GetStyles(claim.Bucket)
	if err != nil {
		PIPA.Log.Error("Get image style with sql err:", err)
		piErrorCode, _ := err.(HandleError)
		return writeErrorResponse(piErrorCode)
	}
	if len(styles.ImageStyle) <= 0 {
		PIPA.Log.Error("Get image style with no row")
		return writeErrorResponse(ErrNoRow)
	}
	response, err := xml.Marshal(styles)
	if err != nil {
		PIPA.Log.Error("Get image style with xml marshal err:", err)
		piErrorCode, _ := err.(HandleError)
		return writeErrorResponse(piErrorCode)
	}
	PIPA.Log.Info("Get image styles successfully!")
	return response, nil
}

func image(data TaskData, ch chan result) {
	var err error
	PIPA.Log.Info("UUID is:", data.Uuid, "Url is:", data.Url)
	res := new(result)
	res.resByte, err = PIPA.Redis.GetImageFromRedis(data.Url)
	if err != nil {
		PIPA.Log.Error("Get Image From Redis err:", err)
	}
	if len(res.resByte) > 0 {
		if err != nil {
			PIPA.Log.Error(data.Url, data.Uuid, err)
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
		PIPA.Log.Error(data.Url, data.Uuid, err)
		res.resByte = nil
		res.resErr = err
		ch <- *res
		return
	}
	err = PIPA.Redis.PushRequest(taskdata)
	if err != nil {
		PIPA.Log.Error("Push Request to redis err:", data.Url, data.Uuid, err)
		res.resByte = nil
		res.resErr = err
		ch <- *res
		return
	}
	res.resByte, err = PIPA.Redis.PopResponse(data.Uuid, data.Url)
	if err != nil {
		PIPA.Log.Error("pop redis Response err:", res.resByte, err)
		res.resErr = err
		ch <- *res
		return
	}
	res.resErr = nil
	ch <- *res
	return
}
