// Package main provides colorized output from piped input.
// It is used like this.
//    $ cat -n FILE | colorize Printf '//.*$'
// Multiple patterns are allowed.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

// The program version.
var version = "0.8.0"
var termHighlight = "\033[1;31m"
var termReset = "\033[0m"
var termBold = "\033[1m"

// main entry point to the program.
func main() {
	rexs := loadRegexps()
	lines := loadPipelineData()
	for _, line := range lines {
		fmt.Printf("%v\n", highlight(rexs, line))
	}
}

// loadRegexps load the regular expressions to
// search for from the command line into a slice.
// It uses standard golang regular expressions.
func loadRegexps() []*regexp.Regexp {
	rexs := make([]*regexp.Regexp, 0)
	literal := false
	for i:=1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "--":
			literal = true
		case "-h","--help":
			if literal == false {
				printHelp()
			}
			fallthrough
		case "-V", "--version":
			if literal == false {
				printVersion()
			}
			fallthrough
		default:
			fmt.Printf("compiling %v\n", arg)
			rex := regexp.MustCompile(arg)
			rexs = append(rexs, rex)
		}
	}
	return rexs
}

// loadPipelineData loads the pipeline data into
// a buffer.
// This assumes non-streaming data. It expects
// file that may have been filtered by other
// tools like cat, tac, grep, etc.
func loadPipelineData() []string {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	return lines
}

// Colorizes the patterns that match in bold-red.
// The color could be parameterized.
func highlight(rexs []*regexp.Regexp, inLine string) string {
	line := inLine
	reset := "\033[0m"
	for _, rex := range rexs {
		matches := rex.FindStringSubmatch(line)
		if matches != nil {
			for _, match := range matches {
				repl := fmt.Sprintf("%v%v%v", termHighlight, match, reset)
				line = strings.Replace(line, match, repl, -1)
			}
		}
	}
	return line
}

// help displays on-line help and exits.
func printHelp() {
	txt := `
USAGE
    %[2]v%[4]v%[3]v PATTERNS

DESCRIPTION
    This program colorizes piped input.

    It is useful for highlighting specific items that may help
    debugging in log files and other unstructured data.

    It is different from colorizing using grep because all data is
    kept.

    Here is a simple example.

        $ %[2]vcat foo.txt%[3]v
        aaa
        bbb
        ccc
        ddd
        EEE
        $ %[2]vcat -n foo.txt | %[4]v 'bbb' '^cc' '(?i)eee'%[3]v
        aaa
        %[1]vbbb%[3]v
        %[1]vcc%[3]vc
        ddd
        %[1]vEEE%[3]v

    Note that you could also specify a single regular expression
    with OR syntax like this and get the same result.

       $ %[2]vcat -n foo.txt | %[4]v 'bbb|^cc|(?i)eee'%[3]v

    The regular expression syntax is the regular golang regular
    expression syntax that is described here:
    https://github.com/google/re2/wiki/Syntax.

OPTIONS
    --              Treat all options as patterns. This allows
                    the search to look for --help or --version.
                    For example the following command would
                    highlight the "%[1]v--help%[3]v" string.
                        %[2]vcat -n %[4]v.go | %[4]v -- --help %[3]v

    -h, --help      Print the program help and exit.

    -v, --version   Print the program version and exit.

AUTHOR
    Joe Linoff

VERSION
    %[5]v

LICENSE
   MIT Open Source
   Copyright (c) 2018

`
	fmt.Printf(txt, termHighlight, termBold, termReset, path.Base(os.Args[0]), version)
	os.Exit(0)
}

// version displays the program version and exits.
func printVersion() {
	fmt.Printf("%v-%v\n", path.Base(os.Args[0]), version)
	os.Exit(0)
}
