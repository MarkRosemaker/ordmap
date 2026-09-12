---
logo:
    alt: 'ordmap logo: a gopher holding a map, surrounded by keys'
    source: logo.jpg
    width: 300
---

![Code Coverage](https://img.shields.io/badge/coverage-98.5%25-brightgreen)


`ordmap` is a Go package that provides a generic ordered map implementation, primarily designed for [JSON v2](https://pkg.go.dev/encoding/json/v2) marshalling and unmarshalling.

An ordered map maintains the order of keys based on insertion, allowing you to iterate over the map in the order in which entries were added. This can be particularly useful for applications where the order of elements is important, such as in JSON serialization or when maintaining the sequence of operations.
