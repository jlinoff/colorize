# colorize
[![Releases](https://img.shields.io/github/release/jlinoff/colorize.svg?style=flat)](https://github.com/jlinoff/colorize/releases)

Colorize unstructured text to hightlight specific strings for an ANSI terminal

This simple little program is useful for highlighting specific items
that may help debugging issues in log files.

It can be used with tools like cat and grep.

Here is how it could be used.

```bash
$ cat -n log | colorize "error:|note:|warning:"
```

The regular expression syntax is the regular golang regular
expression syntax that is described here:
https://github.com/google/re2/wiki/Syntax.

It is written go-1.10. It has been tested on Mac and Linux systems.

### Usage
In the simplest case, simply run the program with a list of patterns to
patch and pipe input to it as shown in the initial example.

The color of each pattern can specified independently using a color
map as this simple example shows where strings in enclosed in double
quotes are highlighted in green and the three keywords are highlighted
in red.

```bash
$ cat -n log | colorize -c "red+bold,green" "error:|note:|warning:" '\"[^\"]*\"'
```

There is more information about color maps in the *COLOR MAP* section.

### Example
```bash
$ bin/native/colorize  -i /tmp/make.log -c red+bold,green,red+bold 'error:.*$|note:.*$|warning:.*$' "'[^']*'" '^.*errors generated'
```
![screen shot 2018-04-12 at 10 55 12 am](https://user-images.githubusercontent.com/2991242/38695130-3c57fce8-3e40-11e8-9c6f-048f8e338df6.png)

### Options

| Short | Long | Description |
| ----- | ---- | ----------- |
| -- | -- | Everything after this is taken as literal. |
| -c COLOR_MAP | --color-map | Specify the pattern color map. |
| -h | --help | Print help message and exit. |
| -i FILE | --input FILE | Read a file. If not specified, read stdin. |
| -v | --verbose | Increase the level of verbosity. |
| -V | --version | Print the version and exit. |

### Color Map
The color map allows different colors to be specified for each
pattern. If this option is not specified the same color is used for
all patterns.

The color map allows an ANSI terminal color code to be specified.
It also has some short cut names to make things easier. The short
cuts are appended using a plus (+) sign for each pattern.

A simple example might make this clearer. If you have two patterns
and want to color matches to the first one in green and matches to
the second one in blue you would do this.

```bash
$ cat logfile | colorize -c green+bold,blue+bold 'patern1' 'pattern2'
```

If the color map (-c) is not specified all matches will be in red+bold
which is the default.

The raw components of the ANSI terminal code can be specified as digits.
For example, `-c red+bold` is the same as `-c 31+1`. They both ultimately
translate to the binary sequence: `\033[31;1m`.

If the number of color map entries is less than the number of
patterns, the last color map entry is extended to all of the
remaining patterns.

The table below presents the list of short cut options.

| Ref  |  Name         | Value     | Description          |
| ---: | ------------- | --------: | -------------------- |
|    1 |  blue         |       34  | Foreground blue. |
|    2 |  blueB        |       44  | Background blue. |
|    3 |  bold         |        1  | Make the color bold. |
|    4 |  cyan         |       36  | Foreground cyan. |
|    5 |  cyanB        |       46  | Background cyan. |
|    6 |  faint        |        2  | Make the color faint. |
|    7 |  gray         | 38;5;245  | Foreground light gray. |
|    8 |  grayB        | 48;5;245  | Background light gray. |
|    9 |  green        |       32  | Foreground green. |
|   10 |  greenB       |       42  | Background green. |
|   11 |  italic       |        3  | Mke the text italic. |
|   12 |  magenta      |       35  | Foreground magenta. |
|   13 |  magentaB     |       45  | Backgroundmagenta. |
|   14 |  normal       |       22  | Make the next normal. |
|   15 |  red          |       31  | Foreground red. |
|   16 |  redB         |       41  | Background red. |
|   17 |  reset        |        0  | Reset. |
|   18 |  reverse      |        7  | Reverse the fore/background colors. |
|   19 |  strike       |        9  | Strike-through the text. |
|   20 |  underline    |        4  | Underline the text. |
|   21 |  white        |       36  | Foreground white. |
|   22 |  whiteB       |       46  | Background white. |
|   23 |  yellow       |       33  | Foreground yellow. |
|   24 |  yellowB      |       43  | Background yellow. |

Here are some more examples of what can be done with the short cuts.

#### -c 'red+greenB+bold,blue+bold'
1. First pattern is bold red on a bold green background.
2. Second pattern is bold blue on the default background.

#### -c 'red+greenB'
1. First pattern is normal red on a normal green background.

### Anlyzing a colorized buffer in emacs
Add the following code to `~/.emacs`.
```lisp
;; colorize ansi term sequences
(require 'ansi-color)
(defun ansi-colorize ()
  "Colorize ansi escape sequences in the current buffer."
  (interactive)
  (ansi-color-apply-on-region (point-min) (point-max)))
```

Write the colorized log data to a separate log.
```bash
$ make 2>&1 /tmp/log
$ colorize -i /tmp/log -c red+bold,blue+bold,green+bol 'error:|warning:' 'note:' "'[^']*'" > /tmp/logc
```

Bring `/tmp/logc` up in emacs and run the `ansi-colorize` function.
