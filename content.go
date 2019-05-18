package markdown

import (
	"context"
	"github.com/lectio/properties"
)

// Content encapsulates the components of a Markdown page
type Content interface {
	FrontMatter() properties.Properties
	HaveFrontMatter() bool
	Body() []byte
	BodyText() string
}

// IdentifiedContent provides a convenience Content for identified content
type IdentifiedContent interface {
	Content
	PrimaryKey() string
}

// DefaultContent is the default content object
type DefaultContent struct {
	id          string
	frontMatter properties.Properties
	body        []byte
}

func newDefaultContent(ctx context.Context, id string, frontMatter properties.Properties, body []byte, options ...interface{}) (*DefaultContent, bool, error) {
	return &DefaultContent{id: id, frontMatter: frontMatter, body: body}, true, nil
}

func (c *DefaultContent) PrimaryKey() string {
	return c.id
}

func (c *DefaultContent) FrontMatter() properties.Properties {
	return c.frontMatter
}

func (c *DefaultContent) HaveFrontMatter() bool {
	return c.frontMatter != nil
}

func (c *DefaultContent) Body() []byte {
	return c.body
}

func (c *DefaultContent) BodyText() string {
	return string(c.body)
}
