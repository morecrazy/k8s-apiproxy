package util

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
	"io"
)

type User struct {
	Name string
	Tel string
	Age string
}

type Info struct {
	Id string
	User []*User
	Collage string
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) (bool) {
	var exist = true;
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false;
	}
	return exist;
}

func marshalJSON() []byte {
	user1 := &User{"刘凯宁", "15200806084", "21"}
	user2 := &User{"lkn", "13244398812", "12"}
	info := Info{"0010", []*User{user1, user2}, "中南大学"}
	info_json, _ := json.Marshal(info)
	//fmt.Printf("JSON格式：%s", info_json)
	//file, err1 := os.OpenFile("./info.json", os.O_CREATE|os.O_WRONLY, 0)
	//checkError(err1)
	//defer file.Close()
	err2 := ioutil.WriteFile("./info.json", info_json, 0666)
	checkError(err2)
	return info_json
}


func unmarshalJSON(data_JSON []byte) {
	data := make(map[string]interface{})
	err := json.Unmarshal(data_JSON, &data)
	checkError(err)
	//data_map := data.(map[string]interface{})
	for k, v := range data {
		switch valueType := v.(type) {
		case string:
			fmt.Println(k, "is string", valueType)
		case int:
			fmt.Println(k, "is int", valueType)
		case []interface{}:
			fmt.Println(k, "is array:")
			for k, u := range valueType {
				fmt.Println("k=:", k, " v=:", u)
			}
		default:
			fmt.Println(k, "is a type that I don't know")
		}
	}
}

func copyFile(sourceFile, resultFile string) {
	//source文件
	source, err := os.Open(sourceFile)
	checkError(err)
	defer source.Close()
	//复制的文件
	result, err := os.OpenFile(resultFile, os.O_WRONLY|os.O_CREATE, 0644)
	checkError(err)
	defer result.Close()
	//复制
	io.Copy(result, source)
	fmt.Println("复制成功!")

}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

