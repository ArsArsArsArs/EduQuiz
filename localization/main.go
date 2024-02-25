package localization

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var enLocalizer *i18n.Localizer
var ruLocalizer *i18n.Localizer

func PhrasesInit() {
	langBundle := i18n.NewBundle(language.English)
	langBundle.MustAddMessages(language.English, &i18n.Message{
		ID:    "ErrorOpeningApp",
		Other: "We got an error while opening the app. Restart may help.",
	},
		&i18n.Message{
			ID:    "NothingToSee",
			Other: "Nothing to see here... yet!",
		},
		&i18n.Message{
			ID:    "Settings",
			Other: "Settings",
		},
		&i18n.Message{
			ID:    "Apply",
			Other: "Apply",
		},
		&i18n.Message{
			ID:    "Cancel",
			Other: "Cancel",
		},
		&i18n.Message{
			ID:    "Language",
			Other: "Language",
		},
		&i18n.Message{
			ID:    "LanguageChanged",
			Other: "You've changed the language",
		},
		&i18n.Message{
			ID:    "RestartAppLanguage",
			Other: "You should restart the app",
		},
		&i18n.Message{
			ID:    "Error",
			Other: "Error",
		},
		&i18n.Message{
			ID:    "NoChangeLang",
			Other: "Failed to change the language",
		},
		&i18n.Message{
			ID:    "NewLibraryTitle",
			Other: "New Library",
		},
		&i18n.Message{
			ID:    "NewLibraryDesc",
			Other: "Create a new library so you can add your question cards there.",
		},
		&i18n.Message{
			ID:    "NewLibraryNameTitle",
			Other: "Library's name",
		},
		&i18n.Message{
			ID:    "NewLibraryErrLength",
			Other: "Must be 1-25 symbols",
		},
		&i18n.Message{
			ID:    "NewLibraryErrSymbols",
			Other: "Invalid symbols",
		},
		&i18n.Message{
			ID:    "NewLibraryErrExists",
			Other: "Already exists",
		},
		&i18n.Message{
			ID:    "Create",
			Other: "Create",
		},
		&i18n.Message{
			ID:    "NewLibraryElseDesc",
			Other: "Or else you can import an already existing EduQuiz library.",
		},
		&i18n.Message{
			ID:    "Import",
			Other: "Import",
		},
		&i18n.Message{
			ID:    "NewLibraryCreateErr",
			Other: "Failed to create a library",
		},
		&i18n.Message{
			ID:    "NewLibrarySuccessTitle",
			Other: "A library has been created",
		},
		&i18n.Message{
			ID:    "NewLibrarySuccessDesc",
			Other: "Access it at the main menu",
		},
		&i18n.Message{
			ID:    "OpenLibraryError",
			Other: "Failed to open the library",
		},
		&i18n.Message{
			ID:    "DaysRow",
			Other: "Days in a row",
		},
		&i18n.Message{
			ID:    "DelLibraryTitle",
			Other: "Library deletion",
		},
		&i18n.Message{
			ID:    "USure",
			Other: "Are you sure?",
		},
		&i18n.Message{
			ID:    "No",
			Other: "No",
		},
		&i18n.Message{
			ID:    "Yes",
			Other: "Yes",
		},
		&i18n.Message{
			ID:    "DeleteLibraryError",
			Other: "Failed to delete the library",
		},
		&i18n.Message{
			ID:    "EditLibraryTitle",
			Other: "Library renaming",
		},
		&i18n.Message{
			ID:    "EditLibraryTitleError",
			Other: "Failed to edit the name",
		},
		&i18n.Message{
			ID:    "Start",
			Other: "Start",
		},
		&i18n.Message{
			ID:    "CardsManagement",
			Other: "New card",
		},
		&i18n.Message{
			ID:    "QAcard",
			Other: "Question & Answer",
		},
		&i18n.Message{
			ID:    "TEXTcard",
			Other: "Text",
		},
		&i18n.Message{
			ID:    "ChooseType",
			Other: "Choose the card's type",
		},
		&i18n.Message{
			ID:    "Select",
			Other: "Select something...",
		},
		&i18n.Message{
			ID:    "QAcardDesc",
			Other: "You set a question and an answer to it. Thus, the system asks you the question and you should write or pick the right answer.",
		},
		&i18n.Message{
			ID:    "TEXTcardDesc",
			Other: "You set a text. Thus, the system gives the text to you with some words being omitted. You should write or pick the correct words so the text is full again.",
		},
		&i18n.Message{
			ID:    "inputQuestionReq",
			Other: "Enter a question",
		},
		&i18n.Message{
			ID:    "inputAnswerReq",
			Other: "Enter the answer",
		},
		&i18n.Message{
			ID:    "inputPTReq",
			Other: "Presentation type",
		},
		&i18n.Message{
			ID:    "inputPTReqExplanation",
			Other: "This parameter determines how you are going to answer while studying",
		},
		&i18n.Message{
			ID:    "jokeQuestion",
			Other: "Why did the chicken cross the road?",
		},
		&i18n.Message{
			ID:    "jokeAnswer",
			Other: "To get to the other side!",
		},
		&i18n.Message{
			ID:    "QA_PTsimple",
			Other: "Showing the answer",
		},
		&i18n.Message{
			ID:    "QA_PTchoice",
			Other: "Choosing the right answer between wrong ones",
		},
		&i18n.Message{
			ID:    "QA_PTtyping",
			Other: "Inputting the answer manually",
		},
		&i18n.Message{
			ID:    "inputTextReq",
			Other: "Enter a text",
		},
		&i18n.Message{
			ID:    "jokeText",
			Other: "Once upon a time...",
		},
		&i18n.Message{
			ID:    "T_PTsimple",
			Other: "Choosing words",
		},
		&i18n.Message{
			ID:    "T_PTtyping",
			Other: "Typing words manually",
		},
		&i18n.Message{
			ID:    "agmLabel",
			Other: "Agressive mode",
		},
		&i18n.Message{
			ID:    "agmExplanation",
			Other: "Omit more words",
		},
		&i18n.Message{
			ID:    "incorrectFields",
			Other: "Some fields may be filled incorrectly",
		},
		&i18n.Message{
			ID:    "AddCardError",
			Other: "Failed to create a card",
		},
		&i18n.Message{
			ID:    "Success",
			Other: "Success!",
		},
		&i18n.Message{
			ID:    "NewCardSuccess",
			Other: "A card has been created successfully",
		},
		&i18n.Message{
			ID:    "DelCardTitle",
			Other: "Card deletion",
		},
		&i18n.Message{
			ID:    "DelCardError",
			Other: "Failed to delete the card",
		},
		&i18n.Message{
			ID:    "TotalCards",
			Other: "Cards in total",
		},
		&i18n.Message{
			ID:    "MATCHcard",
			Other: "Matching",
		},
		&i18n.Message{
			ID:    "MATCHcardDesc",
			Other: "You set several pairs of words. Thus, the system gives you those words randomly. You should pair them right.",
		},
		&i18n.Message{
			ID:    "InputPairsReq",
			Other: "Enter pairs",
		},
		&i18n.Message{
			ID:    "TooManyPairs",
			Other: "Too many pairs!",
		},
		&i18n.Message{
			ID:    "TooFewPairs",
			Other: "Too few pairs!",
		},
		&i18n.Message{
			ID:    "WrongPairError",
			Other: "Your pairs are incorrect",
		},
		&i18n.Message{
			ID:    "NoCardsError",
			Other: "You must have at least one card",
		},
		&i18n.Message{
			ID:    "AlmostReady",
			Other: "Almost ready",
		},
		&i18n.Message{
			ID:    "AlmostReady",
			Other: "Almost ready",
		},
		&i18n.Message{
			ID:    "HowManyPercentReq",
			Other: "Cards percentage",
		},
		&i18n.Message{
			ID:    "HowManyPercentDesc",
			Other: "Specify the percentage of the total number of cards you want to proceed with at this moment.",
		},
		&i18n.Message{
			ID:    "CardsLeft",
			Other: "Cards left",
		},
		&i18n.Message{
			ID:    "Question",
			Other: "Question",
		},
		&i18n.Message{
			ID:    "Answer",
			Other: "Answer",
		},
		&i18n.Message{
			ID:    "ShowAnswer",
			Other: "Show the answer",
		},
		&i18n.Message{
			ID:    "Next",
			Other: "Next",
		},
		&i18n.Message{
			ID:    "ChooseAnswer",
			Other: "Choose the right answer",
		},
		&i18n.Message{
			ID:    "ThatsRight",
			Other: "Right!",
		},
		&i18n.Message{
			ID:    "ThatsWrong",
			Other: "Wrong! Try again",
		},
		&i18n.Message{
			ID:    "WriteAnswer",
			Other: "Write the right answer",
		},
		&i18n.Message{
			ID:    "Check",
			Other: "Check",
		},
		&i18n.Message{
			ID:    "FillInWords",
			Other: "Fill in the missing words",
		},
		&i18n.Message{
			ID:    "FillEverything",
			Other: "You should fill in all the words!",
		},
		&i18n.Message{
			ID:    "SelectedSlot",
			Other: "Selected slot",
		},
		&i18n.Message{
			ID:    "FailedImport",
			Other: "Failed importing",
		},
		&i18n.Message{
			ID:    "WrongFile",
			Other: "The file doesn't seem to be an EQ library",
		},
	)

	langBundle.MustAddMessages(language.Russian, &i18n.Message{
		ID:    "ErrorOpeningApp",
		Other: "Возникла ошибка при открытии приложения. Перезапуск может помочь.",
	},
		&i18n.Message{
			ID:    "NothingToSee",
			Other: "Здесь ничего нет... пока что!",
		},
		&i18n.Message{
			ID:    "Settings",
			Other: "Настройки",
		},
		&i18n.Message{
			ID:    "Apply",
			Other: "Применить",
		},
		&i18n.Message{
			ID:    "Cancel",
			Other: "Отмена",
		},
		&i18n.Message{
			ID:    "Language",
			Other: "Язык",
		},
		&i18n.Message{
			ID:    "LanguageChanged",
			Other: "Язык был изменён",
		},
		&i18n.Message{
			ID:    "RestartAppLanguage",
			Other: "Следует перезапустить приложение",
		},
		&i18n.Message{
			ID:    "Error",
			Other: "Ошибка",
		},
		&i18n.Message{
			ID:    "NoChangeLang",
			Other: "Не удалось поменять язык",
		},
		&i18n.Message{
			ID:    "NewLibraryTitle",
			Other: "Новая библиотека",
		},
		&i18n.Message{
			ID:    "NewLibraryDesc",
			Other: "Создайте новую библиотеку, дабы хранить там свои вопросительные билеты.",
		},
		&i18n.Message{
			ID:    "NewLibraryNameTitle",
			Other: "Название библиотеки",
		},
		&i18n.Message{
			ID:    "NewLibraryErrLength",
			Other: "Должно быть 1-25 символов",
		},
		&i18n.Message{
			ID:    "NewLibraryErrSymbols",
			Other: "Недопустимые символы",
		},
		&i18n.Message{
			ID:    "NewLibraryErrExists",
			Other: "Уже существует",
		},
		&i18n.Message{
			ID:    "Create",
			Other: "Создать",
		},
		&i18n.Message{
			ID:    "NewLibraryElseDesc",
			Other: "Или же Вы можете загрузить готовую библиотеку EduQuiz.",
		},
		&i18n.Message{
			ID:    "Import",
			Other: "Импортировать",
		},
		&i18n.Message{
			ID:    "NewLibraryCreateErr",
			Other: "Не вышло создать библиотеку",
		},
		&i18n.Message{
			ID:    "NewLibrarySuccessTitle",
			Other: "Библиотека создана",
		},
		&i18n.Message{
			ID:    "NewLibrarySuccessDesc",
			Other: "Она уже в главном меню",
		},
		&i18n.Message{
			ID:    "OpenLibraryError",
			Other: "Не вышло открыть библиотеку",
		},
		&i18n.Message{
			ID:    "DaysRow",
			Other: "Дней подряд",
		},
		&i18n.Message{
			ID:    "DelLibraryTitle",
			Other: "Удаление библиотеки",
		},
		&i18n.Message{
			ID:    "USure",
			Other: "Вы уверены?",
		},
		&i18n.Message{
			ID:    "No",
			Other: "Нет",
		},
		&i18n.Message{
			ID:    "Yes",
			Other: "Да",
		},
		&i18n.Message{
			ID:    "DeleteLibraryError",
			Other: "Не вышло удалить библиотеку",
		},
		&i18n.Message{
			ID:    "EditLibraryTitle",
			Other: "Переименование библиотеки",
		},
		&i18n.Message{
			ID:    "EditLibraryTitleError",
			Other: "Не вышло изменить имя",
		},
		&i18n.Message{
			ID:    "Start",
			Other: "Начать",
		},
		&i18n.Message{
			ID:    "CardsManagement",
			Other: "Новый билет",
		},
		&i18n.Message{
			ID:    "QAcard",
			Other: "Вопрос & Ответ",
		},
		&i18n.Message{
			ID:    "TEXTcard",
			Other: "Текст",
		},
		&i18n.Message{
			ID:    "ChooseType",
			Other: "Выберите тип билета",
		},
		&i18n.Message{
			ID:    "Select",
			Other: "Выберите что-нибудь...",
		},
		&i18n.Message{
			ID:    "QAcardDesc",
			Other: "Вы задаёте вопрос и ответ к нему. Таким образом, система спрашивает этот вопрос, и вы должны написать или выбрать правильный ответ.",
		},
		&i18n.Message{
			ID:    "TEXTcardDesc",
			Other: "Вы задаёте текст. Таким образом, система выдаёт этот текст вам, но некоторые слова пропущены. Вы должны написать или выбрать правильные слова, чтобы текст был полон вновь.",
		},
		&i18n.Message{
			ID:    "inputQuestionReq",
			Other: "Введите вопрос",
		},
		&i18n.Message{
			ID:    "inputAnswerReq",
			Other: "Введите ответ",
		},
		&i18n.Message{
			ID:    "inputPTReq",
			Other: "Тип представления",
		},
		&i18n.Message{
			ID:    "inputPTReqExplanation",
			Other: "Этот параметр определяет, каким образом Вы будете отвечать во время обучения",
		},
		&i18n.Message{
			ID:    "jokeQuestion",
			Other: "Сколько весит слон?",
		},
		&i18n.Message{
			ID:    "jokeAnswer",
			Other: "5 тонн!",
		},
		&i18n.Message{
			ID:    "QA_PTsimple",
			Other: "Показ ответа",
		},
		&i18n.Message{
			ID:    "QA_PTchoice",
			Other: "Выбор верного ответа среди неверных",
		},
		&i18n.Message{
			ID:    "QA_PTtyping",
			Other: "Ввод ответа вручную",
		},
		&i18n.Message{
			ID:    "inputTextReq",
			Other: "Введите текст",
		},
		&i18n.Message{
			ID:    "jokeText",
			Other: "Давным-давно, жили-были...",
		},
		&i18n.Message{
			ID:    "T_PTsimple",
			Other: "Выбор",
		},
		&i18n.Message{
			ID:    "T_PTtyping",
			Other: "Ручной ввод",
		},
		&i18n.Message{
			ID:    "agmLabel",
			Other: "Агрессивный режим",
		},
		&i18n.Message{
			ID:    "agmExplanation",
			Other: "Опускать больше слов",
		},
		&i18n.Message{
			ID:    "incorrectFields",
			Other: "Кажется, поля заполнены неправильно",
		},
		&i18n.Message{
			ID:    "AddCardError",
			Other: "Не вышло добавить билет",
		},
		&i18n.Message{
			ID:    "Success",
			Other: "Успешно!",
		},
		&i18n.Message{
			ID:    "NewCardSuccess",
			Other: "Билет создан",
		},
		&i18n.Message{
			ID:    "DelCardTitle",
			Other: "Удаление билета",
		},
		&i18n.Message{
			ID:    "DelCardError",
			Other: "Не вышло удалить билет",
		},
		&i18n.Message{
			ID:    "TotalCards",
			Other: "Всего билетов",
		},
		&i18n.Message{
			ID:    "MATCHcard",
			Other: "Соответствие",
		},
		&i18n.Message{
			ID:    "MATCHcardDesc",
			Other: "Вы задаёте пары слов. Таким образом, система выдаёт эти слова в случайном порядке. Вы должны воссоздать правильные пары.",
		},
		&i18n.Message{
			ID:    "InputPairsReq",
			Other: "Введите пары",
		},
		&i18n.Message{
			ID:    "TooManyPairs",
			Other: "Слишком много пар!",
		},
		&i18n.Message{
			ID:    "TooFewPairs",
			Other: "Слишком мало пар!",
		},
		&i18n.Message{
			ID:    "WrongPairError",
			Other: "Ваши пары некорректны",
		},
		&i18n.Message{
			ID:    "NoCardsError",
			Other: "Надо иметь хотя бы один билет",
		},
		&i18n.Message{
			ID:    "AlmostReady",
			Other: "Почти готово",
		},
		&i18n.Message{
			ID:    "HowManyPercentReq",
			Other: "Процент билетов",
		},
		&i18n.Message{
			ID:    "HowManyPercentDesc",
			Other: "Укажите процент от общего количества билетов, который вы хотите пройти в данный момент.",
		},
		&i18n.Message{
			ID:    "CardsLeft",
			Other: "Билетов осталось",
		},
		&i18n.Message{
			ID:    "Question",
			Other: "Вопрос",
		},
		&i18n.Message{
			ID:    "Answer",
			Other: "Ответ",
		},
		&i18n.Message{
			ID:    "ShowAnswer",
			Other: "Показать ответ",
		},
		&i18n.Message{
			ID:    "Next",
			Other: "Следующее",
		},
		&i18n.Message{
			ID:    "ChooseAnswer",
			Other: "Выберите правильный ответ",
		},
		&i18n.Message{
			ID:    "ThatsRight",
			Other: "Правильно!",
		},
		&i18n.Message{
			ID:    "ThatsWrong",
			Other: "Неправильно! Попробуйте ещё раз",
		},
		&i18n.Message{
			ID:    "WriteAnswer",
			Other: "Впишите правильный ответ",
		},
		&i18n.Message{
			ID:    "Check",
			Other: "Проверить",
		},
		&i18n.Message{
			ID:    "FillInWords",
			Other: "Вставьте пропущенные слова",
		},
		&i18n.Message{
			ID:    "FillEverything",
			Other: "Необходимо вставить все слова!",
		},
		&i18n.Message{
			ID:    "SelectedSlot",
			Other: "Выбранный слот",
		},
		&i18n.Message{
			ID:    "FailedImport",
			Other: "Импорт не удался",
		},
		&i18n.Message{
			ID:    "WrongFile",
			Other: "Файл не похож на библиотеку EQ",
		},
	)

	enLocalizer = i18n.NewLocalizer(langBundle, "en")
	ruLocalizer = i18n.NewLocalizer(langBundle, "ru")
}

func LoadLocalizedPhrase(lang, phrase string) string {
	var output string
	switch lang {
	case "ru":
		output, _ = ruLocalizer.Localize(&i18n.LocalizeConfig{
			MessageID: phrase,
		})
	default:
		output, _ = enLocalizer.Localize(&i18n.LocalizeConfig{
			MessageID: phrase,
		})
	}
	return output
}
