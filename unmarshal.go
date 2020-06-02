package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

func UnmarshalRecord(journal *sdjournal.Journal, to *Record) error {
	entry, err := journal.GetEntry()
	if err != nil {
		return err
	}
	to.TimeUsec = int64(entry.RealtimeTimestamp)
	return unmarshalRecord(entry, reflect.ValueOf(to).Elem())
}

func unmarshalRecord(entry *sdjournal.JournalEntry, toVal reflect.Value) error {
	toType := toVal.Type()

	numField := toVal.NumField()

	// This intentionally supports only the few types we actually
	// use on the Record struct. It's not intended to be generic.

	for i := 0; i < numField; i++ {
		fieldVal := toVal.Field(i)
		fieldDef := toType.Field(i)
		fieldType := fieldDef.Type
		fieldTag := fieldDef.Tag
		fieldTypeKind := fieldType.Kind()

		if fieldTypeKind == reflect.Struct {
			// Recursively unmarshal from the same journal
			unmarshalRecord(entry, fieldVal)
		}

		jdKey := fieldTag.Get("journald")
		if jdKey == "" {
			continue
		}

		value, ok := entry.Fields[jdKey]
		if !ok {
			fieldVal.Set(reflect.Zero(fieldType))
			continue
		}

		if fieldType.Name() == "RawMessage" {
			if !strings.HasPrefix(value, `{"`) {
				jenc, _ := json.Marshal(value)
				value = string(jenc)
			}
			fieldVal.SetBytes(json.RawMessage(value))
			continue
		}

		switch fieldTypeKind {
		case reflect.Int:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				// Should never happen, but not much we can do here.
				fieldVal.Set(reflect.Zero(fieldType))
				continue
			}
			fieldVal.SetInt(int64(intVal))
		case reflect.String:
			fieldVal.SetString(value)
		default:
			// Should never happen
			panic(fmt.Errorf("cannot unmarshal to %s", fieldType))
		}
	}

	return nil
}
