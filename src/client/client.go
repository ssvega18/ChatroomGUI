package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"net"
)

var isLogged bool = false

func main() {

	// connect to server
	a := app.New()
	w2 := a.NewWindow("Login")

	input1 := widget.NewEntry()
	input2 := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Username", Widget: input1}},
		OnSubmit: func() { // optional, handle form submission

			w := a.NewWindow("Chatroom V1.1.0")
			w.Resize(fyne.NewSize(400, 400))

			conn, err := net.Dial("tcp", "localhost:2000")
			if err != nil {
				fmt.Println("ERROR: Could not connect to server.")

			}
			_, err = conn.Write([]byte(input1.Text + " " + input2.Text))
			if err != nil {
				fmt.Println("ERROR: Could not send message to server.")

			}

			ct := container.NewVBox()
			go readMessage(conn, ct)

			input := widget.NewEntry()

			ct.Add(input)

			input.SetPlaceHolder("Enter Message Here: ")
			enterText := widget.NewButton("Send", func() {
				_, err = conn.Write([]byte(input.Text))
				textWidget := canvas.NewText(input1.Text+":", color.RGBA{R: 179, G: 215, B: 255, A: 255})
				textWidget.TextStyle = fyne.TextStyle{Bold: true}
				textWidget2 := canvas.NewText(input.Text, color.White)
				ct.Add(textWidget)
				ct.Add(textWidget2)
			})
			ct.Add(enterText)

			w.SetContent(ct)
			w.Show()
			w2.Close()

		},
	}

	// we can also append items
	form.Append("Password", input2)

	w2.SetContent(form)
	w2.Resize(fyne.NewSize(300, 300))

	w2.Show()
	a.Run()
}

func readMessage(conn net.Conn, ct *fyne.Container) {
	for {
		buffer := make([]byte, 1400)
		dataSize, err := conn.Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}

		username := string(buffer[:dataSize])

		buffer = make([]byte, 1400)
		dataSize, err = conn.Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}

		message := string(buffer[:dataSize])

		textWidget := canvas.NewText(username+":", color.RGBA{R: 0, G: 255, B: 0, A: 255})
		textWidget.TextStyle = fyne.TextStyle{Bold: true}
		textWidget2 := canvas.NewText(message, color.White)
		ct.Add(textWidget)
		ct.Add(textWidget2)
	}
}
