package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// ImgFileInfo contains various information about an image file
type ImgFileInfo struct {
	FileName string
	Format   string
	ImgConf  image.Config
	FileHash []byte
	FileSize int64
	ModTime  time.Time
	doClean  bool
}

func (im *ImgFileInfo) String() string {
	hstring := hex.EncodeToString(im.FileHash)
	if len(hstring) > 8 {
		hstring = hstring[:8]
	}
	sname := im.FileName
	if len(sname) > (20 + 3) {
		sname = "..." + sname[len(sname)-20:]
	}
	return fmt.Sprintf(
		"%s[%dx%d:%s]:<%s>#s:%d#%s",
		sname,
		im.ImgConf.Width, im.ImgConf.Height,
		im.Format,
		hstring,
		im.FileSize,
		im.ModTime,
	)
}

// LongString returns a more detailed/formatted string
func (im *ImgFileInfo) LongString() string {
	return fmt.Sprintf(
		"File      : %s\n"+
			"    Format    : %s\n"+
			"    ImageSize : %dx%d\n"+
			"    FileHash  : %s\n"+
			"    FileSize  : %d\n"+
			"    ModTime   : %s\n",
		im.FileName,
		im.Format,
		im.ImgConf.Width, im.ImgConf.Height,
		hex.EncodeToString(im.FileHash),
		im.FileSize,
		im.ModTime,
	)
}

// Equals checks if another ImgFileInfo instance is equal to this-one, ignoring modification time
func (im *ImgFileInfo) Equals(im2 *ImgFileInfo) bool {
	if im == nil && im2 == nil {
		return true
	}
	if im == nil || im2 == nil {
		return false
	}
	if im.FileName == im2.FileName {
		return true
	}
	if im.FileSize != im2.FileSize {
		return false
	}
	if im.Format != im2.Format {
		return false
	}
	if im.ImgConf.Width != im2.ImgConf.Width || im.ImgConf.Height != im2.ImgConf.Height {
		return false
	}
	if bytes.Compare(im.FileHash, im2.FileHash) != 0 {
		return false
	}
	return true
}

func imageInfo(fn string) *ImgFileInfo {
	ret := ImgFileInfo{}
	ret.doClean = false
	if reader, err := os.Open(fn); err == nil {
		defer reader.Close()
		fstat, statErr := reader.Stat()
		if statErr != nil {
			return nil
		}
		img, format, err := image.DecodeConfig(reader)
		if err != nil {
			return nil
		}

		// Now calculate the hash of the file
		if _, err := reader.Seek(0, os.SEEK_SET); err != nil {
			return nil
		}
		sha := sha256.New()
		if _, err := io.Copy(sha, reader); err != nil {
			return nil
		}
		ret.FileHash = sha.Sum(nil)
		ret.Format = format
		ret.FileName = fn
		ret.ImgConf = img
		ret.FileSize = fstat.Size()
		ret.ModTime = fstat.ModTime()
	} else {
		return nil
	}
	return &ret
}

// GetDirImages gets the information of all images in the specified dir. It does not recurse.
func GetDirImages(tgtdir string) []*ImgFileInfo {
	var fileinfo []*ImgFileInfo
	files, err := ioutil.ReadDir(tgtdir)
	if err != nil {
		panic("Could not open directory '" + tgtdir + "': " + err.Error())
	}

	for _, f := range files {
		fn := filepath.Join(tgtdir, f.Name())
		fi := imageInfo(fn)
		if fi != nil {
			fileinfo = append(fileinfo, fi)
		}
	}
	return fileinfo
}
