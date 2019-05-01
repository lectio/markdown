package markdown

import (
	"io"
)

// Indexer create content addressable hashing or similar key generation for the given Content
type Indexer interface {
	PrimaryKey() string
}

// Content encapsulates the major components of a Markdown page
type Content interface {
	Frontmatter() map[string]interface{}
	Body() string
}

// Reader defines common reader methods
type Reader interface {
	GetContent(Indexer) (Content, error)
	HasContent(Indexer) (bool, error)
}

// Writer defines common writer methods
type Writer interface {
	WriteContent(Indexer, Content) error
	DeleteContent(Indexer) error
}

// Store pulls together all the lifecyle, reader, and writer methods
type Store interface {
	Reader
	Writer
	io.Closer
}
