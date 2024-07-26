package main

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/jung-kurt/gofpdf/v2"
    "github.com/xuri/excelize/v2"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

const (
    host     = "localhost"
    port     = 5433
    user     = "admin"
    password = "mysecretpassword"
    dbname   = "mydatabase"
)

type Ticker struct {
    gorm.Model
    Symbol string `gorm:"unique;not null"`
}

type SP500Data struct {
    gorm.Model
    Date             string
    Open             int
    OpenPrecision    int
    High             int
    HighPrecision    int
    Low              int
    LowPrecision     int
    Close            int
    ClosePrecision   int
    AdjClose         int
    AdjClosePrecision int
    Volume           int
    TickerID         uint
    Ticker           Ticker
}

// convertToIntWithPrecision converts a string to an int and returns the int and its precision
func convertToIntWithPrecision(value string) (int, int) {
    floatValue, err := strconv.ParseFloat(value, 64)
    if err != nil {
        fmt.Println("Error converting string to float:", err)
        return 0, 0
    }

    parts := strings.Split(value, ".")
    precision := 0
    if len(parts) == 2 {
        precision = len(parts[1])
    }

    factor := float64(1)
    for i := 0; i < precision; i++ {
        factor *= 10
    }

    intValue := int(floatValue * factor)
    return intValue, precision
}

func convertToExcel(db *gorm.DB, exportFileName string) {
    var data []SP500Data
    result := db.Preload("Ticker").Find(&data)
    if result.Error != nil {
        fmt.Println("Error querying database:", result.Error)
        return
    }

    f := excelize.NewFile()
    sheetName := "Sheet1"
    f.NewSheet(sheetName)

    headers := []string{"Date", "Open", "Open Precision", "High", "High Precision", "Low", "Low Precision", "Close", "Close Precision", "Adj_Close", "AdjClose Precision", "Volume", "Ticker"}
    for i, header := range headers {
        cell := fmt.Sprintf("%s1", string(rune('A'+i)))
        f.SetCellValue(sheetName, cell, header)
    }

    for i, record := range data {
        f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), record.Date)
        f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), record.Open)
        f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), record.OpenPrecision)
        f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), record.High)
        f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), record.HighPrecision)
        f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), record.Low)
        f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), record.LowPrecision)
        f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), record.Close)
        f.SetCellValue(sheetName, fmt.Sprintf("I%d", i+2), record.ClosePrecision)
        f.SetCellValue(sheetName, fmt.Sprintf("J%d", i+2), record.AdjClose)
        f.SetCellValue(sheetName, fmt.Sprintf("K%d", i+2), record.AdjClosePrecision)
        f.SetCellValue(sheetName, fmt.Sprintf("L%d", i+2), record.Volume)
        f.SetCellValue(sheetName, fmt.Sprintf("M%d", i+2), record.Ticker.Symbol)
    }

    if err := f.SaveAs(exportFileName); err != nil {
        fmt.Println("Error saving Excel file:", err)
        return
    }

    fmt.Println("Data exported to Excel successfully!")
}

