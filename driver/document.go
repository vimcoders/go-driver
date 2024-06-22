package driver

type Document interface {
	DocumentId() string
	DocumentName() string
}
