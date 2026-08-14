package settings

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (p *Parser) setDefaults() error {
	for _, opt := range p.allOpts {
		if opt.DefaultVal == "" {
			continue
		}
		switch opt.Type {
		case OptionType_Boolean:
			if val, err := strconv.ParseBool(opt.DefaultVal); err == nil {
				setFieldByPath(opt.target, opt, val)
			}
		case OptionType_Integer:
			if val, err := strconv.Atoi(opt.DefaultVal); err == nil {
				setFieldByPath(opt.target, opt, val)
			}
		case OptionType_String:
			setFieldByPath(opt.target, opt, opt.DefaultVal)
		case OptionType_Float:
			if val, err := strconv.ParseFloat(opt.DefaultVal, 64); err == nil {
				setFieldByPath(opt.target, opt, val)
			}
		case OptionType_Level:
			if val, err := strconv.Atoi(opt.DefaultVal); err == nil {
				opt.levelValue = val
				setFieldByPath(opt.target, opt, val)
			}
		}
	}
	return nil
}

func (p *Parser) validateRequired() error {
	var missing []string
	for _, opt := range p.allOpts {
		if !opt.Required {
			continue
		}

		v := reflect.ValueOf(opt.target)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		parent := v
		for _, part := range opt.Path {
			parent = parent.FieldByName(part)
		}
		field := parent.FieldByName(opt.Name)

		if field.IsZero() && !opt.wasSet {
			missing = append(missing, formatOptionName(opt))
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("settings: required options not set: %s", strings.Join(missing, ", "))
	}
	return nil
}
