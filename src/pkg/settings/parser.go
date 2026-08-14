package settings

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type OptionType int

const (
	OptionType_Boolean OptionType = iota
	OptionType_Integer
	OptionType_String
	OptionType_Float
	OptionType_Level
)

func getTypeFromReflect(t reflect.Type) OptionType {
	switch t.Kind() {
	case reflect.Bool:
		return OptionType_Boolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return OptionType_Integer
	case reflect.Float32, reflect.Float64:
		return OptionType_Float
	case reflect.String:
		return OptionType_String
	default:
		return OptionType_String
	}
}

// Possible tag options:
//
// arg:"-s,--long", arg:"-s", arg:"--long", arg:"--no-long"
// env:"ENVIRONMENT"
// help:"this is some helpful text"
// default:"<value of n-type>"
// options:"<see below>"
//
// Possible "options" tag settings:
// * required
// * invertable
// * negated
// * level
// * maxlevel=n
func (p *Parser) extractOptions(t reflect.Type, path []string, currentSection string, mod *Module) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue
		}

		// Skip XMLName field (used by encoding/xml)
		if field.Name == "XMLName" {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			sectionName := field.Name
			if custom := field.Tag.Get("section"); custom != "" {
				sectionName = custom
			}

			p.extractOptions(field.Type, append(path, field.Name), sectionName, mod)
			continue
		}

		long, short := splitArgs(field.Tag.Get("arg"))
		env := field.Tag.Get("env")
		help := field.Tag.Get("help")
		defVal := field.Tag.Get("default")

		cOptField := strings.Split(field.Tag.Get("options"), ",")
		cOptions := map[string]bool{}
		maxLevel := 0

		for _, c := range cOptField {
			if strings.HasPrefix(c, "maxlevel") {
				b, _ := strings.CutPrefix(c, "maxlevel=")
				maxLevel, _ = strconv.Atoi(b)
			} else {
				cOptions[c] = true
			}
		}

		required := cOptions["required"]
		invertible := cOptions["invertible"]
		negated := cOptions["negated"]
		isLevel := cOptions["level"]

		if negated && strings.HasPrefix(long, "no-") {
			long = strings.TrimPrefix(long, "no-")
		}

		optType := getTypeFromReflect(field.Type)
		if isLevel && optType == OptionType_Integer {
			optType = OptionType_Level
		}

		opt := &Option{
			Name:       field.Name,
			Path:       append([]string{}, path...),
			Section:    currentSection,
			Short:      short,
			Long:       long,
			Env:        env,
			HelpText:   help,
			DefaultVal: defVal,
			Type:       optType,
			Required:   required,
			Invertible: invertible,
			Negated:    negated,
			MaxLevel:   maxLevel,
			ModuleName: mod.Name,
			target:     mod.Target,
		}

		if isLevel && len(short) > 0 {
			opt.LevelChar = short[0]
		}

		mod.Options = append(mod.Options, opt)
		p.allOpts = append(p.allOpts, opt)
	}
}

func setFieldByPath(s any, opt *Option, value any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	for _, part := range opt.Path {
		if v.Kind() == reflect.Struct {
			v = v.FieldByName(part)
			if !v.IsValid() {
				return fmt.Errorf("settings: invalid path part '%s'", part)
			}
		} else {
			return fmt.Errorf("settings: expected struct at path part '%s', got %s", part, v.Kind())
		}
	}

	field := v.FieldByName(opt.Name)
	if !field.IsValid() {
		return fmt.Errorf("settings: field '%s' not found", opt.Name)
	}
	if !field.CanSet() {
		return fmt.Errorf("settings: field '%s' cannot be set", opt.Name)
	}

	val := reflect.ValueOf(value)
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
		return nil
	}
	return fmt.Errorf("settings: cannot convert %T to %s", value, field.Type())
}

func (p *Parser) parseValueString(opt *Option, value string) error {
	var parsed any
	var err error

	switch opt.Type {
	case OptionType_Boolean:
		parsed, err = strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("settings: invalid bool for '%s': %s", opt.Name, value)
		}
	case OptionType_Integer:
		parsed, err = strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("settings: invalid int for '%s': %s", opt.Name, value)
		}
	case OptionType_String:
		parsed = value
	case OptionType_Float:
		parsed, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("settings: invalid float for '%s': %s", opt.Name, value)
		}
	case OptionType_Level:
		level, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("settings: invalid level for '%s': %s", opt.Name, value)
		}
		if level > opt.MaxLevel {
			level = opt.MaxLevel
		}
		if level < 0 {
			level = 0
		}
		opt.levelValue = level
		parsed = level
	}

	opt.wasSet = true
	return setFieldByPath(opt.target, opt, parsed)
}
