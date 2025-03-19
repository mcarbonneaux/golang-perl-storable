# python_perl_storable

# NAME

golang-perl-storable — Packer/unpacker from/to perl-storable binary format.

# VERSION

0.0.0-prealpha — at the moment only reading of format bytes is ready.

# DESCRIPTION

```go

import 'github.com/darviarush/golang-perl-storable/encoding/storable'

var data int

err := storable.Unmarshal(storable, &data)

var anydata any

err := storable.Unmarshal(storable, &anydata)

anydata.(int)

```

# SYNOPSIS

The perl language has its own binary data format for packing any structures: hashes, lists, objects, regular expressions, scalars, file descriptors, links, globs, etc. It is implemented by the module https://metacpan.org/pod/Storable.

This format is quite popular and data from various perl projects packed into a binary string is stored in external storage: mysql, memcached, tarantool, etc.

This go module is designed to unpack data retrieved from such stores into `go` structures and to pack `go` data to put it into the store.  

# FUNCTIONS

## Unmarshal

### ARGUMENTS

- storable - binary string
- classes - dictionary with classes. Optional parameter
- iconv - function for converting strings not in utf8. Optional parameter

## Marshal

### ARGUMENTS

- data - python data: string, number, dictionary, list, object, etc.
- magic - boolean value. Optional. Adds magic number 'pst0' to output

### RETURNS

Binary string with data in Perl Storable format

# SCRIPT

```sh
# Freeze-defrost:
$ echo '[123, "Let it be!"]' | pypls freeze | pypls thaw

# Transfer frozen data in binary form:
$ echo '[123, "Let it be!"]' | pypls freeze -b | pypls thaw -b

# Pass the code in the parameter:
$ pypls freeze --data '[123, "Let it be!"]' | pypls thaw

# Add a magic number and fade the frozen line:
$ pypls freeze -m -s --data '[123, "Let it be!"]' | pypls thaw

# Recode strings (bytes will remain as is):
$ pypls freeze --data '[123, "Let it be!"]' -i cp1251 | pypls thaw -i cp1251

```

# INSTALL

```sh
$ go get github.com/darviarush/golang-perl-storable
```

# AUTHOR

Yaroslav O. Kosmina <darviarush@mail.ru>

# LICENSE

MIT License

Copyright (c) 2022 Yaroslav O. Kosmina