func convertToPDF(db *gorm.DB, exportFileName string) {
    var data []SP500Data
    result := db.Preload("Ticker").Find(&data)
    if result.Error != nil {
        fmt.Println("Error querying database:", result.Error)
        return
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 12)

    headers := []string{"Date", "Open", "Open Precision", "High", "High Precision", "Low", "Low Precision", "Close", "Close Precision", "Adj_Close", "AdjClose Precision", "Volume", "Ticker"}
    for _, header := range headers {
        pdf.CellFormat(20, 10, header, "1", 0, "C", false, 0, "")
    }
    pdf.Ln(-1)

    pdf.SetFont("Arial", "", 12)
    for _, record := range data {
        pdf.CellFormat(20, 10, record.Date, "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.Open), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.OpenPrecision), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.High), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.HighPrecision), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.Low), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.LowPrecision), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.Close), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.ClosePrecision), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.AdjClose), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.AdjClosePrecision), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, strconv.Itoa(record.Volume), "1", 0, "C", false, 0, "")
        pdf.CellFormat(20, 10, record.Ticker.Symbol, "1", 0, "C", false, 0, "")
        pdf.Ln(-1)
    }

    if err := pdf.OutputFileAndClose(exportFileName); err != nil {
        fmt.Println("Error saving PDF file:", err)
        return
    }

    fmt.Println("Data exported to PDF successfully!")
}


func main() {
    fileName := "sp500_data.xlsx"

    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return
    }
    defer func() {
        sqlDB, err := db.DB()
        if err != nil {
            fmt.Println("Error getting underlying *sql.DB:", err)
            return
        }
        sqlDB.Close()
    }()

    db.AutoMigrate(&SP500Data{}, &Ticker{})

    f, err := excelize.OpenFile(fileName)
    if err != nil {
        fmt.Println("Error opening Excel file:", err)
        return
    }
    defer f.Close()

    sheetName := f.GetSheetName(0)
    rowsIterator, err := f.Rows(sheetName)
    if err != nil {
        fmt.Println("Error reading Excel rows:", err)
        return
    }

    _ = rowsIterator.Next() // Skip header row

    tickerMap := make(map[string]uint)
    var sp500Data []SP500Data

    for rowsIterator.Next() {
        row, err := rowsIterator.Columns()
        if err != nil {
            fmt.Println("Error reading row:", err)
            return
        }

        tickerSymbol := row[7]
        tickerID, ok := tickerMap[tickerSymbol]
        if !ok {
            var ticker Ticker
            result := db.Where("symbol = ?", tickerSymbol).First(&ticker)
            if result.Error != nil {
                if result.Error == gorm.ErrRecordNotFound {
                    newTicker := Ticker{Symbol: tickerSymbol}
                    db.Create(&newTicker)
                    tickerID = newTicker.ID
                } else {
                    fmt.Println("Error finding or creating ticker:", result.Error)
                    return
                }
            } else {
                tickerID = ticker.ID
            }
            tickerMap[tickerSymbol] = tickerID
        }

        // Convert string to float and then to int with precision
        open, openPrecision := convertToIntWithPrecision(row[1])
        high, highPrecision := convertToIntWithPrecision(row[2])
        low, lowPrecision := convertToIntWithPrecision(row[3])
        close, closePrecision := convertToIntWithPrecision(row[4])
        adjClose, adjClosePrecision := convertToIntWithPrecision(row[5])
        volume, _ := strconv.Atoi(row[6])

        sp500Data = append(sp500Data, SP500Data{
            Date:              row[0],
            Open:              open,
            OpenPrecision:     openPrecision,
            High:              high,
            HighPrecision:     highPrecision,
            Low:               low,
            LowPrecision:      lowPrecision,
            Close:             close,
            ClosePrecision:    closePrecision,
            AdjClose:          adjClose,
            AdjClosePrecision: adjClosePrecision,
            Volume:            volume,
            TickerID:          tickerID,
        })

        if len(sp500Data) >= 50 {
            result := db.Create(&sp500Data)
            if result.Error != nil {
                fmt.Println("Error inserting batch of rows:", result.Error)
                return
            }
            sp500Data = nil
        }
    }

    if len(sp500Data) > 0 {
        result := db.Create(&sp500Data)
        if result.Error != nil {
            fmt.Println("Error inserting remaining rows:", result.Error)
            return
        }
    }

    if err := rowsIterator.Close(); err != nil {
        fmt.Println("Error closing rows iterator:", err)
        return
    }

    fmt.Println("Data inserted successfully!")

	var action int


	for action != 3 {
		fmt.Println("Do you want to export to Excel or PDF?")
		fmt.Println("1. Excel")
		fmt.Println("2. PDF")
		fmt.Println("3. None")
		fmt.Scanf("%d", &action)

		switch action {
		case 1:
			convertToExcel(db, "sp500_data_export.xlsx")
		case 2:
			convertToPDF(db, "sp500_data_export.pdf")
		case 3:
            fmt.Println("No export selected.")
        default:
            fmt.Println("Invalid option. Please select 1, 2, or 3.")
	}
}
	
	
}

