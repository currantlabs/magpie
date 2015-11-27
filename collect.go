package magpie

func (m *Magpie) Collect() error {
	for _, c := range m.configs {
		err := collect(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func collect(c *config) error {
	fc, err := findFiles(c)
	if err != nil {
		return err
	}
	return writeNest(fc, c)

}
