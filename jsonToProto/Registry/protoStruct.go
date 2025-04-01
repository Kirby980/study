package Registry

import (
	"fmt"
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

func GenerateProtoFile(registry *TypeRegistry) string {
	var builder strings.Builder
	// 输出所有消息
	for _, msg := range registry.types {
		builder.WriteString(fmt.Sprintf("message %s {\n", msg.Name))
		for _, field := range msg.Fields {
			builder.WriteString(fmt.Sprintf("  %s %s = %d;\n", field.Type, field.Name, field.Tag))
		}
		builder.WriteString("}\n\n")
	}
	return builder.String()
}
