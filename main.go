package main

import (
	// "bytes"
	"fmt"
	"strings"

	// "io/ioutil"
	"log"
	"regexp"

	"github.com/jung-kurt/gofpdf"
	"github.com/rsc/pdf"
)

func readPDFFiles(filePath string) {
	// Read PDF files
	pdfFile, err := pdf.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open PDF file: %v", err)
	}
	// fmt.Println("PDF file opened successfully")
	fmt.Println("Reading PDF file: ", filePath)

	fileNameParts := strings.Split(filePath, "-")
	if len(fileNameParts) < 3 {
		log.Fatalf("Invalid file path format, expected MM-DD-YY format: %v", filePath)
	}
	year := fileNameParts[3][:2] // Taking only the YY part of the date for the year

	numPages := pdfFile.NumPage()
	// fmt.Printf("Number of Pages in PDF: %d\n", numPages)

	for i := 1; i <= numPages; i++ {
		page := pdfFile.Page(i)
		if page.V.IsNull() {
			continue
		}

		content := page.Content()
		rows := make(map[float64]string)
		for _, t := range content.Text {
			rows[t.Y] += t.S + " "
		}
		for _, text := range rows {
			text = strings.ReplaceAll(text, " ", "") // Remove all spaces from the text
			if strings.Contains(text, "Mathison") {
				re := regexp.MustCompile(`TYPE:.{0,13}`)
				text = re.ReplaceAllString(text, "")
				payRegEx := regexp.MustCompile(`MathisonProjectMathisonProject(.{8})`)
				payMatches := payRegEx.FindStringSubmatch(text)
				var pay string
				if len(payMatches) > 1 {
					pay = payMatches[1]
					pay = "$" + pay
				}
				date := text[:5]
				pdfg := gofpdf.New("P", "mm", "A4", "") // Create a new PDF. The "P" argument stands for Portrait mode.
				pdfg.AddPage()
				pdfg.SetFont("Arial", "", 12) // SetFont now takes a style parameter, which can be left empty for normal text.
				// AddText is not directly available in gofpdf, you use CellFormat or Text for positioning text
				pdfg.SetFont("Arial", "B", 12) // Set font to bold and size to 12 for "Mathison Projects Inc"
				pdfg.Text(10, 10, "Mathison Projects Inc")
				pdfg.SetFont("Arial", "", 10) // Set font size back to 10 for the rest of the text
				pdfg.Text(10, 15, "8 The Grn Ste R")
				pdfg.Text(10, 20, "Dover, DE 19901")
				pdfg.Text(10, 25, "United States")
				pdfg.Text(10, 30, "jacob@mathisonprojects.com")
				pdfg.SetDrawColor(0, 0, 0) // Set the color for the divider line to black
				pdfg.Line(10, 35, 200, 35) // Draw a line from x1,y1 to x2,y2 as the divider
				pdfg.Text(10, 40, "Employee: Hunter Koenig-Albert")
				pdfg.Text(10, 45, "Employee ID: HKA-2000")
				pdfg.Ln(60) // Move below the text before starting the table

				// gofpdf does not have a direct AddTable function, but you can create tables using CellFormat in a loop for rows and columns
				header := []string{"Employee #", "Pay Date", "Pay"}
				w := []float64{40.0, 35.0, 45.0, 45.0}
				for j, str := range header {
					pdfg.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
				}
				pdfg.Ln(-1) // Move to the next line
				data := [][]string{
					{"HKA-2000", date + "/" + year, pay},
				}
				for _, row := range data {
					for j, datum := range row {
						if j < len(w) { // Check to prevent index out of range error
							pdfg.CellFormat(w[j], 6, datum, "1", 0, "", false, 0, "")
						}
					}
					pdfg.Ln(-1)
				}

				// Finally, save the PDF
				docDate := date
				docDate = strings.Replace(docDate, "/", "-", -1)
				fileName := fmt.Sprintf("Hunter_Koenig-Albert_Paystub_%s_%s.pdf", docDate, year)
				err := pdfg.OutputFileAndClose(fileName)
				if err != nil {
					log.Fatalf("failed to save PDF: %v", err)
				}
				fmt.Println("PDF document generated successfully")
			}
		}
	}
}

func main() {
	filePaths := []string{
		"C:/Users/Doge2/Downloads/Hunter Koenig-Albert Bank Statement 09-30-23.pdf",
		"C:/Users/Doge2/Downloads/Hunter Koenig-Albert Bank Statement 10-31-23.pdf",
		"C:/Users/Doge2/Downloads/Hunter Koenig-Albert Bank Statement 11-30-23.pdf",
		"C:/Users/Doge2/Downloads/Hunter Koenig-Albert Bank Statement 12-31-23.pdf",
		"C:/Users/Doge2/Downloads/Hunter Koenig-Albert Bank Statement 01-31-24.pdf",
		"C:/Users/Doge2/Downloads/Hunter Koenig-Albert Bank Statement 02-29-24.pdf",
	}
	for _, filePath := range filePaths {
		readPDFFiles(filePath)
	}
}
