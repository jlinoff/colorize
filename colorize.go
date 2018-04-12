// Package main provides colorized output from piped input.
// It is used like this.
//    $ cat -n FILE | colorize Printf '//.*$'
// Multiple patterns are allowed.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// The program version.
var version = "0.8.1"
var termHighlight = "\033[1;31m" // default, can be overwritten by -c
var termReset = "\033[0m"
var termBold = "\033[1m"
var verbose = 0
var inputFile = ""

// main entry point to the program.
func main() {
	rexs, colorMap := loadRegexps()
	lines := loadPipelineData()
	highlight(rexs, colorMap, lines)
	if verbose > 0 {
		fmt.Printf("INFO: done\n")
	}
}

// loadRegexps load the regular expressions to
// search for from the command line into a slice.
// It uses standard golang regular expressions.
func loadRegexps() ([]*regexp.Regexp, []string) {
	rexs := make([]*regexp.Regexp, 0)
	colorMap := make([]string, 0)
	colorMapSpecs :=  make([]string, 0)
	literal := false
	for i:=1; i < len(os.Args); i++ {
		arg := os.Args[i]
		add := false
		switch arg {
		case "--":
			if literal == false {
				// This enables -- -- .
				literal = true
			} else {
				add = true
			}
		case "-c","--color-map":
			if literal == false {
				i++
				if i < len(os.Args) {
					colorMapSpecs = append(colorMapSpecs, os.Args[i])
				} else {
					err := fmt.Errorf("ERROR: missing argument for '%v'.", arg)
					log.Fatal(err)
				}
			} else {
				add = true
			}
		case "-h","--help":
			if literal == false {
				printHelp()
			} else {
				add = true
			}
		case "-i","--input":
			if literal == false {
				i++
				if i < len(os.Args) {
					inputFile = os.Args[i]
				} else {
					err := fmt.Errorf("ERROR: missing argument for '%v'.", arg)
					log.Fatal(err)
				}
			}
		case "-V", "--version":
			if literal == false {
				printVersion()
			} else {
				add = true
			}
		case "-v", "-vv", "--verbose":
			if literal == false {
				verbose++
				if arg == "-vv" {
					verbose++
				}
			} else {
				add = true
			}
		default:
			if literal == false {
				if arg[0] == '-' {
					err := fmt.Errorf("ERROR: unrecognized option '%v'.", arg)
					log.Fatal(err)
				}
			}
			add = true
		}

		// Add a new regular expression.
		if add {
			if verbose > 0 {
				fmt.Printf("INFO: Compiling regular expression '%v'\n", arg)
			}
			rex := regexp.MustCompile(arg)
			rexs = append(rexs, rex)
		}
		add = false
	}
	if verbose > 0 {
		fmt.Printf("INFO: program : %v-%v\n", path.Base(os.Args[0]), version)
		fmt.Printf("INFO: input : '%v'\n", inputFile)
	}
	for _, cms := range colorMapSpecs {
		colorMap = updateColorMap(colorMap, cms)
	}
	colorMap = extendColorMap(colorMap, len(rexs))
	if verbose > 0  {
		fmt.Printf("INFO: num patterns : %2d\n", len(rexs))
		if verbose > 1 {
			for i, rex := range rexs {
				fmt.Printf("INFO: pattern[%d] : '%v'\n", i+1, rex)
			}
		}
		fmt.Printf("INFO: num colorMaps : %2d\n", len(rexs))
		if verbose > 1 {
			for i, cm := range colorMap {
				fmt.Printf("INFO: color[%d] : '%vtest\033[0m'\n", i+1, cm)
			}
		}
	}
	return rexs, colorMap
}

// extendColorMap extends the color map so it is the same length or
// longer than the rexs. This allows index lookups.
func extendColorMap(colorMap []string, size int) []string {
	if len(colorMap) == 0 {
		colorMap = append(colorMap, termHighlight)
	}
	le := colorMap[len(colorMap)-1]
	for len(colorMap) < size {
		colorMap = append(colorMap, le)
	}
	return colorMap
}

