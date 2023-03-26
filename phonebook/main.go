package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"time"
)

const CustomTimeFormat = time.Kitchen + " 02 Jan 2006"

type UserData struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess time.Time
}

var data = map[string]UserData{}

func main() {
	args := os.Args
	if len(args) == 1 {
		exe := path.Base(args[0])
		fmt.Printf("Usage: %s search|list|insert|delete <arguments>\n", exe)
		return
	}
	err := populateData()
	if err != nil {
		fmt.Println("Unable to read the data file", err)
		return
	}

	switch args[1] {
	case "search":
		if len(args) != 3 {
			fmt.Println("Usage: search TelephoneNo.")
			return
		}
		result := search(args[2])
		if result == nil {
			fmt.Println("Entry not found:", args[2])
			return
		}
		fmt.Println(*result)
	case "list":
		err := list()
		if err != nil {
			fmt.Println("Unable to list entries", err)
		}
	case "insert":
		if len(args) != 5 {
			fmt.Println("Usage: insert Name Surname TelephoneNo.")
			return
		}
		err := insert(args[2], args[3], args[4])
		if err != nil {
			fmt.Println("Unable to insert entry", err)
		} else {
			fmt.Println("Inserted successfully!")
		}
	case "delete":
		if len(args) != 3 {
			fmt.Println("Usage: delete TelephoneNo.")
			return
		}
		err := deleteEntry(args[2])
		if err != nil {
			fmt.Println("Unable to delete entry", args[2])
		} else {
			fmt.Println("Entry delete successfully!")
		}
	default:
		fmt.Println("Not a valid option")
	}
}

func search(key string) *UserData {
	userData, ok := data[key]
	if !ok {
		fmt.Printf("Key %s does not exist.", key)
		return nil
	}
	updatedUserData := userData
	updatedUserData.LastAccess = time.Now()
	data[key] = updatedUserData
	err := syncPhoneBook()
	if err != nil {
		fmt.Println("unable to update phone book", err)
	}
	return &userData
}

func list() error {
	for k, v := range data {
		fmt.Println(v)
		v.LastAccess = time.Now()
		data[k] = v
	}
	err := syncPhoneBook()
	if err != nil {
		fmt.Println("unable to update phone book", err)
	}
	return err
}

func insert(name, surname, tel string) error {
	if _, ok := data[tel]; ok {
		return fmt.Errorf("input with telephone no. %s already exists", tel)
	}
	_, err := os.Stat("./data.csv")
	if err != nil {
		fmt.Println("File not found")
		return err
	}

	f, err := os.OpenFile("./data.csv", os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	userData := UserData{Name: name, Surname: surname, Tel: tel, LastAccess: time.Now()}
	data[tel] = userData
	writer := csv.NewWriter(f)
	defer writer.Flush()
	err = writer.Write([]string{userData.Name, userData.Surname, userData.Tel, userData.LastAccess.Format(CustomTimeFormat)})
	if err != nil {
		return err
	}
	return err
}

func deleteEntry(key string) error {
	if _, ok := data[key]; !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	delete(data, key)
	err := syncPhoneBook()
	if err != nil {
		return err
	}
	return err
}

func syncPhoneBook() error {
	_, err := os.Stat("./data.csv")
	if err != nil {
		fmt.Println("File not found")
		return err
	}
	f, err := os.OpenFile("./data.csv", os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error opening file")
		return err
	}
	defer f.Close()
	var dataArray [][]string
	for _, v := range data {
		dataArray = append(dataArray, []string{v.Name, v.Surname, v.Tel, v.LastAccess.Format(CustomTimeFormat)})
	}
	_ = f.Truncate(0)
	writer := csv.NewWriter(f)
	defer writer.Flush()
	err = writer.WriteAll(dataArray)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return err
}

func populateData() error {
	_, err := os.Stat("./data.csv")
	if err != nil {
		fmt.Println("File not found")
		return err
	}
	f, err := os.OpenFile("./data.csv", os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println("Error opening file")
		return err
	}
	defer f.Close()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	for _, v := range records {
		timeRead, err := time.Parse(CustomTimeFormat, v[3])
		if err != nil {
			return err
		}
		user := UserData{
			Name:       v[0],
			Surname:    v[1],
			Tel:        v[2],
			LastAccess: timeRead,
		}
		data[user.Tel] = user
	}
	return nil
}
