package nova

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strings"
	"text/tabwriter"
)

// ActionFunc defines the function signature for CLI actions (both global and command-specific).
// It receives a Context containing parsed flags and arguments.
type ActionFunc func(ctx *Context) error

// CLI represents the command-line application structure.
// It holds the application's metadata, commands, global flags, and the main action.
// Name and Version are required fields when creating a CLI instance via NewCLI.
type CLI struct {
	// Name of the application (Required).
	Name string
	// Version string for the application (Required). Displayed with the --version flag.
	Version string
	// Description provides a brief explanation of the application's purpose, shown in help text.
	Description string
	// Commands is a list of commands the application supports.
	Commands []*Command
	// Action is the default function to run when no command is specified.
	// If nil and no command is given, help is shown.
	Action ActionFunc
	// GlobalFlags are flags applicable to the application as a whole, before any command.
	// Flags named "version" or aliased "v" are reserved and cannot be used.
	// The "help" flag/alias is handled solely by the built-in 'help' command.
	GlobalFlags []Flag
	// Authors lists the application's authors, shown in the main help text.
	Authors string

	// internal storage for the parsed global flag set
	globalSet *flag.FlagSet
	// internal storage for the built-in nova flags
	internalGlobalFlags []Flag
	// Internal storage for flag values (pointers needed by std/flag)
	internalFlagPtrs map[string]any // Map flag name to pointer (*string, *int, etc.)
}

// Command defines a specific action the CLI application can perform.
// It includes metadata, flags specific to the command, and the action function.
// Commands named "help" or aliased "h" are reserved.
type Command struct {
	// Name is the primary identifier for the command (Required).
	Name string
	// Aliases provide alternative names for invoking the command. "h" is reserved.
	Aliases []string
	// Usage provides a short description of the command's purpose, shown in the main help list.
	Usage string
	// Description gives a more detailed explanation of the command, shown in the command's help.
	Description string
	// ArgsUsage describes the expected arguments for the command, shown in the command's help.
	// Example: "<input-file> [output-file]"
	ArgsUsage string
	// Action is the function executed when this command is invoked (Required).
	Action ActionFunc
	// Flags are the options specific to this command.
	// Flags named "help" or aliased "h" are reserved and cannot be used.
	Flags []Flag
}

// Flag defines the interface for command-line flags.
// Concrete types like StringFlag, IntFlag, BoolFlag implement this interface.
type Flag interface {
	fmt.Stringer // For generating help text representation.

	// Apply binds the flag definition to a standard Go flag.FlagSet.
	// It configures the flag's name, aliases, usage, and default value.
	Apply(set *flag.FlagSet, cli *CLI) error
	// GetName returns the primary name of the flag (e.g., "config").
	GetName() string
	// IsRequired indicates whether the flag must be provided by the user.
	IsRequired() bool
	// Validate checks if a required flag was provided correctly after parsing.
	Validate(set *flag.FlagSet) error
	// GetAliases returns the alternative names for the flag.
	GetAliases() []string
}

// Context provides access to parsed flags, arguments, and application/command metadata
// within an ActionFunc.
type Context struct {
	// CLI points to the parent application instance.
	CLI *CLI
	// Command points to the specific command being executed (nil for the global action).
	Command *Command
	// flagSet holds the parsed flag set relevant to the current context
	// (either command-specific or global if no command is running).
	flagSet *flag.FlagSet
	// args holds the non-flag arguments remaining after flag parsing.
	args []string
	// globalSet always holds the parsed global flag set for accessing global flags
	// even within a command context.
	globalSet *flag.FlagSet
}

// StringFlag defines a flag that accepts a string value.
type StringFlag struct {
	// Name is the primary identifier for the flag (e.g., "output"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "o"). Single letters are used as -o.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default string
	// Required indicates if the user must provide this flag.
	Required bool
}

// IntFlag defines a flag that accepts an integer value.
type IntFlag struct {
	// Name is the primary identifier for the flag (e.g., "port"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "p"). Single letters are used as -p.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default int
	// Required indicates if the user must provide this flag.
	Required bool
}

// BoolFlag defines a flag that acts as a boolean switch (true if present, false otherwise).
type BoolFlag struct {
	// Name is the primary identifier for the flag (e.g., "verbose"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "V"). Single letters are used as -V.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value. Note: Presence of the flag usually implies true.
	Default bool
	// Required is ignored for BoolFlag as presence implies the value.
	Required bool
}

