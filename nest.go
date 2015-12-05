package magpie

type Nest func(name string) (Asset, error)

func NewNest(assets map[string]Asset) Nest {
	return func(name string) (a Asset, err error) {
		var ok bool
		if a, ok = assets[name]; ok {
			return
		}
		err = &UnknownAssetError{Name: name}
		return
	}
}
