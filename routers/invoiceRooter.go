package routers

import (
	"github.com/gin-gonic/gin"
	"study/controllers"
)

func InvoiceRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/invoices", controllers.GetInvoices())
	incomingRouters.GET("/invoices/:invoice_id", controllers.GetInvoice())
	incomingRouters.POST("/invoices/", controllers.CreateInvoice())
	incomingRouters.POST("/invoices/:invoice_id", controllers.UpdateInvoice())
}
