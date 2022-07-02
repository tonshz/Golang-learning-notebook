package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	APP_KEY    = "admin"
	APP_SECRET = "go-learning"
)

type AccessToken struct {
	Token string `json:"token"`
}

type API struct {
	URL string
}
type AuthParams struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func NewAPI(url string) *API {
	return &API{URL: url}
}

// 获取所有 API 请求都需要带上的 token
func (a *API) getAccessToken(ctx context.Context) (string, error) {
	body, err := a.httpPost(ctx, "/auth", APP_KEY, APP_SECRET)
	if err != nil {
		return "", err
	}

	var accessToken AccessToken
	_ = json.Unmarshal(body, &accessToken)
	return accessToken.Token, nil
}

// 统一的 HTTP GET 请求方法
func (a *API) httpGet(ctx context.Context, token string, path string) ([]byte, error) {
	// 自定义 HTTPClient
	req, _ := http.NewRequest("GET", a.URL+path, nil)
	req.Header.Set("token", token)

	span, newCtx := opentracing.StartSpanFromContext(
		ctx, "HTTP GET: "+a.URL,
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
	)
	span.SetTag("url", a.URL+path)
	_ = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := ctxhttp.Do(newCtx, http.DefaultClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer span.Finish()

	// 读取消息主体，在实际封装中可以将其抽离
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// 统一的 HTTP POST 请求方法
func (a *API) httpPost(ctx context.Context, path string, appKey string, appSecret string) ([]byte, error) {
	//resp, err := ctxhttp.Post(ctx, http.DefaultClient, a.URL+path, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("app_key=%s&app_secret=%s", appKey, appSecret)))
	// 使用 json 传输数据
	reqParam, _ := json.Marshal(&AuthParams{appKey, appSecret})
	reqBody := strings.NewReader(string(reqParam))
	resp, err := ctxhttp.Post(ctx, http.DefaultClient, a.URL+path, "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// 具体的获取标签列表的方法实现
func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	// 获取AccessToken
	token, err := a.getAccessToken(ctx)
	log.Printf("token: %s", token)
	if err != nil {
		return nil, err
	}

	body, err := a.httpGet(ctx, token, fmt.Sprintf("%s?name=%s", "/api/v1/tags", name))
	if err != nil {
		return nil, err
	}

	return body, nil
}
