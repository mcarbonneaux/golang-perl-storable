# python_perl_storable

# NAME

golang-perl-storable — Паковщик/распаковщик из/в бинарного формата perl-storable.

# VERSION

0.0.0-prealpha — на данный момент готово только считывание байт формата.

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

В языке perl есть свой формат бинарных данных для упаковки любых структур: хешей, списков, объектов, регулярок, скаляров, файловых дескрипторов, ссылок, глобов и т.п. Он реализуется модулем https://metacpan.org/pod/Storable.

Данный формат довольно популярен и запакованные в бинарную строку данные различных проектов на perl хранятся во внешних хранилищах: mysql, memcached, tarantool и т.д.

Данный go-модуль предназначен для распаковки данных, полученных из таких хранилищ, в структуры `go` и для упаковки данных `go`, чтобы поместить их в хранилище.  

# FUNCTIONS

## Unmarshal

### ARGUMENTS

- storable - бинарная строка
- classes - словарь с классами. Необязательный параметр
- iconv - функция для конвертации строк не в utf8. Необязательный параметр

## Marshal

### ARGUMENTS

- data - данные питона: строка, число, словарь, список, объект и т.д.
- magic - булево значение. Необязательно. Добавляет к выводу магическое число 'pst0'

### RETURNS

Бинарная строка с данными в формате Perl Storable

# SCRIPT

```sh
# Заморозить-раморозить:
$ echo '[123, "Хай!"]' | pypls freeze | pypls thaw

# Передавать замороженные данные в бинарном виде:
$ echo '[123, "Хай!"]' | pypls freeze -b | pypls thaw -b

# Передавать код в параметре:
$ pypls freeze --data '[123, "Хай!"]' | pypls thaw

# Добавить магическое число и обесцветить замероженную строку:
$ pypls freeze -m -s --data '[123, "Хай!"]' | pypls thaw

# Перекодировать строки (bytes останутся как есть):
$ pypls freeze --data '[123, "Хай!"]' -i cp1251 | pypls thaw -i cp1251

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

