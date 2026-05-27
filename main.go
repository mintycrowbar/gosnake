package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type PlayerPosition struct {
	posX int
	posY int
}

type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

type moveParams struct {
	posX              int
	posY              int
	direction         Direction
	previousDirection Direction
	snakeLength       int
}

type Logger struct {
	ch chan string
}

func NewLogger(path string) (Logger, error) {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return Logger{}, err
	}

	logger := Logger{
		ch: make(chan string, 128),
	}

	go func() {
		defer func(file *os.File) {
			fileCloseError := file.Close()
			if fileCloseError != nil {
				log.Fatalf("error while closing file: %v", fileCloseError)
			}
		}(file)

		for msg := range logger.ch {
			timestamp := time.Now().Format("15:04:05")
			line := fmt.Sprintf("[%s]: %s\n", timestamp, msg)
			_, writeError := file.WriteString(line)
			if writeError != nil {
				log.Fatalf("Failed to write to file: %v", writeError)
			}
		}
	}()

	return logger, nil
}

func (l Logger) Info(msg string) {
	select {
	case l.ch <- msg:
		// message queued successfully
	default:
		// drop message because channel is full
	}
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	width, height := termbox.Size()
	drawBorders()
	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}

	position := PlayerPosition{width/2 - 1, height / 2}

	events := make(chan termbox.Event, 20)
	go func() {
		for {
			event := termbox.PollEvent()
			events <- event
		}
	}()

	var gameSpeed = 250 * time.Millisecond
	ticker := time.NewTicker(gameSpeed)
	defer ticker.Stop()

	snakeLength := 3
	direction := DirRight
	for {
		select {
		case event := <-events:
			if event.Type == termbox.EventKey {
				var newDir Direction
				validKey := true
				switch {
				case event.Ch == 'w' || event.Key == termbox.KeyArrowUp:
					newDir = DirUp
				case event.Ch == 's' || event.Key == termbox.KeyArrowDown:
					newDir = DirDown
				case event.Ch == 'a' || event.Key == termbox.KeyArrowLeft:
					newDir = DirLeft
				case event.Ch == 'd' || event.Key == termbox.KeyArrowRight:
					newDir = DirRight
				case event.Ch == 'q':
					return
				default:
					validKey = false
				}

				if validKey {
					previousDirection := direction
					direction = newDir
					position = changeDirection(direction, previousDirection, position, ticker, gameSpeed)
				}
			}
		case <-ticker.C:
			position = movePlayer(moveParams{posX: position.posX, posY: position.posY, direction: direction, snakeLength: snakeLength})
		}
	}
}

func movePlayer(params moveParams) PlayerPosition {
	// dx and dy represent the head movement
	// tx and ty represent where the tail was relative to the head
	var dx, dy, tx, ty int

	switch params.direction {
	case DirUp:
		dy, ty = -1, params.snakeLength-1
	case DirDown:
		dy, ty = 1, -params.snakeLength+1
	case DirLeft:
		dx, tx = -1, params.snakeLength-1
	case DirRight:
		dx, tx = 1, -params.snakeLength+1
	}

	// remove the previous tail location
	tailX, tailY := params.posX+tx, params.posY+ty
	if termbox.GetCell(tailX, tailY).Ch == 'O' {
		termbox.SetCell(tailX, tailY, ' ', termbox.ColorBlack, termbox.ColorDefault)
	}

	// set new head location
	head := PlayerPosition{params.posX + dx, params.posY + dy}
	termbox.SetCell(head.posX, head.posY, 'O', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return head
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

func drainChannel(ticker *time.Ticker) {
L:
	for {
		select {
		case <-ticker.C:
		default:
			break L
		}
	}
}

func changeDirection(direction Direction, previousDirection Direction, position PlayerPosition, ticker *time.Ticker, gameSpeed time.Duration) (newPosition PlayerPosition) {
	if direction != previousDirection {
		position = movePlayer(moveParams{posX: position.posX, posY: position.posY, direction: direction, previousDirection: previousDirection})
		drainChannel(ticker)
		ticker.Reset(gameSpeed)
	}
	return position
}
