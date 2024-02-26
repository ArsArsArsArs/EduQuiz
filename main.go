package main

import (
	"EduQuiz/localization"
	"EduQuiz/services"
	"encoding/json"
	"errors"
	"image/color"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var a fyne.App
var masterWindow fyne.Window
var conf services.Config
var successLaunch bool = true

var rootUriPath string
var librariesList []string

var tappedGlobal *widget.Button

func main() {
	a = app.NewWithID("arsarsarsars.EduQuiz")
	masterWindow = a.NewWindow("EduQuiz")
	masterWindow.SetMaster()
	localization.PhrasesInit()
	var err error
	rootUriPath = a.Storage().RootURI().Path()
	conf, err = services.RetrieveConfig(rootUriPath)
	if err != nil {
		successLaunch = false
	}

	if !successLaunch {
		errorWarning := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "ErrorOpeningApp"))
		athing := widget.NewLabel(a.Storage().RootURI().Path())
		cnt := container.NewVBox(errorWarning, athing)
		masterWindow.SetContent(cnt)
		masterWindow.ShowAndRun()
	} else {
		showMainWindow()
		masterWindow.ShowAndRun()
	}
}

func newLibrary() {
	tlb := widget.NewToolbar(widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
		masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
		showMainWindow()
	}))
	title := widget.NewRichTextFromMarkdown("## " + localization.LoadLocalizedPhrase(conf.Language, "NewLibraryTitle"))
	desc := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryDesc"))
	desc.Wrapping = fyne.TextWrapWord
	sep := widget.NewSeparator()
	titleNameOption := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryNameTitle"))
	titleNameOption.TextStyle = fyne.TextStyle{Bold: true}
	titleNameEntry := widget.NewEntry()
	titleNameEntry.SetPlaceHolder("MyLibrary")
	titleNameEntry.Validator = func(s string) error {
		if n := utf8.RuneCountInString(s); n <= 0 || n > 25 {
			return errors.New(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryErrLength"))
		}
		if matched, _ := regexp.MatchString(`^[\wа-яА-ЯёЁ\-_]+$`, s); !matched || strings.Contains(s, "/") || strings.Contains(s, "\\") {
			return errors.New(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryErrSymbols"))
		}
		if doesExist := services.IsDBFileExisting(rootUriPath, "libraries/"+s+".json"); doesExist {
			return errors.New(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryErrExists"))
		}
		return nil
	}
	createButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "Create"), theme.ContentAddIcon(), func() {
		if err := titleNameEntry.Validate(); err != nil {
			showErrorDialog(err.Error())
			return
		}
		err := services.LibraryFileCreate(rootUriPath, titleNameEntry.Text)
		if err != nil {
			showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryCreateErr"))
			return
		}
		successDialog := dialog.NewInformation(localization.LoadLocalizedPhrase(conf.Language, "NewLibrarySuccessTitle"), localization.LoadLocalizedPhrase(conf.Language, "NewLibrarySuccessDesc"), masterWindow)
		successDialog.SetOnClosed(func() {
			masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
			showMainWindow()
		})
		successDialog.Show()
	})
	createButton.Alignment = widget.ButtonAlignCenter
	createButton.Disable()
	titleNameEntry.OnChanged = func(s string) {
		if err := titleNameEntry.Validate(); err == nil {
			createButton.Enable()
		} else {
			createButton.Disable()
		}
	}
	descElse := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryElseDesc"))
	descElse.Wrapping = fyne.TextWrapWord
	importButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "Import"), theme.MenuDropDownIcon(), func() {
		chooseFile := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				return
			}
			if uc == nil {
				return
			}
			defer uc.Close()
			body, errR := io.ReadAll(uc)
			if errR != nil {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "FailedImport"))
				return
			}
			var lib services.Library
			err = json.Unmarshal(body, &lib)
			if err != nil {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "WrongFile"))
				return
			}
			if lib.Name == "" {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "WrongFile"))
				return
			}
			for _, card := range lib.Cards {
				if card.Type <= 0 || card.Type > 3 {
					showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "WrongFile"))
					return
				}
				if card.QuestionAnswer.PresentationType > 3 || card.Text.PresentationType > 2 {
					showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "WrongFile"))
					return
				}
			}

			err = services.LibraryFileImport(rootUriPath, lib)
			if err != nil {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "FailedImport"))
				return
			}

			successDialog := dialog.NewInformation(localization.LoadLocalizedPhrase(conf.Language, "NewLibrarySuccessTitle"), localization.LoadLocalizedPhrase(conf.Language, "NewLibrarySuccessDesc"), masterWindow)
			successDialog.SetOnClosed(func() {
				masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
				showMainWindow()
			})
			successDialog.Show()
		}, masterWindow)
		chooseFile.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		chooseFile.Show()
	})
	importButton.Alignment = widget.ButtonAlignCenter
	newCnt := container.NewVBox(tlb, title, desc, sep, titleNameOption, titleNameEntry, createButton, sep, descElse, importButton)
	masterWindow.SetContent(newCnt)
}

func settingsMenu() {
	languageChoosing := widget.NewRadioGroup([]string{"English", "Русский"}, func(s string) {})
	switch conf.Language {
	case "ru":
		languageChoosing.SetSelected("Русский")
	default:
		languageChoosing.SetSelected("English")
	}
	settingsDialog := dialog.NewForm(localization.LoadLocalizedPhrase(conf.Language, "Settings"), localization.LoadLocalizedPhrase(conf.Language, "Apply"), localization.LoadLocalizedPhrase(conf.Language, "Cancel"), []*widget.FormItem{
		widget.NewFormItem(localization.LoadLocalizedPhrase(conf.Language, "Language"), languageChoosing),
	}, func(b bool) {
		if !b {
			return
		}

		var formattedLanguage string
		switch languageChoosing.Selected {
		case "Русский":
			formattedLanguage = "ru"
		default:
			formattedLanguage = "en"
		}
		if formattedLanguage != conf.Language {
			var err error
			conf, err = services.UpdateConfig(rootUriPath, conf, "lang", formattedLanguage)
			if err != nil {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "NoChangeLang"))
				return
			}
			warningDialog := dialog.NewInformation(localization.LoadLocalizedPhrase(conf.Language, "LanguageChanged"), localization.LoadLocalizedPhrase(conf.Language, "RestartAppLanguage"), masterWindow)
			warningDialog.Show()
		}
	}, masterWindow)
	settingsDialog.Show()
}

