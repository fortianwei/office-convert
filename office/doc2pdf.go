package office

func ConvertDoc2Pdf(fileInputPath string, fileOutputPath string) {

	openArgs := []interface{}{fileInputPath}

	/// https://docs.microsoft.com/zh-cn/office/vba/api/word.document.exportasfixedformat
	exportArgs := []interface{}{fileOutputPath, 17}

	closeArgs := []interface{}{}

	quitArgs := []interface{}{}

	convertHandler := ConvertHandler{
		FileInPath:      fileInputPath,
		FileOutPath:     fileOutputPath,
		ApplicationName: "Word.Application",
		WorkspaceName:   "Documents",
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
