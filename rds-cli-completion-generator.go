package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/jsonutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/sortutil"
	"github.com/essentialkaos/ek/v13/strutil"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic utility info
const (
	APP  = "rds-cli-completion-generator"
	VER  = "0.0.3"
	DESC = "Tool to generate completion for RDS CLI"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	TYPE_ONEOF      = "oneof"
	TYPE_BLOCK      = "block"
	TYPE_PURE_TOKEN = "pure-token"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Info map[string]*Command

type Command struct {
	Container string    `json:"container"`
	Arguments Arguments `json:"arguments"`
}

type Argument struct {
	Name       string    `json:"name"`
	Token      string    `json:"token"`
	Type       string    `json:"type"`
	Display    string    `json:"display"`
	IsMultiple bool      `json:"multiple"`
	IsOptional bool      `json:"optional"`
	Arguments  Arguments `json:"arguments"`
}

type Arguments []*Argument

// ////////////////////////////////////////////////////////////////////////////////// //

type CommandInfo struct {
	Name      string
	Arguments []string
}

type InfoSlice []*CommandInfo

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// colors for usage info
var colorTagApp, colorTagVer, colorTagRel string

// gitrev is short hash of the latest git commit
var gitRev string

// ////////////////////////////////////////////////////////////////////////////////// //

// main is main utility function
func main() {
	preConfigureUI()

	_, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error(" - "))
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	}

	err := process()

	if err != nil {
		terminal.Error(err)
		os.Exit(1)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{#DC382C}", "{#A32422}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#160}", "{#124}"
	default:
		colorTagApp, colorTagVer = "{r*}", "{r}"
	}
}

// process starts processing of JSON files
func process() error {
	files := fsutil.List(
		"redis/src/commands", true,
		fsutil.ListingFilter{MatchPatterns: []string{"*.json"}},
	)

	if len(files) == 0 {
		return fmt.Errorf("Can't find any JSON files with commands info")
	}

	fsutil.ListToAbsolute("redis/src/commands", files)

	var commands InfoSlice

	for _, file := range files {
		info, err := extractCommandInfo(file)

		if err != nil {
			return err
		}

		commands = append(commands, info)
	}

	sort.Sort(sort.Reverse(commands))

	printCommandsCode(commands)

	return nil
}

// extractCommandInfo extract info from JSON file
func extractCommandInfo(file string) (*CommandInfo, error) {
	info := Info{}
	err := jsonutil.Read(file, &info)

	if err != nil {
		return nil, fmt.Errorf("Can't extract info from %s: %v", file, err)
	}

	cmdName, cmdArgs := info.Command()

	return &CommandInfo{cmdName, cmdArgs.Flatten()}, nil
}

// printCommandsCode prints code with commands info
func printCommandsCode(commands InfoSlice) {
	fmtc.NewLine()
	for _, c := range commands {
		if len(c.Arguments) == 0 {
			fmtc.Printfn(`  { {y}"%s"{!}, {*}nil{!}, {*}false{!} },`, c.Name)
		} else {
			args := strutil.JoinFunc(c.Arguments, ", ", func(s string) string {
				return fmt.Sprintf("{y}%s{!}", strconv.Quote(s))
			})

			fmtc.Printfn("  { {y}\"%s\"{!}, {*}[]string{!}{"+args+"}, {*}false{!} },", c.Name)
		}
	}
	fmtc.NewLine()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// String returns string representation of command info
func (i *CommandInfo) String() string {
	if len(i.Arguments) == 0 {
		return i.Name
	}

	return fmt.Sprintf("%s %s", i.Name, strings.Join(i.Arguments, " "))
}

// Command returns command name and arguments
func (i Info) Command() (string, Arguments) {
	for k, v := range i {
		cmdName := k

		if v.Container != "" {
			cmdName = v.Container + " " + k
		}

		return cmdName, v.Arguments
	}

	return "", nil
}

// Flatten returns flat slice with all supported arguments
func (a Arguments) Flatten() []string {
	var result []string

	for _, aa := range a {
		result = append(result, aa.String())
	}

	return result
}

// String returns string representation of command argument
func (a *Argument) String() string {
	v := strutil.Q(a.Display, a.Token, a.Name)

	if a.Token != "" && a.Type != TYPE_PURE_TOKEN {
		v = fmt.Sprintf("%s %s", a.Token, strutil.Q(a.Display, a.Name))
	}

	if a.Type == TYPE_ONEOF {
		var vv []string

		for _, aa := range a.Arguments {
			vv = append(vv, aa.String())
		}

		if a.Token != "" {
			v = fmt.Sprintf("%s {%s}", a.Token, strings.Join(vv, "|"))
		} else {
			v = fmt.Sprintf("{%s}", strings.Join(vv, "|"))
		}
	}

	if a.IsOptional {
		v = "[" + v + "]"
	}

	if a.IsMultiple {
		v += "…"
	}

	return v
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Len is the number of elements in the collection
func (s InfoSlice) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j
func (s InfoSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the element with index i
// must sort before the element with index j
func (s InfoSlice) Less(i, j int) bool {
	return sortutil.NaturalLess(s[i].Name, s[j].Name)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, APP))
	case "fish":
		fmt.Printf(fish.Generate(info, APP))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo()

	info.AppNameColorTag = colorTagApp

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,

		UpdateChecker: usage.UpdateChecker{
			"essentialkaos/rds-cli-completion-generator",
			update.GitHubChecker,
		},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
