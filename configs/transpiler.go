package configs

import (
	"fmt"
	"reflect"
	"strings"
)

type UnitCodec[T any] interface {
	Encode(T) (*SystemdUnit, error)
	Decode(*SystemdUnit) (T, error)
}

type Entry struct {
	Key   string
	Value string
}

type SystemdUnit struct {
	Filename string
	Sections map[string][]Entry
}

func EncodeSystemdSection(v any) ([]Entry, error) {
	var entries []Entry

	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("EncodeSystemdSection: expected struct, got %s", rv.Kind())
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		key := field.Tag.Get("systemd")
		if key == "" || key == "-" {
			continue
		}

		if value.IsZero() {
			continue
		}

		switch value.Kind() {

		case reflect.String:
			entries = append(entries, Entry{
				Key:   key,
				Value: value.String(),
			})

		case reflect.Bool:
			entries = append(entries, Entry{
				Key:   key,
				Value: "yes",
			})

		case reflect.Slice:
			for j := 0; j < value.Len(); j++ {
				entries = append(entries, Entry{
					Key:   key,
					Value: fmt.Sprint(value.Index(j).Interface()),
				})
			}

		case reflect.Map:
			iter := value.MapRange()
			for iter.Next() {
				entries = append(entries, Entry{
					Key:   key,
					Value: fmt.Sprintf("%v=%v", iter.Key(), iter.Value()),
				})
			}

		default:
			return nil, fmt.Errorf(
				"EncodeSystemdSection: unsupported kind %s for field %s",
				value.Kind(),
				field.Name,
			)
		}
	}

	return entries, nil
}

func (su *SystemdUnit) ToString() string {
	var b strings.Builder

	for section, entries := range su.Sections {
		writeSection(&b, section, entries)
	}

	return strings.TrimSpace(b.String()) + "\n"
}

func writeSection(b *strings.Builder, section string, entries []Entry) {
	b.WriteString("[")
	b.WriteString(section)
	b.WriteString("]\n")

	for _, e := range entries {
		b.WriteString(e.Key)
		b.WriteString("=")
		b.WriteString(e.Value)
		b.WriteString("\n")
	}

	b.WriteString("\n")
}

type SystemdUnitBuilder struct {
	unit *SystemdUnit
}

func NewSystemdUnitBuilder() *SystemdUnitBuilder {
	return &SystemdUnitBuilder{
		unit: &SystemdUnit{
			Sections: make(map[string][]Entry),
		},
	}
}

func (b *SystemdUnitBuilder) AddEntry(section, key, value string) *SystemdUnitBuilder {
	b.unit.Sections[section] = append(
		b.unit.Sections[section],
		Entry{Key: key, Value: value},
	)
	return b
}

func (b *SystemdUnitBuilder) AddEntries(section string, entries ...Entry) *SystemdUnitBuilder {
	b.unit.Sections[section] = append(
		b.unit.Sections[section],
		entries...,
	)
	return b
}

func (b *SystemdUnitBuilder) Build(filename string) *SystemdUnit {
	b.unit.Filename = filename

	return b.unit
}
