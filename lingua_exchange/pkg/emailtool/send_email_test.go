package emailtool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"lingua_exchange/configs"
	"lingua_exchange/internal/config"
)

func TestSendEmail(t *testing.T) {
	var configFile = configs.Path("lingua_exchange.yml")
	err := config.Init(configFile)
	if err != nil {
		panic("init config error: " + err.Error())
	}
	// templatePath, err := strutil.GetTemplatePath("register.html")
	assert.Equal(t, nil, SendEmail("yep1895@gmail.com", "885588", "Your Sign Up code", "register.html"))

}
