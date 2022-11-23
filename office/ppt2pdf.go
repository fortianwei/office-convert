package office

func ConvertPPT2Pdf(fileInputPath string, fileOutputPath string) {

	openArgs := []interface{}{fileInputPath, -1, 0, 0}

	exportArgs := []interface{}{fileOutputPath, 32}

	closeArgs := []interface{}{}

	quitArgs := []interface{}{}

	convertHandler := ConvertHandler{
		FileInPath:      fileInputPath,
		FileOutPath:     fileOutputPath,
		ApplicationName: "PowerPoint.Application",
		WorkspaceName:   "Presentations",
		Visible:         true,
		DisplayAlerts:   0,
		OpenFileOp: Operation{
			OpType:    "Open",
			Arguments: openArgs,
		},
		ExportOp: Operation{
			OpType:    "SaveAs",
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
