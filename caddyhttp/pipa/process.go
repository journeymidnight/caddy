package pipa

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/journeymidnight/yig-front-caddy/caddyerrors"
	"net/http"
	"strings"
)

type TaskData struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

type result struct {
	resByte []byte
	resErr  error
}

func processRequest(w http.ResponseWriter, r *http.Request, key string) (response []byte, err error) {
	PIPA.Log.Println(20, "Enter request process!")
	fuc := strings.Split(key, "/")
	switch fuc[0] {
	case "image":
		response, err = processImage(w, r)
		return
	case "style":
		response, err = processStyle(w, r)
		return
	default:
		return nil, caddyerrors.ErrNoRouter
	}
}

func processImage(w http.ResponseWriter, r *http.Request) (response []byte, err error) {
	PIPA.Log.Println(20, "Enter image process!")
	ch := make(chan result)
	taskData := &TaskData{}
	taskData.Uuid = uuid.New().String()
	taskData.Url = r.URL.String()
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

func processStyle(w http.ResponseWriter, r *http.Request) (response []byte, err error) {
	return
}

func image(data TaskData, ch chan result) {
	res := new(result)
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
