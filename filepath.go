package gbase

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const longPathRoot = "static" //todo move to global.conf

type Filepath string

func (this Filepath) AbsPath() (string, error) {
	// 获取当前执行文件的路径
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// 获取执行文件所在的目录
	dir := filepath.Dir(executable)
	return fmt.Sprintf("%s/%s", dir, this.WrapLong().ToString()), nil //fmt.Sprintf("当前执行文件所在的目录:", dir)
}

func (this Filepath) WrapLong() Filepath {
	return Filepath(fmt.Sprintf("%s%s", longPathRoot, this))
}

func (this Filepath) UnwrapLong() Filepath {
	return Filepath(strings.Replace(string(this), longPathRoot, "", 1))
}

func (this Filepath) ToString() string {
	return string(this)
}

func (this Filepath) FromString(filePath string) Filepath {
	return Filepath(filePath)
}

func (this Filepath) SubDir(dirName string) Filepath {
	return Filepath(fmt.Sprintf("%s%s", this, dirName))
}

// img/xx/x.png =>img/xx/
func (this Filepath) RemoveLast() Filepath {
	arr := strings.Split(string(this), "/")
	newStr := "/"
	if len(arr) > 1 {
		newStr = strings.Join(arr[:len(arr)-1], "/") + "/"
	}
	return Filepath(newStr)
}

func (this Filepath) Rm() (err error) {
	err = os.Remove(this.ToString())
	return
}

func (this Filepath) RmDir() (err error) {
	err = os.RemoveAll(this.ToString())
	return
}

// 文件改名
func (this Filepath) Mv(name string) (err error) {
	return os.Rename(this.ToString(), name)
}

