package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	// "github.com/mlonV/dingtalk/config"
)

// 发送到玎钉机器人
func SendMsg(url string, data []byte) (body []byte, err error) {
	header := map[string]string{
		"Content-type":  "application/json;charset=UTF-8",
		"Cache-Control": "no-cache",
		"Connection":    "Keep-Alive",
		"User-Agent":    "ding talk robot send",
	}

	if err != nil {
		return nil, fmt.Errorf("josn marshal err : %s", err)

	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest err : %s", err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, _ := client.Do(req)

	body, err = io.ReadAll(resp.Body)
	if err != nil {

		return nil, fmt.Errorf("read body err : %s", err)
	}
	defer resp.Body.Close()

	return body, nil
}
