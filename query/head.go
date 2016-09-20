package goquery

import (
	"time"
	"strings"
	"os"
	"net/http"
	"image"
	"image/draw"
	"image/png"
	"path/filepath"
)

const separator = string(os.PathSeparator)
const steve = "http://assets.mojang.com/SkinTemplates/steve.png"

type HeadFetcher struct {
	targetDir string
	expire    time.Duration
	extension string
}

func NewHeadFetcher(root string, dir string, expire time.Duration, fileExtension string) HeadFetcher {
	root = strings.TrimRight(root, separator)
	dir = strings.Trim(dir, separator)
	target := filepath.Clean(root + separator + dir) + separator
	os.MkdirAll(target, os.ModeDir)
	return HeadFetcher{target, expire, fileExtension}
}

func (fetcher *HeadFetcher) Fetch(uuid string) string {
	expire := minTime(fetcher.expire)
	path := getAbsolutePath(fetcher.targetDir, uuid, fetcher.extension)
	if file, err := os.Stat(path); os.IsNotExist(err) || checkExpired(file, expire) {
		session := GetSession(uuid)
		return fetcher.FetchHead(session)
	}
	return path
}

func (fetcher *HeadFetcher) FetchHeads(sessions ...Session) {
	for i := range sessions {
		fetcher.FetchHead(sessions[i])
	}
}

func (fetcher *HeadFetcher) FetchHead(session Session) string {
	expire := minTime(fetcher.expire)
	path := getAbsolutePath(fetcher.targetDir, session.Id, fetcher.extension)
	if file, err := os.Stat(path); os.IsNotExist(err) || checkExpired(file, expire) {
		if session.Skin != "" {
			writeSkin(session.Skin, path)
		} else {
			writeSkin(steve, path)
		}
	}
	return path
}

func writeSkin(url string, path string) {
	response, err := http.Get(url)
	if err != nil {
		return
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return
	}

	head := image.NewRGBA(image.Rect(0, 0, 8, 8))
	draw.Draw(head, image.Rect(0, 0, 8, 8), img, image.Pt(8, 8), draw.Src)
	draw.Draw(head, image.Rect(0, 0, 8, 8), img, image.Pt(40, 8), draw.Over)

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
