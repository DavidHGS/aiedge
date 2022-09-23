package main

import (
	"aiedge/draw"
	"aiedge/ffmpeg"
	"aiedge/gocv"
	"aiedge/post"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	imgNum := 300
	//拉流抽帧
	ffmpeg.VideoGetNetImg(imgNum)
	//Post
	jwtToken := post.Signin("liujiaxi", "123")
	println(jwtToken.Msg)
	for i := 1; i <= imgNum; i++ {
		filePath := fmt.Sprintf("./images/frame%d.jpg", i) //源文件路径
		detectionRes := post.DetectionResult{}
		detectionRes = post.PostDetection(filePath, jwtToken.Msg)
		fmt.Printf("%d DetectionObject\n", len(detectionRes.Rects))
		//画框
		fileOutput := fmt.Sprintf("./output/outframe%d.jpg", i) //输出文件路径
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
	gocv.Img2Video(30, 1280, 720, imgNum)
}
