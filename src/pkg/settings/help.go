package settings

import (
	"fmt"
	"strings"
)

type helpLine struct {
	opt    *Option
	selArg string
	help   string
}

type helpModule struct {
	name       string
	prettyName string
	lines      []helpLine
	merged     bool
}

func formatOptionName(opt *Option) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("<%s>", opt.Name))

	if opt.Short != "" {
		parts = append(parts, fmt.Sprintf("-%s", opt.Short))
	}

	if opt.Long != "" {
		if opt.Negated {
			baseLong := strings.TrimPrefix(opt.Long, "--")
			parts = append(parts, fmt.Sprintf("--no-%s", baseLong))
		} else {
			parts = append(parts, fmt.Sprintf("--%s", opt.Long))
		}
	}

	if opt.Invertible && opt.Long != "" && !opt.Negated {
		baseLong := strings.TrimPrefix(opt.Long, "--")
		parts = append(parts, fmt.Sprintf("--no-%s", baseLong))
	}

	if opt.Env != "" {
		parts = append(parts, fmt.Sprintf("[%s]", opt.Env))
	}

	return strings.Join(parts, ", ")
}

func (opt *Option) getDefaultDisplay() string {
	if opt.DefaultVal != "" {
		return opt.DefaultVal
	}
	switch opt.Type {
	case OptionType_Boolean:
		return "false"
	case OptionType_Integer, OptionType_Level:
		return "0"
	case OptionType_Float:
		return "0.0"
	case OptionType_String:
		return "\"\""
	}
	return ""
}

// PrintHelp prints help for a specific module, or all modules if name is ""
func (p *Parser) PrintHelp(moduleName string) {
	longestArg := 0
	longestHelp := 0
	var helpModules []helpModule = []helpModule{}

	exec := fmt.Sprintf("./%s{.exe/.lxb}", p.name)
	fmt.Printf("Usage: %s <... arguments>\n\n", exec)

	// Build help modules, merging modules with the same Name
	for _, mod := range p.modules {
		var lines []helpLine

		for _, opt := range mod.Options {
			selArg := formatOptionName(opt)
			helpText := opt.HelpText
			if opt.Required {
				helpText += " (Required)"
			}
			if opt.Type == OptionType_Level {
				helpText += fmt.Sprintf(" (Max Level: %d)", opt.MaxLevel)
			}
			if opt.Invertible && !opt.Negated && (opt.Type == OptionType_Boolean || opt.Type == OptionType_Level) {
				helpText += " (Invertible)"
			}

			if len(selArg) > longestArg {
				longestArg = len(selArg)
			}
			if len(helpText) > longestHelp {
				longestHelp = len(helpText)
			}

			lines = append(lines, helpLine{opt: opt, selArg: selArg, help: helpText})
		}

		// Find existing helpModule with the same internal Name
		var target *helpModule
		for i := range helpModules {
			if helpModules[i].name == mod.Name {
				target = &helpModules[i]
				break
			}
		}

		// Merge if found, otherwise append as new
		if target != nil {
			target.lines = append(target.lines, lines...)
			target.merged = true
		} else {
			helpModules = append(helpModules, helpModule{
				name:       mod.Name,
				prettyName: mod.PrettyName,
				lines:      lines,
			})
		}
	}

	// If asking for a specific module
	if moduleName != "" {
		for _, mod := range helpModules {
			// Allow lookup by both PrettyName and internal Name
			if mod.prettyName == moduleName || mod.name == moduleName {
				p.printModuleHelp(&mod, longestArg, longestHelp)
				return
			}
		}
	}

	for _, mod := range helpModules {
		p.printModuleHelp(&mod, longestArg, longestHelp)
	}
}

func (p *Parser) printModuleHelp(mod *helpModule, longestArg int, longestHelp int) {
	// Check if we need to print sub-sections (from nested structs)
	// An option belongs to a sub-section if it has a reflection path
	hasSubSections := false
	for _, line := range mod.lines {
		if len(line.opt.Path) > 0 {
			hasSubSections = true
			break
		}
	}

	// Determine the display name for the header
	displayName := mod.prettyName
	if mod.merged {
		displayName = "Core"
	}

	if hasSubSections {
		fmt.Printf("%s options <%s>:\n", displayName, mod.name)
		curSection := ""
		for _, line := range mod.lines {
			if line.opt.Section != curSection {
				if curSection != "" {
					fmt.Println()
				}
				fmt.Printf("  %s:\n", line.opt.Section)
				curSection = line.opt.Section
			}
			fmt.Printf("      %-*s%-*s(Default: %s)\n",
				(longestArg + 5), line.selArg,
				(longestHelp + 5), line.help,
				line.opt.getDefaultDisplay())
		}
	} else {
		fmt.Printf("%s options <%s>:\n", displayName, mod.name)
		for _, line := range mod.lines {
			fmt.Printf("      %-*s%-*s(Default: %s)\n",
				(longestArg + 5), line.selArg,
				(longestHelp + 5), line.help,
				line.opt.getDefaultDisplay())
		}
	}
}