// Float64Flag defines a flag that accepts a float64 value.
type Float64Flag struct {
	// Name is the primary identifier for the flag (e.g., "rate"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "r"). Single letters are used as -r.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default float64
	// Required indicates if the user must provide this flag.
	Required bool
}

// StringSliceFlag defines a flag that accepts multiple string values.
// The flag can be repeated on the command line (e.g., --tag foo --tag bar).
type StringSliceFlag struct {
	// Name is the primary identifier for the flag (e.g., "tag"). Used as --name. (Required)
	Name string
	// Aliases provide alternative names (e.g., "t"). Single letters are used as -t.
	Aliases []string
	// Usage provides a brief description of the flag's purpose for help text. (Required)
	Usage string
	// Default sets the default value if the flag is not provided.
	Default []string
	// Required indicates if the user must provide this flag at least once.
	Required bool
}

// stringSliceValue implements the flag.Value interface for string slices.
type stringSliceValue struct {
	destination *[]string // Pointer to the underlying slice
	hasBeenSet  bool      // Track if the flag was set at all (for default handling)
}

// newStringSliceValue creates a new value for string slice flags.
func newStringSliceValue(dest *[]string, defaults []string) *stringSliceValue {
	if dest != nil && (*dest == nil || len(*dest) == 0) {
		// Create a copy of defaults to avoid modifying the original Default slice
		*dest = make([]string, len(defaults))
		copy(*dest, defaults)
	}
	return &stringSliceValue{destination: dest}
}

// String returns a comma-separated representation of the slice (required by flag.Value).
func (s *stringSliceValue) String() string {
	if s.destination == nil || *s.destination == nil {
		return ""
	}
	return fmt.Sprintf("%q", *s.destination)
}

// Set appends the provided value to the slice. Called by flag package for each occurrence.
func (s *stringSliceValue) Set(value string) error {
	if s.destination == nil {
		return fmt.Errorf("internal error: stringSliceValue destination is nil")
	}
	// On the first call to Set for this flag instance, clear any default values.
	if !s.hasBeenSet {
		*s.destination = []string{}
		s.hasBeenSet = true
	}
	*s.destination = append(*s.destination, value)
	return nil
}

// Get returns the underlying slice.
func (s *stringSliceValue) Get() any {
	if s.destination == nil {
		return []string{} // Should not happen if constructed correctly
	}
	return *s.destination
}

// formatFlagNames combines the primary name and aliases into a comma-separated string
// suitable for help text (e.g., "-n, --name"). It sorts them for consistency.
func formatFlagNames(name string, aliases []string) string {
	names := []string{fmt.Sprintf("--%s", name)}
	for _, alias := range aliases {
		prefix := "--"
		if len(alias) == 1 {
			prefix = "-"
		}
		names = append(names, fmt.Sprintf("%s%s", prefix, alias))
	}
	sort.Slice(names, func(i, j int) bool {
		if len(names[i]) != len(names[j]) {
			return len(names[i]) < len(names[j]) // Prioritize shorter names (aliases)
		}
		return names[i] < names[j] // Then sort alphabetically
	})
	return strings.Join(names, ", ")
}

// validateRequiredFlag provides common validation logic for required flags (String, Int, Float64).
// It checks if the flag was set and, optionally, if its value is non-empty (for strings).
func validateRequiredFlag(set *flag.FlagSet, name string, required, checkEmpty bool) error {
	if !required {
		return nil
	}
	f := set.Lookup(name)
	if f == nil {
		return fmt.Errorf("internal error: flag --%s not found during validation", name)
	}

	wasSet := false
	set.Visit(func(visitedF *flag.Flag) {
		if visitedF.Name == name {
			wasSet = true
		}
	})

	if !wasSet {
		// For slice flags, check if the underlying slice is non-empty even if not explicitly "set"
		// (because defaults might populate it). This requires accessing the value.
		if getter, ok := f.Value.(flag.Getter); ok {
			if sliceVal, okSlice := getter.Get().([]string); okSlice {
				if len(sliceVal) > 0 {
					wasSet = true // Consider it set if it has values (even defaults)
				}
			}
		}
	}

	if !wasSet {
		return fmt.Errorf("required flag --%s was not provided", name)
	}

	if checkEmpty && f.Value.String() == "" {
		return fmt.Errorf("required flag --%s cannot be empty", name)
	}

	return nil
}

