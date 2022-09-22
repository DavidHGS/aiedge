package main

// import (
// 	"aiedge/draw"
// 	"aiedge/post"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// )

// func main() {
// 	filePath := "./images/frame1.jpg" //源文件路径
// 	//Post
// 	jwtToken := post.Signin("liujiaxi", "123")
// 	println(jwtToken.Msg)
// 	detectionRes := post.DetectionResult{}
// 	detectionRes = post.PostDetection(filePath, jwtToken.Msg)
// 	fmt.Printf("%d DetectionObject\n", len(detectionRes.Rects))
// 	//画框
// 	fileOutput := "./output/outframe1.jpg" //输出文件路径
// 	file, _ := os.Open(filePath)
// 	fileBytes, _ := ioutil.ReadAll(file)

// 	//循环读出矩阵
// 	var rectsInfo []post.Rectangle
// 	for _, rect := range detectionRes.Rects {
// 		rectsInfo = append(rectsInfo, rect.Rect)
// 	}
// 	draw.DrawRectOnImageAndSave(fileOutput, fileBytes, rectsInfo)
// 	//拼接成视频

// }

// tutorial01.c
// Code based on a tutorial at http://dranger.com/ffmpeg/tutorial01.html

// A small sample program that shows how to use libavformat and libavcodec to
// read video from a file.
//
// Use
//
// gcc -o tutorial01 tutorial01.c -lavformat -lavcodec -lswscale -lz
//
// to build (assuming libavformat and libavcodec are correctly installed
// your system).
//
// Run using
//
// tutorial01 myvideofile.mpg
//
// to write the first five frames from "myvideofile.mpg" to disk in PPM
// format.
import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"unsafe"

	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/swscale"
	// "github.com/leokinglong/goav/avformat"
)

// SaveFrame writes a single frame to disk as a PPM file
func SaveFramePPM(frame *avutil.Frame, width, height, frameNumber int) {
	// Open file
	fileName := fmt.Sprintf("./images/frame%d.ppm", frameNumber)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error Reading")
	}
	defer file.Close()

	// Write header
	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
	file.Write([]byte(header))

	// Write pixel data
	for y := 0; y < height; y++ {
		data0 := avutil.Data(frame)[0]
		buf := make([]byte, width*3)
		startPos := uintptr(unsafe.Pointer(data0)) + uintptr(y)*uintptr(avutil.Linesize(frame)[0])
		for i := 0; i < width*3; i++ {
			element := *(*uint8)(unsafe.Pointer(startPos + uintptr(i)))
			buf[i] = element
		}
		file.Write(buf)
	}
}

// SaveFrame writes a single frame to disk as a JPG file
func SaveFrameJpg(frame *avutil.Frame, width, height, frameNumber int) {
	// Open file
	fileName := fmt.Sprintf("./images/frame%d.jpg", frameNumber)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error Reading")
	}
	defer file.Close()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		data0 := avutil.Data(frame)[0]
		startPos := uintptr(unsafe.Pointer(data0)) + uintptr(y)*uintptr(avutil.Linesize(frame)[0])
		//fmt.Println("startPos: ", startPos)
		xxx := width * 3
		for x := 0; x < width; x++ {
			var pixel = make([]byte, 3)
			for i := 0; i < 3; i++ {
				element := *(*uint8)(unsafe.Pointer(startPos + uintptr(xxx)))
				pixel[i] = element
				xxx++
			}
			img.SetRGBA(x, y, color.RGBA{pixel[0], pixel[1], pixel[2], 0xff})
		}
	}
	err = jpeg.Encode(file, img, nil)
	if err != nil {
		fmt.Println("jpeg.Encode err: ", err)
		log.Println("Error jpeg.Encode")
	}
}

