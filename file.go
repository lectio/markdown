package markdown

import (
	"context"
	"fmt"
	"github.com/lectio/properties"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"os"
)

// BasePathConfigurator defines where to store results
type BasePathConfigurator interface {
	BasePath(ctx context.Context) string
	BaseFS(ctx context.Context) afero.Fs
	CreatePaths(ctx context.Context) (bool, os.FileMode)
	ComposePath(ctx context.Context, relativePath string) (afero.Fs, error)
}

// FileReaderIndexer is used by content readers to get paths and filenames
type FileReaderIndexer interface {
	ReaderIndexer
	ReadFromPathAndFileName(context.Context, ...interface{}) (afero.Fs, string)
}

// FileWriterIndexer is used by content writers to get paths and filenames
type FileWriterIndexer interface {
	WriterIndexer
	WriteToFileName(context.Context, Content, ...interface{}) (afero.Fs, string)
}

// fileStore satisfies the Store interface for reading/writing markdown
type fileStore struct {
	bpc            BasePathConfigurator
	contentFactory ContentFactory
}

// NewFileStore creates a markdown store which reads/writes from the filesystem
func NewFileStore(contentFactory ContentFactory, bpc BasePathConfigurator) Store {
	result := new(fileStore)
	result.bpc = bpc
	result.contentFactory = contentFactory
	return result
}

func (s fileStore) GetContent(ctx context.Context, indexer ReaderIndexer, checkAllow AllowFrontMatterPropertyFunc, options ...interface{}) (Content, error) {
	fs, fileName := indexer.(FileReaderIndexer).ReadFromPathAndFileName(ctx, options...)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, err
	}

	data, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return nil, err
	}

	var allow properties.AllowAddFunc
	if checkAllow != nil {
		allow = checkAllow(ctx, indexer, options...)
	}

	body, props, _, err := s.contentFactory.PropertiesFactory().MutableFromFrontMatter(ctx, data, allow, options...)
	if err != nil {
		return nil, err
	}

	content, _, err := s.contentFactory.NewContent(ctx, props, body, options...)
	return content, err
}

func (s fileStore) HasContent(ctx context.Context, indexer ReaderIndexer, options ...interface{}) (bool, error) {
	fs, fileName := indexer.(FileReaderIndexer).ReadFromPathAndFileName(ctx, options...)
	return afero.Exists(fs, fileName)
}

func (s fileStore) WriteContent(ctx context.Context, indexer WriterIndexer, content Content, prepare PrepareToWriteFrontMatterFunc, options ...interface{}) error {
	fs, fileName := indexer.(FileWriterIndexer).WriteToFileName(ctx, content, options...)
	file, createErr := fs.Create(fileName)
	if createErr != nil {
		return fmt.Errorf("Unable to create file %q: %v", fileName, createErr)
	}
	defer file.Close()

	if content.HaveFrontMatter() {
		fm := make(map[string]interface{})
		var assign properties.MapAssignFunc
		if prepare != nil {
			assign = prepare(ctx, indexer, content, options...)
		}
		content.FrontMatter().Map(ctx, fm, assign, options...)
		frontMatter, fmErr := yaml.Marshal(fm)
		if fmErr != nil {
			return fmt.Errorf("Unable to marshal front matter %q: %v", fileName, fmErr)
		}

		file.WriteString("---\n")
		_, writeErr := file.Write(frontMatter)
		if writeErr != nil {
			return fmt.Errorf("Unable to write front matter %q: %v", fileName, writeErr)
		}
		file.WriteString("---\n")
	}
	_, writeErr := file.Write(content.Body())
	if writeErr != nil {
		return fmt.Errorf("Unable to write content body %q: %v", fileName, writeErr)
	}

	return nil
}

func (s fileStore) DeleteContent(ctx context.Context, indexer WriterIndexer, content Content, options ...interface{}) error {
	fs, fileName := indexer.(FileWriterIndexer).WriteToFileName(ctx, content, options...)
	err := fs.Remove(fileName)
	if err != nil {
		return fmt.Errorf("Unable to delete file %q: %v", fileName, err)
	}
	return nil
}

func (s fileStore) DeletePrimaryKey(ctx context.Context, indexer ReaderIndexer, options ...interface{}) error {
	fs, fileName := indexer.(FileReaderIndexer).ReadFromPathAndFileName(ctx, options...)
	err := fs.Remove(fileName)
	if err != nil {
		return fmt.Errorf("Unable to delete file %q: %v", fileName, err)
	}
	return nil
}

func (s fileStore) Close() error {
	// not required
	return nil
}
