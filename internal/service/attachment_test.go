package service

import (
	"os"
	"path/filepath"
	"testing"
)

var htmlContent = []byte(`<!DOCTYPE html>
<html>
	<head>
		<title>important announcement</title>
	</head>
	<body>
		<h1>important announcement</h1>
		<span>Dr. Ivo Robotnik, 18:00</span>
		<marquee>i've come to make an announcement</marquee>
	</body>
</html>`)

// https://stackoverflow.com/a/2349470
var jpegContent = []byte(
	"\xff\xd8\xff\xe0\x00\x10\x4a\x46\x49\x46\x00\x01\x01\x01\x00\x48\x00\x48\x00" +
		"\x00\xff\xdb\x00\x43\x00\x03\x02\x02\x02\x02\x02\x03\x02\x02\x02\x03\x03" +
		"\x03\x03\x04\x06\x04\x04\x04\x04\x04\x08\x06\x06\x05\x06\x09\x08\x0a\x0a" +
		"\x09\x08\x09\x09\x0a\x0c\x0f\x0c\x0a\x0b\x0e\x0b\x09\x09\x0d\x11\x0d\x0e" +
		"\x0f\x10\x10\x11\x10\x0a\x0c\x12\x13\x12\x10\x13\x0f\x10\x10\x10\xff\xc9" +
		"\x00\x0b\x08\x00\x01\x00\x01\x01\x01\x11\x00\xff\xcc\x00\x06\x00\x10\x10" +
		"\x05\xff\xda\x00\x08\x01\x01\x00\x00\x3f\x00\xd2\xcf\x20\xff\xd9",
)

func TestFilenameFilter(t *testing.T) {
	// Valid input - typical file name
	ext, passed := filenameValid("myphoto.png")
	if ext != ".png" || !passed {
		t.Errorf("myphoto.png failed filenameValid - false positive")
	}

	// Valid input - Unicode in file name
	ext, passed = filenameValid("🦑.flac")
	if ext != ".flac" || !passed {
		t.Errorf("squid.flac failed filenameValid - false positive")
	}

	// Invalid input - Unicode in extension
	ext, passed = filenameValid("vacation.📷")
	if passed {
		t.Errorf("vacation.camera failed filenameValid - false negative")
	}

	// Invalid input - trailing newline
	ext, passed = filenameValid("gm_construct.bz2\n")
	if passed {
		t.Errorf("gm_construct\\n failed filenameValid - false negative")
	}

	// Invalid input - leading newline
	ext, passed = filenameValid("\nflinstones.mid")
	if passed {
		t.Errorf("\\nflinstones.mid failed filenameValid - false negative")
	}

	// Invalid input - middle newline
	ext, passed = filenameValid("grand\ndad.nes")
	if passed {
		t.Errorf("grand\\ndad.nes failed filenameValid - false negative")
	}

	// Invalid input - null character
	ext, passed = filenameValid("emp\x00ty.nul")
	if passed {
		t.Errorf("emp\\x00ty.nul failed filenameValid - false negative")
	}

	// Invalid input - dir traversal
	ext, passed = filenameValid("../jarjar.bnk")
	if passed {
		t.Errorf("../jarjar.bnk failed filenameValid - false negative")
	}

	// Invalid input - dotfile
	ext, passed = filenameValid(".htaccess")
	if passed {
		t.Errorf(".htaccess failed filenameValid - false negative")
	}

	// Invalid input - double dotfile
	ext, passed = filenameValid("..swag")
	if passed {
		t.Errorf("..swag failed filenameValid - false negative")
	}

	// Invalid input - leading tilde
	ext, passed = filenameValid("~jimmy.bin")
	if passed {
		t.Errorf("~jimmy.bin failed filenameValid - false negative")
	}
}

func TestAttachmentManager_StoreFile(t *testing.T) {
	tempDir := t.TempDir()

	a := AttachmentManager{}
	a.SetAcceptableMimes([]string{"image/jpeg", "application/pdf"})
	setDirErr := a.SetAttachmentDirectory(tempDir)
	if setDirErr != nil {
		t.Errorf("Failed to get temporary directory!")
	}

	// Ensure jpegContent has expected mime
	jpegMime := getContentType(jpegContent)
	if jpegMime != "image/jpeg" {
		t.Error("jpegMime has unexpected mime type:", jpegMime)
	}

	// Ensure htmlContent has expected mime
	htmlMime := getContentType(htmlContent)
	if htmlMime != "text/html" {
		t.Error("htmlMime has unexpected mime type:", htmlMime)
	}

	// Valid file - image/jpeg named photo.jpeg
	err := a.StoreFile("0001", "photo.jpeg", jpegContent)
	if err != nil {
		t.Error("Failed to store valid file photo.jpeg:", err)
	}
	_, err = os.Stat(filepath.Join(tempDir, "0001"))
	if err != nil {
		t.Error("Failed to read valid file photo.jpeg:", err)
	}

	// Valid file - image/jpeg named photo.jpg
	err = a.StoreFile("0002", "photo.jpg", jpegContent)
	if err != nil {
		t.Error("Failed to store valid file photo.jpg:", err)
	}
	_, err = os.Stat(filepath.Join(tempDir, "0002"))
	if err != nil {
		t.Error("Failed to read valid file photo.jpg:", err)
	}

	// Invalid file - image/jpeg named photo.png
	err = a.StoreFile("0003", "photo.png", jpegContent)
	if err == nil {
		t.Error("photo.png of type image/jpeg erroneously stored")
	}

	// Invalid file - text/html named not_suspicious.jpg
	err = a.StoreFile("0004", "not_suspicious.jpg", htmlContent)
	if err == nil {
		t.Error("not_suspicious.jpg of type text/html erroneously stored")
	}

	// Invalid file - text/html named index.html
	err = a.StoreFile("0005", "index.html", htmlContent)
	if err == nil {
		t.Error("index.html of type text/html erroneously stored")
	}
}
