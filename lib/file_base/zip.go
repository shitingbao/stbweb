package base

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

//ZipParse 解析zip文件
//baseURL为zip路径，flagURL目标文件解析路径（就是解析在哪个路径下）
func ZipParse(baseURL, flagURL string) error {
	r, err := zip.OpenReader(baseURL)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		ph := path.Join(flagURL, f.Name)
		os.MkdirAll(path.Dir(ph), os.ModePerm)
		fInfo, err := os.Create(ph)
		if err != nil {
			return err
		}

		if _, err = io.CopyN(fInfo, rc, f.FileInfo().Size()); err != nil {
			return err
		}
		rc.Close()
	}
	return nil
}

// CreateZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
//eg:text.zip   ,   []{"aa.txt","file/bb.txt"},zip中包含两级，第一级有aa.txt和file目录,file目录下有bb.txt
//isDir代表是否带目录结构，true代表包含
func CreateZipFiles(filename string, files []string, isDir bool) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file, isDir); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string, isDir bool) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	//这里的名称就是指的是zip内部的路径，如果取完整路径（file/aa.txt）就会带有目录，直接获取文件名，不带目录（aa.txt），内部就不包含目录文件夹
	if isDir {
		header.Name = filename
	} else {
		header.Name = path.Base(filename)
	}

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