//img2Video
func img2Video() {
	imgPath := "./images/netFrame1.jpg"
	dstPath := "./new.mp4"
	videoIndex := -1

	/*************************打开图片 生成解码器信息******************************************************************************************************/
	// 创建jpg的解封装上下文
	pFormatContext := avformat.AvformatAllocContext()
	if avformat.AvformatOpenInput(&pFormatContext, imgPath, nil, nil) != 0 {
		fmt.Printf("Unable to open file %s\n", imgPath)
		os.Exit(1)
	}
	// Retrieve stream information
	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
		fmt.Println("Couldn't find stream information")
		os.Exit(1)
	}

	pFormatContext.AvDumpFormat(0, imgPath, 0)
	/*************************end***********************************************************************************************************************************/

	var pCodecCtx *avcodec.Context
	//创建解码器
	//初始化解码器上下文用于对jpg进行解码
	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
		stream := pFormatContext.Streams()[i]
		/** 对于jpg图片来说，它里面就是一路视频流，所以媒体类型就是AVMEDIA_TYPE_VIDEO
		 */
		if stream.CodecParameters().AvCodecGetType() == avformat.AVMEDIA_TYPE_VIDEO {
			// Find the decoder for the video stream
			pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(stream.Codec().GetCodecId()))
			if pCodec == nil {
				fmt.Println("Unsupported codec!")
				os.Exit(1)
			}
			// Alloc context
			pCodecCtx := pCodec.AvcodecAllocContext3()
			// 将解码器参数复制给编码器
			if pCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(stream.Codec()))) != 0 {
				fmt.Println("Couldn't copy codec context")
				os.Exit(1)
			}
			// Open codec
			if pCodecCtx.AvcodecOpen2(pCodec, nil) < 0 {
				fmt.Println("Could not open codec")
				os.Exit(1)
			}
			videoIndex = i
			break
		}
	}
	///1 创建编码器 查找编码器
	//创建h264编码器及编码上下文
	encodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	if encodec == nil {
		fmt.Println("Unsupported codec!")
		os.Exit(1)
	}
	// Alloc context
	EnCodecCtx := encodec.AvcodecAllocContext3()
	inStream := pFormatContext.Streams()[videoIndex]
	//配置编码器上下文的成员
	//设置画面宽度 高度 输出像素格式（pix_fmt)
	EnCodecCtx.SetEncodeParams(inStream.Codec().GetWidth(), inStream.Codec().GetHeight(), inStream.Codec().GetPixelFormat())
	//设置帧率，num为分子，den为分母，如果是1/25则表示25帧/s
	EnCodecCtx.SetTimebase(1, 25)
	// 设置全局编码信息
	// EnCodecCtx.Set= EnCodecCtx.Flags() | avcodec.AV_CODEC_FLAG_GLOBAL_HEADER
	//打开编码器
	if EnCodecCtx.AvcodecOpen2(encodec, nil) < 0 {
		fmt.Println("Could not open codec")
		os.Exit(1)
	}
	/*************************创建编码器 生成编码器信息****************************************************************************************************************************************************************************/
	//2. 创建输出视频上下文
	outFormat := avformat.AvformatAllocContext()
	if avformat.AvformatAllocOutputContext2(&outFormat, nil, "mp4", dstPath) < 0 {
		fmt.Println("AvformatAllocOutputContext2 error")
		os.Exit(1)
	}
	/*************************创建输出流生成编码器信息****************************************************************************************************************************************************************************/
	//添加视频流
	stream := outFormat.AvformatNewStream(nil)
	videoOutIndex := stream.Index()
	//设置视频流参数;直接从编码器上下文拷贝 将AVCodecContext信息拷贝到AVCodecParameterst结构体中

	// stream.Codec().SetWidth(inStream.Codec().GetWidth())
	// stream.Codec().SetHeight(inStream.Codec().GetHeight())
	// stream.Codec().SetPixelFormat(inStream.Codec().GetPixelFormat())
	// stream.Codec().SetTimeBase(inStream.TimeBase())

	if (*avcodec.Context)(unsafe.Pointer(stream.Codec())).AvcodecCopyContext(EnCodecCtx) < 0 {
		fmt.Println("Couldn't copy codec context")
		os.Exit(1)
	}
	fmt.Println("*****打印AVFormatContext的内容************************")
	outFormat.AvDumpFormat(0, dstPath, 1)
	// avformat.AvIOOpen(outFormat.Pb(),dstPath,)

	// initialize SWS context for software scaling
	swsCtx := swscale.SwsGetcontext(inStream.Codec().GetWidth(), inStream.Codec().GetHeight(), swscale.PixelFormat(inStream.Codec().GetPixelFormat()),
		EnCodecCtx.Width(), EnCodecCtx.Height(), swscale.PixelFormat(EnCodecCtx.PixFmt()), 0, nil, nil, nil) //avcodec.SWS_BICUBIC

	// 打开或创建输出文件

	pb, err := avformat.AvIOOpen(dstPath, avformat.AVIO_FLAG_WRITE)
	outFormat.SetPb(pb)
	if err != nil {
		fmt.Println("Could not open output file!")
		os.Exit(1)
	}
	//写视频文件头
	if outFormat.AvformatWriteHeader(nil) < 0 {
		fmt.Println("avformat_write_header failed!")
		os.Exit(1)
	}
	// 创建编解码用的AVFrame
	deFrame := avutil.AvFrameAlloc()
	enFrame := avutil.AvFrameAlloc()
	avutil.AvSetFrame(enFrame, EnCodecCtx.Width(), EnCodecCtx.Height(), int(EnCodecCtx.PixFmt()))
	avutil.AvFrameGetBuffer(enFrame, 32)
	avutil.AvFrameMakeWritable(enFrame)

	inPkt := avcodec.AvPacketAlloc()
	for { //while 循环
		if pFormatContext.AvReadFrame(inPkt) == 0 {
			break
		}
		if inPkt.StreamIndex() != videoIndex {
			continue
		}
		//先解码
		doDecode(pCodecCtx, inPkt, deFrame, enFrame, swsCtx, EnCodecCtx, outFormat, videoOutIndex)
		inPkt.AvPacketUnref()
	}
	//刷新解码缓冲区
	doDecode(nil, nil, nil, nil, nil, nil, nil, videoOutIndex)
	outFormat.AvWriteTrailer()
	fmt.Println("转换视频结束")
	os.Exit(1)
}
func doDecode(pCodecCtx *avcodec.Context, inPkt *avcodec.Packet, deFrame *avutil.Frame, enFrame *avutil.Frame,
	swsCtx *swscale.Context, EnCodecCtx *avcodec.Context, outFormat *avformat.Context, videoOutIndex int) {
	num_pts := 0
	//先解码
	pCodecCtx.AvcodecSendPacket(inPkt)
	for {
		ret := pCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(deFrame)))
		fmt.Println("ret:", ret)
		if ret == avutil.AvErrorEOF {
			doEncode(nil, EnCodecCtx, outFormat, videoOutIndex)
			break
		} else if ret < 0 {
			break
		}
		//解码后，先进行格式转化再编码
		if swscale.SwsScale2(swsCtx, avutil.Data(deFrame), avutil.Linesize(deFrame), 0, pCodecCtx.Height(),
			avutil.Data(enFrame), avutil.Linesize(enFrame)) < 0 {
			fmt.Println("swscale.SwsScale2 failed!")
			os.Exit(1)
		}
		// 编码前要设置好pts的值，如果en_ctx->time_base为{1,fps}，那么这里pts的值即为帧的个数值
		// en_frame->pts = num_pts++;
		(*avcodec.Packet)(unsafe.Pointer(enFrame)).SetPts(int64(num_pts))
		num_pts++
		doEncode(enFrame, EnCodecCtx, outFormat, videoOutIndex)
	}
}

