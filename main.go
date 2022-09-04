package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Globals struct {
	Version VersionFlag `short:"v" name:"version" help:"print the version."`
}

var Version string
var Revision = "HEAD"
var embedVersion = "0.0.1"

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	if Version == "" {
		Version = embedVersion
	}
	fmt.Printf("sshcmd version %s (rev:%s)\n", Version, Revision)
	app.Exit(0)

	return nil
}

var CLI struct {
	Globals
	Address  string `required:"" name:"address" short:"a" help:"Specify ssh connection destination"`
	Port     int64  `required:"" name:"port" short:"p" help:"Specify ssh host port"`
	User     string `required:"" name:"user" short:"u" help:"Specify ssh user"`
	Cmd      string `required:"" name:"cmd" short:"c" help:"Specify the command to run on the ssh server"`
	AfterCmd string `name:"after-cmd" short:"f" help:"Specify commands to run on the host server"`
}

func main() {
	_ = kong.Parse(&CLI,
		kong.Name("sshcmd"),
		kong.Description("Execute command to ssh host"),
	)

	// connect ssh server
	fmt.Printf("%s@%s's password: ", CLI.User, CLI.Address)
	pass, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatalf("failed read password: %s\n", err)
	}
	fmt.Println("")

	config := &ssh.ClientConfig{
		User: "pi",
		Auth: []ssh.AuthMethod{
			ssh.Password(string(pass)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", CLI.Address, CLI.Port), config)
	if err != nil {
		log.Fatalf("failed to dial: %s\n", err)
	}
	defer client.Close()

	// execute commands on the ssh connection
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed create session: %s\n", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(CLI.Cmd); err != nil {
		log.Fatalf("failed to run session: %s\n", err)
	}
	defer session.Close()

	fmt.Fprintln(os.Stderr, "--- ssh exec result ---")
	fmt.Println(b.String())

	// execute commands on host server
	if CLI.AfterCmd != "" {
		fmt.Fprintln(os.Stderr, "--- after exec result ---")

		cmd := exec.Command("bash", "-c", CLI.AfterCmd)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Fatalf("faild exec after exec cmd: %s", err)
		}
	}

	os.Exit(0)
}
