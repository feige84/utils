package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpSend struct {
	Debug             bool //是否调试模式
	RequestUrl        string
	RequestNum        int64 //请求次数
	Method            string
	Header            map[string]string
	SendData          interface{}
	Format            string //json，form-data, stream
	XMLHttpRequest    bool
	ProxyStr          string
	ConnectTimeout    int64
	ReadWriteTimeout  int64
	DisableKeepAlives bool
}

type HttpResp struct {
	Data        []byte
	CookieObj   []*http.Cookie
	CookieMap   map[string]string
	RedirectURL string
}

//http请求处理
func HttpHandle(r *HttpSend) (*HttpResp, error) {

	//判断是否是有效URL
	urlInfo, err := url.Parse(r.RequestUrl)
	if err != nil {
		return nil, fmt.Errorf("url parse err: %s", err.Error())
	}

	//请求类型
	if r.SendData == nil {
		r.Method = http.MethodGet
	} else {
		r.Method = http.MethodPost
	}

	//初始化header
	if r.Header == nil {
		r.Header = make(map[string]string)
	}

	//是否异步请求，很多json接口都有这类似的判断。
	if r.XMLHttpRequest {
		r.Header["X-Requested-With"] = "XMLHttpRequest"
	}

	//user-agent
	if value, exist := r.Header["User-Agent"]; !exist || value == "" {
		r.Header["User-Agent"] = "(iPhone; CPU iPhone OS 13_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.8(0x17000820) NetType/WIFI Language/zh_CN"
	}

	var (
		req    *http.Request
		client *http.Client
	)
	if r.Method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, r.RequestUrl, nil)
	} else {
		if value, ok := r.SendData.(map[string]string); ok && len(value) > 0 && strings.ToUpper(r.Format) != "JSON" {
			//form-data
			r.Header["Content-Type"] = "application/x-www-form-urlencoded"
			sendBody := http.Request{}
			err = sendBody.ParseForm()
			if err == nil {
				for k, v := range value {
					sendBody.Form.Add(k, v)
				}
				req, err = http.NewRequest(http.MethodPost, r.RequestUrl, strings.NewReader(sendBody.Form.Encode()))
			}
		} else if value, ok := r.SendData.([]byte); ok && len(value) > 0 && strings.ToUpper(r.Format) == "STREAM" {
			//stream
			r.Header["Content-Type"] = "application/octet-stream;tt-data=a"
			req, err = http.NewRequest(http.MethodPost, r.RequestUrl, bytes.NewBuffer(value))
		} else {
			//json
			r.Header["Content-Type"] = "application/json;charset=utf-8"
			sendBody, err := json.Marshal(r.SendData)
			if err != nil {
				return nil, fmt.Errorf("json encode err: %s", err.Error())
			}
			req, err = http.NewRequest(http.MethodPost, r.RequestUrl, bytes.NewBuffer(sendBody))
		}
	}
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest ["+r.RequestUrl+"] err: %s", err.Error())
	}

	//设置header头
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			req.Header.Set(k, v)
		}
	}

	//设置主机名
	req.Host = urlInfo.Host

	//忽略https的证书
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //解决x509: certificate signed by unknown authority
		//Dial: TimeoutDialer(r.ConnectTimeout*time.Second, 5*time.Second), //设置超时，连接超时，读写超时。官方已不推荐用此方法。
	}
	//设置代理
	// http://username:password@http-dyn.abuyun.com:9020"
	if r.ProxyStr != "" {
		urlProxy, err := url.Parse(r.ProxyStr)
		if err == nil {
			transport.Proxy = http.ProxyURL(urlProxy)
		}
	}
	//是否禁用保持活动链接
	if r.DisableKeepAlives {
		transport.DisableKeepAlives = true
	}
	client = &http.Client{
		Transport: transport,
	}
	if r.ConnectTimeout > 0 {
		client.Timeout = time.Duration(r.ConnectTimeout) * time.Second
	}

	if r.Debug {
		Pr(r)
		Pr(req)
	}

	//开始请求
	resp, err := client.Do(req)
	r.RequestNum++
	if err != nil {
		reqErr := url.Error{Err: err}
		if reqErr.Timeout() {
			if r.RequestNum < 2 {
				//fmt.Println("尝试请求第", r.RequestNum, "次")
				tryResp, tryErr := HttpHandle(r)
				return tryResp, tryErr
			}
		}
		return nil, fmt.Errorf("http response err: %v", err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	if r.Debug {
		fmt.Println("http response code:", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return nil, fmt.Errorf("http response code: %d, %v", resp.StatusCode, err)
	}
	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http response unzip is failed: %s", err)
		}
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if r.Debug {
		fmt.Println("respData, err:", string(respData), err)
	}
	if err != nil {
		return nil, fmt.Errorf("http response data: %s, %v", string(respData), err)
	} else {
		r.RequestNum = 0
		httpResp := &HttpResp{
			Data:      respData,
			CookieObj: resp.Cookies(),
		}
		cookies := make(map[string]string)
		if len(httpResp.CookieObj) > 0 {
			for _, c := range httpResp.CookieObj {
				cookies[c.Name] = c.Value
			}
			httpResp.CookieMap = cookies
		}
		location := resp.Header.Get("Location")
		if location != "" {
			httpResp.RedirectURL = location
		} else {
			if resp.Request != nil && resp.Request.URL != nil {
				httpResp.RedirectURL = resp.Request.URL.String()
			}
		}
		return httpResp, nil
	}
}

//获取字符串内容
func (h *HttpResp) String() string {
	return string(h.Data)
}

//获取字节内容
func (h *HttpResp) Bytes() []byte {
	return h.Data
}

//获取返回的COOKIE对象
func (h *HttpResp) Cookies() []*http.Cookie {
	return h.CookieObj
}

//获取返回的COOKIE MAP
func (h *HttpResp) Cookie() map[string]string {
	return h.CookieMap
}

//获取重定向地址
func (h *HttpResp) Redirect() string {
	return h.RedirectURL
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(GetNow().Add(rwTimeout))
		return conn, err
	}
}

