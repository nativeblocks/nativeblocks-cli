package regionModule

import (
	"errors"

	"github.com/nativeblocks/cli/library/fileutil"
)

const regionFileName = "region"

func GetRegion(fm fileutil.FileManager) (*RegionModel, error) {
	var model RegionModel
	if err := fm.LoadFromFile(regionFileName, &model); err != nil {
		return nil, errors.New("region not set. Please set region first using 'nativeblocks region set <url>'")
	}
	return &model, nil
}

func SetRegion(fm fileutil.FileManager, url string) error {
	_ = fm.DeleteFile(RegionFileName)
	_ = fm.DeleteFile(AuthFileName)
	_ = fm.DeleteFile(OrganizationFileName)
	_ = fm.DeleteFile(ProjectFileName)

	region := RegionModel{Url: url}
	if err := fm.SaveToFile(RegionFileName, region); err != nil {
		return err
	}
	return nil
}
