# colorize
[![Releases](https://img.shields.io/github/colorize/jlinoff/csdiff.svg?style=flat)](https://github.com/jlinoff/colorize/releases)

Colorize piped output to hightlight specific strings

This simple little program is useful for highlighting specific items
that may help debugging in log files and other unstructured text data.

Here is how it could be used.

```bash
$ cat -n log | colorize "error:|note:|warning:"
```

Note that this different that using the `--color` option in tools like `grep`
because nothing is filtered.

The regular expression syntax is the regular golang regular
expression syntax that is described here:
https://github.com/google/re2/wiki/Syntax.

It is written go-1.10. It has been tested on Mac and Linux systems.
