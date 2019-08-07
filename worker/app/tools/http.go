package tools

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpReq struct {
	Ua      string
	Referer string
	Cookie  string
	request *http.Request
}

func NewHttpReq() *HttpReq {
	return &HttpReq{}
}

// http get 请求
func (h *HttpReq) HttpGet(url string) (response *http.Response, bytes []byte, err error) {
	client := &http.Client{
		Timeout: 15 * time.Second, //设置请求超时
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				dialer := net.Dialer{
					Timeout: 5 * time.Second,
				}
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}
	if h.request, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	h.DefaultHeader()
	if response, err = client.Do(h.request); err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		return nil, nil, errors.New("")
	}
	defer response.Body.Close()
	bytes, err = ioutil.ReadAll(response.Body)
	return response, bytes, err
}

func (h *HttpReq) SetUa(ua string) *HttpReq {
	h.Ua = ua
	return h
}

func (h *HttpReq) SetReferer(referer string) *HttpReq {
	h.Referer = referer
	return h
}

func (h *HttpReq) SetCookie(cookieStr string) *HttpReq {
	h.Cookie = cookieStr
	return h
}

//设置默认header
func (h *HttpReq) DefaultHeader() {
	var ua, referer, cookie string
	ua = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"
	referer = "https://www.gushiwen.org/"
	cookie = "Hm_lvt_04660099568f561a75456483228a9516=1564134338,1564385639,1565090591,1565165060; Hm_lpvt_04660099568f561a75456483228a9516=1565165132"
	if h.Ua == "" {
		h.Ua = ua
	}
	if h.Referer == "" {
		h.Referer = referer
	}
	if h.Cookie == "" {
		h.Cookie = cookie
	}
	h.request.Header.Set("User-Agent", h.Ua)
	h.request.Header.Set("referer", h.Referer)
	h.request.Header.Set("cookie", h.Cookie)
}
