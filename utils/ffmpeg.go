package utils

// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// int startCmd(const char* cmd) {
//     return system(cmd);
// }

import "C"

import (
	"errors"
	"fmt"
	"tiktok/config"
)

// ChangeVideoToImage：视频转图片
type ChangeVideoToImage struct {
	InputPath  string // 输入视频路径
	OutputPath string // 输出图片路径
	StartTime  string // 开始时间
	KeepTime   string // 持续时间
	Filter     string // 过滤器
	FrameCount int64  // 帧数
	debug      bool   // 是否打印命令
}

// NewChangeVideoToImage：创建一个ChangeVideoToImage的对象
func NewChangeVideoToImage() *ChangeVideoToImage {
	return &videoChanger
}

var videoChanger ChangeVideoToImage

// ffmpeg的参数
const (
	inputVideoPathOption = "-i"        // 输入视频路径
	startTimeOption      = "-ss"       // 开始时间
	keepTimeOption       = "-t"        // 持续时间
	videoFilterOption    = "-vf"       // 过滤器
	formatToImageOption  = "-f"        // 格式转换
	autoReWriteOption    = "-y"        // 自动覆盖
	framesOption         = "-frames:v" // 帧数
)

var (
	defaultVideoSuffix = ".mp4" // 默认视频后缀
	defaultImageSuffix = ".jpg" // 默认图片后缀
)

// ChangeVideoDefaultSuffix：修改默认视频后缀
func ChangeVideoDefaultSuffix(suffix string) {
	defaultVideoSuffix = suffix
}

// ChangeImageDefaultSuffix：修改默认图片后缀
func ChangeImageDefaultSuffix(suffix string) {
	defaultImageSuffix = suffix
}

// GetDefaultVideoSuffix：获取默认视频后缀
func GetDefaultImageSuffix() string {
	return defaultImageSuffix
}

// paramJoin：连接ffmpeg命令的参数
func paramJoin(s1, s2 string) string {
	return fmt.Sprintf(" %s %s ", s1, s2)
}

func (v *ChangeVideoToImage) Debug() {
	v.debug = true
}

// GetQueryString：构建ffmpeg命令(根据Video2Image对象的字段值来构建命令)
func (v *ChangeVideoToImage) GetQueryString() (ret string, err error) {
	if v.InputPath == "" || v.OutputPath == "" {
		err = errors.New("输入输出路径未指定")
		return
	}
	// 创建ffmpeg的路径ret
	ret = config.Global.FfmpegPath
	// 拼接ffmpeg的参数：inputVideoPathOption、v.InputPath、formatToImageOption
	ret += paramJoin(inputVideoPathOption, v.InputPath)
	ret += paramJoin(formatToImageOption, "image2")
	// 填充ffmpeg的参数：startTimeOption、keepTimeOption、videoFilterOption、framesOption、autoReWriteOption
	if v.Filter != "" {
		ret += paramJoin(videoFilterOption, v.Filter)
	}
	if v.StartTime != "" {
		ret += paramJoin(startTimeOption, v.StartTime)
	}
	if v.KeepTime != "" {
		ret += paramJoin(keepTimeOption, v.KeepTime)
	}
	if v.FrameCount != 0 {
		ret += paramJoin(framesOption, fmt.Sprintf("%d", v.FrameCount))
	}
	ret += paramJoin(autoReWriteOption, v.OutputPath)
	return
}

// ExecCommand：执行ffmpeg命令
func (v *ChangeVideoToImage) ExecCommand(cmd string) error {
	// if v.debug {
	// 	log.Println(cmd)
	// }
	// // 将Go的字符串cmd转换为C语言的字符串
	// cCmd := C.CString(cmd)
	// // 放为C字符串分配的内存，因为Go的垃圾回收器不会自动管理C语言的内存
	// defer C.free(unsafe.Pointer(cCmd))
	// // 调用C语言的函数startCmd执行ffmpeg命令
	// status := C.startCmd(cCmd)
	// if status != 0 {
	// 	return errors.New("视频切截图失败")
	// }
	return nil
}
