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
	dialer := CreatDialer()
	assert.Equal(t, nil, SendEmail(dialer, "89897766@qq.com", "885588", "Your Sign Up code", "register.html"))

}
