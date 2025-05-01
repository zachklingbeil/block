package universe

import "github.com/zachklingbeil/factory"

type Universe struct {
	Factory *factory.Factory
	Map     map[string]*any
}
