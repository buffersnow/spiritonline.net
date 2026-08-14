package settings

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func splitArgs(args string) (long, short string) {
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

func (p *Parser) findOptionByShort(short string) *Option {
	for _, opt := range p.allOpts {
		if opt.Short == short {
			return opt
		}
	}
	return nil
}

func (p *Parser) findOptionByLong(name string) (*Option, bool) {
	shouldInvert := false
	argBase := strings.TrimLeft(name, "-")

	for _, opt := range p.allOpts {
		optLong := strings.TrimLeft(opt.Long, "-")
		if optLong != "" && optLong == argBase {
			// If it's an exact match, it's only valid if the option isn't Negated
			// (because Negated options hide their base form and require --no-X).
			if opt.Negated {
				continue
			}
			return opt, false
		}
	}

	if strings.HasPrefix(argBase, "no-") {
		shouldInvert = true
		argBase = strings.TrimPrefix(argBase, "no-")
	}

	for _, opt := range p.allOpts {
		optLong := strings.TrimLeft(opt.Long, "-")
		if optLong == "" || optLong != argBase {
			continue
		}

		if shouldInvert {
			// --no-X is valid for Invertible or Negated options
			if !opt.Invertible && !opt.Negated {
				continue
			}
			return opt, true
		} else {
			// --X is hidden for Negated options
			if opt.Negated {
				continue
			}
			return opt, false
		}
	}

	return nil, false
}

func (p *Parser) parseLongOption(arg string, idx *int, args []string) error {
	var optName, optValue string
	var hasEquals bool

	if eqIdx := strings.Index(arg, "="); eqIdx != -1 {
		optName = arg[:eqIdx]
		optValue = arg[eqIdx+1:]
		hasEquals = true
	} else {
		optName = arg
	}

	opt, shouldInvert := p.findOptionByLong(optName)
	if opt == nil {
		return fmt.Errorf("settings: unknown option: %s", optName)
	}

	switch opt.Type {
	case OptionType_Boolean:
		var parsed bool
		hasExplicitValue := false

		if hasEquals {
			var err error
			parsed, err = strconv.ParseBool(optValue)
			if err != nil {
				return fmt.Errorf("settings: invalid bool for '%s': %s", optName, optValue)
			}
			hasExplicitValue = true
		} else if *idx+1 < len(args) {
			// Peek at next argument to see if it's explicitly a boolean (e.g., --flag true)
			if val, err := strconv.ParseBool(args[*idx+1]); err == nil {
				*idx++
				parsed = val
				hasExplicitValue = true
			}
		}

		if hasExplicitValue {
			// Explicit value: apply inversion only for --no-flag=value
			if shouldInvert {
				parsed = !parsed
			}
		} else {
			// No explicit value: --flag = true, --no-flag = false
			parsed = !shouldInvert
		}

		opt.wasSet = true
		return setFieldByPath(opt.target, opt, parsed)

	case OptionType_Level:
		if hasEquals {
			// e.g. --level=3
			return p.parseValueString(opt, optValue)
		}

		// Peek at next argument to see if it's a number (e.g. --level 3)
		if *idx+1 < len(args) {
			if _, err := strconv.Atoi(args[*idx+1]); err == nil {
				*idx++
				return p.parseValueString(opt, args[*idx])
			}
		}

		// Otherwise, fallback to standard increment (e.g. --verbose --verbose)
		opt.levelValue++
		if opt.levelValue > opt.MaxLevel {
			opt.levelValue = opt.MaxLevel
		}
		opt.wasSet = true
		return setFieldByPath(opt.target, opt, opt.levelValue)

	default:
		// String, Integer, Float
		if !hasEquals {
			*idx++
			if *idx >= len(args) {
				return fmt.Errorf("settings: option '%s' requires a value", optName)
			}
			optValue = args[*idx]
		}
		return p.parseValueString(opt, optValue)
	}
}

func (p *Parser) parseShortOptions(shortStr string, idx *int, args []string) error {
	i := 0
	for i < len(shortStr) {
		char := string(shortStr[i])
		opt := p.findOptionByShort(char)

		if opt == nil {
			return fmt.Errorf("settings: unknown option: -%s", char)
		}

		if opt.Type == OptionType_Level {
			count := 1
			// Handle stacked levels (e.g. -vvv)
			for i+1 < len(shortStr) && shortStr[i+1] == shortStr[i] {
				count++
				i++
			}

			remaining := shortStr[i+1:]
			if remaining != "" {
				// Handle attached value (e.g. -v3)
				return p.parseValueString(opt, remaining)
			}

			// Handle space-separated value (e.g. -v 3)
			if *idx+1 < len(args) {
				if _, err := strconv.Atoi(args[*idx+1]); err == nil {
					*idx++
					return p.parseValueString(opt, args[*idx])
				}
			}

			// Fallback to incrementing (e.g. -v)
			opt.levelValue += count
			if opt.levelValue > opt.MaxLevel {
				opt.levelValue = opt.MaxLevel
			}
			opt.wasSet = true
			if err := setFieldByPath(opt.target, opt, opt.levelValue); err != nil {
				return err
			}
			i++
			continue
		}

		remaining := shortStr[i+1:]

		switch opt.Type {
		case OptionType_Boolean:
			var val bool = true
			if remaining != "" {
				// Handle attached boolean (e.g. -btrue)
				parsed, err := strconv.ParseBool(remaining)
				if err != nil {
					return fmt.Errorf("settings: invalid bool for '-%s': %s", char, remaining)
				}
				val = parsed
			} else if *idx+1 < len(args) {
				// Handle space-separated boolean (e.g. -b true)
				if parsed, err := strconv.ParseBool(args[*idx+1]); err == nil {
					*idx++
					val = parsed
				}
			}
			opt.wasSet = true
			if err := setFieldByPath(opt.target, opt, val); err != nil {
				return err
			}
			i++
		default:
			// String, Integer, Float
			var value string
			if remaining != "" {
				value = remaining
			} else {
				*idx++
				if *idx >= len(args) {
					return fmt.Errorf("settings: option '-%s' requires a value", char)
				}
				value = args[*idx]
			}
			return p.parseValueString(opt, value)
		}
	}
	return nil
}

func (p *Parser) parseArgs(args []string) error {
	for _, opt := range p.allOpts {
		if opt.Type == OptionType_Level {
			opt.levelValue = 0
		}
		opt.wasSet = false
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		// Help interceptors
		if arg == "--help" || arg == "-h" {
			p.PrintHelp("")
			os.Exit(0)
		}
		if strings.HasPrefix(arg, "--help=") {
			moduleName := arg[7:]
			p.PrintHelp(moduleName)
			os.Exit(0)
		}

		if strings.HasPrefix(arg, "--") {
			if err := p.parseLongOption(arg, &i, args); err != nil {
				return err
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			if err := p.parseShortOptions(arg[1:], &i, args); err != nil {
				return err
			}
		} else {
			i++
		}
		i++
	}
	return nil
}
