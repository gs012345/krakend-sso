package main

import (
	"flag"
	"github.com/devopsfaith/krakend-viper"
	"github.com/devopsfaith/krakend/router"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
	"krakend-sso/krakend-sso/sso"
	"log"

	"context"
	"fmt"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/transport/http/client"
)

func main() {
	port := flag.Int("p", 0, "Port of the service")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "./configuration.json", "Path to the configuration filename")
	flag.Parse()
	//
	parser := viper.New()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}
	//logger, err := gologging.NewLogger(serviceConfig.ExtraConfig)
	//if err != nil {
	//	log.Fatal("ERROR:", err.Error())
	//}
	//
	//logger.Debug("config:", serviceConfig)

	ctx, cancel := context.WithCancel(context.Background())

	backendFactory := sso.SsoNewBackendFactory(nil, client.DefaultHTTPRequestExecutor(client.NewHTTPClient))

	routerFactory := krakendgin.NewFactory(krakendgin.Config{
		Engine:         gin.Default(),
		Logger:         nil,
		HandlerFactory: krakendgin.EndpointHandler,
		ProxyFactory:   proxy.NewDefaultFactory(backendFactory, nil),
		RunServer:      router.RunServer,
		Middlewares:    []gin.HandlerFunc{middle()},
	})

	routerFactory.NewWithContext(ctx).Run(serviceConfig)

	cancel()
}

func middle() func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("middle")
	}
}

/*
import (
	"context"
	"flag"
	"fmt"
	"github.com/devopsfaith/krakend-gologging"
	"log"

	"github.com/devopsfaith/krakend-viper"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/devopsfaith/krakend/transport/http/client"
	"github.com/gin-gonic/gin"

	"github.com/devopsfaith/krakend-martian"
)

func main() {
	port := flag.Int("p", 0, "Port of the service")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "/media/gogo/share/src/krakend-demo/configuration.json", "Path to the configuration filename")
	flag.Parse()
	//
	parser := viper.New()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}

	logger, err := gologging.NewLogger(serviceConfig.ExtraConfig)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	logger.Debug("config:", serviceConfig)

	ctx, cancel := context.WithCancel(context.Background())

	backendFactory := martian.NewBackendFactory(nil, client.DefaultHTTPRequestExecutor(client.NewHTTPClient))

	routerFactory := krakendgin.NewFactory(krakendgin.Config{
		Engine:         gin.Default(),
		Logger:         nil,
		HandlerFactory: krakendgin.EndpointHandler,
		ProxyFactory:   proxy.NewDefaultFactory(backendFactory, nil),
		RunServer: router.RunServer,

	})

	routerFactory.NewWithContext(ctx).Run(serviceConfig)

	cancel()
}



*/

