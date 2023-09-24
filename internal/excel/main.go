package excel

import (
	"fmt"
	"os"
	"strconv"

	// "time"

	"github.com/xuri/excelize/v2"
)

const INPUT_FILE = "./src/input/modelo_importacao_lote.xlsx"
const LOG_FILE = "./src/log/log.xlsx"
const SHEET_NAME = "Solicitação"

type Log struct {
	OS string
	Contact string
	Date string
}

type Info struct {
	Link string
	OS string
}

func CreateExcelFile() (log *excelize.File) {
	if _, err := os.Stat(LOG_FILE); err == nil {
		log, err = excelize.OpenFile(LOG_FILE)
		if err != nil {
			fmt.Println("Erro pra abrir o arquivo?")
			fmt.Println(err)
		}
		return log

	} else {
		fmt.Println("Arquivo nao encontrado, criando um novo")
		log = excelize.NewFile()
		defer func() {
			if err = log.Close(); err != nil {
				panic("Couldn't close log file")
			}
		}()
		
			log.SetCellValue("Sheet1", "A1", "Ordem de Serviço")
			log.SetCellValue("Sheet1", "B1", "Contato com Cliente")
			log.SetCellValue("Sheet1", "C1", "Data")

			if err := log.SaveAs(LOG_FILE); err != nil {
				fmt.Println(err)
			}
			return log
	}
}

func WriteLog(file *excelize.File, log Log) bool {

	logSize := GetLogSize()
	
	if err := file.SetCellStr("Sheet1", fmt.Sprintf("A%d", logSize), log.OS); err != nil {
		panic("Couldn't set OS column")
	}
	if err := file.SetCellStr("Sheet1", fmt.Sprintf("B%d", logSize), log.Contact); err != nil {
		panic("Error while writing file")
	}
	if err := file.SetCellStr("Sheet1", fmt.Sprintf("C%d", logSize), log.Date); err != nil {
		panic("Error while writing file")
	}

	if err := file.SaveAs(LOG_FILE); err != nil {
		fmt.Println("Couldn't save log")
		return false
	}
	fmt.Println("Log saved")
	return true
}

func GetLogSize() int {
	f, err := excelize.OpenFile(LOG_FILE)
    if err != nil {
        panic(err)
    }
	
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println("Erro achando o tamanho do log")
		panic(err)
	}
	
	if err = f.Close(); err != nil {
			fmt.Println(err)
	}

	fmt.Println(len(rows))
	
	return len(rows) + 1
		
}

func GetLinks() []Info {
	var res []Info

	rows := getSheet()
	
	for i, row := range rows {
		if i != 0 {
			name := row[4]
			phone := row[12]
			address := row[10]
			os := row[0]

			msg := strconv.Quote(fmt.Sprintf("Olá, %s. Tudo bem? Queria confirmar se o seu endereço realmente é %s", name, address))

			link := fmt.Sprintf("https://web.whatsapp.com/send?phone=%s&text=%s", phone, msg)
			info := Info{link, os}
			res = append(res, info)
		}
	}
	return res
}


func getSheet() [][]string {
	if f, err := excelize.OpenFile(INPUT_FILE); err != nil {
		panic(err)
	} else {
		defer func() {
			// Close the spreadsheet.
			if err := f.Close(); err != nil {
					fmt.Println(err)
			}
		}()
		if rows, err := f.GetRows(SHEET_NAME); err != nil {
			panic(err)
		} else {
			return rows
		}
	}
}