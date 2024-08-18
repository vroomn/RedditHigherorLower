package main

import (
	datamanipulation "RedditHigherorLowerServer/DataManipulation"
	"RedditHigherorLowerServer/DataManipulation/envfile"
)

func main() {
	envfile.Parse()

	datamanipulation.WriteData()

	//datamanipulation.GetSubreddits()
}