// validateRequiredSliceFlag handles validation for required slice flags.
func validateRequiredSliceFlag(set *flag.FlagSet, name string, required bool) error {
	if !required {
		return nil
	}
	f := set.Lookup(name)
	if f == nil {
		return fmt.Errorf("internal error: flag --%s not found during validation", name)
	}

	// Check if the flag was explicitly set by the user on the command line.
	userSet := false
	set.Visit(func(visitedF *flag.Flag) {
		if visitedF.Name == name {
			userSet = true
		}
	})

	// If required and not set by user, it's missing.
	if !userSet {
		return fmt.Errorf("required flag --%s was not provided", name)
	}

	// Additionally, ensure the resulting slice isn't empty if required.
	// This handles cases where the user might provide `--flag ""` which Set might accept.
	if getter, ok := f.Value.(flag.Getter); ok {
		if sliceVal, okSlice := getter.Get().([]string); okSlice {
			if len(sliceVal) == 0 {
				return fmt.Errorf("required flag --%s cannot be empty", name)
			}
		}
	}

	return nil
}

// findCommand searches the CLI's commands list for a command matching the given name or alias.
// Returns the *Command if found, otherwise nil. It checks commands added by the user
// AND the internal help command.
func (c *CLI) findCommand(name string) *Command {
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return cmd
		}
		if slices.Contains(cmd.Aliases, name) {
			return cmd
		}
	}
	// Check internal help command aliases after user commands.
	if name == "help" || name == "h" {
		return c.createHelpCommand()
	}
	return nil
}

// createHelpCommand generates the internal help command definition.
func (c *CLI) createHelpCommand() *Command {
	return &Command{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Shows a list of commands or detailed help for a command",
		Action: func(ctx *Context) error {
			args := ctx.Args()
			targetCmdName := ""
			if len(args) > 0 {
				targetCmdName = args[0]
			}

			if targetCmdName != "" {
				cmd := ctx.CLI.findCommand(targetCmdName)
				if cmd == nil || cmd.Name == "help" {
					return fmt.Errorf("unknown command %q for help", targetCmdName)
				}
				cmd.ShowHelp(os.Stdout, ctx.CLI.Name)
				return nil
			}
			ctx.CLI.ShowHelp(os.Stdout)
			return nil
		},
	}
}

// checkCommandConflicts validates a command against reserved names/aliases and basic requirements.
func checkCommandConflicts(cmd *Command) error {
	if cmd == nil {
		return fmt.Errorf("cannot process nil command")
	}
	if cmd.Name == "help" {
		return fmt.Errorf("command name 'help' is reserved by nova")
	}
	if cmd.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	if cmd.Action == nil {
		return fmt.Errorf("command '%s' action cannot be nil", cmd.Name)
	}
	if slices.Contains(cmd.Aliases, "h") {
		return fmt.Errorf("command alias 'h' is reserved by nova for command '%s'", cmd.Name)
	}
	for _, flag := range cmd.Flags {
		if err := checkFlagConflicts(flag, true); err != nil {
			return fmt.Errorf("command '%s': %w", cmd.Name, err)
		}
	}
	return nil
}

// checkFlagConflicts validates a flag against reserved names/aliases and basic requirements.
// `isCommandFlag` determines if command-level restrictions apply (reserving help/h).
func checkFlagConflicts(f Flag, isCommandFlag bool) error {
	if f == nil {
		return fmt.Errorf("cannot process nil flag")
	}
	name := f.GetName()
	if name == "" {
		return fmt.Errorf("flag name cannot be empty")
	}

	reservedNames := map[string]struct{}{}
	reservedAliases := map[string]struct{}{}

	if isCommandFlag {
		reservedNames["help"] = struct{}{}
		reservedAliases["h"] = struct{}{}
	} else {
		reservedNames["version"] = struct{}{}
		reservedAliases["v"] = struct{}{}
	}

	if _, reserved := reservedNames[name]; reserved {
		return fmt.Errorf("flag name '%s' is reserved by nova", name)
	}

	for _, alias := range f.GetAliases() {
		if _, reserved := reservedAliases[alias]; reserved {
			return fmt.Errorf("flag alias '%s' for flag '--%s' is reserved by nova", alias, name)
		}
	}

	// Basic validation for flag fields.
	var usage string

	switch ft := f.(type) {
	case *StringFlag:
		usage = ft.Usage
	case *IntFlag:
		usage = ft.Usage
	case *BoolFlag:
		usage = ft.Usage
	case *Float64Flag:
		usage = ft.Usage
	case *StringSliceFlag:
		usage = ft.Usage
	default:
		return fmt.Errorf("flag '--%s': unknown flag type %T", name, f)
	}

	if usage == "" {
		return fmt.Errorf("flag '--%s': Usage cannot be empty", name)
	}

	return nil
}

