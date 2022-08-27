package types

type Config struct {
	GeminiHost     string `figyr:"default=:1965"`
	WebDAVHost     string `figyr:"default=:8888"`
	ContentDir     string `figyr:"required"`
	GeminiCertsDir string `figyr:"required"`

	S3Bucket    string `figyr:"required"`
	S3Namespace string `figyr:"required"`

	Development bool `figyr:"optional"`
	Debug       bool `figyr:"optional"`
}
