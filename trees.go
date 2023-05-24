package main

import (
	// Basic packages
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png" // Importing the PNG package to support loading PNG images

	"github.com/faiface/pixel"          // Importing the Pixel library
	"github.com/faiface/pixel/pixelgl"  // OpenGL from Pixel library
	"github.com/faiface/pixel/text"     // Text from pixel library
	"golang.org/x/image/font/basicfont" // Import basic fonts
)

// loadPicture loads an image from a file and returns a pixel.Picture object.
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// Declare the treeCountLabel variable outside the run function
var treeCountLabel *text.Text

// run is the main game loop where game logic is implemented.
func run() {
	// Window configuration
	cfg := pixelgl.WindowConfig{
		Title:  "Trees!",                 // Window title
		Bounds: pixel.R(0, 0, 1024, 768), // Window size
		VSync:  true,                     // Enable VSync (synchronizes frame rate with monitor refresh rate)
	}
	// Create a new window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Declare some variables
	var (
		windowSize       = pixel.V(1024, 768)     // Window size
		camPos           = windowSize.Scaled(0.5) // Camera position
		camSpeed         = 500.0                  // Camera speed
		camZoom          = 1.0                    // Initial camera zoom level
		minZoom          = 0.2                    // Minimum zoom level
		maxZoom          = 2.0                    // Maximum zoom level
		camZoomSpeed     = 1.2                    // Camera zoom speed
		treesPlanted     = 0                      // Number of trees planted
		initialFontScale = 2.0                    // Initial font scale
		frames           = 0                      // Frames counter initial value
		second           = time.Tick(time.Second) // Tick in seconds
	)

	// Define text fonts
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	// Text position at start
	basicTxt := text.New(pixel.V(windowSize.X/1.20-camPos.X, windowSize.Y/0.90-camPos.Y), basicAtlas)

	// Author variable and print text with fmt
	author := "Jordan"
	fmt.Fprintln(basicTxt, "Controls:")
	fmt.Fprintln(basicTxt, "- Arrows: Move Camera")
	fmt.Fprintln(basicTxt, "- Scroll: Zoom")
	fmt.Fprintln(basicTxt, "- Left Click: Plant Tree")
	fmt.Fprintln(basicTxt, "\nJust have fun planting trees!")
	fmt.Fprintf(basicTxt, "- %s", author)

	// Load the spritesheet image for trees
	spritesheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}

	// First batch (trees)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	// Prepare tree frames from the spritesheet (cut them from the spritesheet)
	var treesFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			treesFrames = append(treesFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	// Enable texture filtering (makes the image smoother) (keep commented)
	// win.SetSmooth(true)

	last := time.Now()

	// Game loop using a for loop
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		// Calculate the position of the tree count label
		countTxtPos := win.Bounds().Min.Add(pixel.V(5, win.Bounds().H()-25))
		countTxtPos = cam.Unproject(countTxtPos)

		// // Declare treeCountLabel variable
		treeCountLabel := text.New(countTxtPos, basicAtlas)

		// Draw tree count label
		fmt.Fprintf(treeCountLabel, "Trees planted: %d", treesPlanted)

		// Escape key to quit
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}

		// Mouse button left to plant tree
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			tree := pixel.NewSprite(spritesheet, treesFrames[rand.Intn(len(treesFrames))])
			mouse := cam.Unproject(win.MousePosition())
			// Draws a random tree from the spritesheet
			tree.Draw(batch, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
			treesPlanted++
		}

		// Arrow key to move camera left
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		// Arrow key to move camera right
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		// Arrow key to move camera down
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		// Arrow key to move camera up
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}

		// Adjust zoom level with mouse wheel
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		// Clamp the zoom level to stay within the specified limits
		camZoom = math.Max(minZoom, math.Min(maxZoom, camZoom))

		// Set the background color to grass green #4F8227
		win.Clear(pixel.RGB(0x4F, 0x82, 0x27).Scaled(1.0 / 255))
		// Draws images in batch 1
		batch.Draw(win)
		// Draw tuto text to screen
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))

		// Draw the treeCountLabel text
		treeCountLabel.Draw(win, pixel.IM.Scaled(treeCountLabel.Orig, initialFontScale/camZoom))

		// Update the game constantly
		win.Update()

		// Check FPS and put it in window frame
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

// Starts the program
func main() {
	pixelgl.Run(run) // Run the game loop defined in the run() function
}