func doEncode(enFrame *avutil.Frame, EnCodecCtx *avcodec.Context, outFormat *avformat.Context, videoOutIndex int) {
	EnCodecCtx.AvcodecSendPacket((*avcodec.Packet)(unsafe.Pointer(enFrame)))
	for {
		out_Pkt := avcodec.AvPacketAlloc()
		if EnCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(out_Pkt))) < 0 {
			out_Pkt.AvPacketUnref()
			break
		}
		//成功编码了，写入前进行时间基的转换
		stream := outFormat.Streams()[videoOutIndex]
		out_Pkt.AvPacketRescaleTs(EnCodecCtx.AvCodecGetPktTimebase(), stream.TimeBase())
		fmt.Printf("video pts %d", out_Pkt.Pts())
		outFormat.AvWriteFrame(out_Pkt)
	}

}

func video2Img() {
	filename := "./src/sample.mp4"
	//视频流地址
	//初始化网络库
	// filename := "rtmp://192.168.20.221:30200/live/2gz8r2nfcg8q8"
	// avformat.AvformatNetworkInit()
	// Open video file
	pFormatContext := avformat.AvformatAllocContext()

	if avformat.AvformatOpenInput(&pFormatContext, filename, nil, nil) != 0 {
		fmt.Printf("Unable to open file %s\n", filename)
		os.Exit(1)
	}

	// Retrieve stream information
	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
		fmt.Println("Couldn't find stream information")
		os.Exit(1)
	}

	// Dump information about file onto standard error
	pFormatContext.AvDumpFormat(0, filename, 0)

	// Find the first video stream
	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
		switch pFormatContext.Streams()[i].CodecParameters().AvCodecGetType() {
		case avformat.AVMEDIA_TYPE_VIDEO:

			// Get a pointer to the codec context for the video stream
			pCodecCtxOrig := pFormatContext.Streams()[i].Codec()
			// Find the decoder for the video stream
			pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(pCodecCtxOrig.GetCodecId()))
			if pCodec == nil {
				fmt.Println("Unsupported codec!")
				os.Exit(1)
			}
			// Copy context
			pCodecCtx := pCodec.AvcodecAllocContext3()
			if pCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig))) != 0 {
				fmt.Println("Couldn't copy codec context")
				os.Exit(1)
			}

			// Open codec
			if pCodecCtx.AvcodecOpen2(pCodec, nil) < 0 {
				fmt.Println("Could not open codec")
				os.Exit(1)
			}

			// Allocate video frame
			pFrame := avutil.AvFrameAlloc()

			// Allocate an AVFrame structure
			pFrameRGB := avutil.AvFrameAlloc()
			if pFrameRGB == nil {
				fmt.Println("Unable to allocate RGB Frame")
				os.Exit(1)
			}

			// Determine required buffer size and allocate buffer
			numBytes := uintptr(avcodec.AvpictureGetSize(avcodec.AV_PIX_FMT_RGB24, pCodecCtx.Width(),
				pCodecCtx.Height()))
			buffer := avutil.AvMalloc(numBytes)

			// Assign appropriate parts of buffer to image planes in pFrameRGB
			// Note that pFrameRGB is an AVFrame, but AVFrame is a superset
			// of AVPicture
			avp := (*avcodec.Picture)(unsafe.Pointer(pFrameRGB))
			avp.AvpictureFill((*uint8)(buffer), avcodec.AV_PIX_FMT_RGB24, pCodecCtx.Width(), pCodecCtx.Height())

			// initialize SWS context for software scaling
			swsCtx := swscale.SwsGetcontext(
				pCodecCtx.Width(),
				pCodecCtx.Height(),
				(swscale.PixelFormat)(pCodecCtx.PixFmt()),
				pCodecCtx.Width(),
				pCodecCtx.Height(),
				avcodec.AV_PIX_FMT_RGB24,
				avcodec.SWS_BILINEAR,
				nil,
				nil,
				nil,
			)

			// Read frames and save first five frames to disk
			frameNumber := 1
			packet := avcodec.AvPacketAlloc()
			for pFormatContext.AvReadFrame(packet) >= 0 { //判断是否还有数据帧
				// Is this a packet from the video stream?
				if packet.StreamIndex() == i {
					// Decode video frame
					response := pCodecCtx.AvcodecSendPacket(packet) //调用avcodec_send_packet，发送packet给ffmpeg解码，然后再从ffmpeg的解码队列中获取解码后的帧数据（AVFrame）。
					if response < 0 {
						fmt.Printf("Error while sending a packet to the decoder: %s\n", avutil.ErrorFromCode(response))
					}
					for response >= 0 {
						response = pCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(pFrame)))
						if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
							break
						} else if response < 0 {
							if response == -11 {
								fmt.Println("response error = - 11") //ReadEagain错误
								break
							}
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return
						}

						if frameNumber <= 30 {
							// Convert the image from its native format to RGB
							swscale.SwsScale2(swsCtx, avutil.Data(pFrame),
								avutil.Linesize(pFrame), 0, pCodecCtx.Height(),
								avutil.Data(pFrameRGB), avutil.Linesize(pFrameRGB))

							// Save the frame to disk
							fmt.Printf("Writing frame %d\n", frameNumber)
							SaveFrameJpg(pFrameRGB, pCodecCtx.Width(), pCodecCtx.Height(), frameNumber)
						} else {
							return
						}
						frameNumber++
					}
				}

				// Free the packet that was allocated by av_read_frame
				packet.AvFreePacket()
			}

			// Free the RGB image
			avutil.AvFree(buffer)
			avutil.AvFrameFree(pFrameRGB)

			// Free the YUV frame
			avutil.AvFrameFree(pFrame)

			// Close the codecs
			pCodecCtx.AvcodecClose()
			(*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig)).AvcodecClose()

			// Close the video file
			pFormatContext.AvformatCloseInput()

			// Stop after saving frames of first video straem
			break
		case avformat.AVMEDIA_TYPE_DATA:
			i = 2
			// Get a pointer to the codec context for the video stream
			pCodecCtxOrig := pFormatContext.Streams()[i].Codec()
			// Find the decoder for the video stream
			pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(pCodecCtxOrig.GetCodecId()))
			if pCodec == nil {
				fmt.Println("Unsupported codec!")
				os.Exit(1)
			}
			// Copy context
			pCodecCtx := pCodec.AvcodecAllocContext3()
			if pCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig))) != 0 {
				fmt.Println("Couldn't copy codec context")
				os.Exit(1)
			}

			// Open codec
			if pCodecCtx.AvcodecOpen2(pCodec, nil) < 0 {
				fmt.Println("Could not open codec")
				os.Exit(1)
			}

			// Allocate video frame
			pFrame := avutil.AvFrameAlloc()

			// Allocate an AVFrame structure
			pFrameRGB := avutil.AvFrameAlloc()
			if pFrameRGB == nil {
				fmt.Println("Unable to allocate RGB Frame")
				os.Exit(1)
			}

			// Determine required buffer size and allocate buffer
			numBytes := uintptr(avcodec.AvpictureGetSize(avcodec.AV_PIX_FMT_RGB24, pCodecCtx.Width(),
				pCodecCtx.Height()))
			buffer := avutil.AvMalloc(numBytes)

			// Assign appropriate parts of buffer to image planes in pFrameRGB
			// Note that pFrameRGB is an AVFrame, but AVFrame is a superset
			// of AVPicture
			avp := (*avcodec.Picture)(unsafe.Pointer(pFrameRGB))
			avp.AvpictureFill((*uint8)(buffer), avcodec.AV_PIX_FMT_RGB24, pCodecCtx.Width(), pCodecCtx.Height())

			// initialize SWS context for software scaling
			swsCtx := swscale.SwsGetcontext(
				pCodecCtx.Width(),
				pCodecCtx.Height(),
				(swscale.PixelFormat)(pCodecCtx.PixFmt()),
				pCodecCtx.Width(),
				pCodecCtx.Height(),
				avcodec.AV_PIX_FMT_RGB24,
				avcodec.SWS_BILINEAR,
				nil,
				nil,
				nil,
			)

			// Read frames and save first five frames to disk
			frameNumber := 1
			packet := avcodec.AvPacketAlloc()
			for pFormatContext.AvReadFrame(packet) >= 0 {
				// Is this a packet from the video stream?
				if packet.StreamIndex() == i {
					// Decode video frame
					response := pCodecCtx.AvcodecSendPacket(packet)
					if response < 0 {
						fmt.Printf("Error while sending a packet to the decoder: %s\n", avutil.ErrorFromCode(response))
					}
					for response >= 0 {
						response = pCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(pFrame)))
						if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
							break
						} else if response < 0 {
							if response == -11 {
								fmt.Println("response error = - 11") //ReadEagain错误
								break
							}
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return
						}

						if frameNumber <= 30 {
							// Convert the image from its native format to RGB
							swscale.SwsScale2(swsCtx, avutil.Data(pFrame),
								avutil.Linesize(pFrame), 0, pCodecCtx.Height(),
								avutil.Data(pFrameRGB), avutil.Linesize(pFrameRGB))

							// Save the frame to disk
							fmt.Printf("Writing frame %d\n", frameNumber)
							SaveFrameJpg(pFrameRGB, pCodecCtx.Width(), pCodecCtx.Height(), frameNumber)
						} else {
							return
						}
						frameNumber++
					}
				}

				// Free the packet that was allocated by av_read_frame
				packet.AvFreePacket()
			}

			// Free the RGB image
			avutil.AvFree(buffer)
			avutil.AvFrameFree(pFrameRGB)

			// Free the YUV frame
			avutil.AvFrameFree(pFrame)

			// Close the codecs
			pCodecCtx.AvcodecClose()
			(*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig)).AvcodecClose()

			// Close the video file
			pFormatContext.AvformatCloseInput()

			// Stop after saving frames of first video straem
			break
		default:
			fmt.Println("Didn't find a video stream")
			os.Exit(1)
		}
	}
}

