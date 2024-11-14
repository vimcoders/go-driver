package dun

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Antispam struct {
	TaskId     string `json:"taskId"`
	DataId     string `json:"dataId"`
	Suggestion int    `json:"suggestion"`
}

type Result struct {
	Antispam `json:"antispam"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Result  `json:"result"`
}

type Client struct {
	apiUrl     string
	version    string
	secretId   string
	secretKey  string
	businessId string
}

func (x *Client) CheckText(params url.Values) error {
	params["secretId"] = []string{secretId}
	params["businessId"] = []string{businessId}
	params["version"] = []string{version}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["nonce"] = []string{strconv.FormatInt(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000), 10)}
	// params["signatureMethod"] = []string{"SM3"} // 签名方法支持国密SM3，默认MD5
	params["signature"] = []string{x.genSignature(params)}
	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(contents))
	var response Response
	if err := json.Unmarshal(contents, &response); err != nil {
		return err
	}
	fmt.Println(response)
	if response.Suggestion != 0 {
		return errors.New("response.Suggestion")
	}
	return nil
}

func (x *Client) genSignature(params url.Values) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paramStr += key + params[key][0]
	}
	paramStr += secretKey
	if params["signatureMethod"] != nil && params["signatureMethod"][0] == "SM3" {
		// sm3Reader := sm3.New()
		// sm3Reader.Write([]byte(paramStr))
		// return hex.EncodeToString(sm3Reader.Sum(nil))
	}
	md5Reader := md5.New()
	md5Reader.Write([]byte(paramStr))
	return hex.EncodeToString(md5Reader.Sum(nil))
}

func NewClient() *Client {
	return &Client{
		apiUrl:     apiUrl,
		version:    version,
		secretId:   secretId,
		secretKey:  secretKey,
		businessId: businessId,
	}
}

const (
	apiUrl     = "http://as.dun.163.com/v5/text/check"
	version    = "v5.2"
	secretId   = "secretId"   //产品密钥ID，产品标识
	secretKey  = "secretKey"  //产品私有密钥，服务端生成签名信息使用，请严格保管，避免泄露
	businessId = "businessId" //业务ID，易盾根据产品业务特点分配
)
