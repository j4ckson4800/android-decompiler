package main

import (
	"flag"
	"fmt"

	"github.com/j4ckson4800/android-decompiler/cmd/proto-gen/internal/app"
	"github.com/j4ckson4800/android-decompiler/cmd/proto-gen/internal/app/codegen"
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

	gen, err := codegen.NewCodegen(*outDir)
	if err != nil {
		fmt.Printf("Error initializing codegen: %v\n", err)
		return
	}

	parser, err := app.NewParser(*apkFile, gen)
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
