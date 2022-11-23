package office

func ConvertXsl2Pdf(fileInputPath string, fileOutputPath string) {

	openArgs := []interface{}{fileInputPath}

	exportArgs := []interface{}{0, fileOutputPath}

	closeArgs := []interface{}{}

	quitArgs := []interface{}{}

	convertHandler := ConvertHandler{
		FileInPath:      fileInputPath,
		FileOutPath:     fileOutputPath,
		ApplicationName: "Excel.Application",
		WorkspaceName:   "Workbooks",
		Visible:         false,
		DisplayAlerts:   0,
		OpenFileOp: Operation{
			OpType:    "Open",
			Arguments: openArgs,
		},
		ExportOp: Operation{
			OpType:    "ExportAsFixedFormat",
			Arguments: exportArgs,
		},
		CloseOp: Operation{

			OpType:    "Close",
			Arguments: closeArgs,
		},
		QuitOp: Operation{

			OpType:    "Quit",
			Arguments: quitArgs,
		},
	}
	convertHandler.Convert()
}
