package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"

	pc "github.com/EdlinOrg/prominentcolor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

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

func NewColors(color []pc.ColorItem) Colors {
	result := Colors{
		Light:  NewColor(color[2]),
		Middle: NewColor(color[1]),
		Dark:   NewColor(color[0]),
	}

	return result
}

func process(k int, arg int, img image.Image) (output []pc.ColorItem, err error) {
	res, err := pc.KmeansWithAll(k, img, arg, uint(pc.DefaultSize), pc.GetDefaultMasks())
	if err != nil {
		log.Println(err)
	}

	if len(res) != 3 {
		return nil, fmt.Errorf("")
	}

	for _, color := range res {
		output = append(output, color)
	}

	return output, nil
}

func analyzeImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	f, err := file.Open()
	if err != nil {
		return err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	str, err := process(3, pc.ArgumentAverageMean, img)
	if err != nil {
		return err
	}

	result := NewColors(str)

	b, err := json.Marshal(&result)
	if err != nil {
		return err
	}

	return c.String(200, string(b))
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/analyze", analyzeImage)
	e.Logger.Fatal(e.Start(":3000"))
}
