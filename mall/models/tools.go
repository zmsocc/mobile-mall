package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//时间戳转换成日期
func UnixToTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

//日期转换成时间戳 2020-05-02 15:04:05
func DateToUnix(str string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

//获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

//获取时间戳-毫秒级
func GetUnixMilli() int64 {
	return time.Now().UnixMilli()
}

//获取当前的日期
func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

//获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// Md5 加密
func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 将字符串转换为 int 类型
func Int(str string) (int, error) {
	n, err := strconv.Atoi(str)
	return n, err
}

// 将 int 类型转换为 string 类型
func String(n int) string {
	str := strconv.Itoa(n)
	return str
}

// 上传图片
func UploadImg(ctx *gin.Context, picName string) (string, error) {
	// 1.获取上传的文件
	file, err := ctx.FormFile(picName)
	if err != nil {
		return "", errors.New("获取文件失败")
	}

	// 2.获取后缀名，判断类型是否正确	.jpg .png .gif .jpeg
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool {
		".jpg": true,
		".png": true,
		".gif": true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("图片后缀名不合法")
	}

	// 3.创建图片保存目录， static/upload/20250717
	day := GetDay()
	dir := "./static/upload/" + day
	// 创建目录
	_ = os.Mkdir(dir, 0666)

	// 4.生成文件名称和文件保存的目录， 111111111.jpg
	fileName := strconv.FormatInt(GetUnix(), 10) + extName 

	// 5.执行上传
	dst := path.Join(dir, fileName)
	ctx.SaveUploadedFile(file, dst)
	return dst, nil
}
