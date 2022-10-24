package gocv

import (
	"fmt"

	"gocv.io/x/gocv"
)

func Img2Video(frame_bfs, frame_width, frame_height, frame_num int) {
	// outputName := "video01.avi" // 尾巴要有数字，不知道哪里少了个库
	// codec := "MJPG"
	outputName := "video02.mp4"
	// codec := "mp4v"
	codec := "X264"
	bfs := float64(frame_bfs)
	frameWidth := frame_width
	frameHeight := frame_height
	// imgs := [...]string{"./images/frame1.jpg", "./images/frame1.jpg"}
	videoWriter, err := gocv.VideoWriterFile(outputName, codec, bfs, frameWidth, frameHeight, true)
	if err != nil {
		fmt.Println("VideoWriterFile error")
		return
	}
	defer videoWriter.Close()
	if !videoWriter.IsOpened() {
		fmt.Println("videoWriter open failed")
	}
	for frameNum := 1; frameNum <= frame_num; frameNum++ {
		inputName := fmt.Sprintf("./output/outframe%d.jpg", frameNum)
		mat := gocv.IMRead(inputName, gocv.IMReadColor)
		// gocv.Resize(mat, &mat, image.Point{X: frameWidth, Y: frameHeight}, 0, 0, gocv.InterpolationArea)
		// gocv.CvtColor(mat, &mat, gocv.ColorBGRToRGBA)
		videoWriter.Write(mat)
		defer mat.Close()
	}
	fmt.Println("Img2Video ok")
}

// func testPushStream() {
// 	camera, err := gocv.VideoCaptureDevice(0)
// 	img := gocv.NewMat()

// 	for {
// 		//从摄像头获取视频数据
// 		camera.Read(&img)
// 		//图像Base64编码
// 		data, err := gocv.IMEncode(".jpg", img)
// 		n := base64.StdEncoding.EncodedLen(len(data))
// 		dst := make([]byte, n)
// 		base64.StdEncoding()
// 	}
// }
