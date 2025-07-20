# front

> Parses frontmatter from a source text file.

## Usage

```
Usage: front [OPTIONS] [INPUT]

Extract and convert YAML frontmatter from text file
INPUT. Outputs a new JSON object `output`:

  output = {
    "path": "-",   // INPUT file path
    "meta": null,  // frontmatter data
    "content": ""  // text content
  }

INPUT defaults to "-", stdout.

Options:
  -d, --delim DELIM   Set frontmatter delimiter
                      to DELIM (default "---")
  -h, --help          Show help
  -o, --out FILE      Output to FILE (default "-")
  -v, --version       Show version
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
