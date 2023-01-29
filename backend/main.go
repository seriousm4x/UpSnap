package main

import (
	"embed"

	"github.com/labstack/echo/v5"
	"github.com/seriousm4x/upsnap/pb"
)

//go:embed all:pb_public
var distDir embed.FS
var distDirFS = echo.MustSubFS(distDir, "pb_public")

func main() {
	pb.StartPocketBase(distDirFS)
}
