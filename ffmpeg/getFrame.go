package ffmpeg

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
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"unsafe"

	"github.com/giorgisio/goav/swscale"
	"gocv.io/x/gocv"

	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
)

// SaveFrame writes a single frame to disk as a PPM file
func SaveFramePPM(frame *avutil.Frame, width, height, frameNumber int) {
	// Open file
	fileName := fmt.Sprintf("./images/frame%d.ppm", frameNumber)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error Reading")
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
		fmt.Println("Error Reading")
		fmt.Println(err)
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
	//img image.Image
	err = jpeg.Encode(file, img, nil)
	if err != nil {
		fmt.Println("jpeg.Encode err: ", err)
		fmt.Println("Error jpeg.Encode")
	}
}

// SaveFrame writes a single frame to gocv.Mat
func SaveFrameTocvMat(frame *avutil.Frame, width, height int) (gocv.Mat, error) {

	bytes := make([]byte, width*height*3)
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
			bytes = append(bytes, byte(pixel[0]>>8), byte(pixel[1]>>8), byte(pixel[2]>>8))
			img.SetRGBA(x, y, color.RGBA{pixel[0], pixel[1], pixel[2], 0xff})
		}
	}
	//img image.Image
	return gocv.NewMatFromBytes(width, height, gocv.MatTypeCV8UC3, bytes)
}

// SaveFrame writes a single frame to buffer
func SaveFrame2Buffer(frame *avutil.Frame, width, height, frameNumber int) (*bytes.Buffer, error) {
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

	buff := new(bytes.Buffer)
	err := jpeg.Encode(buff, img, nil)
	if err != nil {
		fmt.Println("jpeg.Encode err: ", err)
		fmt.Println("Error jpeg.Encode")
	}

	bigBuff := new(bytes.Buffer)
	if bigBuff.Cap() > buff.Len() {

	}
	return buff, err
}

// //img2Video
// func img2Video() {
// 	imgPath := "../output/outframe%3d.jpg"
// 	dstPath := "../new.mp4"
// 	videoIndex := -1
// 	// ??????jpg?????????????????????
// 	pFormatContext := avformat.AvformatAllocContext()
// 	if avformat.AvformatOpenInput(&pFormatContext, imgPath, nil, nil) != 0 {
// 		fmt.Printf("Unable to open file %s\n", imgPath)
// 		os.Exit(1)
// 	}
// 	// Retrieve stream information
// 	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
// 		fmt.Println("Couldn't find stream information")
// 		os.Exit(1)
// 	}
// 	var pCodecCtx *avcodec.Context
// 	//?????????????????????????????????????????????????????????jpg????????????
// 	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
// 		stream := pFormatContext.Streams()[i]
// 		/** ??????jpg????????????????????????????????????????????????????????????????????????AVMEDIA_TYPE_VIDEO
// 		 */
// 		if stream.CodecParameters().AvCodecGetType() == avformat.AVMEDIA_TYPE_VIDEO {
// 			// Find the decoder for the video stream
// 			pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(stream.Codec().GetCodecId()))
// 			if pCodec == nil {
// 				fmt.Println("Unsupported codec!")
// 				os.Exit(1)
// 			}
// 			// Alloc context
// 			pCodecCtx := pCodec.AvcodecAllocContext3()
// 			//Copy context ???????????????????????????????????????AVStream???????????????????????????????????????????????????
// 			if pCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(stream.Codec()))) != 0 {
// 				fmt.Println("Couldn't copy codec context")
// 				os.Exit(1)
// 			}
// 			// Open codec
// 			if pCodecCtx.AvcodecOpen2(pCodec, nil) < 0 {
// 				fmt.Println("Could not open codec")
// 				os.Exit(1)
// 			}
// 			videoIndex = i
// 			break
// 		}
// 	}
// 	///1 ??????????????? ???????????????
// 	//??????h264???????????????????????????
// 	encodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
// 	if encodec == nil {
// 		fmt.Println("Unsupported codec!")
// 		os.Exit(1)
// 	}
// 	// Alloc context
// 	EnCodecCtx := encodec.AvcodecAllocContext3()
// 	inStream := pFormatContext.Streams()[videoIndex]
// 	//?????????????????????????????????
// 	//?????????????????? ?????? ?????????????????????pix_fmt)
// 	EnCodecCtx.SetEncodeParams(inStream.Codec().GetWidth(), inStream.Codec().GetHeight(), inStream.Codec().GetPixelFormat())
// 	//???????????????num????????????den?????????????????????1/25?????????25???/s
// 	EnCodecCtx.SetTimebase(1, 25)
// 	// ????????????????????????
// 	// EnCodecCtx.Set= EnCodecCtx.Flags() | avcodec.AV_CODEC_FLAG_GLOBAL_HEADER
// 	//???????????????
// 	if EnCodecCtx.AvcodecOpen2(encodec, nil) < 0 {
// 		fmt.Println("Could not open codec")
// 		os.Exit(1)
// 	}
// 	//2. ???????????????????????????
// 	outFormat := avformat.AvformatAllocContext()
// 	if avformat.AvformatAllocOutputContext2(&outFormat, nil, "", dstPath) < 0 {
// 		fmt.Println("AvformatAllocOutputContext2 error")
// 		os.Exit(1)
// 	}
// 	//???????????????
// 	stream := outFormat.AvformatNewStream(nil)
// 	videoOutIndex := stream.Index()
// 	//?????????????????????;????????????????????????????????? ???AVCodecContext???????????????AVCodecParameterst????????????
// 	if EnCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(stream.Codec()))) < 0 {
// 		fmt.Println("Couldn't copy codec context")
// 		os.Exit(1)
// 	}
// 	fmt.Println("*****??????AVFormatContext?????????************************")
// 	outFormat.AvDumpFormat(0, dstPath, 1)
// 	// avformat.AvIOOpen(outFormat.Pb(),dstPath,)

