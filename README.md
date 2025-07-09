# front

> Parses frontmatter from a source text file.

## Usage

```
Usage: front [-o OUTPUT] [INPUT]

Parses any YAML frontmatter from INPUT file, converts it
to JSON and outputs it to OUTPUT. Adds the JSON property
"content", which contains the actual non-frontmatter
content from INPUT; and "path", which contains the path
to the input file ("-" if stdin)

If either INPUT or OUTPUT is omitted, or if either of
them is "-", reads from standard input and output.

Options:
  -o, --out PATH   Output file PATH (default "-")
```

## Example

```markdown
---
title: Hello
---
# Hello!
```

```shell
$ front hello.md
{"content":"# Hello!\n","meta":{"title":"Hello"},"path":"hello.md"}
```
