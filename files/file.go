package files

import (
	"crypto/md5"  //nolint:gosec
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/gabriel-vasile/mimetype"
)

// FileInfo describes a file.
type FileInfo struct {
	*Listing
	Fs        afero.Fs          `json:"-"`
	Path      string            `json:"path"`
	Name      string            `json:"name"`
	Size      int64             `json:"size"`
	Extension string            `json:"extension"`
	ModTime   time.Time         `json:"modified"`
	Mode      os.FileMode       `json:"mode"`
	IsDir     bool              `json:"isDir"`
	Type      string            `json:"type"`
	Subtitles []string          `json:"subtitles,omitempty"`
	Content   string            `json:"content,omitempty"`
	Checksums map[string]string `json:"checksums,omitempty"`
	Key       []byte            `json:"key"`
}

// FileOptions are the options when getting a file info.
type FileOptions struct {
	Fs      afero.Fs
	Path    string
	Modify  bool
	Expand  bool
	Checker rules.Checker
}

func generateFileKey(fs afero.Fs, path string) []byte {
	reader, err := fs.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	hash := sha256.New()

	_, err = io.Copy(hash, reader)
	if err != nil {
		log.Fatal(err)
	}

	return hash.Sum(nil)
}

// NewFileInfo creates a File object from a path and a given user. This File
// object will be automatically filled depending on if it is a directory
// or a file. If it's a video file, it will also detect any subtitles.
func NewFileInfo(opts FileOptions) (*FileInfo, error) {
	if !opts.Checker.Check(opts.Path) {
		return nil, os.ErrPermission
	}

	info, err := opts.Fs.Stat(opts.Path)
	if err != nil {
		return nil, err
	}

	var key []byte

	if !info.IsDir() {
		key = generateFileKey(opts.Fs, opts.Path)
	}

	file := &FileInfo{
		Fs:        opts.Fs,
		Path:      opts.Path,
		Name:      info.Name(),
		ModTime:   info.ModTime(),
		Mode:      info.Mode(),
		IsDir:     info.IsDir(),
		Size:      info.Size(),
		Extension: filepath.Ext(info.Name()),
		Key:       key,
	}

	if opts.Expand {
		if file.IsDir {
			if err := file.readListing(opts.Checker); err != nil { //nolint:shadow
				return nil, err
			}
			return file, nil
		}

		err = file.detectType(opts.Modify, true)
		if err != nil {
			return nil, err
		}
	}

	return file, err
}

// Checksum checksums a given File for a given User, using a specific
// algorithm. The checksums data is saved on File object.
func (i *FileInfo) Checksum(algo string) error {
	if i.IsDir {
		return errors.ErrIsDirectory
	}

	if i.Checksums == nil {
		i.Checksums = map[string]string{}
	}

	reader, err := i.Fs.Open(i.Path)
	if err != nil {
		return err
	}
	defer reader.Close()

	var h hash.Hash

	//nolint:gosec
	switch algo {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return errors.ErrInvalidOption
	}

	_, err = io.Copy(h, reader)
	if err != nil {
		return err
	}

	i.Checksums[algo] = hex.EncodeToString(h.Sum(nil))
	return nil
}

//nolint:goconst
//TODO: use constants
func (i *FileInfo) detectType(modify, saveContent bool) error {
	// failing to detect the type should not return error.
	// imagine the situation where a file in a dir with thousands
	// of files couldn't be opened: we'd have immediately
	// a 500 even though it doesn't matter. So we just log it.
	reader, err := i.Fs.Open(i.Path)

	mime, err := mimetype.DetectReader(reader)

	if err != nil {
		i.Type = "blob"
		return nil
	}

	defer reader.Close()

	var reg = regexp.MustCompile(`video|image|audio`)
	mimetype := mime.String()
	isMatch := reg.MatchString(mimetype)

	if isMatch {
		i.Type = strings.Split(mimetype, "/")[0]
		i.detectSubtitles()
		return nil
	}

	i.Type = strings.Replace(mime.Extension(), ".", "", -1)
	if !modify {
		i.Type = "textImmutable"
	}

	if saveContent {
		afs := &afero.Afero{Fs: i.Fs}
		content, err := afs.ReadFile(i.Path)
		if err != nil {
			return err
		}
		i.Content = string(content)
	}
	return nil
}

func (i *FileInfo) detectSubtitles() {
	if i.Type != "video" {
		return
	}

	i.Subtitles = []string{}
	ext := filepath.Ext(i.Path)

	// TODO: detect multiple languages. Base.Lang.vtt

	fPath := strings.TrimSuffix(i.Path, ext) + ".vtt"
	if _, err := i.Fs.Stat(fPath); err == nil {
		i.Subtitles = append(i.Subtitles, fPath)
	}
}

func (i *FileInfo) readListing(checker rules.Checker) error {
	afs := &afero.Afero{Fs: i.Fs}
	dir, err := afs.ReadDir(i.Path)
	if err != nil {
		return err
	}

	listing := &Listing{
		Items:    []*FileInfo{},
		NumDirs:  0,
		NumFiles: 0,
	}

	for _, f := range dir {
		name := f.Name()
		fPath := path.Join(i.Path, name)

		if !checker.Check(fPath) {
			continue
		}

		if strings.HasPrefix(f.Mode().String(), "L") {
			// It's a symbolic link. We try to follow it. If it doesn't work,
			// we stay with the link information instead if the target's.
			info, err := i.Fs.Stat(fPath)
			if err == nil {
				f = info
			}
		}

		var key []byte

		if !f.IsDir() {
			key = generateFileKey(i.Fs, fPath)
		}

		file := &FileInfo{
			Fs:        i.Fs,
			Name:      name,
			Size:      f.Size(),
			ModTime:   f.ModTime(),
			Mode:      f.Mode(),
			IsDir:     f.IsDir(),
			Extension: filepath.Ext(name),
			Path:      fPath,
			Key:       key,
		}

		if file.IsDir {
			listing.NumDirs++
		} else {
			listing.NumFiles++

			err := file.detectType(true, false)
			if err != nil {
				return err
			}
		}

		listing.Items = append(listing.Items, file)
	}

	i.Listing = listing
	return nil
}
