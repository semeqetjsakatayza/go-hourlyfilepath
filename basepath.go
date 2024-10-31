package hourlyfilepath

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type BaseFolder struct {
	baseFolderPath string
	folderMode     os.FileMode
	collectMinutes int
	timeZone       *time.Location
}

func NewBaseFolder(
	baseFolderPath string,
	folderMode os.FileMode,
	collectMinutes int,
	timeZone *time.Location) (b *BaseFolder, err error) {
	if baseFolderPath, err = filepath.Abs(baseFolderPath); nil != err {
		err = fmt.Errorf("cannot have absolute path for base folder: %w", err)
		return
	}
	if err = os.MkdirAll(baseFolderPath, folderMode); nil != err {
		err = fmt.Errorf("cannot setup base folder [%s]: %w", baseFolderPath, err)
		return
	}
	if (collectMinutes < 0) || (collectMinutes > 59) {
		collectMinutes = 0
	}
	if timeZone == nil {
		timeZone = time.UTC
	}
	b = &BaseFolder{
		baseFolderPath: baseFolderPath,
		folderMode:     folderMode,
		collectMinutes: collectMinutes,
		timeZone:       timeZone,
	}
	return
}
