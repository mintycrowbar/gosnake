package main

import (
	"log"

	"github.com/nsf/termbox-go"
)

type PlayerPosition struct {
	posX int
	posY int
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	width, height := termbox.Size()
	drawBorders()
	termbox.SetCell(width/2, height/2, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}

	position := PlayerPosition{width / 2, height / 2}
	downRunCount := 0
	upRunCount := 0
	leftRunCount := 0
	rightRunCount := 0
	for {
		event := termbox.PollEvent()

		if event.Type == termbox.EventKey && (event.Ch == 'w' || event.Key == termbox.KeyArrowUp) {
			upRunCount++
			upRunCount, position = movePlayer(position.posX, position.posY, upRunCount, "up")
		} else if event.Type == termbox.EventKey && (event.Ch == 's' || event.Key == termbox.KeyArrowDown) {
			downRunCount++
			downRunCount, position = movePlayer(position.posX, position.posY, downRunCount, "down")
		} else if event.Type == termbox.EventKey && (event.Ch == 'a' || event.Key == termbox.KeyArrowLeft) {
			leftRunCount++
			leftRunCount, position = movePlayer(position.posX, position.posY, leftRunCount, "left")
		} else if event.Type == termbox.EventKey && (event.Ch == 'd' || event.Key == termbox.KeyArrowRight) {
			rightRunCount++
			rightRunCount, position = movePlayer(position.posX, position.posY, rightRunCount, "right")
		} else if event.Type == termbox.EventKey && event.Ch == 'q' {
			return
		}
	}
}

func movePlayer(posX int, posY int, runIndex int, direction string) (int, PlayerPosition) {
	position := PlayerPosition{}

	switch direction {
	case "up":
		termbox.SetCell(posX, posY-1, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{posX, posY - 1}
		termbox.SetCell(posX, posY+1, ' ', termbox.ColorBlack, termbox.ColorDefault)
	case "down":
		termbox.SetCell(posX, posY+1, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{posX, posY + 1}
		termbox.SetCell(posX, posY-1, ' ', termbox.ColorBlack, termbox.ColorDefault)
	case "left":
		termbox.SetCell(posX-1, posY, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{posX - 1, posY}
		termbox.SetCell(posX+1, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)
	case "right":
		termbox.SetCell(posX+1, posY, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{posX + 1, posY}
		termbox.SetCell(posX-1, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)
	}

	termbox.SetCell(posX, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)
	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return runIndex, position
}

func drawBorders() {
	fg := termbox.ColorBlack
	bg := termbox.ColorDefault

	width, height := termbox.Size()

	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		log.Fatal(err)
	}

	// draw the left border
	for i := 0; i < height; i++ {
		termbox.SetCell(0, i, '#', fg, bg)
	}

	// draw the right border
	for i := 0; i < height; i++ {
		termbox.SetCell(width-1, i, '#', fg, bg)
	}

	// draw the top border
	for i := 0; i < width; i++ {
		termbox.SetCell(i, 0, '#', fg, bg)
	}

	// draw the bottom border
	for i := 0; i < width; i++ {
		termbox.SetCell(i, height-1, '#', fg, bg)
	}
}
