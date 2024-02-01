package builder

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/gosimple/slug"
)

type Scheme struct {
	System      string            `yaml:"system"`
	Name        string            `yaml:"name"`
	Slug        string            `yaml:"slug"`
	Author      string            `yaml:"author"`
	Description string            `yaml:"description"`
	Variant     string            `yaml:"variant"`
	Palette     map[string]string `yaml:"palette"`
}

func BuildTemplate(scheme Scheme, template []byte) (string, error) {
	templateVars := map[string]any{}

	templateVars["scheme-name"] = scheme.Name
	templateVars["scheme-author"] = scheme.Author
	templateVars["scheme-description"] = scheme.Description
	templateVars["scheme-slug"] = scheme.Slug
	templateVars["scheme-slug-underscored"] = strings.ReplaceAll(scheme.Slug, "-", "_")
	templateVars["scheme-system"] = scheme.System
	templateVars["scheme-variant"] = scheme.Variant

	if scheme.Variant != "" {
		templateVars["scheme-is-"+scheme.Variant+"-variant"] = true
	}

	if scheme.Slug == "" {
		newSlug := slug.Make(scheme.Name)
		templateVars["scheme-slug"] = newSlug
		templateVars["scheme-slug-underscored"] = strings.ReplaceAll(newSlug, "-", "_")
	}

	for key, clrString := range scheme.Palette {
		c, err := ParseHexColor(clrString)
		if err != nil {
			return "", err
		}

		templateVars[key+"-hex"] = fmt.Sprintf("%02x%02x%02x", c.R, c.G, c.B)
		templateVars[key+"-hex-bgr"] = fmt.Sprintf("%02x%02x%02x", c.B, c.G, c.R)
		templateVars[key+"-hex-r"] = fmt.Sprintf("%02x", c.R)
		templateVars[key+"-hex-g"] = fmt.Sprintf("%02x", c.G)
		templateVars[key+"-hex-b"] = fmt.Sprintf("%02x", c.B)
		templateVars[key+"-rgb-r"] = c.R
		templateVars[key+"-rgb-g"] = c.G
		templateVars[key+"-rgb-b"] = c.B
		templateVars[key+"-dec-r"] = float32(c.R) / 255
		templateVars[key+"-dec-g"] = float32(c.G) / 255
		templateVars[key+"-dec-b"] = float32(c.B) / 255
	}

	data, err := mustache.Render(string(template), templateVars)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(data), nil
}

func ParseHexColor(hexColor string) (color.RGBA, error) {
	if hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}

	re := regexp.MustCompile("^[0-9a-fA-F]{3}$|^[0-9a-fA-F]{6}$")

	if !re.MatchString(hexColor) {
		return color.RGBA{}, errors.New("Invalid hex color!")
	}

	if len(hexColor) == 3 {
		hexColor = doubleEachChar(hexColor)
	}

	value, err := strconv.ParseUint(hexColor, 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}

	rgba := color.RGBA{
		R: uint8((value >> 16) & 0xFF),
		G: uint8((value >> 8) & 0xFF),
		B: uint8(value & 0xFF),
		A: 255,
	}

	return rgba, nil
}

func doubleEachChar(input string) string {
	result := ""
	for _, char := range input {
		result += string(char) + string(char)
	}
	return result
}

func rgbToXterm256(r, g, b uint8) int {
	xtermIndex := (int(r) * 6 / 256) * 36
	xtermIndex += (int(g) * 6 / 256) * 6
	xtermIndex += int(b) * 6 / 256
	xtermIndex += 16
	return xtermIndex
}
