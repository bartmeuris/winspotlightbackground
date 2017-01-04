package main

import (
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
)

var prefLandscape, prefPortrait, copySmallImages, targetDeDup, targetValidateClean *bool
var ignoreWidth, ignoreHeight *int

const smallWidth = 150
const smallHeight = 150

// Validate if an image matches the specified rules
func (img *ImgFileInfo) Validate() error {
	if !*copySmallImages && ((img.ImgConf.Width <= smallWidth) || (img.ImgConf.Height <= smallHeight)) {
		return newImgError(fmt.Sprintf("Image %dx%d or smaller", smallWidth, smallHeight), img)
	}
	if (*ignoreHeight != 0) && (img.ImgConf.Height != *ignoreHeight) {
		return newImgError(fmt.Sprintf("height %d != expected %d", img.ImgConf.Height, *ignoreHeight), img)
	}
	if (*ignoreWidth != 0) && (img.ImgConf.Width != *ignoreWidth) {
		return newImgError(fmt.Sprintf("width %d != expected %d", img.ImgConf.Width, *ignoreWidth), img)
	}
	if *prefLandscape || *prefPortrait {
		if *prefLandscape && !*prefPortrait && (img.ImgConf.Width < img.ImgConf.Height) {
			return newImgError(fmt.Sprintf("%dx%d not landscape", img.ImgConf.Width, img.ImgConf.Height), img)
		}
		if !*prefLandscape && *prefPortrait && (img.ImgConf.Height > img.ImgConf.Width) {
			return newImgError(fmt.Sprintf("%dx%d not landscape", img.ImgConf.Width, img.ImgConf.Height), img)
		}
	}
	return nil
}

// RemoveDuplicateImages removes all duplicate images from a target directory. It does not recurse
func RemoveDuplicateImages(fileInfo []*ImgFileInfo, removefiles bool) {
	cleancount := 0
	for of := range fileInfo {
		for cf := range fileInfo {
			if of == cf {
				continue
			}
			if fileInfo[of].Equals(fileInfo[cf]) {
				if fileInfo[of].doClean || fileInfo[cf].doClean {
					// Already determined which-one to clean
					continue
				}
				if fileInfo[cf].ModTime.After(fileInfo[of].ModTime) {
					log.Printf("Remove:\n     %s\n  won:\n     %s", fileInfo[of], fileInfo[cf])
					fileInfo[of].doClean = true
					cleancount++
				} else {
					log.Printf("Remove:\n     %s\n  won:\n     %s", fileInfo[cf], fileInfo[of])
					fileInfo[cf].doClean = true
					cleancount++
				}
			}
		}
	}
	if cleancount == 0 {
		return
	}
	if removefiles {
		log.Printf("Removing files:")
	} else {
		log.Printf("Found duplicate files: (not removing)")
	}
	for _, fi := range fileInfo {
		if fi.doClean {
			log.Printf("    %s", fi)
			if removefiles {
				if err := os.Remove(fi.FileName); err != nil {
					log.Printf("Could not remove %s: %s", fi, err)
				}
			}
		}
	}
}

func main() {
	targetDir := flag.String("target", filepath.Join(os.Getenv("USERPROFILE"), "/Pictures/Spotlight"), "the target directory")
	prefLandscape = flag.Bool("landscape", true, "only copy landscape images (width > height)")
	prefPortrait = flag.Bool("portrait", false, "only copy portrait images (height > width)")
	ignoreWidth = flag.Int("width", 0, "only copy files with this width. 0 = ignore width")
	ignoreHeight = flag.Int("height", 0, "only copy files with this height. 0 = ignore height")
	copySmallImages = flag.Bool("copysmall", false, "copy small images (<= 150)")
	targetDeDup = flag.Bool("targetdedup", false, "remove all duplicate images in target directory")
	targetValidateClean = flag.Bool("targetvalidateremove", false, "validate all files in target directory and remove them if they don't match (DANGAROUS)")
	logFile := flag.String("logfile", "", "file to log to")
	flag.Parse()
	if *logFile != "" {
		if f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
			defer f.Close()
			log.SetOutput(f)
		}
	}
	// Construct the spotlight source directory
	spotdir := filepath.Join(os.Getenv("USERPROFILE"), "/AppData/Local/Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy/LocalState/Assets/")

	// Ensure the target directory exists
	os.MkdirAll(*targetDir, 0777)

	files := GetDirImages(spotdir)
	for _, f := range files {
		if err := f.Validate(); err != nil {
			log.Printf("Skipped: %s", err)
			continue
		}
		log.Printf("Copy %s to %s", f, *targetDir)
		if cerr := CopyFile(f.FileName, filepath.Join(*targetDir, filepath.Base(f.FileName)+"."+f.Format)); cerr != nil {
			log.Printf("ERROR: Could not copy %s to %s", f, filepath.Join(*targetDir, filepath.Base(f.FileName)+"."+f.Format))
		}
	}

	log.Printf("Analyzing target folder...")
	rmfiles := GetDirImages(*targetDir)
	RemoveDuplicateImages(rmfiles, *targetDeDup)
	for _, f := range rmfiles {
		if err := f.Validate(); err != nil {
			if *targetValidateClean {
				log.Printf("Image in targetdir not matching requested validation, removing: %s", f)
				if err := os.Remove(f.FileName); err != nil {
					log.Printf("Could not remove %s: %s", f, err)
				}
			} else {
				log.Printf("WARN: Image in targetdir not matching requested validation: %s", f)
			}
		}
	}
}