// updateColorMap updates the color map from the
// -c or --color-map options.
func updateColorMap(colorMap []string, colorSpec string) []string {
	colorTable := map[string]string {
		"blue": "34",
		"blueB": "44",
		"bold": "1",
		"cyan": "36",
		"cyanB": "46",
		"faint": "2",
		"gray": "38;5;245",
		"grayB": "48;5;245",
		"green": "32",
		"greenB": "42",
		"italic": "3",
		"magenta": "35",
		"magentaB": "45",
		"normal": "22",
		"red": "31",
		"redB": "41",
		"reset": "0",
		"reverse": "7",
		"strike": "9",
		"underline": "4",
		"white": "36",
		"whiteB": "46",
		"yellow": "33",
		"yellowB": "43",
	}
	groups := strings.Split(colorSpec, ",")
	for g, group := range groups {
		parts := strings.Split(group, "+")
		var code string
		for i, part := range parts {
			if i > 0 {
				// Prepend the semicolon, if necessary.
				code += ";"
			}
			if _, err := strconv.ParseInt(part, 10, 64); err == nil {
				// It is a  digit.
				code += part
			} else if val, ok := colorTable[part] ; ok {
				// Get the value from the table.
				code += (val)
			} else {
				// Use the raw value.
				// This allows complex expressions like:
				// -c 'red+32;2;255;82;197;48;2;155;106;0'
				code += (part)
			}
		}
		if verbose > 0 {
			fmt.Printf("INFO: colorMap[%d] = '%vm'\n", g+1, code)
		}
		cc := fmt.Sprintf("\033[%vm", code)
		colorMap = append(colorMap, cc)
	}
	return colorMap
}

// loadPipelineData loads the pipeline data into
// a buffer.
// This assumes non-streaming data. It expects
// file that may have been filtered by other
// tools like cat, tac, grep, etc.
func loadPipelineData() []string {
	var data []byte
	var err error
	if len(inputFile) == 0 {
		if verbose > 0 {
			fmt.Printf("INFO: loading from stdin\n")
		}
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		if verbose > 0 {
			fmt.Printf("INFO: loading from file '%v'\n", inputFile)
		}
		data, err = ioutil.ReadFile(inputFile)
	}
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	return lines
}

// highlight the data
func highlight(rexs []*regexp.Regexp, colorMap []string, lines []string) {
	if verbose > 0 {
		fmt.Printf("INFO: highlight\n")
	}
	for _, line := range lines {
		fmt.Printf("%v\n", highlightLine(rexs, colorMap, line))
	}
}

