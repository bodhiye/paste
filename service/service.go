package service

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"paste.org.cn/paste/db"
	"paste.org.cn/paste/proto"
	"paste.org.cn/paste/util"
)

type Paste struct {
	db.Paste
}

func (p *Paste) PostPaste(c *gin.Context) {
	var (
		ctx, log = util.EnsureWithLogger(c)
		req      proto.PostPasteReq
		length   int
	)

	err := c.BindJSON(&req)
	if err != nil {
		log.Errorf("BindJSON failed: %+v", err)
		c.JSON(http.StatusBadRequest, InvalidArgs)
		return
	}

	length = utf8.RuneCountInString(req.Content)
	if length > 10000 {
		log.Errorf("Content is too long: %d", length)
		c.JSON(http.StatusBadRequest, TooManyContent)
		return
	}

	entry := db.PasteEntry{
		Langtype: req.Langtype,
		Content:  req.Content,
		Password: util.String2md5(req.Password),
	}

	key, err := p.Paste.Set(ctx, entry)
	if err != nil {
		log.Errorf("Failed to insert entry into database: %+v", err)
		c.JSON(http.StatusBadRequest, PasteFailed)
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key})
}

func (p *Paste) PostPasteOnce(c *gin.Context) {

}

func (p *Paste) GetPaste(c *gin.Context) {

}