//func SsoNewBackendFactory(logger logging.Logger, re client.HTTPRequestExecutor) proxy.BackendFactory {
//	return NewConfiguredBackendFactory(logger, func(_ *config.Backend) client.HTTPRequestExecutor { return re })
//}
//
//func NewConfiguredBackendFactory(logger logging.Logger, ref func(*config.Backend) client.HTTPRequestExecutor) proxy.BackendFactory {
//	//parse.Register("static.Modifier", staticModifierFromJSON)
//	return func(remote *config.Backend) proxy.Proxy {
//		//logger.Error(result, remote.ExtraConfig)
//		re := ref(remote)
//		return proxy.NewHTTPProxyWithHTTPExecutor(remote, HTTPRequestExecutor(re), remote.Decoder)
//	}
//
//}
//
//func HTTPRequestExecutor(re client.HTTPRequestExecutor) client.HTTPRequestExecutor {
//	return func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
//		if err = modifyRequest(req); err != nil {
//			respErr := modifyResponse(resp, err)
//			if respErr != nil{
//				return
//			}
//			return resp, nil
//		}
//		mctx, ok := req.Context().(*Context)
//		if !ok || !mctx.SkippingRoundTrip() {
//			resp, err = re(ctx, req)
//			if err != nil {
//				return
//			}
//			if resp == nil {
//				err = ErrEmptyResponse
//				return
//			}
//		} else if resp == nil {
//			resp = &http.Response{
//				Request:    req,
//				Header:     http.Header{},
//				StatusCode: http.StatusOK,
//				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
//			}
//		}
//		err = modifyResponse(resp, nil)
//		return
//	}
//}
//
//func modifyRequest(req *http.Request) error {
//	if _, ok := req.Header["X-Sso-Fullticketid"]; !ok{
//		return errors.New("缺少认证header信息")
//	}
//	ticket := req.Header["X-Sso-Fullticketid"][0]
//	//fmt.Printf("获取到的请求头是:%s", ticket)
//	//fmt.Println("获取到的值是:")
//	//fmt.Println(req.Header["X-Sso-Fullticketid"][0])
//	userInfo, err := ssoGetUserModel(ticket)
//	if err != nil{
//		return fmt.Errorf("请求校验sso失败:%s", err.Error())
//	}
//	if userInfo.ErrorCode == 4012{
//		return fmt.Errorf("非法的ticket:%s", userInfo.Message)
//	}
//	if userInfo.Data == nil{
//		return fmt.Errorf("当前用户不存在:%s", userInfo.Message)
//	}
//	req.Header["UserEmail"] = []string{userInfo.Data.LoginEmail}
//	req.Header["AccountGuid"] = []string{userInfo.Data.AccountGuid}
//	return nil
//}
//
//func modifyResponse(resp *http.Response, err error) error {
//	if resp.Header == nil {
//		resp.Header = http.Header{}
//	}
//	if resp.StatusCode == 0 {
//		resp.StatusCode = http.StatusOK
//	}
//	if err != nil{
//		rsp := Response{
//			ErrorId:4001,
//			Reason:fmt.Sprintf("认证失败:%s", err.Error()),
//			Desc:fmt.Sprintf("认证失败:%s", err.Error()),
//		}
//		response, jsonErr := json.Marshal(rsp)
//		if jsonErr != nil{
//			resp.Body  = ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf("数据反序列化失败:%s", jsonErr.Error())))
//		}
//		resp.Body  = ioutil.NopCloser(bytes.NewBuffer(response))
//	}
//	return nil
//}
//
//var (
//	// ErrEmptyValue is the error returned when there is no config under the namespace
//	//ErrEmptyValue = errors.New("getting the extra config for the martian module")
//	//// ErrBadValue is the error returned when the config is not a map
//	//ErrBadValue = errors.New("casting the extra config for the martian module")
//	//// ErrMarshallingValue is the error returned when the config map can not be marshalled again
//	//ErrMarshallingValue = errors.New("marshalling the extra config for the martian module")
//	// ErrEmptyResponse is the error returned when the modifier receives a nil response
//	ErrEmptyResponse = errors.New("getting the http response from the request executor")
//)
//
//type Context struct {
//	context.Context
//	skipRoundTrip bool
//}
//
//// SkipRoundTrip flags the context to skip the round trip
//func (c *Context) SkipRoundTrip() {
//	c.skipRoundTrip = true
//}
//
//// SkippingRoundTrip returns the flag for skipping the round trip
//func (c *Context) SkippingRoundTrip() bool {
//	return c.skipRoundTrip
//}
//
//var _ context.Context = &Context{Context: context.Background()}
//
//
//func ssoGetUserModel(ticket string) (*SsoTicketUserInfoResponse, error) {
//	token_url := "http://10.200.60.36:8800/api/v2/info"
//	reqClient := &http.Client{}
//	v := url.Values{}
//	req, err := http.NewRequest("GET", token_url, strings.NewReader(v.Encode()))
//	req.Header.Set("ticket", ticket)
//	resp, err := reqClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	var token_response SsoTicketUserInfoResponse
//	err = json.Unmarshal(body, &token_response)
//	return &token_response, nil
//
//}
//
//type SsoTicketUserInfoResponse struct {
//	ErrorCode int      `json:"errorCode"`
//	Data      *SsoData `json:"data"`
//	Message   string   `json:"message"`
//}
//
//type SsoData struct {
//	LoginEmail  string         `json:"LoginEmail"`
//	AccountGuid string         `json:"AccountGuid"`
//	DisplayName string         `json:"DisplayName"`
//}
//
//type Response struct {
//	ErrorId int `json:"error_id"`
//	Reason string `json:"reason"`
//	Desc string `json:"desc"`
//}