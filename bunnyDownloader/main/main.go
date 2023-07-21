package main

import "app"

func main() {
	videoURL := "https://www.bilibili.com/video/BV1Bt41167c3"
	engine := app.Default()
	engine.Download(videoURL)
}
