package office

import (
	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	log "github.com/sirupsen/logrus"
)

/// 更多内容请参考官方COM文档 https://docs.microsoft.com/zh-cn/office/vba/api/word.application
type Operation struct {
	OpType    string
	Arguments []interface{}
}

/// 部分应用不允许隐藏 ，比如ppt，所以Visible需要设定下
type ConvertHandler struct {
	FileInPath      string
	FileOutPath     string
	ApplicationName string
	WorkspaceName   string
	Visible         bool
	DisplayAlerts   int
	OpenFileOp      Operation
	ExportOp        Operation
	CloseOp         Operation
	QuitOp          Operation
}

type DomConvertObject struct {
	Application *ole.IDispatch
	Workspace   *ole.IDispatch
	SingleFile  *ole.IDispatch
}

func (handler ConvertHandler) Convert() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	log.Println("handle open start")
	dom := handler.Open()
	log.Println("handle open end")
	log.Println("handler in file path is " + handler.FileInPath)
	log.Println("handler out file path is " + handler.FileOutPath)

	defer dom.Application.Release()
	defer dom.Workspace.Release()
	defer dom.SingleFile.Release()

	handler.Export(dom)
	log.Println("handle export end")

	handler.Close(dom)
	log.Println("handle close end")

	handler.Quit(dom)
	log.Println("handle quit end")

}
func (handler ConvertHandler) Open() DomConvertObject {
	var dom DomConvertObject
	unknown, err := oleutil.CreateObject(handler.ApplicationName)
	if err != nil {
		panic(err)
	}
	dom.Application = unknown.MustQueryInterface(ole.IID_IDispatch)

	oleutil.MustPutProperty(dom.Application, "Visible", handler.Visible)
	oleutil.MustPutProperty(dom.Application, "DisplayAlerts", handler.DisplayAlerts)

	dom.Workspace = oleutil.MustGetProperty(dom.Application, handler.WorkspaceName).ToIDispatch()

	dom.SingleFile = oleutil.MustCallMethod(dom.Workspace, handler.OpenFileOp.OpType, handler.OpenFileOp.Arguments...).ToIDispatch()
	return dom
}

func (handler ConvertHandler) Export(dom DomConvertObject) {
	oleutil.MustCallMethod(dom.SingleFile, handler.ExportOp.OpType, handler.ExportOp.Arguments...)

}

func (handler ConvertHandler) Close(dom DomConvertObject) {
	if handler.ApplicationName == "PowerPoint.Application" {
		oleutil.MustCallMethod(dom.SingleFile, handler.CloseOp.OpType, handler.CloseOp.Arguments...)
	} else {
		oleutil.MustCallMethod(dom.Workspace, handler.CloseOp.OpType, handler.CloseOp.Arguments...)
	}
}

func (handler ConvertHandler) Quit(dom DomConvertObject) {
	oleutil.MustCallMethod(dom.Application, handler.QuitOp.OpType, handler.QuitOp.Arguments...)
}