func showErrorDialog(errText string) {
	errDialog := dialog.NewInformation(localization.LoadLocalizedPhrase(conf.Language, "Error"), errText, masterWindow)
	errDialog.Show()
}

func getStrike(libraryName string) int {
	strike := a.Preferences().Int("strike_" + libraryName)
	strikeUpdated := a.Preferences().String("strikeUpdated_" + libraryName)
	if strikeUpdated == "" {
		return 0
	}
	strikeUpdatedFormat, err := time.Parse(time.DateOnly, strikeUpdated)
	if err != nil {
		return 0
	}
	strikeExpired := strikeUpdatedFormat.AddDate(0, 0, 2)
	if strikeExpired.Before(time.Now()) {
		return 0
	}
	return strike
}

func addStrike(libraryName string) {
	strike := getStrike(libraryName)
	strikeUpdated := a.Preferences().String("strikeUpdated_" + libraryName)
	strikeUpdatedFormat, errTime := time.Parse(time.DateOnly, strikeUpdated)
	if errTime != nil {
		a.Preferences().SetInt("strike_"+libraryName, strike+1)
		a.Preferences().SetString("strikeUpdated_"+libraryName, time.Now().Format(time.DateOnly))
		return
	}
	if strikeNewDay := strikeUpdatedFormat.AddDate(0, 0, 1); strikeNewDay.Before(time.Now()) {
		a.Preferences().SetInt("strike_"+libraryName, strike+1)
		a.Preferences().SetString("strikeUpdated_"+libraryName, time.Now().Format(time.DateOnly))
	}
}

func showMainWindow() {
	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), settingsMenu)
	contentAddButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), newLibrary)
	barCnt := container.NewVBox(widget.NewProgressBarInfinite())
	cnt := container.NewBorder(container.NewGridWithColumns(2, settingsButton, contentAddButton), nil, nil, nil, barCnt)
	masterWindow.SetContent(cnt)
	go func() {
		var isSomething bool
		var err error
		librariesList, isSomething, err = services.GetAllLibraries(rootUriPath)
		if err != nil || !isSomething {
			greeting := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "NothingToSee"))
			cnt.Remove(barCnt)
			cnt.Add(greeting)
			cnt.Refresh()
		} else {
			scrollList := widget.NewList(func() int { return len(librariesList) }, func() fyne.CanvasObject {
				b := widget.NewButton("temp", func() {})
				return b
			}, libraryButton)
			cnt.Remove(barCnt)
			cnt.Add(scrollList)
			cnt.Refresh()
		}
	}()
}

func libraryButton(i widget.ListItemID, o fyne.CanvasObject) {
	b := o.(*widget.Button)
	b.SetText(services.RemoveEndSymbols(librariesList[i], 5))
	b.OnTapped = func() {
		tappedGlobal = b
		showLibrary()
	}
}