// setupInternalGlobalFlags defines and stores the built-in global flag.
func (c *CLI) setupInternalGlobalFlags() {
	c.internalGlobalFlags = []Flag{
		&BoolFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Print the application version",
		},
	}
}

// parseFlagSet handles the common logic of creating a flag.FlagSet, applying flags to it,
// and parsing arguments. It returns the parsed set and any parsing error.
func parseFlagSet(cli *CLI, flags []Flag, args []string, name string, output io.Writer) (*flag.FlagSet, error) {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.SetOutput(output)
	set.Usage = func() {}

	for _, f := range flags {
		if err := f.Apply(set, cli); err != nil {
			return nil, fmt.Errorf("error applying flag --%s for %q: %w", f.GetName(), name, err)
		}
	}

	err := set.Parse(args)
	if err != nil {
		return set, err
	}
	return set, nil
}

// validateFlags iterates through a list of user-defined flags and calls their Validate method
// against the parsed flag set. It collects and returns any validation errors.
func validateFlags(flags []Flag, set *flag.FlagSet) error {
	var errs []string
	for _, f := range flags {
		// Skip validation for internally managed flags.
		if bf, ok := f.(*BoolFlag); ok {
			if bf.Name == "version" && bf.Usage == "Print the application version" {
				continue
			}
		}
		// Validate user-defined flags.
		if err := f.Validate(set); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}
	return nil
}

// printFlagsHelp formats and prints a list of flags to the writer.
func printFlagsHelp(w io.Writer, title string, flags []Flag, isCommandHelp bool) {
	userFlags := []Flag{}
	var internalVersionFlag Flag

	for _, f := range flags {
		name := f.GetName()
		isInternal := false
		if bf, ok := f.(*BoolFlag); ok {
			if name == "version" && bf.Usage == "Print the application version" {
				internalVersionFlag = f
				isInternal = true
			}
		}
		if !isInternal {
			userFlags = append(userFlags, f)
		}
	}

	hasContent := len(userFlags) > 0 || internalVersionFlag != nil
	if !hasContent {
		return
	}

	fmt.Fprintln(w, title)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	sort.Slice(userFlags, func(i, j int) bool {
		return userFlags[i].GetName() < userFlags[j].GetName()
	})

	for _, f := range userFlags {
		fmt.Fprintf(tw, "  %s\n", f.String())
	}

	if internalVersionFlag != nil && !isCommandHelp {
		fmt.Fprintf(tw, "  %s\n", internalVersionFlag.String())
	}

	tw.Flush()
	fmt.Fprintln(w)
}

// ShowHelp prints the main help message for the application to the specified writer.
func (c *CLI) ShowHelp(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [global options] <command> [command options] [arguments...]\n", c.Name)
	fmt.Fprintf(w, "       %s help <command>\n\n", c.Name)
	fmt.Fprintf(w, "Version: %s\n\n", c.Version)

	if c.Description != "" {
		fmt.Fprintf(w, "%s\n\n", strings.TrimSpace(c.Description))
	}

	allCommands := make([]*Command, len(c.Commands)+1)
	copy(allCommands, c.Commands)
	allCommands[len(c.Commands)] = c.createHelpCommand()

	sort.Slice(allCommands, func(i, j int) bool {
		if allCommands[i].Name == "help" {
			return false
		}
		if allCommands[j].Name == "help" {
			return true
		}
		return allCommands[i].Name < allCommands[j].Name
	})

	if len(allCommands) > 0 {
		fmt.Fprintln(w, "Commands:")
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		for _, cmd := range allCommands {
			aliases := ""
			if len(cmd.Aliases) > 0 {
				aliases = fmt.Sprintf(" (%s)", strings.Join(cmd.Aliases, ", "))
			}
			fmt.Fprintf(tw, "  %s%s\t%s\n", cmd.Name, aliases, cmd.Usage)
		}
		tw.Flush()
		fmt.Fprintln(w)
	}

	printFlagsHelp(w, "Global Options:", append(c.GlobalFlags, c.internalGlobalFlags...), false)

	if c.Authors != "" {
		fmt.Fprintf(w, "\nAuthors: %s\n", c.Authors)
	}
}

