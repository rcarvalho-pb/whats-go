package excel

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

const INPUT_FILE = "src/input/modelo_importacao_lote.xlsx"
const LOG_FILE = "src/log/log.xlsx"
const SHEET_NAME = "Solicitação"

type Log struct {
	OS string
	Contact string
}

type Info struct {
	Link string
	OS string
}

func CreateExcelFile() (log *excelize.File) {
	if _, err := os.Stat(fmt.Sprintf("%s-%s", LOG_FILE, time.Now().Format(time.DateTime))); err != nil {
		if os.Remove(LOG_FILE) != nil {
			fmt.Println("Couldn't dele log file")
		}
	}

	log = excelize.NewFile()
	defer func() {
		if err := log.Close(); err != nil {
			panic("Couldn't close log file")
		}
	}()

	log.SetCellValue("Sheet1", "A1", "Ordem de Serviço")
	log.SetCellValue("Sheet1", "B1", "Contato com Cliente")

	return log
}

func WriteLog(logs []Log, file *excelize.File) bool {
	for i, log := range logs {
		if err := file.SetCellStr("Sheet1", fmt.Sprintf("A%d", i+2), log.OS); err != nil {
			panic("Couldn't set OS column")
		}
		if err := file.SetCellStr("Sheet1", fmt.Sprintf("B%d", i+2), log.Contact); err != nil {
			panic("Error while writing file")
		}
	}
	return true
}

func SaveLog(logs []Log) (bool) {
	log := CreateExcelFile()
	if WriteLog(logs, log) {
		if err := log.SaveAs(fmt.Sprintf("%s-%s", LOG_FILE, time.Now().Format(time.DateTime))); err != nil {
			panic(err)
		}
	}
	return true
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
		if rows, err := f.GetRows(SHEET_NAME); err != nil {
			panic(err)
		} else {
			return rows
		}
	}
}