func showLibrary() {
	lib, isGot := services.RetrieveLibraryFile(rootUriPath, tappedGlobal.Text)
	if !isGot {
		showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "OpenLibraryError"))
		return
	}
	masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
	tlb := widget.NewToolbar(widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
		masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
		showMainWindow()
	}))
	renameButton := widget.NewButtonWithIcon("", theme.NewThemedResource(resourceEditSvg), func() {
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("MyLibrary")
		nameEntry.Validator = func(s string) error {
			if n := utf8.RuneCountInString(s); n <= 0 || n > 25 {
				return errors.New(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryErrLength"))
			}
			if matched, _ := regexp.MatchString(`^[\wа-яА-ЯёЁ\-_]+$`, s); !matched || strings.Contains(s, "/") || strings.Contains(s, "\\") {
				return errors.New(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryErrSymbols"))
			}
			if doesExist := services.IsDBFileExisting(rootUriPath, "libraries/"+s+".json"); doesExist {
				return errors.New(localization.LoadLocalizedPhrase(conf.Language, "NewLibraryErrExists"))
			}
			return nil
		}
		renamingDialog := dialog.NewForm(localization.LoadLocalizedPhrase(conf.Language, "EditLibraryTitle"), localization.LoadLocalizedPhrase(conf.Language, "Apply"), localization.LoadLocalizedPhrase(conf.Language, "Cancel"), []*widget.FormItem{
			widget.NewFormItem("", nameEntry),
		}, func(choice bool) {
			if !choice {
				return
			}

			if err := nameEntry.Validate(); err != nil {
				showErrorDialog(err.Error())
				return
			}
			err := services.LibraryFileEdit(rootUriPath, tappedGlobal.Text, nameEntry.Text)
			if err != nil {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "EditLibraryTitleError"))
				return
			}
			_, err = services.UpdateLibrary(rootUriPath, nameEntry.Text, lib, "name", nameEntry.Text)
			if err != nil {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "EditLibraryTitleError"))
				return
			}
			strikeValue := getStrike(lib.Name)
			stikeUpdatedValue := a.Preferences().String("strikeUpdated_" + lib.Name)
			percentValue := a.Preferences().Int("percent_" + lib.Name)
			a.Preferences().RemoveValue("strike_" + lib.Name)
			a.Preferences().RemoveValue("strikeUpdated_" + lib.Name)
			a.Preferences().RemoveValue("percent_" + lib.Name)
			a.Preferences().SetInt("strike_"+nameEntry.Text, strikeValue)
			a.Preferences().SetString("strikeUpdated_"+nameEntry.Text, stikeUpdatedValue)
			a.Preferences().SetInt("percent_"+nameEntry.Text, percentValue)

			showMainWindow()
		}, masterWindow)
		renamingDialog.Show()
	})
	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		deletionDialog := dialog.NewConfirm(localization.LoadLocalizedPhrase(conf.Language, "DelLibraryTitle"), localization.LoadLocalizedPhrase(conf.Language, "USure"), func(choice bool) {
			if !choice {
				return
			}
			err := services.LibraryFileDelete(rootUriPath, tappedGlobal.Text)
			if err != nil {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "DeleteLibraryError"))
				return
			}
			a.Preferences().RemoveValue("strike_" + tappedGlobal.Text)
			a.Preferences().RemoveValue("strikeUpdated_" + tappedGlobal.Text)
			a.Preferences().RemoveValue("percent_" + tappedGlobal.Text)
			showMainWindow()
		}, masterWindow)
		deletionDialog.SetConfirmImportance(widget.DangerImportance)
		deletionDialog.SetDismissText(localization.LoadLocalizedPhrase(conf.Language, "No"))
		deletionDialog.SetConfirmText(localization.LoadLocalizedPhrase(conf.Language, "Yes"))
		deletionDialog.Show()
	})
	exportButton := widget.NewButtonWithIcon("", theme.UploadIcon(), func() {
		chooseLocation := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
			if err != nil {
				return
			}
			if uc == nil {
				return
			}
			defer uc.Close()
			body, errOS := os.ReadFile(filepath.Join(rootUriPath, "libraries", lib.Name+".json"))
			if errOS != nil {
				return
			}
			uc.Write(body)
		}, masterWindow)
		chooseLocation.SetFileName(lib.Name + ".json")
		chooseLocation.Show()
	})
	title := widget.NewRichTextFromMarkdown("## " + lib.Name)
	strikeText := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "DaysRow")+": "+strconv.Itoa(getStrike(lib.Name)), color.NRGBA{R: 255, G: 153, B: 0, A: 255})
	strikeText.TextStyle = fyne.TextStyle{Bold: true}
	strikeText.Alignment = fyne.TextAlignCenter
	totalCardsText := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "TotalCards")+": "+strconv.Itoa(len(lib.Cards)), color.NRGBA{R: 213, G: 38, B: 91, A: 255})
	totalCardsText.TextStyle = fyne.TextStyle{Bold: true}
	totalCardsText.Alignment = fyne.TextAlignCenter
	cardsList := widget.NewList(func() int { return len(lib.Cards) }, func() fyne.CanvasObject {
		return container.NewVBox(widget.NewLabel(""), widget.NewLabel(""), widget.NewButtonWithIcon("", theme.CancelIcon(), func() {}))
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		cardTypeLabel := o.(*fyne.Container).Objects[0].(*widget.Label)
		cardTypeLabel.Alignment = fyne.TextAlignCenter
		cardTypeLabel.TextStyle = fyne.TextStyle{Bold: true}
		cardPreviewLabel := o.(*fyne.Container).Objects[1].(*widget.Label)
		cardPreviewLabel.Wrapping = fyne.TextWrapWord
		deleteButton := o.(*fyne.Container).Objects[2].(*widget.Button)
		deleteButton.OnTapped = func() {
			deletionDialog := dialog.NewConfirm(localization.LoadLocalizedPhrase(conf.Language, "DelLibraryTitle"), localization.LoadLocalizedPhrase(conf.Language, "USure"), func(choice bool) {
				if !choice {
					return
				}
				copyLib := lib
				copyLib.Cards = append(copyLib.Cards[:i], copyLib.Cards[i+1:]...)
				err := services.UpdateLibraryFile(rootUriPath, copyLib)
				if err != nil {
					showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "DelCardError"))
					return
				}
				lib = copyLib
				masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
				showLibrary()
			}, masterWindow)
			deletionDialog.SetConfirmImportance(widget.DangerImportance)
			deletionDialog.SetDismissText(localization.LoadLocalizedPhrase(conf.Language, "No"))
			deletionDialog.SetConfirmText(localization.LoadLocalizedPhrase(conf.Language, "Yes"))
			deletionDialog.Show()
		}

		if lib.Cards[i].Type == 1 {
			cardTypeLabel.SetText(localization.LoadLocalizedPhrase(conf.Language, "QAcard"))
			cutText := services.CutString(lib.Cards[i].QuestionAnswer.Question, 25)
			if cutText == lib.Cards[i].QuestionAnswer.Question {
				cardPreviewLabel.SetText(lib.Cards[i].QuestionAnswer.Question)
			} else {
				cardPreviewLabel.SetText(cutText + "...")
			}
		} else if lib.Cards[i].Type == 2 {
			cardTypeLabel.SetText(localization.LoadLocalizedPhrase(conf.Language, "TEXTcard"))
			cutText := services.CutString(lib.Cards[i].Text.Text, 25)
			if cutText == lib.Cards[i].Text.Text {
				cardPreviewLabel.SetText(lib.Cards[i].Text.Text)
			} else {
				cardPreviewLabel.SetText(cutText + "...")
			}
		} else if lib.Cards[i].Type == 3 {
			cardTypeLabel.SetText(localization.LoadLocalizedPhrase(conf.Language, "MATCHcard"))
			cardPreviewLabel.SetText(lib.Cards[i].Matching.Items[0].FirstString + " — " + lib.Cards[i].Matching.Items[0].SecondString)
		}
	})
	startButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "Start"), theme.MenuExpandIcon(), startCardsButton)
	startButton.Importance = widget.SuccessImportance
	manageButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "CardsManagement"), theme.DocumentIcon(), func() {
		tlb := widget.NewToolbar(widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
			showLibrary()
		}))

		libDisplayName := widget.NewLabel(lib.Name)
		libDisplayName.TextStyle = fyne.TextStyle{Italic: true}
		chooseTheTypeLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "ChooseType"))
		cardCreationCnt := container.NewVBox(container.NewBorder(nil, nil, tlb, libDisplayName), chooseTheTypeLabel)
		typeChoice := widget.NewSelect([]string{localization.LoadLocalizedPhrase(conf.Language, "QAcard"), localization.LoadLocalizedPhrase(conf.Language, "TEXTcard"), localization.LoadLocalizedPhrase(conf.Language, "MATCHcard")}, func(s string) {
			switch s {
			case "Question & Answer", "Вопрос & Ответ":
				typeDesc := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "QAcardDesc"))
				typeDesc.Wrapping = fyne.TextWrapWord
				inputQuestionLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "inputQuestionReq"))
				inputQuestion := widget.NewEntry()
				inputQuestion.MultiLine = true
				inputQuestion.SetPlaceHolder(localization.LoadLocalizedPhrase(conf.Language, "jokeQuestion"))
				inputQuestion.Validator = func(s string) error {
					if n := utf8.RuneCountInString(s); n <= 0 || n > 1024 {
						return errors.New("1-1024")
					}
					return nil
				}
				inputAnswerLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "inputAnswerReq"))
				inputAnswer := widget.NewEntry()
				inputAnswer.MultiLine = true
				inputAnswer.SetPlaceHolder(localization.LoadLocalizedPhrase(conf.Language, "jokeAnswer"))
				inputAnswer.Validator = func(s string) error {
					if n := utf8.RuneCountInString(s); n <= 0 || n > 2048 {
						return errors.New("1-2048")
					}
					return nil
				}
				inputPTLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "inputPTReq") + "\n*" + localization.LoadLocalizedPhrase(conf.Language, "inputPTReqExplanation") + "*")
				inputPTLabel.Wrapping = fyne.TextWrapWord
				var inputPTnum int
				inputPT := widget.NewRadioGroup([]string{localization.LoadLocalizedPhrase(conf.Language, "QA_PTsimple"), localization.LoadLocalizedPhrase(conf.Language, "QA_PTchoice"), localization.LoadLocalizedPhrase(conf.Language, "QA_PTtyping")}, func(s string) {
					if s == localization.LoadLocalizedPhrase(conf.Language, "QA_PTsimple") {
						inputPTnum = 1
					} else if s == localization.LoadLocalizedPhrase(conf.Language, "QA_PTchoice") {
						inputPTnum = 2
					} else if s == localization.LoadLocalizedPhrase(conf.Language, "QA_PTtyping") {
						inputPTnum = 3
					}
				})
				inputPT.Required = true
				inputPT.SetSelected(localization.LoadLocalizedPhrase(conf.Language, "QA_PTsimple"))
				saveCardButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "Create"), theme.ContentAddIcon(), func() {
					if inputQuestion.Validate() != nil || inputAnswer.Validate() != nil {
						showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "incorrectFields"))
						return
					}
					libCopy := lib
					libCopy.Cards = append(lib.Cards, services.LibraryCardBase{
						Type: 1,
						QuestionAnswer: services.LibraryCardQA{
							Question:         inputQuestion.Text,
							Answer:           inputAnswer.Text,
							PresentationType: inputPTnum,
						},
					})
					err := services.UpdateLibraryFile(rootUriPath, libCopy)
					if err != nil {
						showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "AddCardError"))
						return
					}
					lib = libCopy
					successDialog := dialog.NewInformation(localization.LoadLocalizedPhrase(conf.Language, "Success"), localization.LoadLocalizedPhrase(conf.Language, "NewCardSuccess"), masterWindow)
					successDialog.SetOnClosed(func() {
						masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
						showLibrary()
					})
					successDialog.Show()
				})
				saveCardButton.Disable()
				inputQuestion.OnChanged = func(s string) {
					if inputQuestion.Validate() == nil && inputAnswer.Validate() == nil {
						if saveCardButton.Disabled() {
							saveCardButton.Enable()
						}
					} else {
						if !saveCardButton.Disabled() {
							saveCardButton.Disable()
						}
					}
				}
				inputAnswer.OnChanged = func(s string) {
					if inputQuestion.Validate() == nil && inputAnswer.Validate() == nil {
						if saveCardButton.Disabled() {
							saveCardButton.Enable()
						}
					} else {
						if !saveCardButton.Disabled() {
							saveCardButton.Disable()
						}
					}
				}
				for i, obj := range cardCreationCnt.Objects {
					if i <= 2 {
						continue
					}
					cardCreationCnt.Remove(obj)
				}
				cardCreationCnt.Add(typeDesc)
				cardCreationCnt.Add(inputQuestionLabel)
				cardCreationCnt.Add(inputQuestion)
				cardCreationCnt.Add(inputAnswerLabel)
				cardCreationCnt.Add(inputAnswer)
				cardCreationCnt.Add(inputPTLabel)
				cardCreationCnt.Add(inputPT)
				cardCreationCnt.Add(saveCardButton)
				cardCreationCnt.Refresh()

			case "Text", "Текст":
				typeDesc := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "TEXTcardDesc"))
				typeDesc.Wrapping = fyne.TextWrapWord
				inputTextLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "inputTextReq"))
				inputText := widget.NewEntry()
				inputText.MultiLine = true
				inputText.SetPlaceHolder(localization.LoadLocalizedPhrase(conf.Language, "jokeText"))
				inputText.Validator = func(s string) error {
					if n := utf8.RuneCountInString(s); n <= 0 || n > 2048 {
						return errors.New("1-2048")
					}
					return nil
				}
				inputPTLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "inputPTReq") + "\n*" + localization.LoadLocalizedPhrase(conf.Language, "inputPTReqExplanation") + "*")
				inputPTLabel.Wrapping = fyne.TextWrapWord
				var inputPTnum int
				inputPT := widget.NewRadioGroup([]string{localization.LoadLocalizedPhrase(conf.Language, "T_PTsimple"), localization.LoadLocalizedPhrase(conf.Language, "T_PTtyping")}, func(s string) {
					if s == localization.LoadLocalizedPhrase(conf.Language, "T_PTsimple") {
						inputPTnum = 1
					} else if s == localization.LoadLocalizedPhrase(conf.Language, "T_PTtyping") {
						inputPTnum = 2
					}
				})
				inputPT.Required = true
				inputPT.SetSelected(localization.LoadLocalizedPhrase(conf.Language, "T_PTsimple"))
				agroChoiceLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "agmLabel"))
				agroChoice := widget.NewCheck(localization.LoadLocalizedPhrase(conf.Language, "agmExplanation"), func(b bool) {})
				saveCardButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "Create"), theme.ContentAddIcon(), func() {
					if inputText.Validate() != nil {
						showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "incorrectFields"))
						return
					}
					libCopy := lib
					libCopy.Cards = append(lib.Cards, services.LibraryCardBase{
						Type: 2,
						Text: services.LibraryCardText{
							Text:             inputText.Text,
							PresentationType: inputPTnum,
							AgressiveMode:    agroChoice.Checked,
						},
					})
					err := services.UpdateLibraryFile(rootUriPath, libCopy)
					if err != nil {
						showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "AddCardError"))
						return
					}
					lib = libCopy
					successDialog := dialog.NewInformation(localization.LoadLocalizedPhrase(conf.Language, "Success"), localization.LoadLocalizedPhrase(conf.Language, "NewCardSuccess"), masterWindow)
					successDialog.SetOnClosed(func() {
						masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
						showLibrary()
					})
					successDialog.Show()
				})
				saveCardButton.Disable()
				inputText.OnChanged = func(s string) {
					if inputText.Validate() == nil {
						if saveCardButton.Disabled() {
							saveCardButton.Enable()
						}
					} else {
						if !saveCardButton.Disabled() {
							saveCardButton.Disable()
						}
					}
				}
				for i, obj := range cardCreationCnt.Objects {
					if i <= 2 {
						continue
					}
					cardCreationCnt.Remove(obj)
				}
				cardCreationCnt.Add(typeDesc)
				cardCreationCnt.Add(inputTextLabel)
				cardCreationCnt.Add(inputText)
				cardCreationCnt.Add(inputPTLabel)
				cardCreationCnt.Add(inputPT)
				cardCreationCnt.Add(agroChoiceLabel)
				cardCreationCnt.Add(agroChoice)
				cardCreationCnt.Add(saveCardButton)
				cardCreationCnt.Refresh()

			case "Matching", "Соответствие":
				typeDesc := widget.NewLabel(localization.LoadLocalizedPhrase(conf.Language, "MATCHcardDesc"))
				typeDesc.Wrapping = fyne.TextWrapWord
				inputPairsLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "InputPairsReq"))
				pairsContainer := container.NewVBox()
				plusFieldButton := widget.NewButtonWithIcon("", resourceEmojiPlusPng, func() {
					inputFirst := widget.NewEntry()
					inputFirst.SetPlaceHolder("1")
					inputSecond := widget.NewEntry()
					inputSecond.SetPlaceHolder("2")

					inputFirst.Validator = func(s string) error {
						if n := utf8.RuneCountInString(s); n <= 0 || n > 15 {
							return errors.New("1-15")
						}
						return nil
					}
					inputSecond.Validator = func(s string) error {
						if n := utf8.RuneCountInString(s); n <= 0 || n > 15 {
							return errors.New("1-15")
						}
						return nil
					}

					if len(pairsContainer.Objects) >= 10 {
						showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "TooManyPairs"))
						return
					}
					pairsContainer.Add(container.NewGridWithColumns(2, inputFirst, inputSecond))
					pairsContainer.Refresh()
				})
				minusFieldButton := widget.NewButtonWithIcon("", resourceEmojiMinusPng, func() {
					if len(pairsContainer.Objects) <= 1 {
						showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "TooFewPairs"))
						return
					}
					pairsContainer.Remove(pairsContainer.Objects[len(pairsContainer.Objects)-1])
					pairsContainer.Refresh()
				})
				inputFirst := widget.NewEntry()
				inputFirst.SetPlaceHolder("1")
				inputSecond := widget.NewEntry()
				inputSecond.SetPlaceHolder("2")

				inputFirst.Validator = func(s string) error {
					if n := utf8.RuneCountInString(s); n <= 0 || n > 15 {
						return errors.New("1-15")
					}
					return nil
				}
				inputSecond.Validator = func(s string) error {
					if n := utf8.RuneCountInString(s); n <= 0 || n > 15 {
						return errors.New("1-15")
					}
					return nil
				}

				pairsContainer.Add(container.NewGridWithColumns(2, inputFirst, inputSecond))
				saveCardButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "Create"), theme.ContentAddIcon(), func() {
					var inputMatching []services.MatchingItem
					for _, f := range pairsContainer.Objects {
						if f.(*fyne.Container).Objects[0].(*widget.Entry).Validate() != nil || f.(*fyne.Container).Objects[1].(*widget.Entry).Validate() != nil {
							showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "WrongPairError"))
							return
						} else {
							inputMatching = append(inputMatching, services.MatchingItem{
								FirstString:  f.(*fyne.Container).Objects[0].(*widget.Entry).Text,
								SecondString: f.(*fyne.Container).Objects[1].(*widget.Entry).Text,
							})
						}
					}

					libCopy := lib
					libCopy.Cards = append(lib.Cards, services.LibraryCardBase{
						Type: 3,
						Matching: services.LibraryCardMatching{
							Items: inputMatching,
						},
					})
					err := services.UpdateLibraryFile(rootUriPath, libCopy)
					if err != nil {
						showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "AddCardError"))
						return
					}
					lib = libCopy
					successDialog := dialog.NewInformation(localization.LoadLocalizedPhrase(conf.Language, "Success"), localization.LoadLocalizedPhrase(conf.Language, "NewCardSuccess"), masterWindow)
					successDialog.SetOnClosed(func() {
						masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
						showLibrary()
					})
					successDialog.Show()
				})
				for i, obj := range cardCreationCnt.Objects {
					if i <= 2 {
						continue
					}
					cardCreationCnt.Remove(obj)
				}
				cardCreationCnt.Add(typeDesc)
				cardCreationCnt.Add(inputPairsLabel)
				cardCreationCnt.Add(pairsContainer)
				cardCreationCnt.Add(container.NewGridWithColumns(2, plusFieldButton, minusFieldButton))
				cardCreationCnt.Add(saveCardButton)
				cardCreationCnt.Refresh()
			}
		})
		typeChoice.PlaceHolder = localization.LoadLocalizedPhrase(conf.Language, "Select")
		cardCreationCnt.Add(typeChoice)
		masterWindow.SetContent(cardCreationCnt)
	})
	manageButton.Importance = widget.HighImportance
	newCnt := container.NewBorder(container.NewVBox(tlb, container.NewGridWithColumns(3, renameButton, deleteButton, exportButton), title, container.NewGridWithColumns(2, strikeText, totalCardsText)), container.NewVBox(startButton, manageButton), nil, nil, widget.NewSeparator(), cardsList)
	masterWindow.SetContent(newCnt)
}

