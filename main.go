package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"log"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

func main() {
	cli := &xrpc.Client{
		Host: "https://bsky.social",
	}

	input := &atproto.ServerCreateSession_Input{
		Identifier: os.Getenv("HANDDLE"),
		Password:   os.Getenv("PASSWORD"),
	}

	output, err := atproto.ServerCreateSession(context.Background(), cli, input)
	if err != nil {
		log.Fatal(err)
	}

	cli.Auth = &xrpc.AuthInfo{
		AccessJwt:  output.AccessJwt,
		RefreshJwt: output.RefreshJwt,
		Handle:     output.Handle,
		Did:        output.Did,
	}

	post := &bsky.FeedPost{
		LexiconTypeID: "app.bsky.feed.post",
		Text:          fmt.Sprintf("現在時刻は %vです\n", time.Now().Format("15:04")),
		CreatedAt:     time.Now().Local().Format(time.RFC3339),
	}

	feed := &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       cli.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: post,
		},
	}

	resp, err := atproto.RepoCreateRecord(context.Background(), cli, feed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
