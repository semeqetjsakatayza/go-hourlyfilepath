package hourlyfilepath

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var ErrNotReachUnusedFileName = errors.New("cannot reach unused file name")

func (b *BaseFolder) UnuseFilePath(t time.Time, fileNameBase, fileNameSuffix string, flag int, perm os.FileMode, maxAttempts int) (unuseFilePath string, err error) {
	folderPath, err := b.SetupFolder(t)
	if nil != err {
		err = fmt.Errorf("cannot setup folder for time (%v): %w", t, err)
		return
	}
	for attempt := 0; attempt <= maxAttempts; attempt++ {
		fileNameAttemptPart := strconv.FormatInt(int64(attempt), 10)
		fileName := fileNameBase + fileNameAttemptPart + fileNameSuffix
		unuseFilePath = filepath.Join(folderPath, fileName)
		if _, err = os.Stat(unuseFilePath); nil != err {
			if os.IsNotExist(err) {
				err = nil
				return
			}
			log.Printf("WARN: cannot stat path for searching unused file name [%s]: %v", unuseFilePath, err)
		}
	}
	err = ErrNotReachUnusedFileName
	return
}
