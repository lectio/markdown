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

func (i readerIndexer) ReaderPrimaryKey(context.Context) string {
	return i.key
}

func (i readerIndexer) ReadFromPathAndFileName(ctx context.Context) (afero.Fs, string) {
	fileName := fmt.Sprintf("%s.md", i.key)
	return i.bpc.BaseFS(ctx), fileName
}

type MarkdownSuite struct {
	suite.Suite
	bpc BasePathConfigurator
	fs  Store
}

func (suite *MarkdownSuite) SetupSuite() {
	suite.bpc = TheBasePathConfigurator
	suite.fs = NewFileStore(TheContentFactory, suite.bpc)
}

func (suite *MarkdownSuite) WriterPrimaryKey(ctx context.Context, content Content) string {
	ic := content.(IdentifiedContent)
	return ic.PrimaryKey()
}

func (suite *MarkdownSuite) WriteToFileName(ctx context.Context, content Content) (afero.Fs, string) {
	fileName := fmt.Sprintf("%s.md", suite.WriterPrimaryKey(ctx, content))
	return suite.bpc.BaseFS(ctx), fileName
}

func (suite *MarkdownSuite) TearDownSuite() {
}

func (suite *MarkdownSuite) TestNoFrontMatter() {
	frontmatter := make(map[string]interface{})
	bodyBytes, haveFrontMatter, err := ParseYAMLFrontMatter([]byte(noFrontMatter), frontmatter)
	body := string(bodyBytes)
	suite.Nil(err, "Shouldn't have any errors")
	suite.False(haveFrontMatter, "Should not have any front matter")
	suite.Equal(body, noFrontMatter)
}

func (suite *MarkdownSuite) TestValidFrontMatter() {
	ctx := context.Background()
	fmm := make(map[string]interface{})
	bodyBytes, haveFrontMatter, err := ParseYAMLFrontMatter([]byte(validFrontMatter), fmm)
	content, _, err := TheContentFactory.NewIdenfiedContent(ctx, "test01", fmm, haveFrontMatter, bodyBytes)

	suite.Nil(err, "Shouldn't have any errors")
	suite.True(haveFrontMatter, "Should not front matter")

	suite.Equal(content.BodyText(), "test body")

	fm := content.FrontMatter()
	descr, ok := fm.Named(ctx, "description")
	suite.True(ok, "description should be found")
	suite.Equal(descr.AnyValue(ctx), "test description")

	suite.fs.WriteContent(ctx, suite, content)

	ri := &readerIndexer{key: content.PrimaryKey(), bpc: suite.bpc}
	readContent, rcErr := suite.fs.GetContent(ctx, ri)
	suite.Nil(rcErr, "Should not have any read errors")

	fm = readContent.FrontMatter()
	descr, ok = fm.Named(ctx, "description")
	suite.True(ok, "description should be found")
	suite.Equal(descr.AnyValue(ctx), "test description")

	suite.fs.DeleteContent(ctx, ri)

	found, _ := suite.fs.HasContent(ctx, ri)
	suite.False(found, "file should not be found")
}

func (suite *MarkdownSuite) TestInvalidFrontMatter() {
	fmm := make(map[string]interface{})
	_, _, err := ParseYAMLFrontMatter([]byte(invalidFrontMatter1), fmm)

	suite.NotNil(err, "Should have error")
	suite.EqualError(err, "Unexplained front matter parser error; insideFrontMatter: true, yamlStartIndex: 5, yamlEndIndex: 0")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(MarkdownSuite))
}