func main() {
	img2Video2()
	// video2Img()
}
func img2Video4() {
	ofmt_ctx := avformat.AvformatAllocContext()
	out_filename := "out.mp4"
	avformat.AvRegisterAll() //初始化解码器和复用器
	avformat.AvformatAllocOutputContext2(&ofmt_ctx, nil, "mp4", out_filename)
	out_stream := add_vidio_stream(ofmt_ctx)
	ofmt_ctx.AvDumpFormat(0, out_filename, 1)
	flag := ofmt_ctx.Flags() & avcodec.AV_CODEC_ID_MJPEG
	if flag == 0 {
		pb, err := avformat.AvIOOpen(out_filename, avformat.AVIO_FLAG_WRITE)
		ofmt_ctx.SetPb(pb)
		if err != nil {
			fmt.Printf("Could not open output file %s", out_filename)
			os.Exit(1)
		}
	}

	ofmt_ctx.AvformatWriteHeader(nil)

	// numBytes := uintptr(avcodec.AvpictureGetSize(avcodec.AV_PIX_FMT_YUV420P9, out_stream.Codec().GetWidth(),
	// 	out_stream.Codec().GetHeight()))
	// buffer := avutil.AvMalloc(numBytes)
	pkt := avcodec.AvPacketAlloc()
	// pkt.AvInitPacket()
	pkt.SetFlags(pkt.Flags() | avcodec.AV_PKT_FLAG_KEY)
	pkt.SetStreamIndex(out_stream.Index()) //获取视频信息，为压入帧图像做准备
	var fileName string
	for frame_index := 1; frame_index <= 30; frame_index++ {
		fileName = fmt.Sprintf("./images/netFrame%d.jpg", frame_index)
		// 创建jpg的解封装上下文
		in_fmt := avformat.AvformatAllocContext()
		if avformat.AvformatOpenInput(&in_fmt, fileName, nil, nil) != 0 {
			fmt.Printf("Unable to open file %s\n", fileName)
			os.Exit(1)
		}
		// Retrieve stream information
		if in_fmt.AvformatFindStreamInfo(nil) < 0 {
			fmt.Println("Couldn't find stream information")
			os.Exit(1)
		}
		// in_fmt.AvDumpFormat(0, fileName, 0)
		for {
			if ret := in_fmt.AvReadFrame(pkt); ret < 0 {
				break
			}
			// 从instream中获取一些参数
			// inStream := in_fmt.Streams()[0]
			// if pkt.StreamIndex() != 0 {
			// 	fmt.Println("pkt.StreamIndex error")
			// }
			// pkt.SetPts(int64(frame_index*2))
			if ret := ofmt_ctx.AvInterleavedWriteFrame(pkt); ret < 0 { //写入图像到视频
				fmt.Println("pkt write error")
				os.Exit(1)
			}
			pkt.AvPacketUnref()
		}
		fmt.Printf("Write %8d frames to output file\n", frame_index)
		// 关闭输入流上下文
		in_fmt.AvformatCloseInput()
	}
	ofmt_ctx.AvDumpFormat(0, out_filename, 1)
	ofmt_ctx.AvWriteTrailer() //写文件尾
	// ofmt_ctx&&ofmt_ctx.Flags()&avformat.AvGetFrameFilename()
	// avformat.AvformatAllocContext().Pb().Close() //关闭视频文件
	// 关闭输出流上下文
	ofmt_ctx.AvformatCloseInput()
	ofmt_ctx.Pb().Close() //关闭视频文件
	ofmt_ctx.AvformatFreeContext()
}

