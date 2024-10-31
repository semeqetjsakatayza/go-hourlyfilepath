package hourlyfilepath

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func zeroPaddingDecString(n, size int) string {
	var buf [8]byte
	for idx := (size - 1); idx >= 0; idx-- {
		buf[idx] = byte(n%10) + '0'
		n /= 10
	}
	return string(buf[:size])
}

func (b *BaseFolder) SetupFolder(t time.Time) (folderPath string, err error) {
	if t.Location() != b.timeZone {
		t = t.In(b.timeZone)
	}
	yyyymmdd := t.Year()*10000 + int(t.Month())*100 + t.Day()
	folderNameYYYYMMDD := zeroPaddingDecString(yyyymmdd, 8)
	hh := t.Hour()
	folderNameHH := zeroPaddingDecString(hh, 2)
	if b.collectMinutes == 0 {
		folderPath = filepath.Join(b.baseFolderPath, folderNameYYYYMMDD, folderNameHH)
	} else {
		mm := t.Minute()
		mm = int(mm/b.collectMinutes) * b.collectMinutes
		folderNameMM := zeroPaddingDecString(mm, 2)
		folderPath = filepath.Join(b.baseFolderPath, folderNameYYYYMMDD, folderNameHH, folderNameMM)
	}
	if err = os.MkdirAll(folderPath, b.folderMode); nil != err {
		err = fmt.Errorf("cannot setup folder [%s]: %w", folderPath, err)
		return
	}
	return
}
