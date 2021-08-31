package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mlonV/dingtalk/config"
)

func PostMsg(url string, data *config.Msg) (body []byte, err error) {
	header := map[string]string{
		"Content-type":  "application/json;charset=UTF-8",
		"Cache-Control": "no-cache",
		"Connection":    "Keep-Alive",
		"User-Agent":    "ding talk robot send",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("josn marshal err : %s", err)

	}
	fmt.Println(string(jsonData))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest err : %s", err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, _ := client.Do(req)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {

		return nil, fmt.Errorf("read body err : %s", err)
	}
	defer resp.Body.Close()

	return body, nil
}