// ShowHelp prints the help message for a specific command to the specified writer.
func (cmd *Command) ShowHelp(w io.Writer, appName string) {
	usage := fmt.Sprintf("Usage: %s %s", appName, cmd.Name)
	hasFlags := len(cmd.Flags) > 0
	if hasFlags {
		usage += " [command options]"
	}
	if cmd.ArgsUsage != "" {
		usage += " " + cmd.ArgsUsage
	} else if !hasFlags {
		usage += " [arguments...]"
	}
	fmt.Fprintln(w, usage)
	fmt.Fprintln(w)

	if cmd.Description != "" {
		fmt.Fprintln(w, "Description:")
		fmt.Fprintf(w, "  %s\n\n", strings.TrimSpace(cmd.Description))
	}

	printFlagsHelp(w, "Options:", cmd.Flags, true)

	if len(cmd.Aliases) > 0 {
		fmt.Fprintf(w, "Aliases: %s\n", strings.Join(cmd.Aliases, ", "))
	}
}

// GetName returns the primary name of the flag.
func (f *StringFlag) GetName() string { return f.Name }

// GetAliases returns the aliases for the flag.
func (f *StringFlag) GetAliases() []string { return f.Aliases }

// IsRequired returns true if the flag must be provided.
func (f *StringFlag) IsRequired() bool { return f.Required }

// Validate checks if the required flag was provided and is not empty.
func (f *StringFlag) Validate(set *flag.FlagSet) error {
	return validateRequiredFlag(set, f.Name, f.Required, true)
}

// String returns the help text representation of the flag.
func (f *StringFlag) String() string {
	namesStr := formatFlagNames(f.Name, f.Aliases)
	def := ""
	if f.Default != "" {
		def = fmt.Sprintf(" (default: %q)", f.Default)
	}
	req := ""
	if f.Required {
		req = " (Required)"
	}
	return fmt.Sprintf("%s\t%s%s%s", namesStr, f.Usage, def, req)
}

// Apply registers the string flag with the flag.FlagSet.
func (f *StringFlag) Apply(set *flag.FlagSet, cli *CLI) error {
	if _, exists := cli.internalFlagPtrs[f.Name]; exists {
		return nil
	}
	internalDest := f.Default
	cli.internalFlagPtrs[f.Name] = &internalDest
	set.StringVar(&internalDest, f.Name, f.Default, f.Usage)
	for _, alias := range f.Aliases {
		set.StringVar(&internalDest, alias, f.Default, f.Usage+" (alias for --"+f.Name+")")
	}
	return nil
}

// GetName returns the primary name of the flag.
func (f *IntFlag) GetName() string { return f.Name }

// GetAliases returns the aliases for the flag.
func (f *IntFlag) GetAliases() []string { return f.Aliases }

// IsRequired returns true if the flag must be provided.
func (f *IntFlag) IsRequired() bool { return f.Required }

// Validate checks if the required flag was provided.
func (f *IntFlag) Validate(set *flag.FlagSet) error {
	return validateRequiredFlag(set, f.Name, f.Required, false)
}

// String returns the help text representation of the flag.
func (f *IntFlag) String() string {
	namesStr := formatFlagNames(f.Name, f.Aliases)
	def := fmt.Sprintf(" (default: %d)", f.Default)
	req := ""
	if f.Required {
		req = " (Required)"
	}
	return fmt.Sprintf("%s\t%s%s%s", namesStr, f.Usage, def, req)
}

// Apply registers the int flag with the flag.FlagSet.
func (f *IntFlag) Apply(set *flag.FlagSet, cli *CLI) error {
	if _, exists := cli.internalFlagPtrs[f.Name]; exists {
		return nil
	}
	internalDest := f.Default
	cli.internalFlagPtrs[f.Name] = &internalDest
	set.IntVar(&internalDest, f.Name, f.Default, f.Usage)
	for _, alias := range f.Aliases {
		set.IntVar(&internalDest, alias, f.Default, f.Usage+" (alias for --"+f.Name+")")
	}
	return nil
}

// GetName returns the primary name of the flag.
func (f *BoolFlag) GetName() string { return f.Name }