func img2Video3() {
	ofmt_ctx := avformat.AvformatAllocContext()
	out_filename := "out2.mp4"
	avformat.AvRegisterAll() //初始化解码器和复用器
	avformat.AvformatAllocOutputContext2(&ofmt_ctx, nil, "mp4", out_filename)
	out_stream := add_vidio_stream(ofmt_ctx)
	ofmt_ctx.AvDumpFormat(0, out_filename, 1)
	flag := ofmt_ctx.Flags() & avcodec.AV_CODEC_ID_MJPEG
	if flag == 0 {
		pb, err := avformat.AvIOOpen(out_filename, avformat.AVIO_FLAG_WRITE)
		ofmt_ctx.SetPb(pb)
		if err != nil {
			fmt.Printf("Could not open output file %s", out_filename)
			os.Exit(1)
		}
	}

	ofmt_ctx.AvformatWriteHeader(nil)

	frame_index := 1
	// numBytes := uintptr(avcodec.AvpictureGetSize(avcodec.AV_PIX_FMT_YUV420P9, out_stream.Codec().GetWidth(),
	// 	out_stream.Codec().GetHeight()))
	// buffer := avutil.AvMalloc(numBytes)
	pkt := avcodec.AvPacketAlloc()
	pkt.AvInitPacket()
	pkt.SetFlags(pkt.Flags() | avcodec.AV_PKT_FLAG_KEY)
	pkt.SetStreamIndex(out_stream.Index()) //获取视频信息，为压入帧图像做准备
	var fileName string
	for frame_index <= 30 {
		fileName = fmt.Sprintf("./images/netFrame%d.jpg", frame_index)
		file, err := os.Open(fileName)
		fileBytes, _ := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("Open file failed")
		}

		//怎么去读文件指针？
		pkt.AvPacketFromData((*uint8)(unsafe.Pointer(&fileBytes)), len(fileBytes)) //将图像写入包
		if ret := ofmt_ctx.AvWriteFrame(pkt); ret < 0 {                            //写入图像到视频
			fmt.Println("pkt write error")
			os.Exit(1)
		}
		fmt.Printf("Write %8d frames to output file\n", frame_index)
		frame_index++
	}
	ofmt_ctx.AvDumpFormat(0, out_filename, 1)
	// pkt.AvFreePacket()
	ofmt_ctx.AvWriteTrailer() //写文件尾
	// ofmt_ctx&&ofmt_ctx.Flags()&avformat.AvGetFrameFilename()
	ofmt_ctx.Pb().Close()          //关闭视频文件
	ofmt_ctx.AvformatFreeContext() //
}
func add_vidio_stream(oc *avformat.Context) *avformat.Stream {
	st := oc.AvformatNewStream(nil) //AVStream
	//申请AVStream->codec(AVCodecContext对象)空间并设置默认值(由avcodec_get_context_defaults3()设置
	codec := avcodec.AvcodecFindDecoder(avcodec.CodecId(avcodec.AV_CODEC_ID_MJPEG))
	codcCtx := codec.AvcodecAllocContext3()
	(*avcodec.Context)(unsafe.Pointer(st.Codec())).AvcodecCopyContext(codcCtx)
	st.Codec().SetBitRate(400000)
	st.Codec().SetWidth(1280)
	st.Codec().SetHeight(720)
	st.Codec().SetTimeBase(avcodec.NewRational(1, 25))
	st.Codec().SetPixelFormat(avcodec.AV_PIX_FMT_YUV420P9)
	return st
}

