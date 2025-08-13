package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug membuat slug/ alias nama produk di URL dari teks utk seo
func GenerateSlug(input string) string {
	slug := strings.ToLower(input)
	// ganti spasi dan underscore jadi tanda minus
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// hapus karakter non-alfanumerik kecuali tanda minus
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")
	// hilangkan tanda minus ganda
	slug = strings.ReplaceAll(slug, "--", "-")
	// hapus tanda minus di awal/akhir
	slug = strings.Trim(slug, "-")
	return slug
}