// GetAliases returns the aliases for the flag.
func (f *BoolFlag) GetAliases() []string { return f.Aliases }

// IsRequired always returns false for boolean flags.
func (f *BoolFlag) IsRequired() bool { return false }

// Validate always returns nil for boolean flags.
func (f *BoolFlag) Validate(set *flag.FlagSet) error { return nil }

// String returns the help text representation of the flag.
func (f *BoolFlag) String() string {
	namesStr := formatFlagNames(f.Name, f.Aliases)
	def := ""
	if f.Default {
		def = " (default: true)"
	}
	return fmt.Sprintf("%s\t%s%s", namesStr, f.Usage, def)
}

// Apply registers the bool flag with the flag.FlagSet.
func (f *BoolFlag) Apply(set *flag.FlagSet, cli *CLI) error {
	// Special handling for internal --version flag to avoid storing its pointer unnecessarily
	// if we only access it via ctx.Bool("version") which checks the flagset directly.
	// However, for consistency, let's store it too.
	if _, exists := cli.internalFlagPtrs[f.Name]; exists {
		return nil
	}
	internalDest := f.Default
	cli.internalFlagPtrs[f.Name] = &internalDest
	set.BoolVar(&internalDest, f.Name, f.Default, f.Usage)
	for _, alias := range f.Aliases {
		set.BoolVar(&internalDest, alias, f.Default, f.Usage+" (alias for --"+f.Name+")")
	}
	return nil
}

// GetName returns the primary name of the flag.
func (f *Float64Flag) GetName() string { return f.Name }

// GetAliases returns the aliases for the flag.
func (f *Float64Flag) GetAliases() []string { return f.Aliases }

// IsRequired returns true if the flag must be provided.
func (f *Float64Flag) IsRequired() bool { return f.Required }

// Validate checks if the required flag was provided.
func (f *Float64Flag) Validate(set *flag.FlagSet) error {
	return validateRequiredFlag(set, f.Name, f.Required, false)
}

// String returns the help text representation of the flag.
func (f *Float64Flag) String() string {
	namesStr := formatFlagNames(f.Name, f.Aliases)
	def := fmt.Sprintf(" (default: %g)", f.Default) // Use %g for nice float formatting
	req := ""
	if f.Required {
		req = " (Required)"
	}
	return fmt.Sprintf("%s\t%s%s%s", namesStr, f.Usage, def, req)
}

// Apply registers the float64 flag with the flag.FlagSet.
func (f *Float64Flag) Apply(set *flag.FlagSet, cli *CLI) error {
	if _, exists := cli.internalFlagPtrs[f.Name]; exists {
		return nil
	}
	internalDest := f.Default
	cli.internalFlagPtrs[f.Name] = &internalDest
	set.Float64Var(&internalDest, f.Name, f.Default, f.Usage)
	for _, alias := range f.Aliases {
		set.Float64Var(&internalDest, alias, f.Default, f.Usage+" (alias for --"+f.Name+")")
	}
	return nil
}

// GetName returns the primary name of the flag.
func (f *StringSliceFlag) GetName() string { return f.Name }

// GetAliases returns the aliases for the flag.
func (f *StringSliceFlag) GetAliases() []string { return f.Aliases }

// IsRequired returns true if the flag must be provided.
func (f *StringSliceFlag) IsRequired() bool { return f.Required }

// Validate checks if the required flag was provided and resulted in a non-empty slice.
func (f *StringSliceFlag) Validate(set *flag.FlagSet) error {
	return validateRequiredSliceFlag(set, f.Name, f.Required)
}

// String returns the help text representation of the flag.
func (f *StringSliceFlag) String() string {
	namesStr := formatFlagNames(f.Name, f.Aliases)
	def := ""
	if len(f.Default) > 0 {
		// Display default slice using standard Go syntax
		def = fmt.Sprintf(" (default: %q)", f.Default)
	}
	req := ""
	if f.Required {
		req = " (Required)"
	}
	// Indicate that the flag can be specified multiple times
	usage := f.Usage + " (can be specified multiple times)"
	return fmt.Sprintf("%s\t%s%s%s", namesStr, usage, def, req)
}

