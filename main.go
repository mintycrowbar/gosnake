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
			upRunCount, position = moveUp(position.posX, position.posY, upRunCount)
		} else if event.Type == termbox.EventKey && (event.Ch == 's' || event.Key == termbox.KeyArrowDown) {
			downRunCount++
			downRunCount, position = moveDown(position.posX, position.posY, downRunCount)
		} else if event.Type == termbox.EventKey && (event.Ch == 'a' || event.Key == termbox.KeyArrowLeft) {
			leftRunCount++
			leftRunCount, position = moveLeft(position.posX, position.posY, leftRunCount)
		} else if event.Type == termbox.EventKey && (event.Ch == 'd' || event.Key == termbox.KeyArrowRight) {
			rightRunCount++
			rightRunCount, position = moveRight(position.posX, position.posY, rightRunCount)
		} else if event.Type == termbox.EventKey && event.Ch == 'q' {
			return
		}
	}
}

func moveUp(posX int, posY int, runIndex int) (int, PlayerPosition) {
	termbox.SetCell(posX, posY-1, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
	position := PlayerPosition{posX, posY - 1}

	termbox.SetCell(posX, posY+1, ' ', termbox.ColorBlack, termbox.ColorDefault)
	termbox.SetCell(posX, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return runIndex, position
}

func moveDown(posX int, posY int, runIndex int) (int, PlayerPosition) {
	termbox.SetCell(posX, posY+1, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
	position := PlayerPosition{posX, posY + 1}

	termbox.SetCell(posX, posY-1, ' ', termbox.ColorBlack, termbox.ColorDefault)
	termbox.SetCell(posX, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return runIndex, position
}

func moveLeft(posX int, posY int, runIndex int) (int, PlayerPosition) {
	termbox.SetCell(posX-1, posY, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
	position := PlayerPosition{posX - 1, posY}

	termbox.SetCell(posX+1, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)
	termbox.SetCell(posX, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return runIndex, position
}
func moveRight(posX int, posY int, runIndex int) (int, PlayerPosition) {
	termbox.SetCell(posX+1, posY, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
	position := PlayerPosition{posX + 1, posY}

	termbox.SetCell(posX-1, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)
	termbox.SetCell(posX, posY, ' ', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return runIndex, position
}
