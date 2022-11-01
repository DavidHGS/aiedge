package main

import (
	"aiedge/draw"
	"aiedge/ffmpeg"
	"aiedge/gocv"
	"aiedge/post"
	"aiedge/stream"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Username string `json:"username" gorm:"primary_key;not null;unique;comment:'用户名'"`
	Password string `json:"password" gorm:"comment:'密码'"`
}

func main() {
	//系统变量
	//IMG_NUM 默认为30
	//FPS 25
	//PULLSTREAM_URL 默认为 rtmp://aiedge.ndsl-lab.cn:8035/live/stream12
	//PUSHSTREAM_URL 默认为 rtmp://aiedge.ndsl-lab.cn:8035/live/stream
	/*****************************图片参数*************************************************/
	imgNum := 30 //截取图片数量
	var err error
	if IMG_NUM := os.Getenv("IMG_NUM"); IMG_NUM != "" {
		imgNum, err = strconv.Atoi(IMG_NUM)
		fmt.Println(imgNum)
	}
	inputPath := "./images/frame%d.jpg"     //拉流抽帧图片地址
	outputPath := "./output/outframe%d.jpg" //目标检测图片输出地址
	/*****************************输出视频参数**********************************************/
	//根据视频流输入设置
	fps := 25 //视频帧率
	if FPS := os.Getenv("FPS"); FPS != "" {
		fps, err = strconv.Atoi(FPS)
		fmt.Println("fps:", fps)
	}
	if err != nil {
		fmt.Println("Atoi error")
	}

	frameWidth := 1280 //视频帧宽度ls
	frameHeight := 720 //视频帧高度
	/*****************************接口参数*************************************************/
	// pullStreamUrl := "rtmp://aiedge.ndsl-lab.cn:8035/live/2d6tsn0iqy108" //拉流地址
	// pullStreamUrl := "rtmp://aiedge.ndsl-lab.cn:8035/live/stream1" //拉流地址
	// if PULLSTREAM_URL := os.Getenv("PULLSTREAM_URL"); PULLSTREAM_URL != "" {
	// 	pullStreamUrl = os.Getenv("PULLSTREAM_URL") //拉流地址
	// }
	// fmt.Println(pullStreamUrl)
	pushStreamUrl := "rtmp://aiedge.ndsl-lab.cn:8035/live/stream" //推流地址
	if PUSHSTREAM_URL := os.Getenv("PUSHSTREAM_URL"); PUSHSTREAM_URL != "" {
		pushStreamUrl = os.Getenv("PUSHSTREAM_URL") //拉流地址
	}
	// objectDetectionUrl := "http://192.168.1.200:30001/v1/object-detection" //目标检测接口地址
	objectDetectionUrl := "http://192.168.20.150:30123/v1/face/detection" //人脸识别接口地址
	if OBJECTDETECTIONURL := os.Getenv("OBJECTDETECTIONURL"); OBJECTDETECTIONURL != "" {
		objectDetectionUrl = os.Getenv("OBJECTDETECTIONURL") //人脸识别地址
	}
	loginUrl := "http://192.168.20.150:30089/api/v1/auth/signin" //用户登录，获取jwt接口地址
	edge_devname := "edge1-cam1.aiedge-public-device"
	if EDGE_DEVNAME := os.Getenv("EDGE_DEVNAME"); EDGE_DEVNAME != "" {
		edge_devname = os.Getenv("EDGE_DEVNAME")
	}
	start := time.Now()

	// buff := ffmpeg.VideoGetNetImg2buff(imgNum, pullStreamUrl)
	//Post
	//add secreat username password
	var user User
	data, err := ioutil.ReadFile("./config/user.json")
	if err != nil {
		log.Println("load user.json fail: ", err)
		return
	}
	if err := json.Unmarshal(data, &user); err != nil {
		log.Println("Unmarshal user.json fail: ", err)
		return
	}
	jwtToken := post.Signin(user.Username, user.Password, loginUrl)
	println(jwtToken.Msg)

	//环境变量 edge-devname=edge1-cam1
	edge := strings.Split(strings.Split(edge_devname, ".")[0], "-")[0]
	devName := strings.Split(strings.Split(edge_devname, ".")[0], "-")[1]
	namespace := "aiedge-public-device"
	//getStreamCode
	streamcode := stream.GetStream(jwtToken.Msg, edge+"-"+devName, edge, namespace)
	fmt.Println(streamcode)

	//send heartbeat
	var hbsignal = make(chan string)
	//start goruntime
	go stream.Sendheartbeats(jwtToken.Msg, devName, edge, namespace, streamcode.Streamcode, hbsignal)

	//拉流抽帧
	fmt.Println(streamcode.Rtmp)
	ffmpeg.VideoGetNetImg(imgNum, streamcode.Rtmp)

	//stop goruntime
	hbsignal <- "stop"
	close(hbsignal)

	for i := 1; i <= imgNum; i++ {
		filePath := fmt.Sprintf(inputPath, i) //源文件路径
		detectionRes := post.DetectionResult{}
		// detectionRes = post.PostDetectionFormBuffer(filePath, jwtToken.Msg, objectDetectionUrl, buff)
		detectionRes = post.PostDetection(filePath, jwtToken.Msg, objectDetectionUrl)
		fmt.Printf("%d DetectionObject\n", len(detectionRes.Rects))
		//画框
		fileOutput := fmt.Sprintf(outputPath, i) //输出文件路径
		file, _ := os.Open(filePath)
		fileBytes, _ := ioutil.ReadAll(file)

		//循环读出矩阵
		var rectsInfo []post.Rectangle
		for _, rect := range detectionRes.Rects {
			rectsInfo = append(rectsInfo, rect.Rect)
		}
		draw.DrawRectOnImageAndSave(fileOutput, fileBytes, rectsInfo)
	}

	//拼接成视频
	// gocv.Img2Video(30, 1280, 720, imgNum)
	gocv.Img2Video(fps, frameWidth, frameHeight, imgNum)
	elapsed := time.Since(start)
	fmt.Printf("Time Cost %s \n", elapsed)
	// //ffmpeg -re -stream_loop -1 -i video02.mp4 -vcodec copy -acodec copy -f flv -y rtmp://aiedge.ndsl-lab.cn:8035/live/stream

	command := "ffmpeg -re -stream_loop -1 -i video02.mp4 -vcodec copy -acodec copy -f flv -y %s"
	command = fmt.Sprintf(command, pushStreamUrl)
	fmt.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("command output: %q", out.String())
	// for true {
	// 	fmt.Printf("ok")
	// }

}

//docker build . --network="host" --build-arg "HTTP_PROXY=http://127.0.0.1:7890" --build-arg "HTTPS_PROXY=http://127.0.0.1:7890" -t hgs/aiapp:v8
//docker run -it -e IMG_NUM="125" -e FPS="25" -e PULLSTREAM_URL="rtmp://192.168.20.150:30200/live/stream1" -e PUSHSTREAM_URL="rtmp://192.168.20.150:30200/live/stream"  hgs/aiapp:v7 sh -c"cd /tmp/aiedge/&& go run main.go"
