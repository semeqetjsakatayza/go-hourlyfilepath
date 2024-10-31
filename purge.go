package hourlyfilepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func parseDecString(s string) (int, bool) {
	result := 0
	for _, c := range []byte(s) {
		v := int(c - '0')
		if v < 0 || v > 9 {
			return 0, false
		}
		result = result*10 + v
	}
	return result, true
}

func purgeFolderMM(errX []error, pathHH string, mm int) (errS []error) {
	errS = errX
	fsMMs, err1 := os.ReadDir(pathHH)
	if nil != err1 {
		errS = append(errS, fmt.Errorf("cannot read hour folder [%s]: %w", pathHH, err1))
		return
	}
	for _, entryMM := range fsMMs {
		if !entryMM.IsDir() {
			continue
		}
		nameMM := entryMM.Name()
		if len(nameMM) != 2 {
			continue
		}
		if valueMM, ok := parseDecString(nameMM); ok && (valueMM < mm) {
			pathMM := filepath.Join(pathHH, nameMM)
			if err1 := os.RemoveAll(pathMM); nil != err1 {
				errS = append(errS, fmt.Errorf("cannot remove minute collection folder [%s]: %w", pathMM, err1))
			}
			continue
		}
	}
	return
}

func purgeFolderHH(errX []error, pathYYYYMMDD string, hh, mm int) (errS []error) {
	errS = errX
	fsHHs, err1 := os.ReadDir(pathYYYYMMDD)
	if nil != err1 {
		errS = append(errS, fmt.Errorf("cannot read day folder [%s]: %w", pathYYYYMMDD, err1))
		return
	}
	for _, entryHH := range fsHHs {
		if !entryHH.IsDir() {
			continue
		}
		nameHH := entryHH.Name()
		if len(nameHH) != 2 {
			continue
		}
		valueHH, ok := parseDecString(nameHH)
		if !ok {
			continue
		}
		if valueHH > hh {
			continue
		}
		pathHH := filepath.Join(pathYYYYMMDD, nameHH)
		if valueHH < hh {
			if err1 = os.RemoveAll(pathHH); nil != err1 {
				errS = append(errS, fmt.Errorf("cannot remove hour folder [%s]: %w", pathHH, err1))
			}
			continue
		}
		if mm >= 0 {
			errS = purgeFolderMM(errS, pathHH, mm)
		}
	}
	return
}

func purgeFolderYYYYMMDD(baseFolderPath string, fsYYYYMMDDs []os.DirEntry, yyyymmdd, hh, mm int) (errS []error) {
	for _, entryYYYYMMDD := range fsYYYYMMDDs {
		if !entryYYYYMMDD.IsDir() {
			continue
		}
		nameYYYYMMDD := entryYYYYMMDD.Name()
		if len(nameYYYYMMDD) != 8 {
			continue
		}
		valueYYYYMMDD, ok := parseDecString(nameYYYYMMDD)
		if !ok {
			continue
		}
		if valueYYYYMMDD > yyyymmdd {
			continue
		}
		pathYYYYMMDD := filepath.Join(baseFolderPath, nameYYYYMMDD)
		if valueYYYYMMDD < yyyymmdd {
			if err1 := os.RemoveAll(pathYYYYMMDD); nil != err1 {
				errS = append(errS, fmt.Errorf("cannot remove day folder [%s]: %w", pathYYYYMMDD, err1))
			}
			continue
		}
		errS = purgeFolderHH(errS, pathYYYYMMDD, hh, mm)
	}
	return
}

// Purge removes all folders and its content with path name elder than the specified time.
// Only value in path name is considered, not the file modification time.
func (b *BaseFolder) Purge(expireAt time.Time) (err error) {
	if expireAt.Location() != b.timeZone {
		expireAt = expireAt.In(b.timeZone)
	}
	yyyymmdd := expireAt.Year()*10000 + int(expireAt.Month())*100 + expireAt.Day()
	hh := expireAt.Hour()
	mm := -1
	if b.collectMinutes != 0 {
		mm = expireAt.Minute()
		mm = int(mm/b.collectMinutes) * b.collectMinutes
	}
	fsYYYYMMDDs, err := os.ReadDir(b.baseFolderPath)
	if nil != err {
		err = fmt.Errorf("cannot read base folder [%s]: %w", b.baseFolderPath, err)
		return
	}
	if errS := purgeFolderYYYYMMDD(b.baseFolderPath, fsYYYYMMDDs, yyyymmdd, hh, mm); len(errS) > 0 {
		err = errors.Join(errS...)
	}
	return
}
