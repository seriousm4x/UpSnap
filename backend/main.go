package main

import (
	"embed"
	"log"

	"ganigeorgiev.com/pocketbase/pocketbase/apis"
	"github.com/seriousm4x/upsnap/pb"
)

//go:embed all:pb_public
var distDir embed.FS
var distDirFS = apis.MustSubFS(distDir, "pb_public")

func main() {
	if err := pb.StartPocketBase(distDirFS); err != nil {
		log.Fatal(err)
	}
}