func startCardsButton() {
	lib, isGot := services.RetrieveLibraryFile(rootUriPath, tappedGlobal.Text)
	if !isGot {
		showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "OpenLibraryError"))
		return
	}
	if len(lib.Cards) <= 0 {
		showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "NoCardsError"))
		return
	}
	masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
	tlb := widget.NewToolbar(widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
		masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
		showLibrary()
	}))
	var previousPercent int = a.Preferences().IntWithFallback("percent_"+lib.Name, 20)
	almostReadyLabel := widget.NewRichTextFromMarkdown("## " + localization.LoadLocalizedPhrase(conf.Language, "AlmostReady") + "...")
	percentRequirementLabel := widget.NewRichTextFromMarkdown("#### " + localization.LoadLocalizedPhrase(conf.Language, "HowManyPercentReq") + "\n*" + localization.LoadLocalizedPhrase(conf.Language, "HowManyPercentDesc") + "*")
	percentRequirementLabel.Wrapping = fyne.TextWrapWord
	percentRequirement := widget.NewSlider(1, 100)
	percentRequirement.SetValue(float64(previousPercent))
	resultSliderLabel := widget.NewLabel(strconv.Itoa(previousPercent) + "%")
	resultSliderLabel.Alignment = fyne.TextAlignCenter
	resultSliderLabel.TextStyle = fyne.TextStyle{Bold: true}
	percentRequirement.OnChanged = func(f float64) {
		resultSliderLabel.SetText(strconv.Itoa(int(f)) + "%")
	}
	applyButton := widget.NewButtonWithIcon(localization.LoadLocalizedPhrase(conf.Language, "Apply"), theme.NewThemedResource(resourceCheckSvg), func() {
		masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
		a.Preferences().SetInt("percent_"+lib.Name, int(percentRequirement.Value))
		cardsAmountFloat := (float64(len(lib.Cards)) * percentRequirement.Value) / 100.0
		cardsAmount := int(math.Ceil(cardsAmountFloat))

		cardsToShow := make([]services.LibraryCardBase, 0, cardsAmount)
		copyLibCards := make([]services.LibraryCardBase, len(lib.Cards))
		copy(copyLibCards, lib.Cards)
		for i := 0; i < cardsAmount; i++ {
			index := rand.Intn(len(copyLibCards))
			cardsToShow = append(cardsToShow, copyLibCards[index])
			copyLibCards = append(copyLibCards[:index], copyLibCards[index+1:]...)
		}

		tlb := widget.NewToolbar(widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			masterWindow.SetContent(container.NewVBox(widget.NewProgressBarInfinite()))
			showLibrary()
		}))
		tasksProgress := widget.NewProgressBar()
		tasksProgress.Max = float64(cardsAmount)
		tasksProgress.TextFormatter = func() string {
			return localization.LoadLocalizedPhrase(conf.Language, "CardsLeft") + ": " + strconv.Itoa(cardsAmount-int(tasksProgress.Value))
		}
		cardCnt := container.NewVBox()
		var i int = 0
		goNextButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "Next"), func() {})
		goNextButton.OnTapped = func() {
			i += 1
			if i >= len(cardsToShow) {
				addStrike(lib.Name)
				showLibrary()
				return
			}
			tasksProgress.SetValue(tasksProgress.Value + 1)
			callACardShow(cardsToShow[i], cardCnt, goNextButton)
		}
		goNextButton.Disable()
		quizCnt := container.NewBorder(container.NewVBox(tlb, tasksProgress), goNextButton, nil, nil, cardCnt)

		masterWindow.SetContent(quizCnt)

		callACardShow(cardsToShow[i], cardCnt, goNextButton)
	})
	masterWindow.SetContent(container.NewVBox(tlb, almostReadyLabel, percentRequirementLabel, container.NewBorder(nil, nil, widget.NewLabel("1"), widget.NewLabel("100"), percentRequirement), resultSliderLabel, applyButton))
}