func SimpleGet(requestUrl string) (string, error) {
	//判断是否是有效URL
	_, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
		//panic(err.Error())
	}
	//开始请求
	resp, err := http.Get(requestUrl)
	if err != nil {
		return "", err
		//panic(err.Error())
	}
	//用完关闭
	defer resp.Body.Close()
	//不是返回OK。就跳过。
	if resp.StatusCode != http.StatusOK {
		//fmt.Println(resp.StatusCode)
		return "", errors.New("resp code is " + fmt.Sprint(resp.StatusCode))
		//return "http response code: " + string(resp.StatusCode)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
		//panic(err.Error())
	}

	return string(respData), nil
}

func SimplePost(requestUrl string, params map[string]string) (string, error) {
	//contentType := "application/json"
	//参数，多个用&隔开

	//表单方式
	sendBody := &http.Request{}
	//sendBody.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	sendErr := sendBody.ParseForm()
	if sendErr != nil {
		return "parse form err", sendErr
	}
	for k, v := range params {
		sendBody.Form.Add(k, v)
	}
	data := strings.NewReader(strings.TrimSpace(sendBody.Form.Encode()))

	tr := &http.Transport{ //解决x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   3 * time.Second,
		Transport: tr, //解决x509: certificate signed by unknown authority
	}
	req, err := http.NewRequest(http.MethodPost, requestUrl, data)
	if err != nil {
		return "http request err: " + requestUrl, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	reqErr := url.Error{Err: err}
	if reqErr.Timeout() {
		//fmt.Println("尝试请求第", r.RequestNum, "次")
		tryResp, tryErr := SimplePost(requestUrl, params)
		return tryResp, tryErr
	}

	//resp, err := http.Post(requestUrl, "application/x-www-form-urlencoded", data)
	if err != nil {
		return "http response err: " + err.Error(), err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respData), nil
}
