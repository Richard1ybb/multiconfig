package multiconfig

type MultiLoader []Loader

// NewMultiLoader creates a loader that executes the loaders one by one in order
// and returns on the first error.
func NewMultiLoader(loader ...Loader) Loader {
	return MultiLoader(loader)
}

// Load loads the source into the config defined by struct s
func (m MultiLoader) Load(s interface{}) error {
	for _, loader := range m {
		if err := loader.Load(s); err != nil {
			return err
		}
	}

	return nil
}

// MustLoad loads the source into the struct, it panics if gets any error
func (m MultiLoader) MustLoad(s interface{}) {
	if err := m.Load(s); err != nil {
		panic(err)
	}
}
