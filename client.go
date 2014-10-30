// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// ots2
package goots

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"time"

	. "github.com/GiterLab/goots/log"
	. "github.com/GiterLab/goots/otstype"
	"github.com/GiterLab/goots/urllib"
)

var OTSDebugEnable bool = false // OTS调试默认关闭

const (
	DEFAULT_ENCODING       = "utf8"
	DEFAULT_SOCKET_TIMEOUT = 50
	DEFAULT_MAX_CONNECTION = 50
	DEFAULT_LOGGER_NAME    = "ots2-client"
)

var defaultOTSSetting = OTSClient{
	"http://127.0.0.1:8800",     // EndPoint
	"OTSMultiUser177_accessid",  // AccessId
	"OTSMultiUser177_accesskey", // AccessKey
	"TestInstance177",           // InstanceName
	DEFAULT_SOCKET_TIMEOUT,      // SocketTimeout
	DEFAULT_MAX_CONNECTION,      // MaxConnection
	DEFAULT_LOGGER_NAME,         // LoggerName
	DEFAULT_ENCODING,            // Encoding
	&defaultProtocol,
}
var settingMutex sync.Mutex

// Overwrite default settings
func SetDefaultSetting(setting OTSClient) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultOTSSetting = setting
	if defaultOTSSetting.SocketTimeout == 0 {
		defaultOTSSetting.SocketTimeout = 50
	}
	if defaultOTSSetting.MaxConnection == 0 {
		defaultOTSSetting.MaxConnection = 50
	}
}

// 创建一个新的OTSClient实例
func New(end_point, accessid, accesskey, instance_name string, kwargs ...interface{}) (o *OTSClient, err error) {
	o = &defaultOTSSetting
	o.EndPoint = end_point
	o.AccessId = accessid
	o.AccessKey = accesskey
	o.InstanceName = instance_name

	for i, v := range kwargs {
		switch i {
		case 0: // SocketTimeout --> int32
			if _, ok := v.(int); ok {
				o.SocketTimeout = v.(int)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.SocketTimeout should be int type, not %v", reflect.TypeOf(v)))
			}

		case 1: // MaxConnection --> int32
			if _, ok := v.(int); ok {
				o.MaxConnection = v.(int)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.MaxConnection should be int type, not %v", reflect.TypeOf(v)))
			}

		case 2: // LoggerName --> string
			if _, ok := v.(string); ok {
				o.LoggerName = v.(string)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.LoggerName should be string type, not %v", reflect.TypeOf(v)))
			}

		case 3: // Encoding --> string
			if _, ok := v.(string); ok {
				o.Encoding = v.(string)
			} else {
				return nil, (OTSClientError{}.Set("OTSClient.Encoding should be string type, not %v", reflect.TypeOf(v)))
			}
		}
	}

	// parse end point
	end_point_url, err := url.Parse(end_point)
	if err != nil {
		return nil, (OTSClientError{}.Set("url parse error", err))
	}
	if end_point_url.Scheme != "http" && end_point_url.Scheme != "https" {
		return nil, (OTSClientError{}.Set("protocol of end_point must be 'http' or 'https', e.g. http://ots.aliyuncs.com:80."))
	}

	if end_point_url.Host == "" {
		return nil, (OTSClientError{}.Set("host of end_point should be specified, e.g. http://ots.aliyuncs.com:80."))
	}

	// set default setting for urllib
	url_setting := urllib.HttpSettings{
		false,            // ShowDebug
		"GiterLab",       // UserAgent
		60 * time.Second, // ConnectTimeout
		60 * time.Second, // ReadWriteTimeout
		nil,              // TlsClientConfig
		nil,              // Proxy
		nil,              // Transport
		false,            // EnableCookie
	}
	if o.SocketTimeout != 0 {
		url_setting.ConnectTimeout = time.Duration(o.SocketTimeout) * time.Second
		url_setting.ReadWriteTimeout = time.Duration(o.SocketTimeout) * time.Second
	}
	if OTSDebugEnable {
		url_setting.ShowDebug = true
	}
	urllib.SetDefaultSetting(url_setting)

	o.protocol = newProtocol(nil)
	o.protocol.Set(o.AccessId, o.AccessKey, o.InstanceName, o.Encoding, o.LoggerName)

	return o, nil
}

// OTSClient实例
type OTSClient struct {
	// OTS服务的地址（例如 'http://instance.cn-hangzhou.ots.aliyun.com:80'），必须以'http://'开头
	EndPoint string
	// 访问OTS服务的accessid，通过官方网站申请或通过管理员获取
	AccessId string
	// 访问OTS服务的accesskey，通过官方网站申请或通过管理员获取
	AccessKey string
	// 访问的实例名，通过官方网站控制台创建或通过管理员获取
	InstanceName string

	// 连接池中每个连接的Socket超时，单位为秒，可以为int或float。默认值为50
	SocketTimeout int
	// 连接池的最大连接数。默认为50
	// golang http自带连接池管理，此参数留作以后扩展用
	MaxConnection int

	// 用来在请求中打DEBUG日志，或者在出错时打ERROR日志
	LoggerName string

	// 字符编码格式，此参数未用,兼容python
	Encoding string

	// protocol
	protocol *ots_protocol
}

