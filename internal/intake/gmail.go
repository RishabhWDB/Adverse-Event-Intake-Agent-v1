package intake

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
)

func FetchUnseenEmails() ([]string, error) {
	address := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_APP_PASSWORD")

	c, err := imapclient.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer c.Close()

	if err := c.Login(address, password).Wait(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	_, err = c.Select("INBOX", nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to select inbox: %w", err)
	}

	searchData, err := c.Search(&imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
	}, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(searchData.AllSeqNums()) == 0 {
		log.Println("No unseen emails found")
		return []string{}, nil
	}

	var bodies []string
	seqSet := imap.SeqSetNum(searchData.AllSeqNums()...)

	fetchOptions := &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{{}},
	}

	fetchCmd := c.Fetch(seqSet, fetchOptions)
	defer fetchCmd.Close()

	for {
		msg := fetchCmd.Next()
		if msg == nil {
			break
		}
		for {
			item := msg.Next()
			if item == nil {
				break
			}
			if bodySection, ok := item.(imapclient.FetchItemDataBodySection); ok {
				buf := make([]byte, 0)
				chunk := make([]byte, 4096)
				for {
					n, err := bodySection.Literal.Read(chunk)
					if n > 0 {
						buf = append(buf, chunk[:n]...)
					}
					if err != nil {
						break
					}
				}
				bodies = append(bodies, string(buf))
			}
		}
	}

	return bodies, nil
}

func ExtractPlainText(rawEmail string) (string, error) {
	reader := strings.NewReader(rawEmail)
	mr, err := mail.CreateReader(reader)
	if err != nil {
		return "", fmt.Errorf("failed to create mail reader: %w", err)
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			if strings.HasPrefix(contentType, "text/plain") {
				body, err := io.ReadAll(part.Body)
				if err != nil {
					return "", err
				}
				return string(body), nil
			}
		}
	}

	return "", fmt.Errorf("no plain text part found")
}
