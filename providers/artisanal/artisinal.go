package artisanal

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/euank/api.error.technology/errortech"

	irc "github.com/fluffle/goirc/client"
)

type provider struct{}

func New() provider {
	return provider{}
}

func (provider) Name() string {
	return "artisanal"
}

var badErrErr = errortech.Error{
	Short:    "No artisans available",
	Tags:     []string{"404"},
	Language: "english",
}

func (provider) GetError(lang string, tags []string) errortech.Error {
	nick := fmt.Sprintf("err%v", rand.Int())
	c := irc.SimpleClient(nick)
	output := make(chan errortech.Error)

	c.Config().Server = "irc.wobscale.website"

	c.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			conn.Join("#errors")
			conn.Privmsg("#errors", fmt.Sprintf("Please provide me an error in language %v with tags %v. Respond with '%v: <err>'", lang, tags, nick))
		})

	c.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			txt := line.Text()
			if !strings.HasPrefix(txt, "all: ") && !strings.HasPrefix(txt, nick+": ") {
				return
			}
			parts := strings.SplitN(txt, ": ", 2)
			output <- errortech.Error{
				Short:    parts[1],
				Language: lang,
				Tags:     tags,
			}
			conn.Quit("gone")
		})

	if err := c.Connect(); err != nil {
		return badErrErr
	}

	return <-output

}
