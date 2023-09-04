package fabric

import "github.com/lordwestcott/gofabric"

type Fabric struct {
	Fabric *gofabric.App
}

func (f *Fabric) Init() error {
	fab, err := gofabric.InitApp()
	if err != nil {
		return err
	}

	f.Fabric = fab
	return nil
}
