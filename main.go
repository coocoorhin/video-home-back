package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type videolist struct {
	Name     string
	Location string
}

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// router.Static("/", "./public")

	router.GET("/videolist", func(c *gin.Context) {

		var list []videolist

		root := "../videos/"
		err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			if filepath.Ext(info.Name()) == ".mp4" {
				absPath, _ := filepath.Abs(path)
				list = append(list, videolist{Name: info.Name(), Location: absPath})
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error walking the path ")
			return
		}

		if list != nil {
			listJSON, _ := json.Marshal(list)
			c.JSON(http.StatusOK, listJSON)
		} else {
			c.Status(http.StatusNotFound)
		}

	})

	router.POST("/upload", func(c *gin.Context) {
		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		if filepath.Ext(file.Filename) != ".torrent" {
			c.Status(http.StatusUnprocessableEntity)
		} else {
			filenameDst := "../torrents/" + filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, filenameDst); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))

		}

	})
	router.Run(":8080")
}
