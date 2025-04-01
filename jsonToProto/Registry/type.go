package Registry

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TypeRegistry struct {
	Message map[string]string
	types   map[string]ProtoMessage
	counter int
}

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		Message: make(map[string]string),
		types:   make(map[string]ProtoMessage),
	}
}

func ProcessValue(val interface{}, registry *TypeRegistry, suggestedName string) string {
	switch v := val.(type) {
	case []interface{}:
		if len(v) == 0 {
			return "bytes"
		}
		elemType := ProcessValue(v[0], registry, suggestedName)
		return "repeated " + elemType
	case map[string]interface{}:
		return ProcessObject(v, registry, suggestedName)
	case json.Number:
		if _, err := v.Int64(); err == nil {
			return "int32"
		}
		return "double"
	case string:
		return "string"
	case bool:
		return "bool"
	case nil:
		return "google.protobuf.NullValue"
	default:
		return "bytes"
	}
}
func ProcessObject(obj map[string]interface{}, registry *TypeRegistry, suggestedName string) string {
	fields := make([]struct{ Name, Type string }, 0)
	//flagMap := false
	for k, v := range obj {
		if _, err := strconv.Atoi(k); err == nil {
			//		flagMap = true
			return "map<string,string>"
		}
		fieldType := ProcessValue(v, registry, k)
		fields = append(fields, struct{ Name, Type string }{k, fieldType})
	}

	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
	var sig strings.Builder
	for _, f := range fields {
		sig.WriteString(f.Name + ":" + f.Type + ";")
	}
	if name, ok := registry.Message[sig.String()]; ok {
		return name
	}
	messageName := GenerateMessageName(suggestedName, registry)
	protoField := make([]ProtoField, 0, 20)
	tag := 1
	for _, field := range fields {
		protoField = append(protoField, ProtoField{
			Name: strings.ToLower(field.Name),
			Type: field.Type,
			Tag:  tag,
		})
		tag++
	}
	registry.Message[sig.String()] = messageName
	registry.types[messageName] = ProtoMessage{
		Name:   messageName,
		Fields: protoField,
	}
	return messageName
}
func GenerateMessageName(base string, r *TypeRegistry) string {
	name := strings.Title(base)
	if _, exists := r.types[name]; !exists {
		return name
	}
	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s%d", name, i)
		if _, exists := r.types[newName]; !exists {
			return newName
		}
	}
}
