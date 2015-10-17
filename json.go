package gochat

import (
	"encoding/json"
	"fmt"
)

//解析 Json
type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
	//	UserLogo string `json:"user_logo"`
}

func ParseMsg(line string) (*Message, error) {

	var message Message

	data := []byte(line)
	//	fmt.Println("传到解析json里的字符串", data)
	err := json.Unmarshal(data, &message)
	if err != nil {
		fmt.Println("json 解析出错", err)
		return nil, err
	}

	return &message, nil
}

//生成 Json
func (m *Message) JSON() string {

	b, _ := json.Marshal(m)
	/*	if err != nil {
		fmt.Println("Marshal 出错", err)
		return "", err
	}*/

	return string(b)
}
