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

// AllowFrontMatterPropertyFunc allows custom handling of front matter properties for reading
type AllowFrontMatterPropertyFunc func(context.Context, ReaderIndexer, ...interface{}) properties.AllowAddFunc

// Reader defines common reader methods
type Reader interface {
	GetContent(context.Context, ReaderIndexer, AllowFrontMatterPropertyFunc, ...interface{}) (Content, error)
	HasContent(context.Context, ReaderIndexer, ...interface{}) (bool, error)
}

// PrepareToWriteFrontMatterFunc allows custom handling of front matter properties before writing
type PrepareToWriteFrontMatterFunc func(context.Context, WriterIndexer, Content, ...interface{}) properties.MapAssignFunc

// Writer defines common writer methods
type Writer interface {
	WriteContent(context.Context, WriterIndexer, Content, PrepareToWriteFrontMatterFunc, ...interface{}) error
	DeleteContent(context.Context, WriterIndexer, Content, ...interface{}) error
	DeletePrimaryKey(context.Context, ReaderIndexer, ...interface{}) error
}

// Store pulls together all the lifecyle, reader, and writer methods
type Store interface {
	Reader
	Writer
	io.Closer
}
