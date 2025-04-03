package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Kirby980/study/jsonToProto/Registry"
)

func Test_jsonToProto(t *testing.T) {

	Value, err := os.ReadFile("test.json")
	if err != nil {
		panic("读取文件失败")
	}
	value := Registry.RemoveJSONComment(string(Value))
	//fmt.Printf("清理后的 JSON:\n%q\n", value)
	orderedValue, err := Registry.ParseJson(value)
	if err != nil {
		fmt.Println(err)
		return
	}
	registry := Registry.NewTypeRegistry()
	processValue := Registry.ProcessValueV2(orderedValue, registry, "json_to_proto")
	fmt.Println(processValue)
	file := Registry.GenerateProtoFile(registry, "json_to.proto")
	fmt.Println(file)

}
func UnOrderV1(value string) {
	fmt.Println(value)
	var data any
	decoder := json.NewDecoder(bytes.NewReader([]byte(value)))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			log.Printf("语法错误位置：%d,上下文：%s", syntaxErr.Offset, value[syntaxErr.Offset-10:syntaxErr.Offset+10])
		}
		log.Fatal(err)
	}
	fmt.Printf("%T\n", data)
	registry := Registry.NewTypeRegistry()
	Registry.ProcessValueV1(data, registry, "json_to_proto")
	Registry.GenerateProtoFile(registry, "json_to.proto")
}
