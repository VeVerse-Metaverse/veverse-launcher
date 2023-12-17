package config

// SessionEncryptionKey is the encryption key used to encrypt the session.
var (
	SessionEncryptionKey = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX="
)

// Supported build configurations.
const (
	Development = "Development"
	Test        = "Test"
	Shipping    = "Shipping"
)

const (
	LauncherPort = "13730" // Port used by the main launcher instance to communicate with the subsequent processes.
	GamePort     = "13731" // Port used by the game instance to communicate with the main launcher instance.
)