// Colorizes the patterns that match in bold-red.
// The color could be parameterized.
func highlightLine(rexs []*regexp.Regexp, colorMap []string, inLine string) string {
	line := inLine
	reset := "\033[0m"
	for i, rex := range rexs {
		code := colorMap[i]
		matches := rex.FindStringSubmatch(line)
		if matches != nil {
			for _, match := range matches {
				repl := fmt.Sprintf("%v%v%v", code, match, reset)
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
    %[2]v%[4]v%[3]v [OPTIONS] PATTERNS

DESCRIPTION
    This program colorizes unstructured text data.

    It is useful for highlighting specific items that may help
    debugging issues in log files.

    It can be used with tools like cat and grep.

    Here is a simple example.

        $ %[2]vcat foo.txt%[3]v
        aaa
        bbb
        ccc
        ddd
        EEE
        $ %[2]vcat foo.txt | %[4]v 'bbb' '^cc' '(?i)eee'%[3]v
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

       %[2]vhttps://github.com/google/re2/wiki/Syntax%[3]v

OPTIONS
    --              Treat all options as patterns. This allows
                    the search to look for --help or --version.
                    For example the following command would
                    highlight the "%[1]v--help%[3]v" string.
                        %[2]vcat -n %[4]v.go | %[4]v -- --help %[3]v

    -c MAP, --color-map MAP
                    Specify a color map.
                    The color map allows different colors to be
                    specified for each pattern. If this option
                    is not specified the same color is used
                    for all patterns.

                    See the COLOR MAP section for detailed
                    information.

    -h, --help      Print the program help and exit.

    -i FILE, --input FILE
                    Read from a file instead of stdin.

    -v, --verbose   Increase the level of verbosity.
                    It can be specified multiple times.
                    This is not normally useful. It is only
                    intended for debugging.

    -V, --version   Print the program version and exit.

COLOR MAP
    By default each pattern has the same highlight color but that can
    be customized by using a color map.

    The color map allows an ANSI terminal color code to be specified.
    It also has some short cut names to make things easier. The short
    cuts are appended using a plus (+) sign for each pattern.

    A simple example might make this clearer. If you have two patterns
    and want to color matches to the first one in green and matches to
    the second one in blue you would do this.

        $ cat logfile | %[4]v -c green+bold,blue+bold 'patern1' 'pattern2'

    If the color map (-c) is not specified all matches will be in red+bold
    which is the default.

    The raw components of the ANSI terminal code can be specified as digits.
    For example, "-c red+bold" is the same as "-c 31+1". They both ultimately
    translate to the binary sequence: "\033[31;1m".

    If the number of color map entries is less than the number of
    patterns, the last color map entry is extended to all of the
    remaining patterns.

    The table below presents the list of short cut options.

         #  Name         Value  Description
         1  blue            34  Foreground blue.
         2  blueB           44  Background blue.
         3  bold             1  Make the color bold.
         4  cyan            36  Foreground cyan.
         5  cyanB           46  Background cyan.
         6  faint            2  Make the color faint.
         7  gray      38;5;245  Foreground light gray.
         8  grayB     48;5;245  Background light gray.
         9  green           32  Foreground green.
        10  greenB          42  Background green.
        11  italic           3  Mke the text italic.
        12  magenta         35  Foreground magenta.
        13  magentaB        45  Backgroundmagenta.
        14  normal          22  Make the next normal.
        15  red             31  Foreground red.
        16  redB            41  Background red.
        17  reset            0  Reset.
        18  reverse          7  Reverse the fore/background colors.
        19  strike           9  Strike-through the text.
        20  underline        4  Underline the text.
        21  white           36  Foreground white.
        22  whiteB          46  Background white.
        23  yellow          33  Foreground yellow.
        24  yellowB         43  Background yellow.

    Here are some examples of what can be done with the short cuts.

        -c 'red+greenB+bold,blue+bold'
            First pattern is bold red on a bold green background.
            Second pattern is bold blue on the default background.

        -c 'red+greenB'
            First pattern is normal red on a normal green background.

EXAMPLES
    # Example 1: help
    $ %[2]v%[4]v -h%[3]v

    # Example 2: simple
    #            highlight "error:", "warning:" and "note:" in default red+bold
    #            ignore case
    $ %[2]vcat -n logfile | %[4]v '(?i)error:|warning:|note:'%[3]v

    # Example 3: more complex
    #            highlight "error:" in red
    #            highlight "warning:" in blue
    #            highlight strings in double quotes in green.
    $ %[2]vcat -n logfile | %[4]v -c red+bold,blue+bold,green+bold '(?i)error:' '(?i)warning:' '"[^"]*"'%[3]v

    # Example 4: change the background color
    $ %[2]vcat -n logfile | %[4]v -c greenB+red+bold '(?i)error:|warning:|note:'%[3]v

    # Example 5: use full blown RGB.
    $ %[2]vcat -n logfile | %[4]v -c '38;2;255;82;197;48;2;155;106;10' '(?i)error:|warning:|note:'%[3]v

    # Example 5: use full blown RGB and reverse it.
    $ %[2]vcat -n logfile | %[4]v -c 'reverse+38;2;255;82;197;48;2;155;106;10' '(?i)error:|warning:|note:'%[3]v

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
