package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"study/Registry"
)

func main() {

	Value, err := os.ReadFile("test.json")
	if err != nil {
		panic("读取文件失败")
	}
	value := Registry.RemoveJSONComment(string(Value))
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
	Registry.ProcessValue(data, registry, "json_to_proto")
	Registry.GenerateProtoFile(registry, "json_to.proto")
}
