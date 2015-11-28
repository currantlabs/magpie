package magpie

type UnknownAssetError struct {
	Name string
}

func (uae *UnknownAssetError) Error() string {
	return "magpie: unknown asset " + uae.Name
}
