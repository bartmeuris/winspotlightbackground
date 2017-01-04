package main

import "fmt"

type imgError struct {
	Err     string
	ImgInfo *ImgFileInfo
}

func newImgError(err string, info *ImgFileInfo) imgError {
	return imgError{
		err,
		info,
	}
}

func (e imgError) Error() string {
	return fmt.Sprintf("%s: %s", e.ImgInfo, e.Err)
}