// Apply registers the string slice flag with the flag.FlagSet using a custom Value.
func (f *StringSliceFlag) Apply(set *flag.FlagSet, cli *CLI) error {
	if _, exists := cli.internalFlagPtrs[f.Name]; exists {
		return nil
	}
	// Create the internal slice variable
	internalDest := make([]string, len(f.Default))
	copy(internalDest, f.Default)
	cli.internalFlagPtrs[f.Name] = &internalDest // Store pointer to the slice

	// Use the custom stringSliceValue, pointing it to our internal slice
	value := newStringSliceValue(&internalDest, f.Default) // Pass pointer to internal slice

	set.Var(value, f.Name, f.Usage)
	for _, alias := range f.Aliases {
		set.Var(value, alias, f.Usage+" (alias for --"+f.Name+")")
	}
	return nil
}

// Args returns the non-flag arguments remaining after parsing for the current context.
func (c *Context) Args() []string {
	return c.args
}

// String returns the string value of a flag specified by name.
// It checks command flags first, then global flags. Returns "" if not found or type mismatch.
func (c *Context) String(name string) string {
	if ptr, ok := c.CLI.internalFlagPtrs[name]; ok {
		if stringPtr, ok := ptr.(*string); ok {
			fs := c.flagSet // Command set or global set
			if fs.Lookup(name) != nil {
				return *stringPtr
			}
			// If not in current context's flagset, check global if different
			if c.globalSet != fs && c.globalSet.Lookup(name) != nil {
				return *stringPtr // Should be the same pointer, but check existence
			}
		}
	}
	return ""
}

// Int returns the integer value of a flag specified by name.
// It checks command flags first, then global flags. Returns 0 if not found or type mismatch.
func (c *Context) Int(name string) int {
	if ptr, ok := c.CLI.internalFlagPtrs[name]; ok {
		if intPtr, ok := ptr.(*int); ok {
			fs := c.flagSet
			if fs.Lookup(name) != nil {
				return *intPtr
			}
			if c.globalSet != fs && c.globalSet.Lookup(name) != nil {
				return *intPtr
			}
		}
	}
	return 0
}

// Bool returns the boolean value of a flag specified by name.
// It checks command flags first, then global flags. Returns false if not found or type mismatch.
func (c *Context) Bool(name string) bool {
	if ptr, ok := c.CLI.internalFlagPtrs[name]; ok {
		if boolPtr, ok := ptr.(*bool); ok {
			fs := c.flagSet
			if fs.Lookup(name) != nil {
				// For bools, presence matters. Check if it was set.
				wasSet := false
				fs.Visit(func(f *flag.Flag) {
					if f.Name == name {
						wasSet = true
					}
				})
				// Also check global if different
				if !wasSet && c.globalSet != fs {
					c.globalSet.Visit(func(f *flag.Flag) {
						if f.Name == name {
							wasSet = true
						}
					})
				}
				// Return the value pointed to, which reflects default or parsed value
				return *boolPtr
			}
		}
	}
	return false // Default if not found or wrong type
}

// Float64 returns the float64 value of a flag specified by name.
// It checks command flags first, then global flags. Returns 0.0 if not found or type mismatch.
func (c *Context) Float64(name string) float64 {
	if ptr, ok := c.CLI.internalFlagPtrs[name]; ok {
		if floatPtr, ok := ptr.(*float64); ok {
			fs := c.flagSet
			if fs.Lookup(name) != nil {
				return *floatPtr
			}
			if c.globalSet != fs && c.globalSet.Lookup(name) != nil {
				return *floatPtr
			}
		}
	}
	return 0.0
}

// StringSlice returns the []string value of a flag specified by name.
// It checks command flags first, then global flags. Returns nil if not found or type mismatch.
func (c *Context) StringSlice(name string) []string {
	if ptr, ok := c.CLI.internalFlagPtrs[name]; ok {
		if slicePtr, ok := ptr.(*[]string); ok {
			fs := c.flagSet
			if fs.Lookup(name) != nil {
				// Return a copy to prevent modification
				ret := make([]string, len(*slicePtr))
				copy(ret, *slicePtr)
				return ret
			}
			if c.globalSet != fs && c.globalSet.Lookup(name) != nil {
				ret := make([]string, len(*slicePtr))
				copy(ret, *slicePtr)
				return ret
			}
		}
	}
	return nil
}

