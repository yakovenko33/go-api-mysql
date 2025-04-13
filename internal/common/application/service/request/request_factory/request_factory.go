package request_factory

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetLanguage(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")

	//("en-US,en;q=0.9,ru;q=0.8")
	langs := strings.Split(lang, ",")
	if len(langs) > 0 {
		lang = strings.Split(langs[0], ";")[0]
	}

	if lang == "" {
		lang = "en"
	}
	return lang
}
