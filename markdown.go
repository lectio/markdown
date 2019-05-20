package markdown

import (
	"context"
	"github.com/lectio/properties"
	"io"
)

// A ReaderIndexer instance provides the indexing keys for something that needs to be read
type ReaderIndexer interface {
	ReaderPrimaryKey(context.Context, ...interface{}) string
}

// A WriterIndexer instance provides the indexing keys for something that needs to be written
type WriterIndexer interface {
	WriterPrimaryKey(context.Context, Content, ...interface{}) string
}

// Reader defines common reader methods
type Reader interface {
	GetContent(context.Context, ReaderIndexer, properties.AllowAddFunc, ...interface{}) (Content, error)
	HasContent(context.Context, ReaderIndexer, ...interface{}) (bool, error)
}

// Writer defines common writer methods
type Writer interface {
	WriteContent(context.Context, WriterIndexer, Content, properties.MapAssignFunc, ...interface{}) error
	DeleteContent(context.Context, WriterIndexer, Content, ...interface{}) error
	DeletePrimaryKey(context.Context, ReaderIndexer, ...interface{}) error
}

// Store pulls together all the lifecyle, reader, and writer methods
type Store interface {
	Reader
	Writer
	io.Closer
}
