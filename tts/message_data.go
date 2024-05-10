package tts

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star"
)

type messageData struct {
	Name     string
	Doc      string
	Fields   []*messageField
	Optional bool
}

type messageField struct {
	Name       string
	JSONName   string
	ZeroValue  string
	Type       string
	IsRepeated bool
	IsEnum     bool
	IsMessage  bool
}

func createField(field pgs.Field) *messageField {
	var f messageField
	f.Name = field.Name().LowerCamelCase().String()
	f.JSONName = field.Name().String()
	f.Type = protoTypeToTSType(field.Type())
	f.IsRepeated = field.Type().IsRepeated()
	f.IsEnum = field.Type().ProtoType() == pgs.EnumT
	f.IsMessage = field.Type().ProtoType() == pgs.MessageT
	f.ZeroValue = f.populateZeroValue()
	return &f
}

func protoTypeToTSType(typ pgs.FieldType) string {
	switch typ.ProtoType() {
	case pgs.EnumT:
		if typ.IsRepeated() {
			return typ.Element().Enum().Name().String()
		}
		return typ.Enum().Name().String()
	case pgs.MessageT:
		if typ.IsRepeated() {
			return typ.Element().Embed().Name().String()
		}
		return typ.Embed().Name().String()
	case pgs.BoolT:
		return "boolean"
	case pgs.Int32T, pgs.Int64T,
		pgs.SInt32, pgs.SInt64,
		pgs.UInt32T, pgs.UInt64T,
		pgs.DoubleT, pgs.FloatT:
		return "number"
	case pgs.StringT:
		return "string"
	case pgs.BytesT:
		return "string"
	}
	panic("unknown type: " + typ.Element().ProtoType().String())
}

func (mf *messageField) populateZeroValue() string {
	switch mf.Type {
	case "boolean":
		return "false"
	case "number":
		return "0"
	case "string":
		return "''"
	}
	if mf.IsRepeated {
		return "[]"
	}
	if mf.IsEnum {
		return "''"
	}
	return "{}"
}

func (mf messageField) ResolveType() string {
	t := mf.Type

	if mf.IsRepeated {
		switch t {
		case "string", "number", "boolean":
			return fmt.Sprintf("(props['%s']! || []).map((v) => { return %s(v)})", mf.JSONName, title.String(t))
		}
		if mf.IsEnum {
			return fmt.Sprintf("(props['%s']! || []).map((v) => { return (%s)[v] })", mf.JSONName, t)
		}
		return fmt.Sprintf("(props['%s']! || []).map((v) => { return %s.fromJSON(v) })", mf.JSONName, t)
	}

	switch t {
	case "string", "number", "boolean":
		return fmt.Sprintf("%s(props['%s'] || %s)!", title.String(t), mf.JSONName, mf.ZeroValue)
	}

	if mf.IsEnum {
		return fmt.Sprintf("(%s)[props['%s']! || '']!", t, mf.JSONName)
	}

	return fmt.Sprintf("%s.fromJSON(props['%s']!)", t, mf.JSONName)
}

func (mf messageField) PrintType() string {
	resp := mf.Type
	if mf.IsRepeated {
		resp += "[]"
	}
	return resp
}

func (mf messageField) PrintTypeProperties() string {
	resp := mf.Type
	if mf.isClass() {
		resp += "Properties"
	}
	if mf.IsRepeated {
		resp += "[]"
	}
	return resp
}

func (mf messageField) SetConstructorProp(optional bool) string {
	if mf.isBasic() || mf.IsEnum {
		return fmt.Sprintf("props.%s!", mf.Name)
	}
	if mf.IsRepeated {
		return fmt.Sprintf("(props.%s! || []).map((v) => { return new %s(v!) })", mf.Name, mf.Type)
	}
	if optional {
		return fmt.Sprintf("props.%s && new %s(props.%s!)", mf.Name, mf.Type, mf.Name)
	}
	return fmt.Sprintf("new %s(props.%s!)", mf.Type, mf.Name)
}

func (mf messageField) SetToObjectProp(optional bool) string {
	if mf.isBasic() || mf.IsEnum {
		return fmt.Sprintf("this.%s", mf.Name)
	}
	if mf.IsRepeated {
		return fmt.Sprintf("(this.%s || []).map((v) => { return v.toObject() })", mf.Name)
	}
	if optional {
		return fmt.Sprintf("this.%s?.toObject()", mf.Name)
	}
	return fmt.Sprintf("this.%s.toObject()", mf.Name)
}

func (mf messageField) isBasic() bool {
	switch mf.Type {
	case "string", "number", "boolean":
		return true
	}
	return false
}

func (mf messageField) isClass() bool {
	return !mf.isBasic() && !mf.IsEnum
}
