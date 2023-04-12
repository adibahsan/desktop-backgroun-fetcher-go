package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/getlantern/systray"
)

var (
	a fyne.App
	w fyne.Window
)

func main() {
	// Create a new Fyne app
	a = app.New()
	w = a.NewWindow("MyApp")

	// Create a label to display status
	statusLabel := widget.NewLabel("GET request not started")

	// Add the status label to the window
	w.SetContent(statusLabel)

	// Initialize systray
	systray.Run(onSystrayReady, onSystrayExit)
}

func onSystrayReady() {
	// Set systray icon
	// systray.SetIcon(getIcon())

	// Set systray tooltip
	systray.SetTooltip("MyApp")

	// Add systray menu item to open window
	mOpen := systray.AddMenuItem("Open", "Open MyApp")

	// Add systray menu item to exit
	mQuit := systray.AddMenuItem("Quit", "Quit MyApp")

	// Start a timer to make GET request every minute
	go func() {
		for {
			makeGETRequest()
			time.Sleep(10 * time.Second)
		}
	}()

	// Handle systray menu item click events
	for {
		select {
		case <-mOpen.ClickedCh:
			// Show window when "Open" menu item is clicked
			w.Show()
		case <-mQuit.ClickedCh:
			systray.Quit()
			os.Exit(0)
		}
	}
}

func onSystrayExit() {
	// Clean up resources when systray is exited
	w.Close()
	a.Quit()
}

func makeGETRequest() {
	// Make GET request
	resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Update status label
	statusLabel := w.Content().(*widget.Label)
	fmt.Println("Getting Success" + time.Now().Format("2006-01-02 15:02:45"))
	status := statusLabel.Text + "\nGET request successful" // Append to previous text
	statusLabel.SetText(status)

	// Play beep sound on successful GET request
	playBeep()
}

func playBeep() {
	// Load beep sound from file
	f, err := os.Open("beep.mp3") // Replace with your own beep sound file
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer streamer.Close()

	// Initialize speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Play the beep sound
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the beep sound to finish
	<-done

	// Show a notification while playing the beep sound
	notification := fyne.NewNotification("New Order", "You have a new Order")

	// notification.SetIcon(getIcon()) // Set the icon for the notification

	// Configure the notification to open a URL when clicked
	fyne.CurrentApp().SendNotification(notification)
}

func getIcon() []byte {
	// Return an icon image as byte array
	// Replace this with your own icon image
	return []byte{}
}
