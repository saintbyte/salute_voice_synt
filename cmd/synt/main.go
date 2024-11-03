package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/aerth/playwav"
	"github.com/joho/godotenv"
	salutespeech_api "github.com/saintbyte/salute_speech_api"
	"github.com/skratchdot/open-golang/open"
	"io"
	"log"
	"os"
)

func main() {
	godotenv.Load()
	myApp := app.New()
	myWindow := myApp.NewWindow("Salute Text to Speech")
	openFileButton := widget.NewButton("Open File", func() {
		openFileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println("Error:", err)
				return
			}
			if reader == nil {
				log.Println("Canceled")
				return
			}
			log.Println("File opened:", reader.URI().String())
		}, myWindow)
		openFileDialog.Show()
	})

	saveFileButton := widget.NewButton("Save File", func() {
		saveFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				log.Println("Error:", err)
				return
			}
			if writer == nil {
				log.Println("Canceled")
				return
			}
			log.Println("File saved:", writer.URI().String())
			writer.Close()
		}, myWindow)
		saveFileDialog.Show()
	})

	toolbar := container.NewHBox(openFileButton, saveFileButton)

	input := widget.NewMultiLineEntry()

	voiceButton := widget.NewButton("Озвучить", func() {
		text := input.Text
		if text == "" {
			return
		}
		err := createVoiceFile(text, "output.wav")
		if err != nil {
			fmt.Println("Ошибка при создании голосового файла:", err)
			return
		}
		fmt.Println("Голосовой файл создан: output.wav")
	})

	playButton := widget.NewButton("Послушать", func() {
		err := playAudioFile("output.wav")
		if err != nil {
			fmt.Println("Ошибка при воспроизведении аудио:", err)
			return
		}
	})
	openFolderButton := widget.NewButton("Открыть папку", func() {
		open.Start("./")
	})
	content := container.New(layout.NewVBoxLayout(), toolbar, input, voiceButton, playButton, openFolderButton)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()
}
func createVoiceFile(text2speech_or_ssml, fileName string) error {
	log.Println(text2speech_or_ssml)
	s := salutespeech_api.NewSaluteSpeechApi()
	s.Debug = false
	s.AudioType = salutespeech_api.SaluteSpeechApi_OutputAudioTypeWAV16
	s.Voice = s.GetVoiceById("Pon_24000")
	data, err := s.Synthesize(text2speech_or_ssml)
	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	if data == nil {
		panic("data is null")
	}
	for {
		n, err := data.Read(buf)
		if err != nil {
			log.Println(err)
		}
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := file.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func playAudioFile(fileName string) error {
	_, err := playwav.FromFile(fileName)
	return err
}
