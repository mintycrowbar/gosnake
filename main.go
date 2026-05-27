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

type moveParams struct {
	posX              int
	posY              int
	direction         string
	previousDirection string
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
	direction := "right"
	for {
		select {
		case event := <-events:
			if event.Type == termbox.EventKey {
				var newDir string
				switch {
				case event.Ch == 'w' || event.Key == termbox.KeyArrowUp:
					newDir = "up"
				case event.Ch == 's' || event.Key == termbox.KeyArrowDown:
					newDir = "down"
				case event.Ch == 'a' || event.Key == termbox.KeyArrowLeft:
					newDir = "left"
				case event.Ch == 'd' || event.Key == termbox.KeyArrowRight:
					newDir = "right"
				case event.Ch == 'q':
					return
				}

				if newDir != "" {
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
	position := PlayerPosition{}

	// TODO – handle turns: when the player turns, compare the amount of cells before and after the angle, meaning when the head reaches n tiles ahead, remove that n tiles from the tail
	// TODO – keep track of where the tail is and count the number of cells between the tail and the head, for example N. Then, every time the movePlayer function is called, compare
	switch params.direction {
	case "up":
		// clear the tail cell
		if termbox.GetCell(params.posX, params.posY+params.snakeLength-1).Ch == 'O' {
			termbox.SetCell(params.posX, params.posY+params.snakeLength-1, ' ', termbox.ColorBlack, termbox.ColorDefault)
		}
		//if termbox.GetCell(params.posX, params.posY-params.snakeLength-2).Ch == ' ' && termbox.GetCell(params.posX-1, params.pos-params.snakeLength-1)
		// draw the head a cell one cell ahead of where it currently is
		termbox.SetCell(params.posX, params.posY-1, 'O', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{params.posX, params.posY - 1}
	case "down":
		if termbox.GetCell(params.posX, params.posY-params.snakeLength+1).Ch == 'O' {
			termbox.SetCell(params.posX, params.posY-params.snakeLength+1, ' ', termbox.ColorBlack, termbox.ColorDefault)
		}
		termbox.SetCell(params.posX, params.posY+1, 'O', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{params.posX, params.posY + 1}
	case "left":
		if termbox.GetCell(params.posX+params.snakeLength-1, params.posY).Ch == 'O' {
			termbox.SetCell(params.posX+params.snakeLength-1, params.posY, ' ', termbox.ColorBlack, termbox.ColorDefault)
		}
		termbox.SetCell(params.posX-1, params.posY, 'O', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{params.posX - 1, params.posY}
	case "right":
		if termbox.GetCell(params.posX-params.snakeLength+1, params.posY).Ch == 'O' {
			termbox.SetCell(params.posX-params.snakeLength+1, params.posY, ' ', termbox.ColorBlack, termbox.ColorDefault)
		}
		termbox.SetCell(params.posX+1, params.posY, 'O', termbox.ColorBlack, termbox.ColorDefault)
		position = PlayerPosition{params.posX + 1, params.posY}
	}

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return position
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

func changeDirection(direction string, previousDirection string, position PlayerPosition, ticker *time.Ticker, gameSpeed time.Duration) (newPosition PlayerPosition) {
	if direction != previousDirection {
		position = movePlayer(moveParams{posX: position.posX, posY: position.posY, direction: direction, previousDirection: previousDirection})
		drainChannel(ticker)
		ticker.Reset(gameSpeed)
	}
	return position
}
