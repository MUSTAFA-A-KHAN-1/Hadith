package image

import (
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/abdullahdiaa/garabic"
	"github.com/fogleman/gg"
)

type Generator struct {
	fontDir string
}

func NewGenerator(fontDir string) *Generator {
	return &Generator{
		fontDir: fontDir,
	}
}

func (g *Generator) GenerateHadithImage(title, narrator, arabicText, englishText, reference string) ([]byte, error) {
	const W, H = 1080, 1080
	dc := gg.NewContext(W, H)

	g.drawBackground(dc)

	arabicFontPath := g.getFontPath("Amiri-Regular.ttf")
	englishFontPath := g.getFontPath("Caveat-Regular.ttf")

	// --- 1. Title (Top, Green, Uppercase) ---
	dc.SetHexColor("#558B2F") // Olive Green
	if err := dc.LoadFontFace(englishFontPath, 110); err != nil {
		return nil, fmt.Errorf("failed to load title font: %w", err)
	}
	titleY := 150.0
	dc.DrawStringAnchored(strings.ToUpper(title), float64(W)/2, titleY, 0.5, 0.5)

	// --- 2. Attribution (Black, smaller) ---
	dc.SetHexColor("#1a1a1a") // Black
	attributionY := titleY + 80

	// If narrator is empty, use default. If provided, use it.
	// We want to handle the "Prophet Muhammad ﷺ" part intelligently if present in the narrator string,
	// but that requires complex parsing.
	// For now, if the narrator string contains "Prophet" or "Messenger", we just print it as is (using Caveat).
	// If we want to support the symbol, we'd need to split the string.
	// Given the constraint "narrator is extracted from JSON", let's print the narrator string.

	// However, the user wants "include all these details in the generated image as shown in the image".
	// The image shows "The Prophet Muhammad ﷺ said:".
	// If the JSON `Narrator` field is just "Umar bin Al-Khattab", then the text should probably be:
	// "Umar bin Al-Khattab narrated that The Prophet Muhammad ﷺ said:"?
	// Or maybe the `Narrator` field in the JSON *already* contains the full chain?
	// Let's assume `narrator` is the text to be displayed.

	// We will try to replace "Prophet Muhammad" with "Prophet Muhammad ﷺ" if it doesn't have it,
	// OR just render the symbol if we can detect the placeholder.
	// But `Caveat` font doesn't have the symbol. `Amiri` does.
	// So we need to render the symbol with Amiri.

	// Heuristic: If narrator string is provided, use it.
	// If it contains "Prophet Muhammad" or "Messenger of Allah", we can try to inject the symbol.
	// But simply rendering the narrator string using Caveat is the safest first step.
	// The user said "the narrator is not hardcoded".

	displayText := narrator
	if displayText == "" {
		displayText = "The Prophet Muhammad ﷺ said:"
	} else {
		// Ensure it ends with a colon if it looks like an intro
		if !strings.HasSuffix(displayText, ":") && !strings.HasSuffix(displayText, ".") {
			displayText += ":"
		}
	}

	// Check if we need to inject the symbol (U+FDFA)
	// If the string contains "Pbuh" or "SAW", replace it?
	displayText = strings.ReplaceAll(displayText, "(saw)", "ﷺ")
	displayText = strings.ReplaceAll(displayText, "(pbuh)", "ﷺ")

	// Drawing logic with potential mixed fonts
	if strings.Contains(displayText, "ﷺ") {
		parts := strings.Split(displayText, "ﷺ")

		totalWidth := 0.0
		// Measure width first
		for i, part := range parts {
			dc.LoadFontFace(englishFontPath, 50)
			w, _ := dc.MeasureString(part)
			totalWidth += w
			if i < len(parts)-1 {
				dc.LoadFontFace(arabicFontPath, 50)
				w, _ = dc.MeasureString("ﷺ")
				totalWidth += w
			}
		}

		startX := (float64(W) - totalWidth) / 2
		currentX := startX

		for i, part := range parts {
			dc.LoadFontFace(englishFontPath, 50)
			dc.DrawStringAnchored(part, currentX, attributionY, 0, 0.5)
			w, _ := dc.MeasureString(part)
			currentX += w

			if i < len(parts)-1 {
				dc.LoadFontFace(arabicFontPath, 50)
				dc.DrawStringAnchored("ﷺ", currentX, attributionY, 0, 0.5)
				w, _ = dc.MeasureString("ﷺ")
				currentX += w
			}
		}
	} else {
		// Just render plain text
		if err := dc.LoadFontFace(englishFontPath, 50); err != nil {
			return nil, fmt.Errorf("failed to load attribution font: %w", err)
		}
		dc.DrawStringAnchored(displayText, float64(W)/2, attributionY, 0.5, 0.5)
	}

	// --- 3. Arabic Text (Centered, Large) ---
	dc.SetHexColor("#000000") // Black
	if err := dc.LoadFontFace(arabicFontPath, 70); err != nil {
		return nil, fmt.Errorf("failed to load arabic font: %w", err)
	}

	shapedArabic := garabic.Shape(arabicText)
	maxWidth := float64(W) - 160
	lines := dc.WordWrap(shapedArabic, maxWidth)
	lineHeight := dc.FontHeight() * 1.5
	arabicHeight := float64(len(lines)) * lineHeight

	// Position Arabic below attribution with some gap
	arabicStartY := attributionY + 100 + (arabicHeight / 2)

	for i, line := range lines {
		reversedLine := g.reversePreservingCombiningMarks(line)
		dc.DrawStringAnchored(reversedLine, float64(W)/2, arabicStartY+float64(i)*lineHeight - (arabicHeight/2), 0.5, 0.5)
	}

	// --- 4. English Translation (Centered, Caveat) ---
	dc.SetHexColor("#1a1a1a")
	if err := dc.LoadFontFace(englishFontPath, 60); err != nil {
		return nil, fmt.Errorf("failed to load english font: %w", err)
	}

	englishLines := dc.WordWrap(englishText, maxWidth)
	englishHeight := float64(len(englishLines)) * (dc.FontHeight() * 1.2)

	// Position English below Arabic with gap
	englishStartY := arabicStartY + (arabicHeight/2) + 80 + (englishHeight/2)

	for i, line := range englishLines {
		dc.DrawStringAnchored(line, float64(W)/2, englishStartY+float64(i)*(dc.FontHeight()*1.2) - (englishHeight/2), 0.5, 0.5)
	}

	// --- 5. Reference (Bottom, Smaller) ---
	dc.SetHexColor("#4a4a4a") // Dark Gray
	if err := dc.LoadFontFace(englishFontPath, 40); err != nil {
		return nil, fmt.Errorf("failed to load ref font: %w", err)
	}
	refY := float64(H) - 100
	dc.DrawStringAnchored(reference, float64(W)/2, refY, 0.5, 0.5)

	var buf bytes.Buffer
	if err := dc.EncodePNG(&buf); err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *Generator) drawBackground(dc *gg.Context) {
	// Light blue/white tint similar to reference image
	dc.SetHexColor("#F0F8FF")
	dc.Clear()

	// Add subtle noise/texture
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	width := dc.Width()
	height := dc.Height()

	// Faint blobs
	for i := 0; i < 5; i++ {
		x := rnd.Float64() * float64(width)
		y := rnd.Float64() * float64(height)
		r := 100 + rnd.Float64()*200

		// Very light pink/orange/blue pastel blobs
		rCol := 200 + rnd.Intn(55)
		gCol := 200 + rnd.Intn(55)
		bCol := 200 + rnd.Intn(55)

		dc.SetRGBA255(rCol, gCol, bCol, 20) // Very transparent
		dc.DrawCircle(x, y, r)
		dc.Fill()
	}
}

func (g *Generator) getFontPath(fontName string) string {
	return filepath.Join(g.fontDir, fontName)
}

func (g *Generator) reversePreservingCombiningMarks(s string) string {
	runes := []rune(s)
	var clusters [][]rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if isCombiningMark(r) && len(clusters) > 0 {
			clusters[len(clusters)-1] = append(clusters[len(clusters)-1], r)
		} else {
			clusters = append(clusters, []rune{r})
		}
	}

	for i, j := 0, len(clusters)-1; i < j; i, j = i+1, j-1 {
		clusters[i], clusters[j] = clusters[j], clusters[i]
	}

	var res []rune
	for _, cluster := range clusters {
		res = append(res, cluster...)
	}
	return string(res)
}

func isCombiningMark(r rune) bool {
	return (r >= 0x064B && r <= 0x065F) ||
		r == 0x0670 ||
		(r >= 0x0610 && r <= 0x061A) ||
		(r >= 0x06D6 && r <= 0x06DC) ||
		(r >= 0x06DF && r <= 0x06E4) ||
		(r >= 0x06E7 && r <= 0x06E8) ||
		(r >= 0x06EA && r <= 0x06ED)
}
