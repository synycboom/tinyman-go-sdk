package utils

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/synycboom/tinyman-go-sdk/types"
)

// Program returns a program byte array to be used in LogicSig
func Program(definition types.Logic, variables map[string]uint64) ([]byte, error) {
	template, err := base64.StdEncoding.DecodeString(definition.Bytecode)
	if err != nil {
		return nil, err
	}

	copiedVariables := make([]types.Variable, len(definition.Variables))
	copy(copiedVariables, definition.Variables)

	sort.SliceStable(copiedVariables, func(i, j int) bool {
		return copiedVariables[i].Index < copiedVariables[j].Index
	})

	offset := 0
	for _, v := range copiedVariables {
		ss := strings.Split(v.Name, "TMPL_")
		name := strings.ToLower(ss[len(ss)-1])
		value := variables[name]
		start := v.Index - offset
		end := start + v.Length
		valueEncoded, err := EncodeValue(value, v.Type)
		if err != nil {
			return nil, err
		}

		diff := v.Length - len(valueEncoded)
		offset += diff
		template = updateByteRange(template, valueEncoded, start, end)
	}

	return template, nil
}

func updateByteRange(src []byte, update []byte, start, end int) (out []byte) {
	if start >= 0 {
		if start < len(src) {
			out = append(out, src[0:start]...)
		} else {
			out = append(out, src...)
		}

		out = append(out, update...)
		if end <= start && start < len(src) {
			out = append(out, src[start:]...)
		} else if end > start && end < len(src) {
			out = append(out, src[end:]...)
		}
	}

	return
}

// EncodeVarInt encodes value to be used in program
func EncodeValue(value any, valueType string) ([]byte, error) {
	if valueType != "int" {
		return nil, fmt.Errorf("unsupported value type %s", valueType)
	}

	if v, ok := value.(uint64); ok {
		return EncodeVarInt(v), nil
	}

	return nil, fmt.Errorf("unsuported value %v", value)
}

// EncodeVarInt encodes 64-bit unsigned integer value
func EncodeVarInt(number uint64) []byte {
	var buf []byte
	for {
		toWrite := number & 0x7f
		number >>= 7
		if number != 0 {
			buf = append(buf, byte(toWrite|0x80))
		} else {
			buf = append(buf, byte(toWrite))
			break
		}
	}

	return buf
}
