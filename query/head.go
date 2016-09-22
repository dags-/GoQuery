package goquery

import (
	"time"
	"strings"
	"os"
	"path/filepath"
	"net/http"
	"image/png"
	"image"
	"image/draw"
	"math"
)

const separator = string(os.PathSeparator)
const steve = "http://assets.mojang.com/SkinTemplates/steve.png"

type HeadFetcher struct {
	targetDir string
	expire    time.Duration
	extension string
	scale     int
}

func NewHeadFetcher(root string, dir string, expire time.Duration, fileExtension string, scale int) HeadFetcher {
	root = strings.TrimRight(root, separator)
	dir = strings.Trim(dir, separator)
	target := filepath.Clean(root + separator + dir) + separator
	os.MkdirAll(target, os.ModeDir)
	return HeadFetcher{target, expire, fileExtension, int(math.Max(float64(scale), 1))}
}

func (fetcher *HeadFetcher) Path(uuid string) string {
	return getAbsolutePath(fetcher.targetDir, uuid, fetcher.extension)
}

func (fetcher *HeadFetcher) Fetch(uuid string) string {
	path := fetcher.Path(uuid)
	expire := minTime(fetcher.expire)
	if file, err := os.Stat(path); os.IsNotExist(err) || checkExpired(file, expire) {
		session := GetSession(uuid)
		return fetcher.FetchHead(session)
	}
	return path
}

func (fetcher *HeadFetcher) FetchHead(session Session) string {
	expire := minTime(fetcher.expire)
	path := getAbsolutePath(fetcher.targetDir, session.Id, fetcher.extension)
	if file, err := os.Stat(path); os.IsNotExist(err) || checkExpired(file, expire) {
		if session.Skin != "" {
			writeHead(session.Skin, path, fetcher.scale)
		} else {
			writeHead(steve, path, fetcher.scale)
		}
	}
	return path
}

func writeHead(url string, path string, scale int) {
	response, err := http.Get(url)
	if err != nil {
		return
	}

	img, err := png.Decode(response.Body)
	if err != nil {
		return
	}

	head := scaleImage(img, image.Rect(8, 8, 16, 16), scale)
	helmet := scaleImage(img, image.Rect(40, 8, 48, 16), scale)
	draw.Over.Draw(head, head.Bounds(), helmet, image.Pt(0, 0))

	if out, err := os.Create(path); err == nil {
		png.Encode(out, head)
		out.Close()
	}
}

func getAbsolutePath(dir string, id string, extension string) string {
	return dir + id + extension
}

func checkExpired(file os.FileInfo, expire time.Duration) bool {
	return time.Now().Sub(file.ModTime()) > expire
}

func minTime(expire time.Duration) time.Duration {
	if expire < time.Duration(1 * time.Minute) {
		return time.Duration(1 * time.Minute)
	}
	return expire
}
