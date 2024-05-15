package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/trzsz/go-arg"
)

type sshArgs struct {
	Ver         bool     `arg:"-V,--version" help:"show program's version number and exit"`
	Destination string   `arg:"positional" help:"alias in ~/.ssh/config, or [user@]hostname[:port]"`
	Command     string   `arg:"positional" help:"command to execute instead of a login shell"`
	Argument    []string `arg:"positional" help:"command arguments separated by spaces"`
	Debug       bool     `arg:"-v,--debug" help:"verbose mode for debugging, same as ssh's -vvv"`
}

func (sshArgs) Description() string {
	return "Simple ssh client with trzsz ( trz / tsz ) support.\n"
}

func (sshArgs) Version() string {
	return fmt.Sprintf("trzsz ssh %s", kTsshVersion)
}

type Parser struct {
	*arg.Parser
}

func (p *Parser) WriteHelp(w io.Writer) {
	var b bytes.Buffer
	p.Parser.WriteHelp(&b)
	s := strings.Replace(b.String(), " [-v]", "", 1)
	s = strings.Replace(s, "  -v, --version          show program's version number and exit\n", "", 1)
	fmt.Fprint(w, s)

}

func (p *Parser) WriteUsage(w io.Writer) {
	var b bytes.Buffer
	p.Parser.WriteUsage(&b)
	s := strings.Replace(b.String(), " [-v]", "", 1)
	fmt.Fprint(w, s)

}

func NewParser(config arg.Config, dests ...interface{}) (*Parser, error) {
	p, err := arg.NewParser(config, dests...)
	return &Parser{p}, err
}

const kTsshVersion = "0.1.19"

func main() {

	log.SetFlags(log.Lshortfile)
	log.Println("Override built in option `-v, --version` of package `arg`")
	var (
		args sshArgs
	)

	// Like `parser := arg.MustParse(args)` but override built in option `-v, --version` of package `arg`
	parser, err := NewParser(arg.Config{}, &args)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(-1)
	}
	a2s := make([]string, 0) // without built in option
	deb := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "-help", "--help":
			parser.WriteHelp(os.Stderr)
			os.Exit(0)
		case "-version", "--version":
			fmt.Fprintln(os.Stderr, args.Version())
			os.Exit(0)
		case "-v":
			deb = true
		default:
			a2s = append(a2s, arg)
		}
	}
	err = parser.Parse(a2s)
	if err != nil {
		parser.WriteUsage(os.Stdout)
		fmt.Fprintln(os.Stdout, err)
		os.Exit(0)
	}
	if args.Ver {
		fmt.Fprintln(os.Stderr, args.Version())
		os.Exit(0)
	}
	args.Debug = args.Debug || deb
	log.Println("Debug", args.Debug)
	log.Println("Ver", args.Ver)
	log.Println("Destination", args.Destination)
	log.Println("Command", args.Command)
	log.Println("Argument", args.Argument)
}