func (o *OTSClient) String() string {
	r := ""
	r = r + fmt.Sprintln("#### OTSClinet Config ####")
	r = r + fmt.Sprintln("API_VERSION  :", API_VERSION)
	r = r + fmt.Sprintln("DebugEnable  :", OTSDebugEnable)
	r = r + fmt.Sprintln("EndPoint     :", o.EndPoint)
	r = r + fmt.Sprintln("AccessId     :", o.AccessId)
	r = r + fmt.Sprintln("AccessKey    :", o.AccessKey)
	r = r + fmt.Sprintln("InstanceName :", o.InstanceName)
	r = r + fmt.Sprintln("SocketTimeout:", o.SocketTimeout)
	r = r + fmt.Sprintln("MaxConnection:", o.MaxConnection)
	r = r + fmt.Sprintln("LoggerName   :", o.LoggerName)
	// r = r + fmt.Sprintln("Encoding:", o.Encoding)
	r = r + fmt.Sprintln("##########################")

	return r
}

// 在OTSClinet创建后（既调用了New函数），需要重新修改OTSClinet的参数时
// 可以调用此函数进行设置，参数使用字典方式，可以使用的字典如下：
// Debug --> bool
// EndPoint --> string
// AccessId --> string
// AccessKey --> string
// InstanceName --> string
// SocketTimeout --> int
// MaxConnection --> int
// LoggerName --> string
// Encoding --> string
// 注：具体参数意义请查看OTSClinet定义处的注释
func (o *OTSClient) Set(kwargs DictString) *OTSClient {
	if len(kwargs) != 0 {
		for k, v := range kwargs {
			switch k {
			case "Debug":
				if v1, ok := v.(bool); ok {
					setting := urllib.GetDefaultSetting()
					setting.ShowDebug = v1
				} else {
					panic(OTSClientError{}.Set("Debug should be bool, not %v", reflect.TypeOf(v)))
				}
			case "EndPoint":
				if v1, ok := v.(string); ok {
					o.EndPoint = v1
				} else {
					panic(OTSClientError{}.Set("EndPoint should be string, not %v", reflect.TypeOf(v)))
				}
				// parse end point
				end_point_url, err := url.Parse(v.(string))
				if err != nil {
					panic(OTSClientError{}.Set("url parse error", err))
				}
				if end_point_url.Scheme != "http" && end_point_url.Scheme != "https" {
					panic(OTSClientError{}.Set("protocol of end_point must be 'http' or 'https', e.g. http://ots.aliyuncs.com:80."))
				}

				if end_point_url.Host == "" {
					panic(OTSClientError{}.Set("host of end_point should be specified, e.g. http://ots.aliyuncs.com:80."))
				}

			case "AccessId":
				if v1, ok := v.(string); ok {
					o.AccessId = v1
				} else {
					panic(OTSClientError{}.Set("AccessId should be string, not %v", reflect.TypeOf(v)))
				}

			case "AccessKey":
				if v1, ok := v.(string); ok {
					o.AccessKey = v1
				} else {
					panic(OTSClientError{}.Set("AccessKey should be string, not %v", reflect.TypeOf(v)))
				}

			case "InstanceName":
				if v1, ok := v.(string); ok {
					o.InstanceName = v1
				} else {
					panic(OTSClientError{}.Set("InstanceName should be string, not %v", reflect.TypeOf(v)))
				}

			case "SocketTimeout":
				if v1, ok := v.(int); ok {
					o.SocketTimeout = v1
				} else {
					panic(OTSClientError{}.Set("SocketTimeout should be int, not %v", reflect.TypeOf(v)))
				}

			case "MaxConnection":
				if v1, ok := v.(int); ok {
					o.MaxConnection = v1
				} else {
					panic(OTSClientError{}.Set("MaxConnection should be int, not %v", reflect.TypeOf(v)))
				}

			case "LoggerName":
				if v1, ok := v.(string); ok {
					o.LoggerName = v1
				} else {
					panic(OTSClientError{}.Set("LoggerName should be string, not %v", reflect.TypeOf(v)))
				}

			case "Encoding":
				if v1, ok := v.(string); ok {
					o.Encoding = v1
				} else {
					panic(OTSClientError{}.Set("Encoding should be string, not %v", reflect.TypeOf(v)))
				}

			default:
				panic(OTSClientError{}.Set("Unknown param %s", k))
			}
		}
	}

	return o
}

func _request_helper(api_name string, args ...interface{}) []interface{} {
	return nil
}
