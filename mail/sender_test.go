package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thetherington/simplebank/util"
)

func TestSendEmailWithMailHog(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewMailHogSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword, config.EmailServerHost, config.EmailServerPort)

	subject := "A test email"
	content := `
		<h1>Hello world</h1>
		<p>This is a test message from <a href="http://hetheringtons.org">Hetheringtons</a></p>
	`

	to := []string{"thomas@hetheringtons.org"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
