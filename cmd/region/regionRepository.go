package region

import (
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
)

const regionFileName = "region"

func GetRegion(fm fileutil.FileManager) (*RegionModel, error) {
	var model RegionModel
	if err := fm.LoadFromFile(regionFileName, &model); err != nil {
		return nil, fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
	}
	return &model, nil
}

func SetRegion(fm fileutil.FileManager, url string) error {
	fm.DeleteFile(RegionFileName)
	fm.DeleteFile(AuthFileName)
	fm.DeleteFile(OrganizationFileName)
	fm.DeleteFile(ProjectFileName)

	region := RegionModel{Url: url}
	if err := fm.SaveToFile(RegionFileName, region); err != nil {
		return err
	}
	return nil
}