func img2Video2() {
	imgPath := "./images/netFrame1.jpg"
	dstPath := "./new.mp4"
	videoIndex := -1

	/*************************1.图片解封装******************************************************************************************************/
	// 创建jpg的解封装上下文
	in_fmt := avformat.AvformatAllocContext()
	if avformat.AvformatOpenInput(&in_fmt, imgPath, nil, nil) != 0 {
		fmt.Printf("Unable to open file %s\n", imgPath)
		os.Exit(1)
	}
	// Retrieve stream information
	if in_fmt.AvformatFindStreamInfo(nil) < 0 {
		fmt.Println("Couldn't find stream information")
		os.Exit(1)
	}

	in_fmt.AvDumpFormat(0, imgPath, 0)
	/*************************end***********************************************************************************************************************************/
	/*************************2.图片解码******************************************************************************************************/
	//创建解码器
	//初始化解码器上下文用于对jpg进行解码
	var de_ctx *avcodec.Context
	for i := 0; i < int(in_fmt.NbStreams()); i++ {
		stream := in_fmt.Streams()[i]
		/** 对于jpg图片来说，它里面就是一路视频流，所以媒体类型就是AVMEDIA_TYPE_VIDEO
		 */
		if stream.CodecParameters().AvCodecGetType() == avformat.AVMEDIA_TYPE_VIDEO {
			// 用输出流上下文创建输出流
			// outStream = outFormat.AvformatNewStream(nil)
			// Find the decoder for the video stream
			codec := avcodec.AvcodecFindDecoder(avcodec.CodecId(stream.Codec().GetCodecId()))
			if codec == nil {
				fmt.Println("Unsupported codec!")
				os.Exit(1)
			}
			de_ctx = codec.AvcodecAllocContext3()
			// 设置解码参数;文件解封装的AVStream中就包括了解码参数，这里直接流中拷贝即可
			if de_ctx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(stream.Codec()))) != 0 {
				fmt.Println("Couldn't copy codec context")
				os.Exit(1)
			}
			// 初始化解码器及解码器上下文

			if de_ctx.AvcodecOpen2(codec, nil) < 0 {
				fmt.Println("avcodec_open2() fail")
				os.Exit(1)
			}
			videoIndex = i
			break
		}
	}
	/*************************3.创建mp4文件封装器****************************************************************************************************************************************************************************/
	//2. 创建输出视频上下文
	ou_fmt := avformat.AvformatAllocContext()
	if avformat.AvformatAllocOutputContext2(&ou_fmt, nil, "mp4", dstPath) < 0 {
		fmt.Println("AvformatAllocOutputContext2 error")
		os.Exit(1)
	}
	// 添加视频流
	ou_stream := ou_fmt.AvformatNewStream(nil)
	video_ou_index := ou_stream.Index()
	ou_fmt.AvDumpFormat(0, dstPath, 1)
	/*************************4.创建h264的编码器及编码器上下文****************************************************************************************************************************************************************************/
	///1 创建编码器 查找编码器
	//创建h264编码器及编码上下文
	en_codec := avcodec.AvcodecFindEncoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	if en_codec == nil {
		fmt.Println("Unsupported codec!")
		os.Exit(1)
	}
	// Alloc context
	en_ctx := en_codec.AvcodecAllocContext3()
	// 设置h264编码参数
	inStream := in_fmt.Streams()[videoIndex]
	//配置编码器上下文的成员
	//设置画面宽度 高度 输出像素格式（pix_fmt)
	en_ctx.SetEncodeParams(inStream.Codec().GetWidth(), inStream.Codec().GetHeight(), inStream.Codec().GetPixelFormat())
	//设置帧率，num为分子，den为分母，如果是1/25则表示25帧/s
	en_ctx.SetTimebase(1, 25)
	// 设置全局编码信息
	// EnCodecCtx.Set= EnCodecCtx.Flags() | avcodec.AV_CODEC_FLAG_GLOBAL_HEADER
	//打开编码器

	if en_ctx.AvcodecOpen2(en_codec, nil) < 0 {
		fmt.Println("Could not open codec")
		os.Exit(1)
	}
	/*************************4.初始化mp4文件封装器****************************************************************************************************************************************************************************/
	//设置视频流参数;直接从编码器上下文拷贝 将AVCodecContext信息拷贝到AVCodecParameterst结构体中
	if (*avcodec.Context)(unsafe.Pointer(ou_stream.Codec())).AvcodecCopyContext(en_ctx) < 0 {
		fmt.Println("Couldn't copy codec context")
		os.Exit(1)
	}
	fmt.Println("*****打印AVFormatContext的内容************************")
	ou_fmt.AvDumpFormat(0, dstPath, 1)
	// 初始化视频封装器输出缓冲区
	pb, err := avformat.AvIOOpen(dstPath, avformat.AVIO_FLAG_WRITE)
	ou_fmt.SetPb(pb)
	if err != nil {
		fmt.Println("Could not open output file!")
		os.Exit(1)
	}

	// initialize SWS context for software scaling
	swsCtx := swscale.SwsGetcontext(de_ctx.Width(), de_ctx.Height(), swscale.PixelFormat(de_ctx.PixFmt()),
		en_ctx.Width(), en_ctx.Height(), swscale.PixelFormat(en_ctx.PixFmt()), 0, nil, nil, nil) //avcodec.SWS_BICUBIC

	//写视频文件头
	if ou_fmt.AvformatWriteHeader(nil) < 0 {
		fmt.Println("avformat_write_header failed!")
		os.Exit(1)
	}

	// 创建编解码用的AVFrame
	deFrame := avutil.AvFrameAlloc()
	enFrame := avutil.AvFrameAlloc()
	avutil.AvSetFrame(enFrame, en_ctx.Width(), en_ctx.Height(), int(en_ctx.PixFmt()))
	avutil.AvFrameGetBuffer(enFrame, 0)
	avutil.AvFrameMakeWritable(enFrame)

	//开始写入

	for frameNumber := 1; frameNumber < 30; frameNumber++ {
		//重新创建输入流
		fileName := fmt.Sprintf("./images/netFrame%d.jpg", frameNumber)
		// 创建jpg的解封装上下文
		in_fmt := avformat.AvformatAllocContext()
		if avformat.AvformatOpenInput(&in_fmt, fileName, nil, nil) != 0 {
			fmt.Printf("Unable to open file %s\n", fileName)
			os.Exit(1)
		}
		// Retrieve stream information
		if in_fmt.AvformatFindStreamInfo(nil) < 0 {
			fmt.Println("Couldn't find stream information")
			os.Exit(1)
		}

		in_fmt.AvDumpFormat(0, imgPath, 0)
		////////////////////////////
		for { //while 循环
			in_pkt := avcodec.AvPacketAlloc()
			if in_fmt.AvReadFrame(in_pkt) != 0 { //是否是MP4中的帧
				break
			}
			if in_pkt.StreamIndex() != videoIndex { //是否是mp4中的视频流
				continue
			}
			//先解码
			//doDecode
			doDecode(de_ctx, in_pkt, deFrame, enFrame, swsCtx, en_ctx, ou_fmt, video_ou_index)
			in_pkt.AvPacketUnref()
		}
		//刷新解码缓冲区
		//doDecode
		// doDecode(nil, nil, nil, nil, nil, nil, nil, video_ou_index)
		ou_fmt.AvWriteTrailer()
		fmt.Println("转换视频结束")
		os.Exit(1)
	}
}
