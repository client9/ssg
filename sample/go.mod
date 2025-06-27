module github.com/client9/ssg/sample

go 1.24.4

require (
	github.com/client9/ssg v0.0.0-local
	github.com/client9/ssg/htmlcontent v0.0.0-local
	github.com/yosssi/gohtml v0.0.0-20201013000340-ee4748c638f4
)

require golang.org/x/net v0.41.0 // indirect

replace github.com/client9/ssg v0.0.0-local => ../

replace github.com/client9/ssg/htmlcontent v0.0.0-local => ../htmlcontent
