package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/version"
	"github.com/luxploit/red"
)

var (
	shouldDisplayHelp red.ContextKey = "options_printhelp"
	helpDisplayTopic  red.ContextKey = "options_helptopic"
)

// Option defines a single option parsed from struct tags
type Option struct {
	Name       string     // The struct field name (e.g., "Cache")
	Path       []string   // Reflection path to the field (e.g., ["Server", "Cache"])
	Section    string     // Help menu section name derived from parent struct or module
	Short      string     // Short flag character (e.g., "v" becomes -v)
	Long       string     // Long flag name (e.g., "--cache" becomes --cache)
	Env        string     // Environment variable name to check (e.g., "APP_CACHE")
	HelpText   string     // Description shown in the help menu
	DefaultVal string     // Default value represented as a string
	Type       OptionType // The data type of the option (Bool, Int, etc.)
	Required   bool       // If true, parsing fails if the option remains zero-valued
	Invertible bool       // If true, both --flag and --no-flag are valid
	Negated    bool       // If true, the base --flag is hidden and ONLY --no-flag is accepted
	LevelChar  byte       // The character used for counting levels (e.g., 'v' for -vvv)
	MaxLevel   int        // The maximum level allowed for level-based options
	ModuleName string     // The name of the module this option belongs to
	target     any        // Internal: pointer to the parent module struct for fast reflection setting

	levelValue int
	wasSet     bool
}

// Module represents a registered options struct
type Module struct {
	Name       string
	PrettyName string
	Target     any
	Options    []*Option
}

type Parser struct {
	modules   []*Module
	allOpts   []*Option // Flat list for fast lookups
	name      string
	moduleXML map[string][]byte
}

func New(r *red.Context, v *version.BuildTag) (*Parser, error) {
	p := &Parser{
		name:      v.GetService(),
		moduleXML: make(map[string][]byte),
	}

	isHelp, topic := util.ParseHelpRequest()

	r.Set(shouldDisplayHelp, isHelp)
	r.Set(helpDisplayTopic, topic)

	return p, nil
}

// RegisterModule adds an options struct as a named module
func (p *Parser) RegisterModule(name, prettyName string, target any) error {
	mod := &Module{
		Name:       name,
		PrettyName: prettyName,
		Target:     target,
	}

	t := reflect.TypeOf(target)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("settings: module target must be a pointer to a struct, got %T", target)
	}

	// Extract options using the module name as the base section
	p.extractOptions(t, []string{}, prettyName, mod)
	p.modules = append(p.modules, mod)
	return nil
}

func ShouldDisplayHelp(r *red.Context) bool {
	return r.Get(shouldDisplayHelp).(bool)
}

// Responsible for printing the menu
func HandleHelpMenu(r *red.Context, p *Parser) error {

	if ShouldDisplayHelp(r) {
		p.PrintHelp(r.Get(helpDisplayTopic).(string))
		os.Exit(0)
	}

	return nil
}

func ParseSettings(r *red.Context, v *version.BuildTag, p *Parser) error {
	configDir := filepath.Join("cfg", v.GetService())

	for _, mod := range p.modules {
		path := filepath.Join(configDir, mod.Name+".xml")

		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // module has no on-disk config; rely on defaults/env/args
			}
			return fmt.Errorf("settings: read config for module '%s': %w", mod.Name, err)
		}

		if err := p.SetModuleXML(mod.Name, data); err != nil {
			return err
		}
	}

	return p.Parse(os.Args[1:], ShouldDisplayHelp(r))
}
