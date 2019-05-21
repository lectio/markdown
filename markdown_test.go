package markdown

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"testing"

	"github.com/stretchr/testify/suite"
)

const validFrontMatter = `
---
description: test description
---
test body
`

const noFrontMatter = `test body without front matter`

const invalidFrontMatter1 = `
---
description: test description

test body
`

type readerIndexer struct {
	key string
	bpc BasePathConfigurator
}

func (i readerIndexer) ReaderPrimaryKey(context.Context, ...interface{}) string {
	return i.key
}

func (i readerIndexer) ReadFromPathAndFileName(ctx context.Context, options ...interface{}) (afero.Fs, string) {
	fileName := fmt.Sprintf("%s.md", i.key)
	return i.bpc.BaseFS(ctx), fileName
}

type MarkdownSuite struct {
	suite.Suite
	bpc BasePathConfigurator
	fs  *fileStore
}

func (suite *MarkdownSuite) SetupSuite() {
	suite.bpc = TheBasePathConfigurator
	suite.fs = NewFileStore(TheContentFactory, suite.bpc).(*fileStore)
}

func (suite *MarkdownSuite) WriterPrimaryKey(ctx context.Context, content Content, options ...interface{}) string {
	ic := content.(IdentifiedContent)
	return ic.PrimaryKey()
}

func (suite *MarkdownSuite) WriteToFileName(ctx context.Context, content Content, options ...interface{}) (afero.Fs, string) {
	fileName := fmt.Sprintf("%s.md", suite.WriterPrimaryKey(ctx, content))
	return suite.bpc.BaseFS(ctx), fileName
}

func (suite *MarkdownSuite) TearDownSuite() {
}

func (suite *MarkdownSuite) TestNoFrontMatter() {
	ctx := context.Background()
	bodyBytes, props, _, err := suite.fs.contentFactory.PropertiesFactory().MutableFromFrontMatter(ctx, []byte(noFrontMatter), nil)
	body := string(bodyBytes)
	suite.Nil(err, "Shouldn't have any errors")
	suite.Nil(props, "Should not have any front matter")
	suite.Equal(body, noFrontMatter)
}

func (suite *MarkdownSuite) TestValidFrontMatter() {
	ctx := context.Background()
	bodyBytes, props, _, err := suite.fs.contentFactory.PropertiesFactory().MutableFromFrontMatter(ctx, []byte(validFrontMatter), nil)
	suite.NotNil(props, "Should have front matter")

	content, _, err := TheContentFactory.NewIdentifiedContent(ctx, "test01", props, bodyBytes)

	suite.Nil(err, "Shouldn't have any errors")
	suite.Equal(content.BodyText(), "test body")

	fm := content.FrontMatter()
	descr, ok := fm.Named(ctx, "description")
	suite.True(ok, "description should be found")
	suite.Equal(descr.AnyValue(ctx), "test description")

	suite.fs.WriteContent(ctx, suite, content, nil)

	ri := &readerIndexer{key: content.PrimaryKey(), bpc: suite.bpc}
	readContent, rcErr := suite.fs.GetContent(ctx, ri, nil)
	suite.Nil(rcErr, "Should not have any read errors")

	fm = readContent.FrontMatter()
	descr, ok = fm.Named(ctx, "description")
	suite.True(ok, "description should be found")
	suite.Equal(descr.AnyValue(ctx), "test description")

	suite.fs.DeleteContent(ctx, suite, content)

	found, _ := suite.fs.HasContent(ctx, ri)
	suite.False(found, "file should not be found")
}

func (suite *MarkdownSuite) TestInvalidFrontMatter() {
	ctx := context.Background()
	_, _, _, err := suite.fs.contentFactory.PropertiesFactory().MutableFromFrontMatter(ctx, []byte(invalidFrontMatter1), nil)
	suite.NotNil(err, "Should have error")
	suite.EqualError(err, "Unexplained front matter parser error; insideFrontMatter: true, yamlStartIndex: 5, yamlEndIndex: 0")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(MarkdownSuite))
}
