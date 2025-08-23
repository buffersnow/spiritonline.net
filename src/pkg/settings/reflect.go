package settings

import (
	"fmt"
	"reflect"
	"strings"
)

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

func (o *Options) helpText(val reflect.Value, envprefix, section string) {

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	fmt.Printf("  %s options:\n", section)

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

			o.helpText(fiVal.Addr(), envPrefix, argSection)
			continue
		}

		helpTag := field.Tag.Get("help")
		envName := field.Tag.Get("env")
		argTag := field.Tag.Get("arg")

		isRequired := false
		if strings.Contains(envName, ",required") {
			isRequired = true
		}

		selArg := ""
		if len(envName) != 0 {
			placeholder := field.Type.Kind().String()
			envnamep := strings.Replace(envName, ",required", "", 1)
			selArg = fmt.Sprintf("%s <%s>", (envprefix + envnamep), placeholder)
		} else if len(argTag) != 0 {
			fullFlag, shortFlag := getArgs(argTag)
			placeholder := field.Type.Kind().String()

			if len(shortFlag) != 0 && len(fullFlag) != 0 {
				selArg = fmt.Sprintf("-%s,--%s <%s>", shortFlag, fullFlag, placeholder)
			} else if len(shortFlag) != 0 {
				selArg = fmt.Sprintf("-%s <%s>", shortFlag, placeholder)
			} else if len(fullFlag) != 0 {
				selArg = fmt.Sprintf("--%s <%s>", fullFlag, placeholder)
			}

		}

		defaultVal := ""
		switch fiVal.Kind() {
		case reflect.String:
			defaultTag := field.Tag.Get("default")
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

		if !isRequired {
			fmt.Printf("      %-*s%-*s(Default: %s)\n", (longestArg + 20), selArg, (longestHelp + 5), helpTag, defaultVal)
		} else {
			fmt.Printf("      %-*s%-*s(Required)\n", (longestArg + 20), selArg, (longestHelp + 5), helpTag)
		}
	}
}

func (o *Options) parseArgs(val reflect.Value) error {

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for idx := 0; idx < val.NumField(); idx++ {
		field, fiVal := typ.Field(idx), val.Field(idx)

		if !fiVal.CanSet() {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			if err := o.parseArgs(fiVal.Addr()); err != nil {
				return err
			}

			continue
		}

		helpTag := field.Tag.Get("help")
		if len(helpTag) != 0 && len(helpTag) > longestHelp {
			longestHelp = len(helpTag)
		}

		defaultVal := field.Tag.Get("default")
		if fiVal.Kind() != reflect.String && len(defaultVal) != 0 {
			return fmt.Errorf("default values are currently only supported on string options! (please fix..)")
		}

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
			continue //& this is an env variable
		}

		fullFlag, shortFlag := getArgs(argTag)

		switch fiVal.Kind() {
		case reflect.String:
			if len(defaultVal) == 0 {
				fs.StringVarP(fiVal.Addr().Interface().(*string), fullFlag, shortFlag, fiVal.String(), "")
			} else {
				fs.StringVarP(fiVal.Addr().Interface().(*string), fullFlag, shortFlag, defaultVal, "")
			}
		case reflect.Int:
			fs.IntVarP(fiVal.Addr().Interface().(*int), fullFlag, shortFlag, int(fiVal.Int()), "")
		case reflect.Bool:
			fs.BoolVarP(fiVal.Addr().Interface().(*bool), fullFlag, shortFlag, fiVal.Bool(), "")
		default:
			return fmt.Errorf("unsupported type %s on option %s", fiVal.Kind(), field.Name)
		}

	}

	return nil
}
