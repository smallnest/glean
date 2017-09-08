# Glean

*A go plugin framework that can reload variables and functions from plugins automatically.*


[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/smallnest/glean?status.png)](http://godoc.org/github.com/smallnest/glean)  [![travis](https://travis-ci.org/smallnest/glean.svg?branch=master)](https://travis-ci.org/smallnest/glean) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/glean)](https://goreportcard.com/report/github.com/smallnest/glean) 


## Installation

```sh
go get -u github.com/smallnest/glean
```

## Feature

- load symbol and you don't worry about errors
- load/reload exported variables and funtions from plugins
- watch plugins' changes and reload pointer of variables and function in applications

## Examples

see [Examples](https://github.com/smallnest/glean/tree/master/_example)