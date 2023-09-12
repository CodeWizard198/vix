package vix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// Context 此context不是线程安全的
// 因为一个context应该是对应着给一个goroutine使用
// 不应该多个goroutine使用同一个context
type Context struct {
	// Req 请求request
	Req *http.Request
	// Resp 响应response
	Resp http.ResponseWriter
	// PathParam 路径参数
	// "/shop/:id" PathParam[id]
	PathParam          map[string]string
	ResponseData       []byte
	ResponseStatusCode int
	MatchRouter        string
	Form               url.Values
	// 请求参数缓存
	paramCache url.Values
}

// BindJSON 全局JSON解析器
// 可以通过UseNumber和DisallowUnknownFields方法来控制全局的json解析的力度
func (c *Context) BindJSON(target any) error {
	resolver := json.NewDecoder(c.Req.Body)
	if useNumber {
		resolver.UseNumber()
	}
	if disallowUnknownFields {
		resolver.DisallowUnknownFields()
	}
	return resolver.Decode(target)
}

// BindJSONbyOpt 控制单个json解析的力度
// numberUse：是否开启number模式，开启后以Number作为数字类型，默认是float64
// disallow：是否开启json检测，当json中有结构体未定义的对象时报错
func (c *Context) BindJSONbyOpt(target any, numberUse bool, disallow bool) error {
	resolver := json.NewDecoder(c.Req.Body)
	if numberUse {
		resolver.UseNumber()
	}
	if disallow {
		resolver.DisallowUnknownFields()
	}
	return resolver.Decode(target)
}

// GetFormValue 获取表单数据
// parsForm会自动缓存，无需当心每次解析表单数据时都要parse
func (c *Context) GetFormValue(key string) StringValue {
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{
			Err: err,
		}
	}
	form := c.Req.Form
	return StringValue{
		Value: form.Get(key),
	}
}

// GetParamValue 获取请求的请求参数
// 在第一次获取后，会自动缓存，下次获取时就可以不用再解析路径参数了
func (c *Context) GetParamValue(key string) StringValue {
	if c.paramCache == nil {
		c.paramCache = c.Req.URL.Query()
	}
	val, ok := c.paramCache[key]
	if !ok {
		return StringValue{
			Err: errors.New(fmt.Sprintf("请求参数没有该key:[%s]", key)),
		}
	}
	return StringValue{
		Value: val[0],
	}
}

// GetMoreParamValues 批量获取请求参数，可以传入多个key
// key string[] 不可以为空
// paramMap：key（为输入的key） value（为请求路径对于key的值） 类型默认为string
func (c *Context) GetMoreParamValues(key ...string) (map[string]StringValue, error) {
	if len(key) == 0 {
		return nil, errors.New("key不可以为空")
	}
	paramMap := make(map[string]StringValue)
	for _, k := range key {
		val := c.GetParamValue(k)
		if val.Err != nil {
			return nil, val.Err
		}
		paramMap[k] = val
	}
	return paramMap, nil
}

// FormFile 获得传输的文件
func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Req.FormFile(key)
}

// SaveFile 将文件保存到本地，采用断点续传的方式
// 若文件夹中存在同名文件，默认为文件信息一样，无需再次保存
// 默认文件名为传输时文件的文件名
// 若需要改变可以传入name，此name需要包含文件类型
func (c *Context) SaveFile(path string, file multipart.File, header *multipart.FileHeader, name ...string) error {
	// 创建path中不存在的文件夹
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	// 获取文件名 优先使用传入的文件名，若没有则使用原来的文件名
	fileName := header.Filename
	if len(name) != 0 {
		fileName = name[0]
	}
	// 构建文件路径 path + 文件名
	filePath := path + fileName
	// 判断是否存在该文件
	// 若不存在就创建一个
	exists, err := isFileExists(filePath)
	if !exists || err != nil {
		_, err = os.Create(filePath)
		if err != nil {
			return err
		}
	}
	// 获取创建/已有的文件信息
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	// 如果大小一样且文件名也一样的，就认为是同一个文件
	// 如果你不希望是这样，请配置不同的文件名
	if info.Size() == header.Size {
		return nil
	}
	fi, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {
			panic(err)
		}
	}(fi)
	// 定位到文件没写入的部分
	_, _ = file.Seek(info.Size(), 0)
	// 每次写入1MB
	data := make([]byte, 1024*1024)
	for {
		total, err := file.Read(data)
		if err == io.EOF {
			break
		}
		_, err = fi.Write(data[:total])
		if err != nil {
			return err
		}
	}
	return nil
}

// SetCookie 设置cookie
// expires：cookie的过期时间
func (c *Context) SetCookie(name, value string, expires ...int) {
	expireTime := 0
	if len(expires) != 0 {
		expireTime = expires[0]
	}
	http.SetCookie(c.Resp, &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: expireTime,
	})
}

// GetCookie 获取请求中的cookie
func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return "", err
	}
	value, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// DelCookie 删除请求中的cookie
func (c *Context) DelCookie(name string) {
	http.SetCookie(c.Resp, &http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
}

// JSON 相应json类型的数据
func (c *Context) JSON(code int, response any) {
	c.ResponseStatusCode = code
	data, err := json.Marshal(response)
	if err != nil {
		c.ResponseData = []byte("")
		c.ResponseStatusCode = http.StatusNotFound
		return
	}
	c.setHeaderJSON(data)
	c.ResponseData = data
	c.ResponseStatusCode = code
}

// STRING 响应string数据
func (c *Context) STRING(code int, response string) {
	c.ResponseStatusCode = code
	c.setHeaderSTRING(response)
	c.ResponseData = []byte(response)
}

// BYTE 响应[]byte数据
func (c *Context) BYTE(code int, response []byte) {
	c.ResponseStatusCode = code
	c.setHeaderBYTE(response)
	c.ResponseData = response
}
