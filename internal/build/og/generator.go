package og

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"site/internal/models"

	"github.com/fogleman/gg"
)

const (
	imageWidth  = 1200
	imageHeight = 630
)

type Generator struct {
	outputDir    string
	profilePhoto image.Image
	fontPath     string
	fontBoldPath string
	siteName     string
	siteURL      string
}

func NewGenerator(outputDir, profilePhotoPath, fontPath, fontBoldPath, siteName, siteURL string) (*Generator, error) {
	g := &Generator{
		outputDir:    outputDir,
		fontPath:     fontPath,
		fontBoldPath: fontBoldPath,
		siteName:     siteName,
		siteURL:      siteURL,
	}

	if profilePhotoPath != "" {
		photo, err := loadImage(profilePhotoPath)
		if err != nil {
			return nil, err
		}
		g.profilePhoto = photo
	}

	return g, nil
}

func (g *Generator) Generate(post *models.Post, collection *models.Collection) (string, error) {
	dc := gg.NewContext(imageWidth, imageHeight)

	dc.SetColor(color.White)
	dc.Clear()

	badgeText := "POST"
	seriesName := ""
	if collection != nil {
		if collection.Type == models.CollectionTypeTopic {
			badgeText = "DOCS"
			seriesName = collection.Name
		} else if collection.Type == models.CollectionTypeSeries {
			if collection.Slug != "blog" {
				seriesName = collection.Name
			}
		}
	}

	contentX := 80.0
	if g.profilePhoto != nil {
		contentX = 240.0
		g.drawProfilePhoto(dc, 80, 200)
	}

	g.drawBadge(dc, badgeText, contentX, 180)
	g.drawTitle(dc, post.Title, contentX, 260)

	if seriesName != "" {
		g.drawSubtitle(dc, seriesName, contentX, 360)
	}

	g.drawFooter(dc)

	outputPath := g.getOutputPath(post, collection)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", err
	}

	if err := dc.SavePNG(outputPath); err != nil {
		return "", err
	}

	return g.getWebPath(post, collection), nil
}

func (g *Generator) drawBadge(dc *gg.Context, text string, x, y float64) {
	if err := dc.LoadFontFace(g.fontBoldPath, 14); err != nil {
		return
	}

	textW, _ := dc.MeasureString(text)
	padX := 14.0
	badgeW := textW + padX*2
	badgeH := 28.0
	radius := 4.0

	dc.SetColor(color.RGBA{30, 30, 30, 255})
	dc.DrawRoundedRectangle(x, y, badgeW, badgeH, radius)
	dc.Fill()

	dc.SetColor(color.White)
	dc.DrawString(text, x+padX, y+20)
}

func (g *Generator) drawTitle(dc *gg.Context, title string, x, y float64) {
	if err := dc.LoadFontFace(g.fontBoldPath, 56); err != nil {
		return
	}

	dc.SetColor(color.RGBA{20, 20, 20, 255})

	maxWidth := float64(imageWidth) - x - 80
	lines := dc.WordWrap(title, maxWidth)

	lineHeight := 68.0
	maxLines := 2

	for i, line := range lines {
		if i >= maxLines {
			break
		}
		if i == maxLines-1 && len(lines) > maxLines {
			line = strings.TrimSuffix(line, " ")
			if len(line) > 3 {
				line = line[:len(line)-3] + "..."
			}
		}
		dc.DrawString(line, x, y+float64(i)*lineHeight)
	}
}

func (g *Generator) drawSubtitle(dc *gg.Context, text string, x, y float64) {
	if err := dc.LoadFontFace(g.fontPath, 26); err != nil {
		return
	}

	dc.SetColor(color.RGBA{100, 100, 100, 255})
	dc.DrawString(text, x, y)
}

func (g *Generator) drawProfilePhoto(dc *gg.Context, x, y float64) {
	size := 120.0
	centerX := x + size/2
	centerY := y + size/2

	dc.SetColor(color.RGBA{230, 230, 230, 255})
	dc.DrawCircle(centerX, centerY, size/2+2)
	dc.Fill()

	dc.DrawCircle(centerX, centerY, size/2)
	dc.Clip()

	bounds := g.profilePhoto.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())
	scale := size / min(imgW, imgH)

	dc.Push()
	dc.Translate(centerX, centerY)
	dc.Scale(scale, scale)
	dc.DrawImageAnchored(g.profilePhoto, 0, 0, 0.5, 0.5)
	dc.Pop()

	dc.ResetClip()
}

func (g *Generator) drawFooter(dc *gg.Context) {
	dc.SetColor(color.RGBA{230, 230, 230, 255})
	dc.SetLineWidth(1)
	dc.DrawLine(80, imageHeight-90, imageWidth-80, imageHeight-90)
	dc.Stroke()

	if err := dc.LoadFontFace(g.fontPath, 20); err != nil {
		return
	}

	dc.SetColor(color.RGBA{100, 100, 100, 255})

	url := g.siteURL
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, "/")

	text := g.siteName + "  ·  " + url
	dc.DrawString(text, 80, imageHeight-45)
}

func (g *Generator) getOutputPath(post *models.Post, collection *models.Collection) string {
	if collection != nil {
		return filepath.Join(g.outputDir, "og", collection.Slug, post.Slug+".png")
	}
	return filepath.Join(g.outputDir, "og", post.Slug+".png")
}

func (g *Generator) getWebPath(post *models.Post, collection *models.Collection) string {
	if collection != nil {
		return "/og/" + collection.Slug + "/" + post.Slug + ".png"
	}
	return "/og/" + post.Slug + ".png"
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		f.Seek(0, 0)
		img, _, err = image.Decode(f)
		if err != nil {
			return nil, err
		}
	}
	return img, nil
}
