package main

import (
	"flag"
	"fmt"

	"github.com/j4ckson4800/android-decompiler/cmd/proto-gen/internal/app"
)

func main() {
	apkFile := flag.String("apk", "", "APK file to parse")
	outDir := flag.String("o", "", "Output directory")
	flag.Parse()

	if *apkFile == "" {
		flag.Usage()
		return
	}

	if *outDir == "" {
		flag.Usage()
		return
	}

	parser, err := app.NewParser(*apkFile, *outDir)
	if err != nil {
		fmt.Printf("Error initializing parser: %v\n", err)
		return
	}

	if err := parser.Parse(); err != nil {
		fmt.Printf("Error parsing APK: %v\n", err)
		return
	}

	if err := parser.GenerateProtoDefs(); err != nil {
		fmt.Printf("Error generating protos: %v\n", err)
		return
	}
}
