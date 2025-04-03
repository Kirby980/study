package Registry

import (
	"fmt"
	"os"
	"strings"
)

type ProtoMessage struct {
	Name   string
	Fields []ProtoField
}
type ProtoField struct {
	Name string
	Type string
	Tag  int
}

func GenerateProtoFile(registry *TypeRegistry, suggestedName string) string {
	var builder strings.Builder
	builder.WriteString("syntax = \"proto3\";\n\n")
	builder.WriteString("package example;\n\n")
	names := make([]string, 0, len(registry.types))
	for name := range registry.types {
		names = append(names, name)
	}
	sort.Strings(names)
	// 输出所有消息
	for _, msg := range registry.types {
		builder.WriteString(fmt.Sprintf("message %s {\n", msg.Name))
		for _, field := range msg.Fields {
			builder.WriteString(fmt.Sprintf("  %s %s = %d;\n", field.Type, field.Name, field.Tag))
		}
		builder.WriteString("}\n\n")
	}

	file, err := os.OpenFile(suggestedName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic("打开文件失败")
	}
	file.WriteString(builder.String())
	return builder.String()
}
