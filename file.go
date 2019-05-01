package markdown

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileIndexer indexes Content for a file
type FileIndexer interface {
	Indexer
	Path() (string, bool)
	PathAndFileName() string
}

// NewContent creates a new instance of Content based on the given parameters
type NewContent func(frontmatter map[string]interface{}, haveFrontmatter bool, body []byte) (Content, error)

// fileStore satisfies the Store interface for reading/writing markdown
type fileStore struct {
	newContent NewContent
}

// NewFileStore creates a markdown store which reads/writes from the filesystem
func NewFileStore(newContent NewContent) Store {
	result := new(fileStore)
	result.newContent = newContent
	return result
}

func (s fileStore) GetContent(indexer Indexer) (Content, error) {
	fileName := indexer.(FileIndexer).PathAndFileName()
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, err
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	frontmatter := make(map[string]interface{})
	body, haveFrontmatter, err := ParseYAMLFrontMatter(data, &frontmatter)
	if err != nil {
		return nil, err
	}

	return s.newContent(frontmatter, haveFrontmatter, body)
}

func (s fileStore) HasContent(indexer Indexer) (bool, error) {
	fileName := indexer.(FileIndexer).PathAndFileName()
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func (s fileStore) WriteContent(indexer Indexer, content Content) error {
	path, createPath := indexer.(FileIndexer).Path()
	if createPath {
		_, err := s.CreateDirIfNotExist(path)
		if err != nil {
			return err
		}
	}

	fileName := indexer.(FileIndexer).PathAndFileName()
	file, createErr := os.Create(filepath.Join(path, fileName))
	if createErr != nil {
		return fmt.Errorf("Unable to create file %q: %v", fileName, createErr)
	}
	defer file.Close()

	frontMatter, fmErr := yaml.Marshal(content.Frontmatter())
	if fmErr != nil {
		return fmt.Errorf("Unable to marshal front matter %q: %v", fileName, fmErr)
	}

	file.WriteString("---\n")
	_, writeErr := file.Write(frontMatter)
	if writeErr != nil {
		return fmt.Errorf("Unable to write front matter %q: %v", fileName, writeErr)
	}

	_, writeErr = file.WriteString("---\n" + content.Body())
	if writeErr != nil {
		return fmt.Errorf("Unable to write content body %q: %v", fileName, writeErr)
	}

	return nil
}

// CreateDirIfNotExist creates a path if it does not exist. It is similar to mkdir -p in shell command,
// which also creates parent directory if not exists.
func (s fileStore) CreateDirIfNotExist(dir string) (bool, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		return true, err
	}
	return false, nil
}

func (s fileStore) DeleteContent(indexer Indexer) error {
	fileName := indexer.(FileIndexer).PathAndFileName()
	err := os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("Unable to delete file %q: %v", fileName, err)
	}
	return nil
}

func (s fileStore) Close() error {
	// not required
	return nil
}