func callACardShow(card services.LibraryCardBase, cnt *fyne.Container, nextB *widget.Button) {
	switch card.Type {
	case 1:
		qaCardShow(card.QuestionAnswer, cnt, nextB)
	case 2:
		textCardShow(card.Text, cnt, nextB)
	case 3:
		matchingCardShow(card.Matching, cnt, nextB)
	}
}

func qaCardShow(data services.LibraryCardQA, cnt *fyne.Container, nextB *widget.Button) {
	nextB.Disable()
	cnt.RemoveAll()
	loadingBar := widget.NewProgressBarInfinite()
	cnt.Add(loadingBar)
	cnt.Refresh()
	switch data.PresentationType {
	case 1:
		questionLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "Question"))
		question := widget.NewRichTextFromMarkdown(data.Question)
		question.Wrapping = fyne.TextWrapWord
		answerLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "Answer"))
		showAnswerButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "ShowAnswer"), func() {})
		showAnswerButton.OnTapped = func() {
			answer := widget.NewRichTextFromMarkdown(data.Answer)
			answer.Wrapping = fyne.TextWrapWord
			cnt.Remove(showAnswerButton)
			cnt.Add(answer)
			cnt.Refresh()
			nextB.Enable()
		}
		cnt.Remove(loadingBar)
		cnt.Add(questionLabel)
		cnt.Add(question)
		cnt.Add(answerLabel)
		cnt.Add(showAnswerButton)
		cnt.Refresh()

	case 2:
		lib, isGot := services.RetrieveLibraryFile(rootUriPath, tappedGlobal.Text)
		if !isGot {
			showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "OpenLibraryError"))
			return
		}
		questionLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "Question"))
		question := widget.NewRichTextFromMarkdown(data.Question)
		question.Wrapping = fyne.TextWrapWord
		answerLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "Answer"))
		var answersList []string
		for _, card := range lib.Cards {
			if card.QuestionAnswer.PresentationType != 2 || card.QuestionAnswer.Answer == data.Answer {
				continue
			}
			answersList = append(answersList, card.QuestionAnswer.Answer)
		}
		rand.Shuffle(len(answersList), func(i, j int) {
			answersList[i], answersList[j] = answersList[j], answersList[i]
		})
		if len(answersList) > 3 {
			answersList = answersList[:3]
		}
		answersList = append(answersList, data.Answer)
		rand.Shuffle(len(answersList), func(i, j int) {
			answersList[i], answersList[j] = answersList[j], answersList[i]
		})
		answersChoosing := widget.NewSelect(answersList, func(s string) {})
		answersChoosing.OnChanged = func(s string) {
			for i, obj := range cnt.Objects {
				if i <= 3 {
					continue
				}
				cnt.Remove(obj)
			}
			if s == data.Answer {
				answersChoosing.Disable()
				successLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsRight"), color.NRGBA{R: 52, G: 201, B: 36, A: 255})
				successLabel.Alignment = fyne.TextAlignCenter
				fullAnswer := widget.NewLabel(data.Answer)
				fullAnswer.Wrapping = fyne.TextWrapWord
				cnt.Add(successLabel)
				cnt.Add(widget.NewSeparator())
				cnt.Add(fullAnswer)
				cnt.Refresh()
				nextB.Enable()
			} else {
				failLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsWrong"), color.NRGBA{R: 255, G: 43, B: 43, A: 255})
				showAnswerButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "ShowAnswer"), func() {})
				showAnswerButton.OnTapped = func() {
					answersChoosing.Disable()
					showAnswerButton.Disable()
					fullAnswer := widget.NewLabel(data.Answer)
					fullAnswer.Wrapping = fyne.TextWrapWord
					cnt.Add(widget.NewSeparator())
					cnt.Add(fullAnswer)
					cnt.Refresh()
					nextB.Enable()
				}
				cnt.Add(container.NewBorder(nil, nil, failLabel, showAnswerButton))
				cnt.Refresh()
			}
		}
		answersChoosing.PlaceHolder = localization.LoadLocalizedPhrase(conf.Language, "ChooseAnswer")
		cnt.Remove(loadingBar)
		cnt.Add(questionLabel)
		cnt.Add(question)
		cnt.Add(answerLabel)
		cnt.Add(answersChoosing)
		cnt.Refresh()
	case 3:
		questionLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "Question"))
		question := widget.NewRichTextFromMarkdown(data.Question)
		question.Wrapping = fyne.TextWrapWord
		answerLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "Answer"))
		inputField := widget.NewEntry()
		inputField.SetPlaceHolder(localization.LoadLocalizedPhrase(conf.Language, "WriteAnswer"))
		inputField.MultiLine = true
		checkButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "Check"), func() {})
		checkButton.OnTapped = func() {
			for i, obj := range cnt.Objects {
				if i <= 4 {
					continue
				}
				cnt.Remove(obj)
			}
			if strings.EqualFold(inputField.Text, data.Answer) {
				inputField.Disable()
				checkButton.Disable()
				successLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsRight"), color.NRGBA{R: 52, G: 201, B: 36, A: 255})
				successLabel.Alignment = fyne.TextAlignCenter
				fullAnswer := widget.NewLabel(data.Answer)
				fullAnswer.Wrapping = fyne.TextWrapWord
				cnt.Add(successLabel)
				cnt.Add(widget.NewSeparator())
				cnt.Add(fullAnswer)
				cnt.Refresh()
				nextB.Enable()
			} else {
				failLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsWrong"), color.NRGBA{R: 255, G: 43, B: 43, A: 255})
				showAnswerButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "ShowAnswer"), func() {})
				showAnswerButton.OnTapped = func() {
					inputField.Disable()
					checkButton.Disable()
					showAnswerButton.Disable()
					fullAnswer := widget.NewLabel(data.Answer)
					fullAnswer.Wrapping = fyne.TextWrapWord
					cnt.Add(widget.NewSeparator())
					cnt.Add(fullAnswer)
					cnt.Refresh()
					nextB.Enable()
				}
				cnt.Add(container.NewBorder(nil, nil, failLabel, showAnswerButton))
				cnt.Refresh()
			}
		}
		cnt.Remove(loadingBar)
		cnt.Add(questionLabel)
		cnt.Add(question)
		cnt.Add(answerLabel)
		cnt.Add(inputField)
		cnt.Add(checkButton)
		cnt.Refresh()
	}
}

