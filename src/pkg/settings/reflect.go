package settings

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	optionType_Environment = iota
	optionType_Argument
)

type parsedOption struct {
	optionType  int
	helpText    string
	tagContent  string
	defaultVal  string
	placeholder string
}

type parsedSection struct {
	sectionName string
	options     []parsedOption
}

var cachedOptions []parsedSection

var longestArg int
var longestHelp int

func getArgs(args string) (long, short string) {
	argParts := strings.Split(args, ",")
	argNameA, argNameB := argParts[0], ""

	if len(argParts) >= 2 {
		argNameB = argParts[1]
	}

	fullFlag, shortFlag := "", ""
	if s, ok := strings.CutPrefix(argNameA, "--"); ok {
		fullFlag = s
		argNameA = s
	} else if s, ok := strings.CutPrefix(argNameB, "--"); ok {
		fullFlag = s
		argNameB = s
	}

	if s, ok := strings.CutPrefix(argNameA, "-"); ok {
		shortFlag = s
		argNameA = s
	} else if s, ok := strings.CutPrefix(argNameB, "-"); ok {
		shortFlag = s
		argNameB = s
	}

	return fullFlag, shortFlag
}

func (o *Options) parseOptions(val reflect.Value, envprefix, section string) error {
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	typ := val.Type()

	psection := parsedSection{sectionName: section}

	for idx := 0; idx < val.NumField(); idx++ {
		field, fiVal := typ.Field(idx), val.Field(idx)

		if !fiVal.CanSet() {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			envPrefix := field.Tag.Get("envPrefix")
			argSection := field.Tag.Get("section")
			if len(argSection) == 0 {
				argSection = field.Name
			}

			if err := o.parseOptions(fiVal.Addr(), envPrefix, argSection); err != nil {
				return err
			}

			continue
		}

		option := parsedOption{placeholder: field.Type.Kind().String()}

		helpTag := field.Tag.Get("help")
		if len(helpTag) != 0 && len(helpTag) > longestHelp {
			longestHelp = len(helpTag)
		}
		option.helpText = helpTag

		defaultVal := field.Tag.Get("default")
		if fiVal.Kind() != reflect.String && len(defaultVal) != 0 {
			return fmt.Errorf("default values are currently only supported on string options! (please fix..)")
		}

		defaultTag := field.Tag.Get("default")
		if len(defaultVal) == 0 {
			switch fiVal.Kind() {
			case reflect.String:
				if len(defaultTag) != 0 {
					defaultVal = defaultTag
				} else {
					defaultVal = "\"\""
				}
			case reflect.Int:
				defaultVal = "0"
			case reflect.Bool:
				defaultVal = "false"
			case reflect.Map:
				defaultVal = "[empty map]"
			}

			envDefaultTag := field.Tag.Get("envDefault")
			if len(envDefaultTag) != 0 {
				defaultVal = envDefaultTag
			}
		}
		option.defaultVal = defaultVal

		envName := field.Tag.Get("env")
		argTag := field.Tag.Get("arg")

		if len(argTag) == 0 && len(envName) == 0 {
			return fmt.Errorf("option %s is neither an arg nor an env setting", field.Name)
		}

		if len(argTag) != 0 && len(argTag) > longestArg {
			longestArg = len(argTag)
		} else if len(envName) != 0 && len(envName) > longestArg {
			longestArg = len(envName)
		}

		if len(envName) != 0 {
			option.optionType = optionType_Environment
			option.tagContent = envprefix + envName
			psection.options = append(psection.options, option)
			continue //& this is an env variable
		}

		option.optionType = optionType_Argument
		option.tagContent = argTag
		fullFlag, shortFlag := getArgs(argTag)

		switch fiVal.Kind() {
		case reflect.String:
			if len(defaultTag) == 0 {
				fs.StringVarP(fiVal.Addr().Interface().(*string), fullFlag, shortFlag, fiVal.String(), "")
			} else {
				fs.StringVarP(fiVal.Addr().Interface().(*string), fullFlag, shortFlag, defaultTag, "")
			}
		case reflect.Int:
			fs.IntVarP(fiVal.Addr().Interface().(*int), fullFlag, shortFlag, int(fiVal.Int()), "")
		case reflect.Bool:
			fs.BoolVarP(fiVal.Addr().Interface().(*bool), fullFlag, shortFlag, fiVal.Bool(), "")
		default:
			return fmt.Errorf("unsupported type %s on option %s", fiVal.Kind(), field.Name)
		}

		psection.options = append(psection.options, option)
	}

	if psection.sectionName != "" {
		cachedOptions = append(cachedOptions, psection)
	}

	return nil
}

func (o *Options) helpText() {

	for _, sections := range cachedOptions {
		fmt.Printf("  %s options:\n", sections.sectionName)

		for _, option := range sections.options {
			isRequired := false
			if strings.Contains(option.tagContent, ",required") {
				isRequired = true
			}

			selArg := ""
			if option.optionType == optionType_Environment {
				envnamep := strings.Replace(option.tagContent, ",required", "", 1)
				selArg = fmt.Sprintf("%s <%s>", envnamep, option.placeholder)
			} else {
				fullFlag, shortFlag := getArgs(option.tagContent)
				if len(shortFlag) != 0 && len(fullFlag) != 0 {
					selArg = fmt.Sprintf("-%s,--%s <%s>", shortFlag, fullFlag, option.placeholder)
				} else if len(shortFlag) != 0 {
					selArg = fmt.Sprintf("-%s <%s>", shortFlag, option.placeholder)
				} else if len(fullFlag) != 0 {
					selArg = fmt.Sprintf("--%s <%s>", fullFlag, option.placeholder)
				}

			}

			if !isRequired {
				fmt.Printf("      %-*s%-*s(Default: %s)\n", (longestArg + 20), selArg, (longestHelp + 5), option.helpText, option.defaultVal)
			} else {
				fmt.Printf("      %-*s%-*s(Required)\n", (longestArg + 20), selArg, (longestHelp + 5), option.helpText)
			}
		}

	}
}
