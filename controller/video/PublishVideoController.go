package video

import (
	"net/http"
	"path/filepath"
	"tiktok/config"
	"tiktok/models"
	"tiktok/utils"

	"tiktok/service/video"

	"github.com/gin-gonic/gin"
)

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
	//文件后缀为图片
	pictureIndexMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)

// PublishVideoController：发布视频，并截取一帧画面作为封面
func PublishVideoController(c *gin.Context) {
	//准备参数
	//1.从上下文中获取userId
	rawId, _ := c.Get("user_id")
	//检查userId是否为int64类型
	userId, ok := rawId.(int64)
	if !ok {
		PublishVideoError(c, "解析UserId出错")
		return
	}
	//2.从上下文中获取video_title
	title := c.PostForm("title")
	//3.从上下文中获取 多部分表单
	form, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}

	//支持多文件上传
	files := form.File["data"]
	//遍历文件
	for _, file := range files {
		suffix := filepath.Ext(file.Filename)    //得到后缀
		if _, ok := videoIndexMap[suffix]; !ok { //判断是否为视频格式
			PublishVideoError(c, "不支持的视频格式")
			continue
		}
		name := utils.NewFileName(userId) //根据userId得到唯一的文件名
		filename := name + suffix
		// 保存文件的路径
		savePath := filepath.Join(config.Global.StaticSourcePath, filename)
		// 保存文件到本地
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		// 调用utils层，截取一帧画面作为封面
		err = utils.SaveImageFromVideo(name, false)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		// 调用service层，将视频信息保存到数据库
		err := video.PostVideo(userId, filename, name+utils.GetDefaultImageSuffix(), title)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		PublishVideoOk(c, file.Filename+"上传成功")
	}
}

// 发送错误
func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.StatusCodeResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// 发送成功
func PublishVideoOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.StatusCodeResponse{
		StatusCode: 0,
		StatusMsg:  msg,
	})
}