func textCardShow(data services.LibraryCardText, cnt *fyne.Container, nextB *widget.Button) {
	nextB.Disable()
	cnt.RemoveAll()
	loadingBar := widget.NewProgressBarInfinite()
	cnt.Add(loadingBar)
	cnt.Refresh()
	switch data.PresentationType {
	case 1:
		var omittedWords []string
		var textAfter string
		if data.AgressiveMode {
			textAfter, omittedWords = services.RemoveRandomWords(data.Text, 50)
		} else {
			textAfter, omittedWords = services.RemoveRandomWords(data.Text, 20)
		}
		textLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "TEXTcard"))
		text := widget.NewLabel(textAfter)
		text.Wrapping = fyne.TextWrapWord
		fillInLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "FillInWords"))
		omittedPasting := container.NewHBox()
		scrollOmitted := container.NewHScroll(omittedPasting)
		copyOmitted := make([]string, len(omittedWords))
		copy(copyOmitted, omittedWords)
		rand.Shuffle(len(copyOmitted), func(i, j int) {
			copyOmitted[i], copyOmitted[j] = copyOmitted[j], copyOmitted[i]
		})
		for i := range copyOmitted {
			selectOmitted := widget.NewSelect(copyOmitted, func(s string) {})
			selectOmitted.PlaceHolder = strconv.Itoa(i + 1)
			omittedPasting.Add(selectOmitted)
		}
		checkButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "Check"), func() {})
		checkButton.OnTapped = func() {
			var answersSlice []string
			for _, obj := range omittedPasting.Objects {
				selectWidget := obj.(*widget.Select)
				if selectWidget.Selected == "" {
					showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "FillEverything"))
					return
				}
				answersSlice = append(answersSlice, selectWidget.Selected)
			}
			for i, obj := range cnt.Objects {
				if i <= 4 {
					continue
				}
				cnt.Remove(obj)
			}
			if services.CompareStringSlices(answersSlice, omittedWords) {
				omittedPasting.Hide()
				checkButton.Disable()
				successLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsRight"), color.NRGBA{R: 52, G: 201, B: 36, A: 255})
				successLabel.Alignment = fyne.TextAlignCenter
				fullAnswer := widget.NewLabel(data.Text)
				fullAnswer.Wrapping = fyne.TextWrapWord
				cnt.Add(successLabel)
				cnt.Add(widget.NewSeparator())
				cnt.Add(fullAnswer)
				cnt.Refresh()
				nextB.Enable()
			} else {
				failLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsWrong"), color.NRGBA{R: 255, G: 43, B: 43, A: 255})
				showAnswerButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "ShowAnswer"), func() {})
				showAnswerButton.OnTapped = func() {
					omittedPasting.Hide()
					checkButton.Disable()
					showAnswerButton.Disable()
					fullAnswer := widget.NewLabel(data.Text)
					fullAnswer.Wrapping = fyne.TextWrapWord
					cnt.Add(widget.NewSeparator())
					cnt.Add(fullAnswer)
					cnt.Refresh()
					nextB.Enable()
				}
				cnt.Add(container.NewBorder(nil, nil, failLabel, showAnswerButton))
				cnt.Refresh()
			}
		}
		cnt.Remove(loadingBar)
		cnt.Add(textLabel)
		cnt.Add(text)
		cnt.Add(fillInLabel)
		cnt.Add(scrollOmitted)
		cnt.Add(checkButton)
		cnt.Refresh()

	case 2:
		var omittedWords []string
		var textAfter string
		if data.AgressiveMode {
			textAfter, omittedWords = services.RemoveRandomWords(data.Text, 50)
		} else {
			textAfter, omittedWords = services.RemoveRandomWords(data.Text, 20)
		}
		textLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "TEXTcard"))
		text := widget.NewLabel(textAfter)
		text.Wrapping = fyne.TextWrapWord
		fillInLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "FillInWords"))
		omittedPasting := container.NewGridWithColumns(len(omittedWords))
		scrollOmitted := container.NewHScroll(omittedPasting)
		for i := range omittedWords {
			omittedWords[i] = services.OmitPunctuation(strings.ToLower(omittedWords[i]))
			entryOmitted := widget.NewEntry()
			entryOmitted.PlaceHolder = strconv.Itoa(i + 1)
			omittedPasting.Add(entryOmitted)
		}
		checkButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "Check"), func() {})
		checkButton.OnTapped = func() {
			var answersSlice []string
			for _, obj := range omittedPasting.Objects {
				entryWidget := obj.(*widget.Entry)
				if entryWidget.Text == "" {
					showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "FillEverything"))
					return
				}
				answersSlice = append(answersSlice, services.OmitPunctuation(strings.ToLower(entryWidget.Text)))
			}
			for i, obj := range cnt.Objects {
				if i <= 4 {
					continue
				}
				cnt.Remove(obj)
			}
			if services.CompareStringSlices(answersSlice, omittedWords) {
				omittedPasting.Hide()
				checkButton.Disable()
				successLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsRight"), color.NRGBA{R: 52, G: 201, B: 36, A: 255})
				successLabel.Alignment = fyne.TextAlignCenter
				fullAnswer := widget.NewLabel(data.Text)
				fullAnswer.Wrapping = fyne.TextWrapWord
				cnt.Add(successLabel)
				cnt.Add(widget.NewSeparator())
				cnt.Add(fullAnswer)
				cnt.Refresh()
				nextB.Enable()
			} else {
				failLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsWrong"), color.NRGBA{R: 255, G: 43, B: 43, A: 255})
				showAnswerButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "ShowAnswer"), func() {})
				showAnswerButton.OnTapped = func() {
					omittedPasting.Hide()
					checkButton.Disable()
					showAnswerButton.Disable()
					fullAnswer := widget.NewLabel(data.Text)
					fullAnswer.Wrapping = fyne.TextWrapWord
					cnt.Add(widget.NewSeparator())
					cnt.Add(fullAnswer)
					cnt.Refresh()
					nextB.Enable()
				}
				cnt.Add(container.NewBorder(nil, nil, failLabel, showAnswerButton))
				cnt.Refresh()
			}
		}
		cnt.Remove(loadingBar)
		cnt.Add(textLabel)
		cnt.Add(text)
		cnt.Add(fillInLabel)
		cnt.Add(scrollOmitted)
		cnt.Add(checkButton)
		cnt.Refresh()
	}
}

