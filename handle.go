// //opencv读取视频 —> 将视频分割为帧 —> 将每一帧进行需求加工后 —> 将此帧写入pipe管道 —> 利用ffmpeg进行推流直播
package main

import (
	. "aiedge/algorithm"
	"aiedge/ffmpeg"
	"aiedge/post"
	"aiedge/stream"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"gocv.io/x/gocv"
)

// type Myimage struct {
// 	img    []byte
// 	width  int
// 	height int
// }

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
	/*****************************去重参数*************************************************/
	deduplicationInterval := 1 //相隔多少秒去重
	/*****************************图片参数*************************************************/
	imgNum := 5 //截取图片数量
	var err error
	if IMG_NUM := os.Getenv("IMG_NUM"); IMG_NUM != "" {
		imgNum, err = strconv.Atoi(IMG_NUM)
		fmt.Println(imgNum)
	}
	inputPath := "./images/frame%d.jpg" //拉流抽帧图片地址
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

	// objectDetectionUrl := "http://192.168.1.200:30001/v1/object-detection" //目标检测接口地址
	fmtinUrl := "http://192.168.20.150:30089/api/v1/auth/signin" //用户登录，获取jwt接口地址
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
		fmt.Println("load user.json fail: ", err)
		return
	}
	if err := json.Unmarshal(data, &user); err != nil {
		fmt.Println("Unmarshal user.json fail: ", err)
		return
	}
	jwtToken := post.Signin(user.Username, user.Password, fmtinUrl)
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

	var imgMap map[string]([]ImgType)
	imgMap = make(map[string]([]ImgType)) //map[i]-----[]image

	t1 := time.Now()
	//创建目标跟踪任务
	auth := "Bearer " + jwtToken.Msg
	missonID := post.CreateMission("http://192.168.1.84:30002/v1/tracking/multiple-object-tracking/create", auth)
	//创建client 发送请求
	client := resty.New()
	for i := 1; i <= imgNum; i++ {
		filePath := fmt.Sprintf(inputPath, i) //源文件路径
		mat := gocv.IMRead(filePath, gocv.IMReadColor)
		buf, err := gocv.IMEncode(".jpg", mat)
		// 检测
		// client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "image/jpeg").
			SetHeader("Authorization", auth).
			SetBody(buf.GetBytes()).
			Post("http://192.168.20.150:30003/v1/face/detection")
		if err != nil {
			fmt.Println(resp)
			fmt.Println("can't run detection ", err.Error())
			return
		}
		fmt.Println(resp.StatusCode())
		// 人脸检测 格式与 object-detection 统一
		detectionResult := post.DetectionResult{}
		json.Unmarshal([]byte(resp.String()), &detectionResult)

		// label 0 人类 1 ... 根据 label 做筛选
		// *****人脸检测 格式要跟 object-detection 统一*****
		detection_result, err := json.Marshal(&detectionResult)
		if err != nil {
			fmt.Println("detectionResult jsonify failed", err.Error())
		}
		// 跟踪
		// 跟踪的输入是检测的输出
		// 跟踪的输出其实就是给每个检测框打了一个 ID
		resp, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", auth).
			SetBody(detection_result).
			SetQueryParam("trajectory_misson_id", strconv.FormatInt(missonID, 10)).
			Post("http://192.168.1.84:30002/v1/tracking/multiple-object-tracking")

		if err != nil {
			fmt.Println(resp)
			fmt.Println("can't run tracking ", err.Error())
			return
		}
		fmt.Println(resp.StatusCode())
		tracking_result := post.TrackingResult{}
		json.Unmarshal([]byte(resp.String()), &tracking_result)
		/*___________________________________对框处理_______________________________________________________*/
		fmt.Printf("%d TrajectoryObject\n", len(tracking_result.Trajectorys))

		for _, trajectory := range tracking_result.Trajectorys {
			//裁剪
			fmt.Println(int(trajectory.Left))
			if int(trajectory.Left) < 0 {
				continue
			}
			croppedMat := mat.Region(image.Rect(int(trajectory.Left), int(trajectory.Top), int(trajectory.Right), int(trajectory.Bottom)))
			resultMat := croppedMat.Clone()

			//路径拼接
			srcPath := "./result/%d_%d_%.3f.jpg"                                                   //id_时间戳_置信度
			srcPathStr := fmt.Sprintf(srcPath, trajectory.Id, time.Now().Unix(), trajectory.Score) //输出文件路径

			// serveFrames(resultMat.ToBytes(), srcPathStr)
			buf_cutted, _ := gocv.IMEncode(".jpg", resultMat)

			//图片信息添加到map
			if _, ok := imgMap[strconv.FormatInt(trajectory.Id, 10)]; !ok { //map[key]为nil 添加
				var tmpArr []ImgType
				tmp := ImgType{
					buf_cutted.GetBytes(),
					srcPathStr,
				}
				tmpArr = append(tmpArr, tmp)
				imgMap[strconv.FormatInt(trajectory.Id, 10)] = tmpArr
				tmpArr = nil
			} else { //map[key]不为nil 添加
				tmpArr := imgMap[strconv.FormatInt(trajectory.Id, 10)]
				tmp := ImgType{
					buf_cutted.GetBytes(),
					srcPathStr,
				}
				tmpArr = append(tmpArr, tmp)
				imgMap[strconv.FormatInt(trajectory.Id, 10)] = tmpArr
			}

			// fmt.Println(resultMat.ToBytes())
			// fmt.Println(buf.GetBytes())
			// fmt.Println(bytes.Equal(resultMat.ToBytes(), buf.GetBytes()))
			// if resultMat.ToBytes()==buf.GetBytes(){}
			gocv.IMWrite(srcPathStr, resultMat)
		}
		/*___________________________________对框处理end_______________________________________________________*/
		//不断获取最新 `deduplicationInterval` 秒的最好图片
		if time.Now().After(t1.Add(time.Second * time.Duration(deduplicationInterval))) {
			t1 = time.Now() //重新计算秒数
			// 过去重
		}
	}
	//——————————————————————————————————————————————————去重————————————————————————————————————————————
	Deduplication(imgMap, auth) //去重
	post.RemoveMission("http://192.168.1.84:30002/v1/tracking/multiple-object-tracking/remove", auth, missonID)
	elapsed := time.Since(start)
	fmt.Printf("Time Cost %s \n", elapsed)
	//websocket推送图片
}

func serveFrames(imgByte []byte, filename string) {

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}
	out, _ := os.Create(filename)
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	//jpeg.Encode(out, img, nil)
	if err != nil {
		log.Println(err)
	}

}
