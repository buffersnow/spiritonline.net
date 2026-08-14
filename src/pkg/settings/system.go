package settings

import (
	"encoding/xml"
	"fmt"
	"os"
)

func (p *Parser) SetModuleXML(moduleName string, data []byte) error {
	for _, mod := range p.modules {
		if mod.Name == moduleName {
			p.moduleXML[moduleName] = data
			return nil
		}
	}
	return fmt.Errorf("settings: cannot set XML for unregistered module '%s'", moduleName)
}

func (p *Parser) parseEnv() error {
	for _, opt := range p.allOpts {
		if opt.Env == "" {
			continue
		}
		envVal := os.Getenv(opt.Env)
		if envVal == "" {
			continue
		}
		if err := p.parseValueString(opt, envVal); err != nil {
			return fmt.Errorf("settings: env error for '%s': %w", opt.Env, err)
		}
	}
	return nil
}

// Parse executes the parsing pipeline (Priority: Args > Env > XML > Defaults)
func (p *Parser) Parse(args []string, helpmenu bool) error {
	if err := p.setDefaults(); err != nil {
		return err
	}

	for _, mod := range p.modules {
		data, ok := p.moduleXML[mod.Name]
		if !ok || len(data) == 0 {
			continue
		}
		if err := xml.Unmarshal(data, mod.Target); err != nil {
			return fmt.Errorf("settings: XML error in module '%s': %w", mod.Name, err)
		}
	}

	if err := p.parseEnv(); err != nil {
		return err
	}

	if args != nil {
		if err := p.parseArgs(args); err != nil {
			return err
		}
	}

	if !helpmenu {
		if err := p.validateRequired(); err != nil {
			return err
		}
	}

	return nil
}
