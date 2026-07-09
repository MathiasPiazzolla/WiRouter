package main

import (
	"fmt"

	"github.com/MathiasPiazzolla/WiRouter/model"

	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println(windows.ComputerName())
	fmt.Print(model.GetAllAdapters())
	//fmt.Print(net.Interfaces())
	//fmt.Print(windows.GetAdaptersAddresses())
	list, err := model.ListWindowsAdapters()

	if err != nil {
		panic(err)
	}

	for _, a := range list {

		fmt.Println("----------------")
		fmt.Println("Nome:", a.Name)
		fmt.Println("Descrizione:", a.Description)
		fmt.Println("MAC:", a.MAC)
		fmt.Println("Tipo:", a.Type)

	}
}
