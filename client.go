package gofaceplus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type FaceClient struct {
	ApiServer string `json:"api_server"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
}

type DetectResult struct {
	Faces     []*Face `json:"face"`
	ImgHeight int64   `json:"img_height"`
	ImgWidth  int64   `json:"img_width"`
	ImgId     string  `json:"img_id"`
	SessionId string  `json:"session_id"`
	ImgUrl    string  `json:"url"`
}

type Face struct {
	Id        string    `json:"face_id"`
	Attrs     FaceAttrs `json:"attribute"`
	Positions Positions `json:"position"`
	Tag       string    `json:"tag"`
}

type FaceAttrs struct {
	Age     Age     `json:"age"`
	Gender  Gender  `json:"gender"`
	Glass   Glass   `json:"glass"`
	Pose    Pose    `json:"pose"`
	Race    Race    `json:"race"`
	Smiling Smiling `json:"smiling"`
}

type Pose struct {
	PitchAngle PitchAngle `json:"pitch_angle"`
	RollAngle  RollAngle  `json:"roll_angle"`
	YawAngle   YawAngle   `json:"yaw_angle"`
}

type Age struct {
	Range int64 `json:"range"`
	Value int64 `json:"value"`
}

type Gender struct {
	Confidence float64 `json:"confidence"`
	Value      string  `json:"value"`
}

type Glass struct {
	Confidence float64 `json:"confidence"`
	Value      string  `json:"value"`
}

type Race struct {
	Confidence float64 `json:"confidence"`
	Value      string  `json:"value"`
}

type Smiling struct {
	Value float64 `json:"value"`
}

type PitchAngle struct {
	Value float64 `json:"value"`
}

type RollAngle struct {
	Value float64 `json:"value"`
}

type YawAngle struct {
	value float64 `json:"value"`
}

type Img struct {
	Id     string `json:"id"`
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
	Url    string `json:"url"`
}
type Positions struct {
	Center     Point   `json:"center"`
	EyeLeft    Point   `json:"eye_left"`
	EyeRight   Point   `json:"eye_right"`
	Height     float64 `json:"height"` //0~100之间的实数，表示检出的脸的高度在图片中百分比
	Width      float64 `json:"width"`  //0~100之间的实数，表示检出的脸的宽度在图片中百分比
	MouseLeft  Point   `json:"mouse_left"`
	MouseRight Point   `json:"mouse_right"`
	Nose       Point   `json:"nose"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (fc *FaceClient) DetectImg(imgPath string) (sessionId string, faces []*Face, img *Img, err error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return
	}
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	fw, err := writer.CreateFormFile("img", imgPath)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, file); err != nil {
		return
	}
	writer.Close()
	address := fmt.Sprintf("%s/detection/detect?api_key=%s&api_secret=%s&attribute=%s",
		fc.ApiServer, fc.ApiKey, fc.ApiSecret, "gender,age,race,smiling,glass,pose")
	fmt.Printf("Address: %s \n", address)
	req, err := http.NewRequest("POST", address, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	fmt.Printf("Response status code: %d", resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	obj := &DetectResult{}
	err = json.Unmarshal(body, obj)
	sessionId = obj.SessionId
	faces = obj.Faces
	img = &Img{
		Id:     obj.ImgId,
		Height: obj.ImgHeight,
		Width:  obj.ImgWidth,
		Url:    obj.ImgUrl,
	}
	return
}