func matchingCardShow(data services.LibraryCardMatching, cnt *fyne.Container, nextB *widget.Button) {
	nextB.Disable()
	cnt.RemoveAll()
	loadingBar := widget.NewProgressBarInfinite()
	cnt.Add(loadingBar)
	cnt.Refresh()

	var firstRow, secondRow []string
	for _, item := range data.Items {
		firstRow = append(firstRow, item.FirstString)
		secondRow = append(secondRow, item.SecondString)
	}
	rand.Shuffle(len(firstRow), func(i, j int) {
		firstRow[i], firstRow[j] = firstRow[j], firstRow[i]
	})
	rand.Shuffle(len(secondRow), func(i, j int) {
		secondRow[i], secondRow[j] = secondRow[j], secondRow[i]
	})

	matchingLabel := widget.NewRichTextFromMarkdown("### " + localization.LoadLocalizedPhrase(conf.Language, "MATCHcard"))
	rows := container.NewGridWithRows(len(data.Items))
	for i := range data.Items {
		selectAnswer := widget.NewSelect(secondRow, func(s string) {})
		selectAnswer.PlaceHolder = "..."
		rows.Add(container.NewGridWithColumns(2, widget.NewLabel(firstRow[i]), selectAnswer))
	}
	checkButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "Check"), func() {})
	checkButton.OnTapped = func() {
		for i, obj := range cnt.Objects {
			if i <= 2 {
				continue
			}
			cnt.Remove(obj)
		}

		for _, row := range rows.Objects {
			secR := row.(*fyne.Container).Objects[1].(*widget.Select)
			if secR.Selected == "" {
				showErrorDialog(localization.LoadLocalizedPhrase(conf.Language, "FillEverything"))
				return
			}
		}

		for _, row := range rows.Objects {
			firR := row.(*fyne.Container).Objects[0].(*widget.Label)
			secR := row.(*fyne.Container).Objects[1].(*widget.Select)
			for _, item := range data.Items {
				if item.FirstString == firR.Text && item.SecondString != secR.Selected {
					failLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsWrong"), color.NRGBA{R: 255, G: 43, B: 43, A: 255})
					showAnswerButton := widget.NewButton(localization.LoadLocalizedPhrase(conf.Language, "ShowAnswer"), func() {})
					showAnswerButton.OnTapped = func() {
						rows.Hide()
						checkButton.Disable()
						showAnswerButton.Disable()
						var fullAnswStr string
						for _, item := range data.Items {
							fullAnswStr += item.FirstString + " — " + item.SecondString + "\n"
						}
						fullAnswer := widget.NewLabel(fullAnswStr)
						fullAnswer.Wrapping = fyne.TextWrapWord
						cnt.Add(widget.NewSeparator())
						cnt.Add(fullAnswer)
						cnt.Refresh()
						nextB.Enable()
					}
					cnt.Add(container.NewBorder(nil, nil, failLabel, showAnswerButton))
					cnt.Refresh()
					return
				}
			}
		}
		rows.Hide()
		checkButton.Disable()
		successLabel := canvas.NewText(localization.LoadLocalizedPhrase(conf.Language, "ThatsRight"), color.NRGBA{R: 52, G: 201, B: 36, A: 255})
		successLabel.Alignment = fyne.TextAlignCenter
		var fullAnswStr string
		for _, item := range data.Items {
			fullAnswStr += item.FirstString + " — " + item.SecondString + "\n"
		}
		fullAnswer := widget.NewLabel(fullAnswStr)
		fullAnswer.Wrapping = fyne.TextWrapWord
		cnt.Add(successLabel)
		cnt.Add(widget.NewSeparator())
		cnt.Add(fullAnswer)
		cnt.Refresh()
		nextB.Enable()
	}

	cnt.Remove(loadingBar)
	cnt.Add(matchingLabel)
	cnt.Add(rows)
	cnt.Add(checkButton)
	cnt.Refresh()
}
