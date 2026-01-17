package sound

// WebSoundManager handles sound for web version
type WebSoundManager struct {
	sounds map[SoundType]string
}

// NewWebSoundManager creates a new web sound manager
func NewWebSoundManager() *WebSoundManager {
	return &WebSoundManager{
		sounds: map[SoundType]string{
			SoundMove:      "/move.mp3",
			SoundCapture:   "/capture.mp3",
			SoundCastle:    "/castle.mp3",
			SoundCheck:     "/check.mp3",
			SoundCheckmate: "/checkmate.mp3",
			SoundIllegal:   "/illegal.mp3",
			SoundPromote:   "/promote.mp3",
		},
	}
}

// GetSoundURL returns the URL for a sound type
func (w *WebSoundManager) GetSoundURL(soundType SoundType) string {
	if url, ok := w.sounds[soundType]; ok {
		return url
	}
	return ""
}

// GetAllSoundURLs returns all sound URLs
func (w *WebSoundManager) GetAllSoundURLs() map[SoundType]string {
	return w.sounds
}
