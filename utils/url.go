package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func ComputeSignature(timestamp int64, secret string) string {
	b := &[]byte{}

	*b = append(*b, strconv.FormatInt(timestamp, 10)...)
	*b = append(*b, '\n')
	*b = append(*b, secret...)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(*b)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))

}

func ComputeHmacSha256(timestamp int64, secret string) string {
	data := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	// sha := hex.EncodeToString(h.Sum(nil))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))

}

// 获取发送数据的url
func GetSendUrl(DingUrl, Secret string) string {
	timestamp := time.Now().UnixNano() / 1e6
	// var timestamp int64
	// timestamp = 1630131322950
	sign := ComputeSignature(timestamp, Secret)
	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(timestamp, 10))
	params.Add("sign", sign)
	PostUrl := DingUrl + "&" + params.Encode()
	return PostUrl
}
