package region

import (
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
)

const regionFileName = "region"

type RegionModel struct {
	Url string `json:"url"`
}

func (regionModel *RegionModel) GetRegion(fm fileutil.FileManager) (*RegionModel, error) {
	var model RegionModel
	if err := fm.LoadFromFile(regionFileName, &model); err != nil {
		return nil, fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
	}
	return &model, nil
}
