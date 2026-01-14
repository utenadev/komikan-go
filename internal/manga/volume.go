package manga

import (
	"regexp"
	"strconv"
	"strings"
)

// VolumeInfo represents extracted volume information
type VolumeInfo struct {
	Title       string
	Volume      int
	HasVolume   bool
	IsSpecial   bool // true for special editions like "ダイズカン"
}

// volumePattern matches volume numbers at the end of titles
var volumePattern = regexp.MustCompile(`\s+(\d+)$`)
// specialPattern matches special edition keywords
var specialPattern = regexp.MustCompile(`(ダイズカン|ガイド|ファンブック|イラスト集|公式ブック|設定資料集|カラー版|完全版|総編集|愛蔵版)`)

// ExtractVolumeInfo extracts volume information from title
func ExtractVolumeInfo(title string) VolumeInfo {
	// Check for special edition
	if specialPattern.MatchString(title) {
		return VolumeInfo{
			Title:     title,
			HasVolume: false,
			IsSpecial: true,
		}
	}

	// Extract volume number from end of title
	matches := volumePattern.FindStringSubmatch(title)
	if len(matches) >= 2 {
		volume, err := strconv.Atoi(matches[1])
		if err == nil {
			// Remove volume number from base title
			baseTitle := volumePattern.ReplaceAllString(title, "")
			return VolumeInfo{
				Title:     strings.TrimSpace(baseTitle),
				Volume:    volume,
				HasVolume: true,
				IsSpecial: false,
			}
		}
	}

	return VolumeInfo{
		Title:     title,
		HasVolume: false,
		IsSpecial: false,
	}
}

// NormalizeTitleForSearch normalizes title for API search
func NormalizeTitleForSearch(title string) string {
	info := ExtractVolumeInfo(title)
	if info.HasVolume {
		return info.Title
	}
	return title
}

// ShouldIncludeAsSeriesVolume determines if a book should be included in series volume tracking
func ShouldIncludeAsSeriesVolume(title, author string) bool {
	info := ExtractVolumeInfo(title)

	// Exclude special editions from volume tracking
	if info.IsSpecial {
		return false
	}

	// Include if it has a volume number
	if info.HasVolume {
		return true
	}

	return false
}
