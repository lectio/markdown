package markdown

import (
	"fmt"
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

type testContent struct {
	id          string
	frontmatter map[string]interface{}
	body        string
}

func (c testContent) PrimaryKey() string {
	return c.id
}

func (c testContent) Path() (string, bool) {
	return "./", false
}

func (c testContent) PathAndFileName() string {
	return fmt.Sprintf("./%s.md", c.id)
}

func (c testContent) Frontmatter() map[string]interface{} {
	return c.frontmatter
}

func (c testContent) HasFrontmatter() bool {
	return c.frontmatter != nil
}

func (c testContent) Body() string {
	return c.body
}

type FrontMatterSuite struct {
	suite.Suite
	fs Store
}

func (suite *FrontMatterSuite) SetupSuite() {
	suite.fs = NewFileStore(
		func(frontmatter map[string]interface{}, haveFrontmatter bool, body []byte) (Content, error) {
			return &testContent{id: "test01", frontmatter: frontmatter, body: string(body)}, nil
		})
}

func (suite *FrontMatterSuite) TearDownSuite() {
}

func (suite *FrontMatterSuite) TestNoFrontMatter() {
	content := testContent{frontmatter: make(map[string]interface{})}
	body, haveFrontMatter, err := ParseYAMLFrontMatter([]byte(noFrontMatter), content.frontmatter)
	content.body = string(body)
	suite.Nil(err, "Shouldn't have any errors")
	suite.False(haveFrontMatter, "Should not have any front matter")
	suite.Equal(content.body, noFrontMatter)
}

func (suite *FrontMatterSuite) TestValidFrontMatter() {
	content := testContent{id: "test01", frontmatter: make(map[string]interface{})}
	body, haveFrontMatter, err := ParseYAMLFrontMatter([]byte(validFrontMatter), content.frontmatter)
	content.body = string(body)

	suite.Nil(err, "Shouldn't have any errors")
	suite.True(haveFrontMatter, "Should not front matter")

	suite.Equal(content.body, "test body")

	fm := content.Frontmatter()
	descr, ok := fm["description"]
	suite.True(ok, "description should be found")
	suite.Equal(descr, "test description")

	suite.fs.WriteContent(content, content)

	readContent, rcErr := suite.fs.GetContent(content)
	suite.Nil(rcErr, "Should not have any read errors")

	fm = readContent.Frontmatter()
	descr, ok = fm["description"]
	suite.True(ok, "description should be found")
	suite.Equal(descr, "test description")

	suite.fs.DeleteContent(readContent.(Indexer))
}

func (suite *FrontMatterSuite) TestInvalidFrontMatter() {
	content := testContent{frontmatter: make(map[string]interface{})}
	_, _, err := ParseYAMLFrontMatter([]byte(invalidFrontMatter1), content.frontmatter)

	suite.NotNil(err, "Should have error")
	suite.EqualError(err, "Unexplained front matter parser error; insideFrontMatter: true, yamlStartIndex: 5, yamlEndIndex: 0")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(FrontMatterSuite))
}
