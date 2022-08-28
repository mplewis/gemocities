package types

type Config struct {
	GeminiHost     string `figyr:"default=:1965"`
	WebDAVHost     string `figyr:"default=:8888"`
	ContentDir     string `figyr:"required"`
	GeminiCertsDir string `figyr:"required"`

	// TODO: Integrate with S3 for ez3
	// TODO: Dev mode uses a local file system for ez3
	S3Bucket    string `figyr:"required"`
	S3Namespace string `figyr:"required"`

	// TODO: Mail credentials

	Development bool `figyr:"optional"`
	Debug       bool `figyr:"optional"`
}