func (this Filepath) CopyTo(distPath Filepath) (err error) {
	if err = distPath.MkDir(); err != nil {
		return err
	}
	sourceFile, err := os.Open(this.ToString())
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(distPath.ToString())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

//func (this Filepath) RemoveFileName() Filepath {
//	directory := filepath.Dir(path)
//	strings.Replace()filepath.Base(this.ToString())
//}

func (this Filepath) GetDirPath() Filepath {
	return Filepath(filepath.Dir(this.ToString()))
	//tmpPath := this
	//if this.HasFileName() {
	//	tmpPath = this.RemoveLast()
	//}
	//return tmpPath
}

// 生成与文件名同名的子文件夹，与当前文件在同一目录下
func (this Filepath) MkSameNameSubDir() (subDir Filepath, err error) {
	subDir = Filepath(fmt.Sprintf("%s/%s/", this.GetDirPath().ToString(), this.GetFileNameWithoutExt()))
	hlog.Debugf("subDir:%s", subDir)
	err = subDir.MkDir()
	return
}

// 此时this必须是相对路径或绝对路径，不能为短路径
func (this Filepath) ExistDir() bool {
	_, err := os.Stat(this.GetDirPath().ToString())
	if err == nil {
		return true
	}
	return false
}

func (this Filepath) ExistFile() bool {
	_, err := os.Stat(this.ToString())
	if err == nil {
		return true
	}
	return false
}

func (this Filepath) MkDir() error {
	//exist := this.ExistDir()
	if !this.ExistDir() {
		return os.MkdirAll(this.GetDirPath().ToString(), os.ModePerm)
	}
	return nil
}

func (this Filepath) HasFileName() bool {
	return path.Ext(this.ToString()) != ""
}

// 返回的是 xxx.png
func (this Filepath) GetFileName() string {
	return filepath.Base(this.ToString())
}

// return .png
func (this Filepath) GetFileExt() string {
	return path.Ext(this.ToString())
}

// 同时适用于文件夹目录和URL，其它方法也是
func (this Filepath) GetFileNameWithoutExt() string {
	return strings.TrimSuffix(this.GetFileName(), this.GetFileExt())
}

func (this Filepath) ModFileName(name string) Filepath {
	return Filepath(fmt.Sprintf("%s/%s%s", this.GetDirPath().ToString(), name, this.GetFileExt()))
}

func (this Filepath) ModFileNameWithExt(name string, ext string) Filepath {
	return Filepath(fmt.Sprintf("%s/%s%s", this.GetDirPath().ToString(), name, ext))
}

func (this Filepath) CheckOwn(uid int64) bool {
	//path := "/path/to/123/directory"
	regex := regexp.MustCompile(`/(\d+)/`)
	match := regex.FindStringSubmatch(this.ToString())
	if len(match) > 1 {
		return match[1] == fmt.Sprintf("%d", uid)
	}
	return false
}

// 调用后需要defer file.Close()
func (this Filepath) CreateFile() (file *os.File, err error) {
	return os.Create(this.ToString())
}

// 从当前路径读取文件,只读模式
// 调用后需要defer file.Close()
func (this Filepath) GetFile() (file *os.File, err error) {
	if err = this.MkDir(); err != nil {
		return
	}
	if !this.ExistFile() {
		return this.CreateFile()
	} else {
		return os.Open(this.ToString())
		//return os.OpenFile(this.ToString(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}
}

// 以可写模式打开文件
// 调用后需要defer file.Close()
func (this Filepath) OpenFileForWrite() (file *os.File, err error) {
	if err = this.MkDir(); err != nil {
		return
	}
	if !this.ExistFile() {
		return this.CreateFile()
	} else {
		return os.OpenFile(this.ToString(), os.O_WRONLY|os.O_CREATE, 0644)
	}
}

func (this Filepath) ReadFile() (content string, err error) {
	var data []byte
	if data, err = os.ReadFile(this.ToString()); err == nil {
		content = string(data)
	}
	return

}

// 从当前路径读取图片文件到base64字符串里
func (this Filepath) GetBase64() (b64 string, err error) {
	// 读取图片文件
	data, err := os.ReadFile(this.ToString())
	if err != nil {
		hlog.Errorf("Filepath.LoadB64Img 读取文件失败")
		return
	}

	// 转为 Base64 编码
	b64 = base64.StdEncoding.EncodeToString(data)
	return
}

func (this Filepath) GetB64Chunks(chunkSize int) (chunks []string, err error) {
	chunks = make([]string, 0)
	//file, err := this.GetFile()
	//if err != nil {
	//	return
	//}
	//defer file.Close()

	file, err := os.Open(this.ToString())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffer reader
	reader := bufio.NewReader(file)

	// Create a buffer to store chunk data

	for {
		chunk := make([]byte, chunkSize)
		// Read chunk data from the file
		n, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			return nil, err
		}

		// If EOF is reached, break the loop
		if err == io.EOF {
			break
		}

		chunks = append(chunks, base64.StdEncoding.EncodeToString(chunk[:n]))
	}
	return
}

func (this Filepath) SaveBase64Img(base64Image string) (err error) {
	var (
		data []byte
		file *os.File
	)
	if err = this.MkDir(); err != nil {
		return errors.New(fmt.Sprintf("make dir failed,err:%s", err.Error()))
	}
	// 解码Base64字符串
	if data, err = base64.StdEncoding.DecodeString(base64Image); err != nil {
		return
	}

	// 创建文件
	file, err = os.Create(this.ToString())
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入数据
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	// 确保数据被写入磁盘
	err = file.Sync()
	if err != nil {
		return err
	}

	// 将字节数组写入本地文件
	//err = os.WriteFile(localPath, data, 0644)
	//err = os.WriteFile(localPath, data, 0666)
	//if err != nil {
	//	return err
	//}

	return
}

// 将内容写入到文件
func (this Filepath) SaveContent(content string) error {
	var (
		err error
	)

	// 创建文件
	file, err := this.GetFile()
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入数据
	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	// 确保数据被写入磁盘
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

// 保存网络文件到自身当前路径
func (this Filepath) SaveFileFromUrl(urlStr string) error {
	var (
		err      error
		resp     *http.Response
		body     []byte
		fileName string
	)

	if resp, err = http.Get(urlStr); err != nil {
		return err
	} else if body, err = io.ReadAll(resp.Body); err != nil {
		return err
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		hlog.Errorf("Filepath.SaveImgFromUrl get file ext failed,err:%s", err.Error())
		return err
	}
	ext := path.Ext(parsedURL.Path)
	if err = this.MkDir(); err != nil {
		return errors.New(fmt.Sprintf("make dir failed,err:%s", err.Error()))
	}
	if !this.HasFileName() {
		fileName = fmt.Sprintf("%s%s", Sn(6), ext)
	}

	if out, err := os.Create(fmt.Sprintf("%s%s", this.ToString(), fileName)); err != nil {
		return err
	} else if _, err = io.Copy(out, bytes.NewReader(body)); err != nil {
		return err
	}

	return nil
}

//// 保存网络图片到自身当前路径
//func (this Filepath) SaveImgFromHttBody(urlStr string) error {
//	var (
//		err      error
//		resp     *http.Response
//		body     []byte
//		fileName string
//	)
//
//	if resp, err = http.Get(urlStr); err != nil {
//		return err
//	} else if body, err = io.ReadAll(resp.Body); err != nil {
//		return err
//	}
//	parsedURL, err := url.Parse(urlStr)
//	if err != nil {
//		hlog.Errorf("Filepath.SaveImgFromUrl get file ext failed,err:%s", err.Error())
//		return err
//	}
//	ext := path.Ext(parsedURL.Path)
//	if err = this.MkDir(); err != nil {
//		return errors.New(fmt.Sprintf("make dir failed,err:%s", err.Error()))
//	}
//	if !this.HasFileName() {
//		fileName = fmt.Sprintf("%s%s", util.Sn(6), ext)
//	}
//
//	if out, err := os.Create(fmt.Sprintf("%s%s", this.ToString(), fileName)); err != nil {
//		return err
//	} else if _, err = io.Copy(out, bytes.NewReader(body)); err != nil {
//		return err
//	}
//
//	return nil
//}

// 图片转成RGBA格式
func (this Filepath) ToRgbaPng(alpha uint8) (dstFilePath Filepath, err error) {
	var (
		file, outFile *os.File
		img           image.Image
	)
	// 打开原始RGB图片文件
	if file, err = os.Open(this.ToString()); err != nil {
		return
	}
	defer file.Close()

	// 解码图片
	img, _, err = image.Decode(file)
	if err != nil {
		return
	}

	// 创建一个新的RGBA图片，大小与原始图片相同
	rgba := image.NewRGBA(img.Bounds())

	// 将RGB图片转换为RGBA图片
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// 遍历每个像素，设置全透明度
	for i := 0; i < len(rgba.Pix); i += 4 {
		rgba.Pix[i+3] = alpha // 设置 Alpha 通道为 0，表示全透明
	}

	// 创建或覆盖目标文件
	dstFilePath = Filepath(fmt.Sprintf("%s%s_rgba.png", this.GetDirPath(), this.GetFileNameWithoutExt()))
	if outFile, err = os.Create(dstFilePath.ToString()); err != nil {
		return
	}
	defer outFile.Close()

	// 将转换后的RGBA图片编码并保存到文件中
	if err = png.Encode(outFile, rgba); err != nil {
		return
	}

	return
}

// 图片转成灰度图像
func (this Filepath) ToGrayJpg() (dstFilePath Filepath, err error) {
	var (
		file, outFile *os.File
		img           image.Image
	)
	// 打开原始RGB图片文件
	if file, err = os.Open(this.ToString()); err != nil {
		return
	}
	defer file.Close()

	// 解码图片
	img, _, err = image.Decode(file)
	if err != nil {
		return
	}

	// 创建一个新的RGBA图片，大小与原始图片相同
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// 遍历每个像素进行转换
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	// 创建或覆盖目标文件
	dstFilePath = Filepath(fmt.Sprintf("%s%s_gray.png", this.GetDirPath(), this.GetFileNameWithoutExt()))
	if outFile, err = os.Create(dstFilePath.ToString()); err != nil {
		return
	}
	defer outFile.Close()

	// 将转换后的灰度图片编码并保存到文件中
	if err = png.Encode(outFile, grayImg); err != nil {
		return
	}
	//if err = jpeg.Encode(outFile, grayImg, nil); err != nil {
	//	return
	//}

	return
}

//-----------------网址类方法

func (this Filepath) GetHttpFileSize() (int64, error) {
	// 发送 HTTP HEAD 请求
	resp, err := http.Head(this.ToString())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 检查请求是否成功
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	// 获取文件大小
	size := resp.ContentLength
	if size < 0 {
		return 0, fmt.Errorf("Unable to determine file size")
	}

	return size, nil
}
