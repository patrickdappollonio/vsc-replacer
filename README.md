# `vsc-replacer` a poor version of VSCode Regexp Search&Replace

`vsc-replacer` is a poor version of VSCode Regexp Search&Replace, which is used to replace the string in the specified file with the specified regular expression.

Works pretty much the same, for example:

```bash
./vsc-replacer --regex "(Hello) (World)" --replacement "$1 Everyone" --dir "/path/to/files" --dry-run
```

This will replace all the `Hello World` in the files under `/path/to/files` with `Hello Everyone`.

**Differences:**

* `vsc-replacer` is a command line tool, while VSCode is a GUI tool.
* VSCode uses Javascript regexp, while `vsc-replacer` uses Golang regexp, consider it when you write your regexp.
