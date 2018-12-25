package main

import (
	"encoding/json"
	"image"
	"log"
	"mime/multipart"
	"strconv"

	pc "github.com/EdlinOrg/prominentcolor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// A Color is a struct containing a ColorItem and a Hex value
type Color struct {
	pc.ColorItem
	Hex string
}

func NewColor(color pc.ColorItem) Color {
	return Color{Hex: "#" + color.AsString(), ColorItem: color}
}

func process(k int, arg int, img image.Image) (output []pc.ColorItem, err error) {
	res, err := pc.KmeansWithAll(k, img, arg, uint(pc.DefaultSize), pc.GetDefaultMasks())
	if err != nil {
		log.Println(err)
	}

	for _, color := range res {
		output = append(output, color)
	}

	return output, nil
}

func analyzeImage(img image.Image, n int) ([]Color, error) {
	colorItems, err := process(n, pc.ArgumentAverageMean, img)
	if err != nil {
		return nil, err
	}

	var colors []Color
	for _, v := range colorItems {
		colors = append(colors, NewColor(v))
	}

	return colors, nil
}

func openImage(fileHeader *multipart.FileHeader) (image.Image, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func analyzeImageRoute(c echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return err
	}

	img, err := openImage(fileHeader)

	result, err := analyzeImage(img, 3)

	b, err := json.Marshal(&result)
	if err != nil {
		return err
	}

	c.String(200, string(b))

	return nil
}

func analyzeNColorsImageRoute(c echo.Context) error {
	fileHeader, err := c.FormFile("image")
	n := c.Param("n")
	i, err := strconv.Atoi(n)
	if err != nil {
		return err
	}

	img, err := openImage(fileHeader)

	result, err := analyzeImage(img, i)

	b, err := json.Marshal(&result)
	if err != nil {
		return err
	}

	c.String(200, string(b))

	return nil
}

func main() {
	// TODO: Replace echo with fasthttp or net/http
	port := ":3000"

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/analyze", analyzeImageRoute)
	e.POST("/analyze/:n", analyzeNColorsImageRoute)
	e.Logger.Fatal(e.Start(port))
}
