package store

import "fmt"

type DownloadHistory struct{}

func (dh *DownloadHistory) Save(title, url string, images []string) {
	fmt.Println("[√] title:", title, "source:", url, "images", images)
}

func (dh *DownloadHistory) Failed(title, url string, images []string) {
	fmt.Println("[×] title:", title, "source:", url, "images", images)
}
