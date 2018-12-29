// go run bottest.go
package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"go.uber.org/zap"

	"github.com/stevenxie/rgv/pkg/bot"
	"github.com/stevenxie/rgv/pkg/reddit"
	ess "github.com/unixpickle/essentials"
)

const (
	target = "uwaterloo"
	rate   = 3 * time.Second
)

func main() {
	fmt.Printf("Creating Reddit bot that watches /r/%s...\n", target)

	// Create logger and bot.Receiver.
	var (
		sp          = &spewer{os.Stdout}
		logger, err = zap.NewDevelopment()
	)
	if err != nil {
		ess.Die("Failed to create Zap logger:", err)
	}

	b, err := bot.New(sp, logger.Sugar())
	if err != nil {
		ess.Die("Failed to create bot:", err)
	}

	if err = b.Run(target, rate); err != nil {
		ess.Die("Error running bot:", err)
	}
	fmt.Println("Running bot...")

	if err = b.Wait(); err != nil {
		ess.Die("Bot errored out:", err)
	}
}

type spewer struct{ out io.Writer }

func (p *spewer) print(a ...interface{}) {
	fmt.Fprint(p.out, a...)
}

func (p *spewer) println(a ...interface{}) {
	fmt.Fprintln(p.out, a...)
}

func (p *spewer) dump(a ...interface{}) {
	spew.Fdump(p.out, a...)
}

func (p *spewer) ReceivePost(post *reddit.Post) error {
	p.print("Received post: ")
	p.dump(post)
	p.println()
	return nil
}

func (p *spewer) ReceiveComment(c *reddit.Comment) error {
	p.print("Received comment: ")
	p.dump(c)
	p.println()
	return nil
}