// 	// initialize SWS context for software scaling
// 	swsCtx := swscale.SwsGetcontext(pCodecCtx.Width(), pCodecCtx.Height(), swscale.PixelFormat(pCodecCtx.PixFmt()),
// 		EnCodecCtx.Width(), EnCodecCtx.Height(), swscale.PixelFormat(EnCodecCtx.PixFmt()), avcodec.SWS_BICUBIC, nil, nil, nil)
// 	//??????????????????
// 	if outFormat.AvformatWriteHeader(nil) < 0 {
// 		fmt.Println("avformat_write_header failed!")
// 		os.Exit(1)
// 	}
// 	// ?????????????????????AVFrame
// 	deFrame := avutil.AvFrameAlloc()
// 	enFrame := avutil.AvFrameAlloc()
// 	avutil.AvSetFrame(enFrame, EnCodecCtx.Width(), EnCodecCtx.Height(), int(EnCodecCtx.PixFmt()))
// 	avutil.AvFrameGetBuffer(enFrame, 32)
// 	avutil.AvFrameMakeWritable(enFrame)

// 	inPkt := avcodec.AvPacketAlloc()
// 	for { //while ??????
// 		if pFormatContext.AvReadFrame(inPkt) == 0 {
// 			break
// 		}
// 		if inPkt.StreamIndex() != videoIndex {
// 			continue
// 		}
// 		//?????????
// 		doDecode(pCodecCtx, inPkt, deFrame, enFrame, swsCtx, EnCodecCtx, outFormat, videoOutIndex)
// 		inPkt.AvPacketUnref()
// 	}
// 	//?????????????????????
// 	doDecode(nil, nil, nil, nil, nil, nil, nil, videoOutIndex)
// 	outFormat.AvWriteTrailer()
// 	fmt.Println("??????????????????")
// 	os.Exit(1)
// }
// func doDecode(pCodecCtx *avcodec.Context, inPkt *avcodec.Packet, deFrame *avutil.Frame, enFrame *avutil.Frame,
// 	swsCtx *swscale.Context, EnCodecCtx *avcodec.Context, outFormat *avformat.Context, videoOutIndex int) {
// 	num_pts := 0
// 	//?????????
// 	pCodecCtx.AvcodecSendPacket(inPkt)
// 	for {
// 		ret := pCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(deFrame)))
// 		fmt.Println("ret:", ret)
// 		if ret == avutil.AvErrorEAGAIN {
// 			doEncode(nil, EnCodecCtx, outFormat, videoOutIndex)
// 			break
// 		} else if ret < 0 {
// 			break
// 		}
// 		//??????????????????????????????????????????
// 		if swscale.SwsScale2(swsCtx, avutil.Data(deFrame), avutil.Linesize(deFrame), 0, pCodecCtx.Height(),
// 			avutil.Data(enFrame), avutil.Linesize(enFrame)) < 0 {
// 			fmt.Println("swscale.SwsScale2 failed!")
// 			os.Exit(1)
// 		}
// 		// ?????????????????????pts???????????????en_ctx->time_base???{1,fps}???????????????pts???????????????????????????
// 		// en_frame->pts = num_pts++;
// 		(*avcodec.Packet)(unsafe.Pointer(enFrame)).SetPts(int64(num_pts))
// 		num_pts++
// 		doEncode(enFrame, EnCodecCtx, outFormat, videoOutIndex)
// 	}
// }

