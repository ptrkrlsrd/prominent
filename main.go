package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"mime/multipart"

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

type Colors struct {
	Light  Color
	Middle Color
	Dark   Color
}

func NewColorsFromSlice(colorSlice []pc.ColorItem) (Colors, error) {
	if len(colorSlice) != 3 {
		return Colors{}, fmt.Errorf("colorSlice contains too few elements: %v", len(colorSlice))
	}

	result := Colors{
		Light:  NewColor(colorSlice[2]),
		Middle: NewColor(colorSlice[1]),
		Dark:   NewColor(colorSlice[0]),
	}

	return result, nil
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

func analyzeImage(img image.Image) (Colors, error) {
	str, err := process(3, pc.ArgumentAverageMean, img)
	if err != nil {
		return Colors{}, err
	}

	return NewColorsFromSlice(str)
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

	result, err := analyzeImage(img)

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
	e.Logger.Fatal(e.Start(port))
}
