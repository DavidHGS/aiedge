package gocv

import (
	"fmt"

	"gocv.io/x/gocv"
)

func Img2Video(frame_bfs, frame_width, frame_height, frame_num int) {
	outputName := "video02.mp4"
	codec := "mp4v"
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
	fmt.Println("ok")
}
