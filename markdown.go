package markdown

import (
	"context"
	"io"
)

// A ReaderIndexer instance provides the indexing keys for something that needs to be read
type ReaderIndexer interface {
	ReaderPrimaryKey(context.Context) string
}

// A WriterIndexer instance provides the indexing keys for something that needs to be written
type WriterIndexer interface {
	WriterPrimaryKey(context.Context, Content) string
}

// Reader defines common reader methods
type Reader interface {
	GetContent(context.Context, ReaderIndexer) (Content, error)
	HasContent(context.Context, ReaderIndexer) (bool, error)
}

// Writer defines common writer methods
type Writer interface {
	WriteContent(context.Context, WriterIndexer, Content) error
	DeleteContent(context.Context, ReaderIndexer) error
}

// Store pulls together all the lifecyle, reader, and writer methods
type Store interface {
	Reader
	Writer
	io.Closer
}
