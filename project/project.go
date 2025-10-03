package project

import (
	"encoding/json"
	"fmt"

	"github.com/adm87/finch-core/hashset"
)

// ======================================================
// Tiled Project Format
// ======================================================

type TiledPropertyType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type TiledEnumPropertyType struct {
	TiledPropertyType

	StorageType   string   `json:"storageType"`
	Values        []string `json:"values"`
	ValuesAsFlags bool     `json:"valuesAsFlags"`
}

type TiledClassPropertyType struct {
	TiledPropertyType

	Color    string             `json:"color"`
	DrawFill bool               `json:"drawFill"`
	Members  []TiledClassMember `json:"members"`
	UseAs    []TiledClassUseAs  `json:"useAs"`
}

type TiledClassMember struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	PropertyType string `json:"propertyType,omitempty"`
	Value        any    `json:"value"`
}

type TiledClassUseAs string

const (
	TiledClassUseAsProperty  TiledClassUseAs = "property"
	TiledClassUseAsMap       TiledClassUseAs = "map"
	TiledClassUseAsLayer     TiledClassUseAs = "layer"
	TiledClassUseAsObject    TiledClassUseAs = "object"
	TiledClassUseAsTile      TiledClassUseAs = "tile"
	TiledClassUseAsTileset   TiledClassUseAs = "tileset"
	TiledClassUseAsWangColor TiledClassUseAs = "wangcolor"
	TiledClassUseAsWangSet   TiledClassUseAs = "wangset"
	TiledClassUseAsProject   TiledClassUseAs = "project"
)

func (t TiledClassUseAs) IsValid() bool {
	return hashset.From(
		TiledClassUseAsProperty,
		TiledClassUseAsMap,
		TiledClassUseAsLayer,
		TiledClassUseAsTile,
		TiledClassUseAsTileset,
		TiledClassUseAsWangColor,
		TiledClassUseAsWangSet,
		TiledClassUseAsProject,
	).Contains(t)
}

type TiledProject struct {
	AutomappingRulesFile string                   `json:"automappingRulesFile"`
	Commands             []any                    `json:"commands"`
	CompatibilityVersion int                      `json:"compatibilityVersion"`
	ExtensionsPath       string                   `json:"extensionsPath"`
	Folders              []string                 `json:"folders"`
	Properties           []any                    `json:"properties"`
	EnumPropertyTypes    []TiledEnumPropertyType  `json:"-"` // Ignored
	ClassPropertyTypes   []TiledClassPropertyType `json:"-"` // Ignored
}

func (p *TiledProject) MarshalJSON() ([]byte, error) {
	type Alias TiledProject
	wrapper := &struct {
		*Alias
		PropertyTypes []any `json:"propertyTypes"`
	}{
		Alias: (*Alias)(p),
	}

	for _, enumType := range p.EnumPropertyTypes {
		wrapper.PropertyTypes = append(wrapper.PropertyTypes, enumType)
	}
	for _, classType := range p.ClassPropertyTypes {
		wrapper.PropertyTypes = append(wrapper.PropertyTypes, classType)
	}

	return json.Marshal(wrapper)
}

func (p *TiledProject) UnmarshalJSON(data []byte) error {
	type Alias TiledProject
	wrapper := &struct {
		*Alias
		PropertyTypes []json.RawMessage `json:"propertyTypes"`
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	for _, rawType := range wrapper.PropertyTypes {
		var baseType struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawType, &baseType); err != nil {
			return err
		}

		switch baseType.Type {
		case "enum":
			var enumType TiledEnumPropertyType
			if err := json.Unmarshal(rawType, &enumType); err != nil {
				return err
			}
			p.EnumPropertyTypes = append(p.EnumPropertyTypes, enumType)
		case "class":
			var classType TiledClassPropertyType
			if err := json.Unmarshal(rawType, &classType); err != nil {
				return err
			}
			p.ClassPropertyTypes = append(p.ClassPropertyTypes, classType)
		default:
			return fmt.Errorf("unknown property type: %s", baseType.Type)
		}
	}

	return nil
}
