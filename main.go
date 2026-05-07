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

	upRunCount := 0
	for {
		event := termbox.PollEvent()
		startingPosition := PlayerPosition{width / 2, height / 2}

		if event.Type == termbox.EventKey && event.Ch == 'w' || event.Key == termbox.KeyArrowUp {
			upRunCount++
			moveUp(startingPosition.posX, startingPosition.posY, upRunCount)
		} else if event.Type == termbox.EventKey && event.Ch == 's' || event.Key == termbox.KeyArrowDown {
			moveDown()
		} else if event.Type == termbox.EventKey && event.Ch == 'a' || event.Key == termbox.KeyArrowLeft {
			moveLeft()
		} else if event.Type == termbox.EventKey && event.Ch == 'd' || event.Key == termbox.KeyArrowRight {
			moveRight()
		} else if event.Type == termbox.EventKey && event.Ch == 'q' {
			return
		}
	}
}

func moveUp(posX int, posY int, runIndex int) int {
	termbox.SetCell(posX, posY-runIndex, 'ඞ', termbox.ColorBlack, termbox.ColorDefault)
	termbox.SetCell(posX, posY-runIndex+1, ' ', termbox.ColorBlack, termbox.ColorDefault)
	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return runIndex
}

func moveDown() {
	// TODO
}

func moveLeft() {
	// TODO
}

func moveRight() {
	// TODO
}