// func doEncode(enFrame *avutil.Frame, EnCodecCtx *avcodec.Context, outFormat *avformat.Context, videoOutIndex int) {
// 	EnCodecCtx.AvcodecSendPacket((*avcodec.Packet)(unsafe.Pointer(enFrame)))
// 	for {
// 		out_Pkt := avcodec.AvPacketAlloc()
// 		if EnCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(out_Pkt))) < 0 {
// 			out_Pkt.AvPacketUnref()
// 			break
// 		}
// 		//???????????????????????????????????????????????????
// 		stream := outFormat.Streams()[videoOutIndex]
// 		out_Pkt.AvPacketRescaleTs(EnCodecCtx.AvCodecGetPktTimebase(), stream.TimeBase())
// 		fmt.Printf("video pts %d", out_Pkt.Pts())
// 		outFormat.AvWriteFrame(out_Pkt)
// 	}

// }

//???????????????????????? ??????????????????
func VideoGetLocalImg(frameNum int, url string) {
	// url = "../src/sample.mp4"
	filename := url
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
			avp.AvpictureFill((*uint8)(buffer), avcodec.AV_PIX_FMT_RGB24,
				pCodecCtx.Width(), pCodecCtx.Height())

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
								fmt.Println("response error = - 11") //ReadEagain??????
								break
							}
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return
						}

						if frameNumber <= frameNum {
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

//?????????????????? ??????????????????
func VideoGetNetImg(frameNum int, url string) {
	//???????????????
	// url="rtmp://192.168.20.221:30200/live/2dfp52anvad7g"
	//??????????????????
	filename := url
	avformat.AvformatNetworkInit()
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
			avp.AvpictureFill((*uint8)(buffer), avcodec.AV_PIX_FMT_RGB24,
				pCodecCtx.Width(), pCodecCtx.Height())

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
								// fmt.Println("response error = - 11 ReadEagain") //ReadEagain??????
								break
							}
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return
						}

						if frameNumber <= frameNum {
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
								// fmt.Println("response error = - 11 ReadEagain") //ReadEagain??????
								break
							}
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return
						}

						if frameNumber <= frameNum {
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

//?????????????????? ?????????
func VideoGetNetImg2buff(frameNum int, url string) *bytes.Buffer {
	//???????????????
	// url="rtmp://192.168.20.221:30200/live/2dfp52anvad7g"
	//??????????????????
	filename := url
	avformat.AvformatNetworkInit()
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
			avp.AvpictureFill((*uint8)(buffer), avcodec.AV_PIX_FMT_RGB24,
				pCodecCtx.Width(), pCodecCtx.Height())

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
								// fmt.Println("response error = - 11 ReadEagain") //ReadEagain??????
								break
							}
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return nil
						}

						if frameNumber <= frameNum {
							// Convert the image from its native format to RGB
							swscale.SwsScale2(swsCtx, avutil.Data(pFrame),
								avutil.Linesize(pFrame), 0, pCodecCtx.Height(),
								avutil.Data(pFrameRGB), avutil.Linesize(pFrameRGB))

							// SaveFrameJpg(pFrameRGB, pCodecCtx.Width(), pCodecCtx.Height(), frameNumber)
							buff, err := SaveFrame2Buffer(pFrameRGB, pCodecCtx.Width(), pCodecCtx.Height(), frameNumber)
							if err != nil {
								fmt.Printf("SaveFrame2Buffer error")
								return nil
							}
							// Save the frame to disk
							fmt.Printf("Writing frame %d\n", frameNumber)
							return buff
						} else {
							return nil
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
								// fmt.Println("response error = - 11 ReadEagain") //ReadEagain??????
								break
							}
							fmt.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
							return nil
						}

						if frameNumber <= frameNum {
							// Convert the image from its native format to RGB
							swscale.SwsScale2(swsCtx, avutil.Data(pFrame),
								avutil.Linesize(pFrame), 0, pCodecCtx.Height(),
								avutil.Data(pFrameRGB), avutil.Linesize(pFrameRGB))

							// Save the frame to disk
							fmt.Printf("Writing frame %d\n", frameNumber)
							buff, err := SaveFrame2Buffer(pFrameRGB, pCodecCtx.Width(), pCodecCtx.Height(), frameNumber)
							//post
							//draw
							//????????????
							if err != nil {
								fmt.Printf("SaveFrame2Buffer error")
								return nil
							}
							return buff
						} else {
							return nil
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
	return nil
}
