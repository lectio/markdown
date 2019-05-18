package markdown

import (
	"context"
	"fmt"
	"github.com/lectio/properties"
	"github.com/spf13/afero"
	"os"
)

var (
	// TheContentFactory is primary content factory for common use cases
	TheContentFactory = &DefaultContentFactory{PropsFactory: properties.ThePropertiesFactory}

	// TheBasePathConfigurator can be used for the current directory
	TheBasePathConfigurator = newDefaultBasePathConfigurator("")
)

// NewContentFunc creates a new instance of Content based on the given parameters
type NewContentFunc func(ctx context.Context, frontmatter map[string]interface{}, haveFrontmatter bool, body []byte, options ...interface{}) (Content, bool, error)

// ContentFactory creates content instances
type ContentFactory interface {
	PropertiesFactory() properties.Factory
	NewContent(ctx context.Context, frontmatter map[string]interface{}, haveFrontmatter bool, body []byte, options ...interface{}) (Content, bool, error)
	NewIdenfiedContent(ctx context.Context, id string, frontmatter map[string]interface{}, haveFrontmatter bool, body []byte, options ...interface{}) (IdentifiedContent, bool, error)
}

// DefaultContentFactory is the default instance
type DefaultContentFactory struct {
	PropsFactory properties.Factory
}

// PropertiesFactory returns the factory used to create Properties instances
func (f *DefaultContentFactory) PropertiesFactory() properties.Factory {
	return f.PropsFactory
}

// NewContent takes a front matter plus body text and creates a Content instance
func (f *DefaultContentFactory) NewContent(ctx context.Context, frontMatter map[string]interface{}, havefrontMatter bool, body []byte, options ...interface{}) (Content, bool, error) {
	if havefrontMatter {
		props, _, err := f.PropsFactory.MutableFromStringMap(ctx, frontMatter, options...)
		if err != nil {
			return nil, false, err
		}
		return newDefaultContent(ctx, "", props, body, options...)
	}
	return newDefaultContent(ctx, "", nil, body, options...)
}

// NewIdenfiedContent takes a front matter plus body text and creates a Content instance with an identity attached
func (f *DefaultContentFactory) NewIdenfiedContent(ctx context.Context, id string, frontMatter map[string]interface{}, havefrontMatter bool, body []byte, options ...interface{}) (IdentifiedContent, bool, error) {
	if havefrontMatter {
		props, _, err := f.PropsFactory.MutableFromStringMap(ctx, frontMatter, options...)
		if err != nil {
			return nil, false, err
		}
		return newDefaultContent(ctx, id, props, body, options...)
	}
	return newDefaultContent(ctx, id, nil, body, options...)
}

// DefaultBasePathConfigurator is the default instance
type DefaultBasePathConfigurator struct {
	rootFS   afero.Fs
	basePath string
	baseFS   afero.Fs
}

func newDefaultBasePathConfigurator(basePath string) *DefaultBasePathConfigurator {
	result := &DefaultBasePathConfigurator{
		rootFS:   afero.NewOsFs(),
		basePath: basePath,
	}
	result.baseFS = result.rootFS
	if baseFS, err := result.newAferoBasePathFs(context.Background(), result.rootFS, basePath); err == nil {
		result.baseFS = baseFS
	} else {
		fmt.Printf("Unable to create DefaultBasePathConfigurator.baseFS in %q, defaulting to rootFS: +v", err)
	}
	return result
}

// BasePath returns the basePath of the configurator
func (bpc *DefaultBasePathConfigurator) BasePath(ctx context.Context) string {
	return bpc.basePath
}

// BaseFS returns the afero Fs of the configurator's base path
func (bpc *DefaultBasePathConfigurator) BaseFS(ctx context.Context) afero.Fs {
	return bpc.baseFS
}

// CreatePaths returns true if paths should be created when afero file systems are requested
func (bpc *DefaultBasePathConfigurator) CreatePaths(ctx context.Context) (bool, os.FileMode) {
	return true, os.FileMode(0755)
}

// ComposePath creates a new afero file system off baseFS given the relative path
func (bpc *DefaultBasePathConfigurator) ComposePath(ctx context.Context, relativePath string) (afero.Fs, error) {
	return bpc.newAferoBasePathFs(ctx, bpc.BaseFS(ctx), relativePath)
}

func (bpc *DefaultBasePathConfigurator) newAferoBasePathFs(ctx context.Context, parent afero.Fs, relativePath string) (afero.Fs, error) {
	create, createMode := bpc.CreatePaths(ctx)
	if create {
		err := parent.MkdirAll(relativePath, createMode)
		if err != nil {
			return nil, fmt.Errorf("Unable to create path %q in newBasePathFs: %v", relativePath, err.Error())
		}
	}
	return afero.NewBasePathFs(parent, relativePath), nil
}
