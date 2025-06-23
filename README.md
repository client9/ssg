# ssg
A toolkit for Static Site Generators


## The Document Pipeline -- What API to use?

THe SSG is a series of transformations on a document.

|-------------------|-----------|-------------|
| Transformer       | Output    | Input       |
|-------------------|-----------|-------------|
| Golang template   | io.Writer | []byte      |
| Golang gzip, tar  | io.Writer | io.Writer   |
| Goldmark Markdown | io.Writer | []byte      |
| JSON Encoder      | io.Writer | any         |
| tdewolff/minify   | io.Writer | io.Reader   |
| GoHTML            | []byte    | []byte      | 
| enescakir/emoji   | string    | string      |


Ok output is easy: Use io.Writer.

For input there are more choices.

If potentially very large (mega or gigebytes or more) or if it could come in vyer slowly (over a network), making a pure io.Write (or io.WriteCloser), like gzip and tar do.   This is not the use case here.

Many transformers take a []byte or string.  However if use io.Reader:

* We can use streaming transformers as is.
* We can do a nifty buffer swap

```go

src :=  bytes.Buffer{}
dest := bytes.Buffer{}

for _,tx := range pipeline {
	tx(dest, src)
	src, dest = dest, src
	dest.reset()
}
```

No need for channels or go routines.


