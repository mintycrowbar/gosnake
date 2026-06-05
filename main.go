package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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

type Counter struct {
	score int
}

type terminalSize struct {
	width  int
	height int
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
	termSize := terminalSize{width, height}

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

	runGame(events, termSize)
}

func runGame(events <-chan termbox.Event, termSize terminalSize) {
	drawBorders()
	makeScoreLabel(termSize)

	var gameSpeed = 250 * time.Millisecond
	ticker := time.NewTicker(gameSpeed)
	defer ticker.Stop()

	snake := &Queue{}
	snake.Push(PlayerPosition{termSize.width/2 - 1, termSize.height / 2})

	snakeLength := 3
	direction := DirRight
	counter := &Counter{8}
	gameOver := false

	drawPointRandom(termSize)

GameLoop:
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
					termbox.Close()
					os.Exit(1)
				default:
					validKey = false
				}

				if validKey {
					previousDirection := direction
					direction = newDir

					snake, gameOver = changeDirection(changeDirParams{
						moveParams: moveParams{
							posX:        snake.GetHead().posX,
							posY:        snake.GetHead().posY,
							direction:   direction,
							snakeLength: snakeLength,
						},
						previousDirection: previousDirection,
					}, termSize, snake, counter, ticker, gameSpeed)

					if gameOver == true {
						break GameLoop
					}
				}
			}
		case <-ticker.C:
			snake, snakeLength, gameOver = movePlayer(moveParams{
				posX:        snake.GetHead().posX,
				posY:        snake.GetHead().posY,
				direction:   direction,
				snakeLength: snakeLength,
			}, snake, counter, termSize)

			if gameOver == true {
				break GameLoop
			}
		}
	}
	drawGameOver(events, termSize)
}

func movePlayer(params moveParams, moveQueue *Queue, counter *Counter, termSize terminalSize) (*Queue, int, bool) {
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

	// check if the next cell is a point or the walls or snake body
	if termbox.GetCell(newPos.posX, newPos.posY).Fg == termbox.ColorCyan {
		params.snakeLength++
		drawPointRandom(termSize)
		counter.Increment()
		updateScoreLabel(counter, termSize)
	} else if termbox.GetCell(newPos.posX, newPos.posY).Ch == '#' || termbox.GetCell(newPos.posX, newPos.posY).Ch == 'O' {
		return moveQueue, params.snakeLength, true
	}

	termbox.SetCell(newPos.posX, newPos.posY, 'O', termbox.ColorBlack, termbox.ColorDefault)

	if len(moveQueue.snakeBody) > params.snakeLength+1 {
		moveQueue.Pop()
	}
	termbox.SetCell(moveQueue.snakeBody[0].posX, moveQueue.snakeBody[0].posY, ' ', termbox.ColorBlack, termbox.ColorDefault)

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}

	return moveQueue, params.snakeLength, false
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

func changeDirection(params changeDirParams, termSize terminalSize, moveQueue *Queue, counter *Counter, ticker *time.Ticker, gameSpeed time.Duration) (*Queue, bool) {
	gameOver := false
	if params.direction != params.previousDirection {
		moveQueue, _, gameOver = movePlayer(moveParams{params.posX, params.posY, params.direction, params.snakeLength}, moveQueue, counter, termSize)
		drainChannel(ticker)
		ticker.Reset(gameSpeed)
	}
	return moveQueue, gameOver
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

func drawPointRandom(termSize terminalSize) {
	randomX := rand.Intn((termSize.width-3)-3+1) + 3
	randomY := rand.Intn((termSize.height-3)-3+1) + 3

	for termbox.GetCell(randomX, randomY).Ch != ' ' {
		randomX = rand.Intn((termSize.width-3)-3+1) + 3
		randomY = rand.Intn((termSize.height-3)-3+1) + 3
	}

	termbox.SetCell(randomX, randomY, 'P', termbox.ColorCyan, termbox.ColorDefault)
}

func makeScoreLabel(termSize terminalSize) {
	chars := []rune{'s', 'c', 'o', 'r', 'e', ':', ' ', '0'}
	for i := 0; i < len(chars); i++ {
		termbox.SetCell(i+1, termSize.height-1, chars[i], termbox.ColorLightGray, termbox.ColorBlack)
	}
}

/*func updateScoreLabel(counter *Counter, termSize terminalSize) {
	score := strconv.Itoa(counter.score)
	firstDigit := score[0:1]
	middleDigit := score[1:2]
	lastDigit := string(score[len(score)-1])

	if counter.score <= 9 {
		termbox.SetCell(8, termSize.height-1, rune(firstDigit[0]), termbox.ColorLightGray, termbox.ColorBlack)
	} else if counter.score >= 11 && counter.score <= 99 {
		termbox.SetCell(8, termSize.height-1, rune(firstDigit[0]), termbox.ColorLightGray, termbox.ColorBlack)
		termbox.SetCell(9, termSize.height-1, rune(middleDigit[1]), termbox.ColorLightGray, termbox.ColorBlack)
	} else if counter.score >= 100 {
		termbox.SetCell(8, termSize.height-1, rune(firstDigit[0]), termbox.ColorLightGray, termbox.ColorBlack)
		termbox.SetCell(9, termSize.height-1, rune(middleDigit[1]), termbox.ColorLightGray, termbox.ColorBlack)
		termbox.SetCell(9, termSize.height-1, rune(lastDigit[2]), termbox.ColorLightGray, termbox.ColorBlack)
	}
}*/

func updateScoreLabel(counter *Counter, termSize terminalSize) {
	score := strconv.Itoa(counter.score)

	if counter.score <= 9 {
		termbox.SetCell(8, termSize.height-1, rune(score[0]), termbox.ColorLightGray, termbox.ColorBlack)
	} else if counter.score >= 10 && counter.score <= 99 {
		termbox.SetCell(8, termSize.height-1, rune(score[0]), termbox.ColorLightGray, termbox.ColorBlack)
		termbox.SetCell(9, termSize.height-1, rune(score[1]), termbox.ColorLightGray, termbox.ColorBlack)
	} else if counter.score >= 100 {
		termbox.SetCell(8, termSize.height-1, rune(score[0]), termbox.ColorLightGray, termbox.ColorBlack)
		termbox.SetCell(9, termSize.height-1, rune(score[1]), termbox.ColorLightGray, termbox.ColorBlack)
		termbox.SetCell(9, termSize.height-1, rune(score[2]), termbox.ColorLightGray, termbox.ColorBlack)
	}
}

func (counter *Counter) Increment() {
	counter.score++
}

func drawGameOver(events <-chan termbox.Event, termSize terminalSize) {
	topMsg := "GAME OVER"
	bottomMsg := "Press R to restart or Q to quit"

	clearErr := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if clearErr != nil {
		log.Fatal(clearErr)
	}

	x := termSize.width / 2
	y := termSize.height / 2

	for i, ch := range topMsg {
		termbox.SetCell(-len(topMsg)/2+x+i, y-2, ch, termbox.ColorLightRed, termbox.ColorDefault)
	}

	for i, ch := range bottomMsg {
		termbox.SetCell(-len(bottomMsg)/2+x+i, y, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	flushErr := termbox.Flush()
	if flushErr != nil {
		log.Fatal(flushErr)
	}

	choice := termbox.PollEvent()

	switch choice.Ch {
	case 'q':
		return
	case 'r':
		runGame(events, termSize)
	}
}
