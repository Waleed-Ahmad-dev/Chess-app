package sound

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// SoundType represents different chess sounds
type SoundType int

const (
	SoundMove SoundType = iota
	SoundCapture
	SoundCastle
	SoundCheck
	SoundCheckmate
	SoundIllegal
	SoundPromote
)

var (
	soundEnabled = true
	assetPath    = "internal/server/assets"
)

// EnableSound enables sound playback
func EnableSound() {
	soundEnabled = true
}

// DisableSound disables sound playback
func DisableSound() {
	soundEnabled = false
}

// SetAssetPath sets the custom path for sound assets
func SetAssetPath(path string) {
	assetPath = path
}

// PlaySound plays a specific chess sound
func PlaySound(soundType SoundType) error {
	if !soundEnabled {
		return nil
	}

	var soundFile string
	switch soundType {
	case SoundMove:
		soundFile = "move.mp3"
	case SoundCapture:
		soundFile = "capture.mp3"
	case SoundCastle:
		soundFile = "castle.mp3"
	case SoundCheck:
		soundFile = "check.mp3"
	case SoundCheckmate:
		soundFile = "checkmate.mp3"
	case SoundIllegal:
		soundFile = "illegal.mp3"
	case SoundPromote:
		soundFile = "promote.mp3"
	default:
		return fmt.Errorf("unknown sound type")
	}

	// Get the absolute path to the sound file
	exePath, err := os.Executable()
	if err != nil {
		// Fallback to relative path
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	soundPath := filepath.Join(exeDir, assetPath, soundFile)

	// Check if file exists
	if _, err := os.Stat(soundPath); os.IsNotExist(err) {
		// Try relative to current directory
		soundPath = filepath.Join(assetPath, soundFile)
	}

	// Play sound based on OS
	return playSystemSound(soundPath)
}

// playSystemSound plays a sound file using system commands
func playSystemSound(soundPath string) error {
	if _, err := os.Stat(soundPath); os.IsNotExist(err) {
		return fmt.Errorf("sound file not found: %s", soundPath)
	}

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd := exec.Command("afplay", soundPath)
		return cmd.Run()
	case "linux":
		// Try various Linux audio players
		players := []string{"aplay", "mpg123", "mpg321", "ffplay", "paplay"}
		for _, player := range players {
			if _, err := exec.LookPath(player); err == nil {
				args := []string{soundPath}
				if player == "ffplay" {
					args = []string{"-nodisp", "-autoexit", soundPath}
				}
				cmd := exec.Command(player, args...)
				return cmd.Run()
			}
		}
		return fmt.Errorf("no compatible audio player found on Linux")
	case "windows":
		cmd := exec.Command("powershell", "-c",
			fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", soundPath))
		return cmd.Run()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
