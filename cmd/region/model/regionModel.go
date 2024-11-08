package model

import (
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
)

const RegionFileName = "region"

type RegionModel struct {
	URL string `json:"url"`
}

func (regionModel *RegionModel) RegionGet(fm fileutil.FileManager) (*RegionModel, error) {
	var model RegionModel
	if err := fm.LoadFromFile(RegionFileName, &model); err != nil {
		return nil, fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
	}
	return &model, nil
}
