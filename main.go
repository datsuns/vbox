package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
)

var DEFAULT_TOOL_PATH = fmt.Sprintf("C:%cProgram Files%cOracle%cVirtualBox%cVBoxManage.exe", filepath.Separator, filepath.Separator, filepath.Separator, filepath.Separator)

const toolPathOption = "tool"

var OptionFlags = []cli.Flag{
	cli.StringFlag{
		Name:  fmt.Sprintf("%s, t", toolPathOption),
		Value: DEFAULT_TOOL_PATH,
		Usage: "specify path to VBoxManage.exe",
	},
	cli.BoolFlag{
		Name:  "verbose, V",
		Usage: "verbose mode",
	},
}

var CommandList = []cli.Command{
	{
		Name:    "now",
		Aliases: []string{},
		Usage:   "list current VM status",
		Action:  cmdNow,
	},
	{
		Name:    "start",
		Aliases: []string{"r"},
		Usage:   "wakeup VM",
		Action:  cmdStart,
		After:   cmdNow,
	},
	{
		Name:    "gui",
		Aliases: []string{},
		Usage:   "wakeup VM gui mode",
		Action:  cmdStartGui,
		After:   cmdNow,
	},
	{
		Name:    "stop",
		Aliases: []string{},
		Usage:   "stop VM",
		Action:  cmdStop,
		After:   cmdNow,
	},
	{
		Name:    "restart",
		Aliases: []string{},
		Usage:   "restart VM",
		Action:  cmdStop,
		After: func(c *cli.Context) error {
			cmdStart(c)
			cmdNow(c)
			return nil
		},
	},
	{
		Name:    "help",
		Aliases: []string{},
		Usage:   "exec VBoxManage's help",
		Action:  cmdHelp,
	},
	{
		Name:    "cmd",
		Aliases: []string{},
		Usage:   "pass all params directly to VBoxManage",
		Action:  cmdCmd,
	},
}

func getGlobalContext(c *cli.Context) *cli.Context {
	parent := c.Parent()
	if parent != nil {
		return parent
	} else {
		return c
	}
}

func loadVbox(c *cli.Context) *Vbox {
	ctx := getGlobalContext(c)
	return NewVbox(ctx.String(toolPathOption), ctx.Bool("verbose"))
}

func cmdNow(c *cli.Context) error {
	vbox := loadVbox(c)
	all := vbox.AllVms()
	running := vbox.RunningVms()
	space := 0
	for k, _ := range all {
		if space < len(k) {
			space = len(k)
		}
	}
	fmt.Println("\nVM status:")
	for k, _ := range all {
		if _, exists := running[k]; exists {
			fmt.Printf(fmt.Sprintf("%%%ds: Run\n", space+1), k)
		} else {
			fmt.Printf(fmt.Sprintf("%%%ds: stop\n", space+1), k)
		}
	}
	return nil
}

// if VM name contains while space, such as "VM name",
// target name may be surrounded by "
func cmdStart(c *cli.Context) error {
	vbox := loadVbox(c)
	if c.NArg() == 0 {
		fmt.Print(" please specify VM image name")
		return nil
	}
	target := c.Args()[0]
	fmt.Printf(">> start [%s]\n", target)
	vbox.StartVm(target)
	return nil
}

func cmdStartGui(c *cli.Context) error {
	vbox := loadVbox(c)
	if c.NArg() == 0 {
		fmt.Print(" please specify VM image name")
		return nil
	}
	target := c.Args()[0]
	fmt.Printf(">> start [%s]\n", target)
	vbox.StartVmGui(target)
	return nil
}

func cmdStop(c *cli.Context) error {
	vbox := loadVbox(c)
	if c.NArg() == 0 {
		fmt.Print(" please specify VM image name")
		return nil
	}
	target := c.Args()[0]
	if target == "all" {
		for k, _ := range vbox.RunningVms() {
			fmt.Printf(">> stop [%s]\n", k)
			vbox.StopVm(k)
		}
	} else {
		fmt.Printf(">> stop [%s]\n", target)
		vbox.StopVm(target)
	}
	return nil
}

func cmdHelp(c *cli.Context) error {
	vbox := loadVbox(c)
	vbox.Help(c.Args())
	return nil
}

func cmdCmd(c *cli.Context) error {
	vbox := loadVbox(c)
	vbox.Command(c.Args())
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "vbox"
	app.Usage = "Virtual Box operation Tool"
	app.Version = "1.0.0"
	app.Commands = CommandList
	app.Action = nil
	app.Flags = OptionFlags

	app.Run(os.Args)
}
