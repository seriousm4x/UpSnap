package main

import (
	"embed"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/seriousm4x/upsnap/pb"
)

//go:embed all:pb_public
var distDir embed.FS
var distDirFS = apis.MustSubFS(distDir, "pb_public")

func main() {
	pb.StartPocketBase(distDirFS)
}
