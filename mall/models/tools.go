package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
	. "github.com/hunterhug/go_image"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
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

//获取时间戳-纳秒级
func GetUnixNano() int64 {
	return time.Now().UnixNano()
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

// 把字符串解析成 html
func Str2Html(str string) template.HTML {
	return template.HTML(str)
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

// 将字符串转换为 float64 类型
func Float64(str string) (float64, error) {
	n, err := strconv.ParseFloat(str, 64)
	return n, err
}

// 通过列获系统设置里面的值
func GetSettingFromColum(columnName string) string {
	setting := Setting{}
	DB.First(&setting)
	// 反射来获取
	v := reflect.ValueOf(setting)
	val := v.FieldByName(columnName).String()
	return val
}

// 获取 Oss 的状态
func GetOssStatus() int{
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}
	ossStatus := config.Section("oss").Key("status").String()
	status, _ := Int(ossStatus)
	return status
}

// 上传图片
func UploadImg(ctx *gin.Context, picName string) (string, error) {
	if GetOssStatus() == 1 {
		return OssUploadImg(ctx, picName)
	}
	return LocalUploadImg(ctx, picName)
}

// 格式化图片， 判断是否开启了 Oss
func FormatImg(str string) string {
	if GetOssStatus() == 1 {
		return GetSettingFromColum("OssDomain") + str
	}
	return "/" + str
}

func Sub(a int, b int) int {
	return a - b
}

func LocalUploadImg(ctx *gin.Context, picName string) (string, error) {
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
	fileName := strconv.FormatInt(GetUnixMilli(), 10) + extName 

	// 5.执行上传
	dst := path.Join(dir, fileName)
	ctx.SaveUploadedFile(file, dst)
	return dst, nil
}

func OssUploadImg(ctx *gin.Context, picName string) (string, error) {
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
	dir := "static/upload/" + day

	// 4.生成文件名称和文件保存的目录， 111111111.jpg
	fileName := strconv.FormatInt(GetUnixMilli(), 10) + extName 

	// 5.执行上传
	dst := path.Join(dir, fileName)
	return OssUpload(file, dst)
}

// 封装 Oss 上传的方法
func OssUpload(file *multipart.FileHeader, dst string) (string, error) {
	// 创建 OssClient 实例
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "", "")
	if err != nil {
		return "", err
	}
	// 获取存储空间
	bucket, err := client.Bucket("digital-mall-occ")
	if err != nil {
		return "", err
	}
	// 读取本地文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 上传文件流
	err = bucket.PutObject(dst, src)
	if err != nil {
		return "", err
	}
	return dst, nil
}

// 生成商品缩略图
func ResizeGoodsImage(filename string) {
	extname := path.Ext(filename)
	thumbnailSize := GetSettingFromColum("ThumbnailSize")
	thumbnailSizeSlice := strings.Split(thumbnailSize, ",")
	for i := 0; i < len(thumbnailSizeSlice); i++ {
		savepath := filename + "_" + thumbnailSizeSlice[i] + "x" + thumbnailSizeSlice[i] + extname
		w, _ := Int(thumbnailSizeSlice[i])
		err := ThumbnailF2F(filename, savepath, w, w)
		if err != nil {
			// 写个日志模块，处理日志
			fmt.Println(err)
		}
	}
}