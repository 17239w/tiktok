package video

//utils "tiktok/utils"

var (
	// 文件后缀为视频
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
	// 文件后缀为图片
	pictureIndexMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)
