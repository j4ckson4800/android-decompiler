package decompiler

type Option func(*ParseConfig)

type ParseConfig struct {
	/*
		SanitizeAnnotations removes annotations from the dex files.

		If we have a lot of heavy patterns or apk has a lot of annotations
		we might sanitize them to reduce string count.
		Average app returns ~24k strings with sanitization instead of ~34k,
		but since annotation parsing is heavy operation we might want to skip it
		for lightweight patterns sanitization performs 2x slower than without it

		b.Run("without annotations", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				apk, _ := decompiler.NewApk(data, decompiler.WithSanitizeAnnotations())
				_ = apk.GetConstrStrings()  // hacky method to get all the const strings
			}
		})

		b.Run("with annotations", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				apk, _ := decompiler.NewApk(data)
				_ = apk.GetConstrStrings()  // hacky method to get all the const strings
			}
		})

		Benchmark on average apk with single pattern:
		BenchmarkExtractStrings/without_annotations-32                 5        1149425920 ns/op
		BenchmarkExtractStrings/with_annotations-32                    9         610055044 ns/op
	*/
	SanitizeAnnotations bool

	/*
		FailOnInvalidDex stops parsing dex files if any of them is invalid.

		Dex file may be invalid only if it's encrypted.
		The most common case of this is when we have a dex file stored in assets folder.

		Example of app with encrypted dex: com.conquer.domino v1.1.8.0
	*/
	FailOnInvalidDex bool

	/*
		FailOnInvalidResource stops parsing apk if .arsc file is invalid in some way
	*/
	FailOnInvalidResource bool
}

func WithSanitizeAnnotations() Option {
	return func(cfg *ParseConfig) {
		cfg.SanitizeAnnotations = true
	}
}

func WithFailOnInvalidDex() Option {
	return func(cfg *ParseConfig) {
		cfg.FailOnInvalidDex = true
	}
}

func WithFailOnInvalidResource() Option {
	return func(cfg *ParseConfig) {
		cfg.FailOnInvalidResource = true
	}
}
