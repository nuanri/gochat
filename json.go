package gochat

import (
	"encoding/json"
	"fmt"
)

type HiMsg struct {
	Cmd  string `json:"cmd"`
	Body []byte `json:"body"`
	//	UserLogo string `json:"user_logo"`
}

type SigBody struct {
	Name string `json:"name"`
}

type SendBody struct {
	To   string `json:"to"`
	From string `json:"from"`
	Msg  string `json:"msg"`
}

//解析 Json
func ParseMsg(data []byte, m interface{}) error {

	err := json.Unmarshal(data, m)
	if err != nil {
		fmt.Println("json 解析出错", err)
		return err
	}

	return nil
}

//生成 Json
func GetJson(m interface{}) []byte {

	b, _ := json.Marshal(m)
	/*	if err != nil {
		fmt.Println("Marshal 出错", err)
		return "", err
	}*/

	return b
}
