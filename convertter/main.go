package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jonas747/dca"
)

func main() {
	files, err := ioutil.ReadDir("../rawmeinkampf")
	if err != nil {
		log.Fatal(err)
	}

	options := &dca.EncodeOptions{
		Volume:        256,
		Channels:      2,
		FrameRate:     48000,
		FrameDuration: 20,
		Bitrate:       128,
		RawOutput:     false,
		Application:   dca.AudioApplication("audio"),
		CoverFormat:   "jpeg",
		VBR:           true,
		Comment:       "",
		Threads:       0,
	}

	for _, f := range files {
		splitted := strings.Split(f.Name(), ".")

		println(splitted[1])

		if splitted[1] == "mp3" {
			encodeSession, err := dca.EncodeFile("../rawmeinkampf/"+f.Name(), options)

			if err != nil {
				fmt.Println("err while encoding: ", err)
			}

			output, err := os.Create("../meinkampf/" + splitted[0] + ".dca")

			if err != nil {
				fmt.Println("err while file create: ", err)
			}

			_, err = io.Copy(output, encodeSession)

			if err != nil {
				fmt.Println("error while copying ", err)
			}
		}

		fmt.Println(f.Name())
	}
}
