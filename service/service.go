package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"paste.org.cn/paste/proto"
	"paste.org.cn/paste/util"
)

type Paste struct{}

func (p *Paste) PostPaste(c *gin.Context) {
	var (
		_, log = util.EnsureWithLogger(c)
		req    proto.PostPasteReq
	)

	err := c.BindJSON(&req)
	if err != nil {
		log.Errorf("BindJSON failed: %+v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}

func (p *Paste) PostPasteOnce(c *gin.Context) {

}

func (p *Paste) GetPaste(c *gin.Context) {

}
