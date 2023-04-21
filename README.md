# go-amf0

[![ci](https://github.com/yutopp/go-amf0/actions/workflows/ci.yml/badge.svg)](https://github.com/yutopp/go-amf0/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/yutopp/go-amf0/branch/master/graph/badge.svg?token=01TbR2Rwue)](https://codecov.io/gh/yutopp/go-amf0)
[![GoDoc](https://godoc.org/github.com/yutopp/go-amf0?status.svg)](http://godoc.org/github.com/yutopp/go-amf0)
[![Go Report Card](https://goreportcard.com/badge/github.com/yutopp/go-amf0)](https://goreportcard.com/report/github.com/yutopp/go-amf0)
[![license](https://img.shields.io/github/license/yutopp/go-amf0.svg)](https://github.com/yutopp/go-amf0/blob/master/LICENSE_1_0.txt)

AMF0 encoder/decoder library written in Go.

- [ ] Decoder
  - [x] Number
  - [x] Boolean
  - [x] String
  - [x] Object
  - [ ] Movieclip
  - [x] null
  - [ ] undefined
  - [ ] Reference
  - [x] ECMA Array
  - [x] Object End
  - [x] Strict Array
  - [x] Date
  - [ ] Long String
  - [ ] Unsupported
  - [ ] RecordSet
  - [ ] XMLDocument
  - [ ] Typed Object
- [ ] Encoder
  - [x] Number
  - [ ] Boolean
  - [x] String
  - [x] Object
  - [ ] Movieclip
  - [x] null
  - [ ] undefined
  - [ ] Reference
  - [x] ECMA Array
  - [x] Object End
  - [x] Strict Array
  - [x] Date
  - [ ] Long String
  - [ ] Unsupported
  - [ ] RecordSet
  - [ ] XMLDocument
  - [ ] Typed Object
- [ ] Documents
- [ ] Optimize

## Installation

```
go get github.com/yutopp/go-amf0
```

## Licence

[Boost Software License - Version 1.0](./LICENSE_1_0.txt)

## References

- [AMF0 specification](https://rtmp.veriskope.com/pdf/amf0-file-format-specification.pdf)
