# pointer

[![Go Report Card](https://goreportcard.com/badge/github.com/borderzero/border0-go/lib/types/pointer)](https://goreportcard.com/report/github.com/borderzero/border0-go/lib/types/pointer)
[![Documentation](https://godoc.org/github.com/borderzero/border0-go/lib/types/pointer?status.svg)](https://godoc.org/github.com/borderzero/border0-go/lib/types/pointer)
[![license](https://img.shields.io/github/license/borderzero/border0-go.svg)](https://github.com/borderzero/border0-go/blob/master/LICENSE)

TL;DR smplifies dealing with pointers in Go...

In Go, taking the address of a literal (string, number, etc) is illegal because it has ambiguous semantics. This fact (and dealing with that ambiguity) often makes Go code longer than it needs to be.

For example, assume you have a struct with all pointer fields... To populate it, you'd have to first define all the values outside of the struct definition:

```
myString := "string"
myBool := true
myInt := 999

s := MyStruct{
    MyString: &myString,
    MyBool:   &myBool,
    MyInt:    &myInt,
}
```

Using this library, this can be simplified to:

```
s := MyStruct{
    MyString: pointer.To("string"),
    MyBool:   pointer.To(true),
    MyInt:    pointer.To(999),
}
```
