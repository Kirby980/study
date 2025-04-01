package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"study/Registry"
)

const Value = `{
  "ret":0,   // 返回状态码，0表示正式，其他表示错误
  "msg":"success", // 额外的信息，当有错误时为错误信息
  "study":{
    "day1": 1,
    "day2": 2,
    "day3": "test"
  },
    "test":{
    "1" : "1",
    "2" : "2",
    },
    "slice":[
    {
    "test1":1,
    "test2":2
    }
    ]
}`

func main() {
	value := Registry.RemoveJSONComment(Value)
	//fmt.Printf("清理后的 JSON:\n%q\n", value)

	fmt.Println(value)
	var data any
	decoder := json.NewDecoder(bytes.NewReader([]byte(value)))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			log.Printf("语法错误位置：%d，上下文：%s", syntaxErr.Offset, value[syntaxErr.Offset-10:syntaxErr.Offset+10])
		}
		log.Fatal(err)
	}
	fmt.Printf("%T\n", data)
	registry := Registry.NewTypeRegistry()
	Registry.ProcessValue(data, registry, "json_to_proto ")
	file := Registry.GenerateProtoFile(registry)
	fmt.Println(file, "----------file")
}
