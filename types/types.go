package types

type Config struct {
	AppDomain      string `figyr:"required,description=The domain name where the app is hosted"`
	GeminiHost     string `figyr:"default=:1965,description=The address to listen on for Gemini requests"`
	WebDAVHost     string `figyr:"default=:8888,description=The address to listen on for WebDAV requests"`
	HTTPHost       string `figyr:"default=:8080,description=The address to listen on for HTTP requests"`
	ContentDir     string `figyr:"required,description=The directory on disk where user content is stored"`
	GeminiCertsDir string `figyr:"required,description=The directory on disk where the server's Gemini certificates are stored"`

	SMTPHost     string `figyr:"required,description=The hostname of the SMTP email server"`
	SMTPPort     int    `figyr:"default=587,description=The port of the SMTP email server"`
	SMTPUsername string `figyr:"required,description=The username to use when connecting to the SMTP email server"`
	SMTPPassword string `figyr:"required,description=The password to use when connecting to the SMTP email server"`

	Development bool `figyr:"optional,description=Configure logging for local development"`
	Debug       bool `figyr:"optional,description=Set log level to debug\\, not info"`
}
