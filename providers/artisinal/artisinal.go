package artisinal

import (
	"fmt"
	"math/rand"

	"github.com/euank/api.error.technology/errortech"

	irc "github.com/fluffle/goirc/client"
)

type provider struct{}

func New() provider {
	return provider{}
}

func (provider) Name() string {
	return "artisinal"
}

var badErrErr = errortech.Error{
	Short:    "No artisans available",
	Tags:     []string{"404"},
	Language: "english",
}

func (provider) GetError(lang string, tags []string) errortech.Error {
	c := irc.SimpleClient(fmt.Sprintf("err%v", rand.Int()))
	output := make(chan errortech.Error)

	c.Config().Server = "irc.wobscale.website"

	c.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			conn.Join("#errors")
			conn.Privmsg("#errors", fmt.Sprintf("Please provide me an error in language %v with tags %v", lang, tags))
		})

	c.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			output <- errortech.Error{
				Short:    line.Text(),
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