// NewCLI creates and validates a new CLI application instance based on the provided configuration.
// The Name and Version fields in the input CLI struct are required.
// It checks for conflicts with reserved names/aliases (help command, h alias, version flag, v alias)
// and basic flag/command requirements.
// Returns the validated CLI instance or an error if validation fails.
func NewCLI(cli *CLI) (*CLI, error) {
	if cli == nil {
		return nil, fmt.Errorf("nova: CLI configuration cannot be nil")
	}
	if strings.TrimSpace(cli.Name) == "" {
		return nil, fmt.Errorf("nova: CLI Name cannot be empty")
	}
	if strings.TrimSpace(cli.Version) == "" {
		return nil, fmt.Errorf("nova: CLI Version cannot be empty")
	}

	// Validate user-defined commands.
	for _, cmd := range cli.Commands {
		if err := checkCommandConflicts(cmd); err != nil {
			return nil, fmt.Errorf("nova: validation error: %w", err)
		}
	}

	// Validate user-defined global flags.
	for _, f := range cli.GlobalFlags {
		// Global flags cannot use 'version', or 'v'. 'help'/'h' are not reserved for global flags.
		if err := checkFlagConflicts(f, false); err != nil {
			return nil, fmt.Errorf("nova: validation error: %w", err)
		}
	}

	cli.internalFlagPtrs = make(map[string]any)
	cli.setupInternalGlobalFlags()

	return cli, nil
}

// AddCommand registers a new command with the application *after* initial NewCLI() validation.
// It performs the same conflict checks as NewCLI(). It is generally recommended to define
// all commands and flags within the struct passed to NewCLI().
func (c *CLI) AddCommand(cmd *Command) error {
	if err := checkCommandConflicts(cmd); err != nil {
		return fmt.Errorf("nova: cannot add command: %w", err)
	}
	c.Commands = append(c.Commands, cmd)
	return nil
}

// Run executes the CLI application based on the provided command-line arguments.
// Call NewCLI() to create and validate the CLI instance before calling Run.
// Run parses flags, handles the built-in version flag and help command, validates required flags,
// and executes the appropriate action (global or command-specific).
func (c *CLI) Run(arguments []string) error {
	if c == nil || c.Name == "" || c.Version == "" || c.internalFlagPtrs == nil {
		return fmt.Errorf("nova: CLI not properly initialized; use NewCLI() first")
	}

	allGlobalFlags := append(c.GlobalFlags, c.internalGlobalFlags...)

	// Split off leading global flags up to the first non-flag (the command name).
	raw := arguments[1:]
	split := len(raw)
	for i, tok := range raw {
		if !strings.HasPrefix(tok, "-") || tok == "-" {
			split = i
			break
		}
	}
	globalArgs := raw[:split]
	restArgs := raw[split:]

	// Parse only the global flags.
	globalSet, err := parseFlagSet(c, allGlobalFlags, globalArgs, c.Name, io.Discard)
	c.globalSet = globalSet

	// Initial context uses the global flag set.
	ctx := &Context{CLI: c, globalSet: globalSet, flagSet: globalSet}

	// Handle global --version.
	if ctx.Bool("version") {
		fmt.Printf("%s version %s\n", c.Name, c.Version)
		return nil
	}

	// Handle global parsing errors.
	if err != nil {
		return err
	}

	// Validate required user-defined global flags.
	if err := validateFlags(c.GlobalFlags, globalSet); err != nil {
		return err
	}

	// Determine command (if any) and its args.
	args := restArgs
	var cmd *Command
	cmdArgs := args
	if len(args) > 0 {
		cmd = c.findCommand(args[0])
		if cmd != nil {
			cmdArgs = args[1:]
		}
	}

	// Execute command or default action.
	if cmd != nil {
		if cmd.Name == "help" {
			ctx.args = cmdArgs
			return cmd.Action(ctx)
		}

		// Parse command-specific flags.
		cmdSet, cmdErr := parseFlagSet(c, cmd.Flags, cmdArgs, cmd.Name, io.Discard)
		ctx.Command = cmd
		ctx.flagSet = cmdSet

		if cmdErr != nil {
			return cmdErr
		}
		if err := validateFlags(cmd.Flags, cmdSet); err != nil {
			return err
		}

		ctx.args = cmdSet.Args()
		return cmd.Action(ctx)

	} else if c.Action != nil {
		// No command, but default action defined.
		ctx.args = args
		return c.Action(ctx)

	} else {
		// No command and no default action: show help or error.
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q", args[0])
		}
		c.ShowHelp(os.Stdout)
		return nil
	}
}
