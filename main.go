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
	posX        int
	posY        int
	direction   Direction
	snakeLength int
}

type changeDirParams struct {
	moveParams
	previousDirection Direction
}

type Queue struct {
	snakeBody []PlayerPosition
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

	snake := &Queue{}
	snake.Push(PlayerPosition{width/2 - 1, height / 2})
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

					snake = changeDirection(changeDirParams{
						moveParams: moveParams{
							posX:        snake.GetHead().posX,
							posY:        snake.GetHead().posY,
							direction:   direction,
							snakeLength: snakeLength,
						},
						previousDirection: previousDirection,
					}, snake, ticker, gameSpeed)
				}
			}
		case <-ticker.C:
			snake, snakeLength = movePlayer(moveParams{
				posX:        snake.GetHead().posX,
				posY:        snake.GetHead().posY,
				direction:   direction,
				snakeLength: snakeLength,
			}, snake)
		}
	}
}

func movePlayer(params moveParams, moveQueue *Queue) (*Queue, int) {
	// dx and dy represent the head movement
	var dx, dy int

	switch params.direction {
	case DirUp:
		dy = -1
	case DirDown:
		dy = 1
	case DirLeft:
		dx = -1
	case DirRight:
		dx = 1
	}

	// set new head location
	newPos := PlayerPosition{params.posX + dx, params.posY + dy}
	moveQueue.Push(newPos)
	termbox.SetCell(newPos.posX, newPos.posY, 'O', termbox.ColorBlack, termbox.ColorDefault)

	if len(moveQueue.snakeBody) > params.snakeLength+1 {
		moveQueue.Pop()
	}
	termbox.SetCell(moveQueue.snakeBody[0].posX, moveQueue.snakeBody[0].posY, ' ', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
	return moveQueue, params.snakeLength
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

func changeDirection(params changeDirParams, moveQueue *Queue, ticker *time.Ticker, gameSpeed time.Duration) *Queue {
	if params.direction != params.previousDirection {
		moveQueue, _ = movePlayer(moveParams{params.posX, params.posY, params.direction, params.snakeLength}, moveQueue)
		drainChannel(ticker)
		ticker.Reset(gameSpeed)
	}
	return moveQueue
}

// add an element to the end of the original slice
func (queue *Queue) Push(position PlayerPosition) {
	queue.snakeBody = append(queue.snakeBody, position)
}

// remove the first element of the slice
func (queue *Queue) Pop() {
	if len(queue.snakeBody) > 0 {
		queue.snakeBody = queue.snakeBody[1:]
	}
}

// get head coords
func (queue *Queue) GetHead() PlayerPosition {
	if len(queue.snakeBody) == 0 {
		return PlayerPosition{}
	}
	return queue.snakeBody[len(queue.snakeBody)-1]
}
