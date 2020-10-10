package common

import (
	"fmt"
	"os"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func pdfHelp() {
	begin := time.Now()
	m := pdf.NewMaroto(consts.Portrait, consts.Letter)
	//m.SetBorder(true)//是否设置边框
	//Row 第一个参数为高度，增加一行
	m.Row(100, func() {
		//col 增加一列，这里的都会排在同一行，第一个参数是宽，这里的宽和高，都是指的是他自己生成的一个背景
		//这里的图片是依附于他的背景上的，相同的这里的宽会影响左右位置，和上面的高同事i设置会影响显示总大小
		//比如上面设置100，这里4和下面的8
		m.Col(4, func() {
			_ = m.FileImage("./bb.png", props.Rect{
				Center:  true,
				Percent: 80, //显示原来图片的百分比
			})
		})
		m.Col(8, func() {
			m.Text("Gopher International Shipping, Inc.", props.Text{
				Top:         12,
				Size:        20,
				Extrapolate: true,
			})
			m.Text("1000 Shipping Gopher Golang TN 3691234 GopherLand (GL)", props.Text{
				Size: 12,
				Top:  22,
			})
		})
		m.ColSpace(4)
	})
	m.Line(10)
	m.Row(40, func() {
		m.Col(4, func() {
			m.Text("João Sant'Ana 100 Main Street Stringfield TN 39021 United Stats (USA)", props.Text{
				Size: 15,
				Top:  12,
			})
		})
		m.ColSpace(4)
		m.Col(4, func() {
			m.QrCode("https://github.com/johnfercher/maroto", props.Rect{
				Center:  true,
				Percent: 75,
			})
		})
	})
	m.Line(10)
	m.Row(100, func() {
		m.Col(12, func() {
			_ = m.Barcode("https://github.com/johnfercher/maroto", props.Barcode{
				Center:  true,
				Percent: 70,
			})
			m.Text("https://github.com/johnfercher/maroto", props.Text{
				Size:  20,
				Align: consts.Center,
				Top:   65,
			})
		})
	})
	m.SetBorder(true)
	m.Row(40, func() {
		m.Col(6, func() {
			m.Text("CODE: 123412351645231245564 DATE: 20-07-1994 20:20:33", props.Text{
				Size: 15,
				Top:  14,
			})
		})
		m.Col(6, func() {
			m.Text("CA", props.Text{
				Top:   1,
				Size:  85,
				Align: consts.Center,
			})
		})
	})
	m.SetBorder(false)
	err := m.OutputFileAndClose("./file/aa.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))
}
