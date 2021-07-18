# unicode symbol dumper

This program dumps unicode symbol.

UTF16, UTF-8, UTF-32 supported.

# Example

```bash
$ printf "🛀🐧" | usd
+-----------+------------+---------+---------------------+
| CHARACTER | CODE POINT |  NAME   |         HEX         |
+-----------+------------+---------+---------------------+
| 🛀        | U+1F6C0    | BATH    | 0xF0 0x9F 0x9B 0x80 |
| 🐧        | U+1F427    | PENGUIN | 0xF0 0x9F 0x90 0xA7 |
+-----------+------------+---------+---------------------+
```

# Usage

```bash
$ usd -help
Usage of usd:
  usd <subcommand>
  usd -help
Sub commands:
  utf8
        dump UTF-8
  utf16
        dump UTF-16
  utf32
        dump UTF-32
Options:
  -help
       show help
  -fileType value
        output file type. default is None (value: CSV|TSV|None)
  -noHeader
        no header
  -version
        show version
$ usd utf8 -help
Usage of utf8:
  utf8 [option]
Options:
  -help
        show help
$ usd utf16 -help
Usage of utf16:
  utf16 [option]
Options:
  -help
        show help
  -endian endian
        UTF16 endian. default is 'Big' (value: Big|Little)
$ usd utf32 -help
Usage of utf32:
  utf32 [option]
Options:
  -help
        show help
  -endian endian
        UTF32 endian. default is 'Big' (value: Big|Little)
```
