//go:build darwin
// +build darwin

package platform

// Browser history paths
func GetSafariHistoryPath() (string, error) {
	return ExpandPath("~/Library/Safari/History.db")
}

func GetChromeHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Google/Chrome/Default/History")
}

func GetFirefoxProfilesPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Firefox/Profiles")
}

func GetEdgeHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Microsoft Edge/Default/History")
}

func GetBraveHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/BraveSoftware/Brave-Browser/Default/History")
}

func GetArcHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Arc/User Data/Default/History")
}

// Shell history paths
func GetBashHistoryPath() (string, error) {
	return ExpandPath("~/.bash_history")
}

func GetZshHistoryPath() (string, error) {
	return ExpandPath("~/.zsh_history")
}

func GetZshSessionsPath() (string, error) {
	return ExpandPath("~/.zsh_sessions")
}

func GetFishHistoryPath() (string, error) {
	return ExpandPath("~/.local/share/fish/fish_history")
}

// Application cache paths
func GetCachesPath() (string, error) {
	return ExpandPath("~/Library/Caches")
}

// Recent files paths
func GetRecentDocumentsPath() (string, error) {
	return ExpandPath("~/Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.ApplicationRecentDocuments")
}


