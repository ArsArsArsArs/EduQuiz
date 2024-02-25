package services

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Language string `json:"lang"`
}

type Library struct {
	Name  string            `json:"name"`
	Cards []LibraryCardBase `json:"cards"`
}

type LibraryCardBase struct {
	//1 - QA, 2 - Text, 3 - Matching
	Type           int                 `json:"type"`
	QuestionAnswer LibraryCardQA       `json:"qa"`
	Text           LibraryCardText     `json:"txt"`
	Matching       LibraryCardMatching `json:"matching"`
}

type LibraryCardQA struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	//1 - Simple showing the answer. 2 - Choosing the right answer between wrong ones. 3 - Writing the right answer
	PresentationType int `json:"present"`
}

type LibraryCardText struct {
	Text string `json:"text"`
	//1 - Choose the words from the list, 2 - Type the words
	PresentationType int `json:"present"`
	//If true, the most words will be removed
	AgressiveMode bool `json:"agressive"`
}

type LibraryCardMatching struct {
	Items []MatchingItem `json:"matchingitems"`
}

type MatchingItem struct {
	FirstString  string `json:"firststr"`
	SecondString string `json:"secondstr"`
}

func IsDBFileExisting(path, fileName string) bool {
	if _, err := os.Stat(filepath.Join(path, fileName)); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func DBFileCreate(path, fileName string) error {
	strg, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}
	defer strg.Close()
	return nil
}

func LibraryFileCreate(path, name string) error {
	if !IsDBFileExisting(path, "libraries") {
		err := os.Mkdir(filepath.Join(path, "libraries"), 0777)
		if err != nil {
			return err
		}
	}
	strg, err := os.Create(filepath.Join(path, "libraries", name+".json"))
	if err != nil {
		return err
	}
	defer strg.Close()
	newLib := Library{
		Name: name,
	}
	b, err := json.Marshal(newLib)
	if err != nil {
		return err
	}
	_, err = strg.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func LibraryFileImport(path string, lib Library) error {
	if !IsDBFileExisting(path, "libraries") {
		err := os.Mkdir(filepath.Join(path, "libraries"), 0777)
		if err != nil {
			return err
		}
	}
	if IsDBFileExisting(filepath.Join(path, "libraries"), lib.Name+".json") {
		return errors.New("already exists")
	}
	strg, err := os.Create(filepath.Join(path, "libraries", lib.Name+".json"))
	if err != nil {
		return err
	}
	defer strg.Close()
	b, err := json.Marshal(lib)
	if err != nil {
		return err
	}
	_, err = strg.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func LibraryFileDelete(path, name string) error {
	err := os.Remove(filepath.Join(path, "libraries", name+".json"))
	if err != nil {
		return err
	}
	return nil
}

func LibraryFileEdit(path, oldName, newName string) error {
	err := os.Rename(filepath.Join(path, "libraries", oldName+".json"), filepath.Join(path, "libraries", newName+".json"))
	if err != nil {
		return err
	}
	return nil
}

func RetrieveConfig(path string) (Config, error) {
	if !IsDBFileExisting(path, "config.json") {
		err := DBFileCreate(path, "config.json")
		if err != nil {
			return Config{}, err
		}
	}
	body, err := os.ReadFile(filepath.Join(path, "config.json"))
	if err != nil {
		return Config{}, err
	}

	var conf Config
	err = json.Unmarshal(body, &conf)
	if err != nil {
		return Config{}, nil
	}
	return conf, nil
}

func UpdateConfig(path string, conf Config, value string, key any) (Config, error) {
	switch value {
	case "lang":
		copyConf := conf
		copyConf.Language = key.(string)
		err := updateJSONConfig(path, copyConf)
		if err != nil {
			return conf, err
		}
		return copyConf, nil
	default:
		return conf, nil
	}
}

func updateJSONConfig(path string, conf Config) error {
	confFile, err := os.OpenFile(filepath.Join(path, "config.json"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer confFile.Close()
	err = json.NewEncoder(confFile).Encode(conf)
	if err != nil {
		return err
	}
	return nil
}

func RetrieveLibraryFile(path, libraryName string) (Library, bool) {
	if !IsDBFileExisting(path, "libraries/"+libraryName+".json") {
		return Library{}, false
	}
	body, err := os.ReadFile(filepath.Join(path, "libraries", libraryName+".json"))
	if err != nil {
		log.Println(err)
		return Library{}, false
	}

	var lib Library
	err = json.Unmarshal(body, &lib)
	if err != nil {
		log.Println(err)
		return Library{}, false
	}
	return lib, true
}

func UpdateLibrary(path, libName string, lib Library, value string, key any) (Library, error) {
	switch value {
	case "name":
		copyLib := lib
		copyLib.Name = key.(string)
		err := UpdateLibraryFile(path, copyLib)
		if err != nil {
			return lib, err
		}
		return copyLib, nil
	default:
		return lib, nil
	}
}

func UpdateLibraryFile(path string, lib Library) error {
	confFile, err := os.OpenFile(filepath.Join(path, "libraries", lib.Name+".json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer confFile.Close()
	data, err := json.Marshal(lib)
	if err != nil {
		return err
	}
	_, err = confFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func GetAllLibraries(path string) ([]string, bool, error) {
	dirStr := filepath.Join(path, "libraries")
	if !IsDBFileExisting(dirStr, "") {
		return []string{}, false, nil
	}

	dir, err := os.Open(dirStr)
	if err != nil {
		return []string{}, false, err
	}
	defer dir.Close()

	filesStr, err := dir.Readdirnames(0)
	if err != nil {
		return []string{}, false, err
	}
	if len(filesStr) == 0 {
		return []string{}, false, nil
	}

	var result []string
	for _, fileStr := range filesStr {
		if !strings.HasSuffix(fileStr, ".json") {
			continue
		}
		result = append(result, fileStr)
	}
	return result, true, nil
}
