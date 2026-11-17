// Command hangprofile sets a "dual" (hangi) profile picture —
// پروفایل هنگی — with ONE call: photo1 and photo2 end up on the object
// at the same time (two different pictures, not stacked on each other),
// and every avatar the run did not create is removed so the object ends
// with exactly the pair. All the work (tolerant upload, simultaneous
// dual bind, clean-up) lives inside the library now — see
// parsrubika.Hangi / account.Account.HangiProfile.
//
// Usage:
//
//	go run . --phone +989123456789
//	go run . --phone +989123456789 --guid c0DBedk0525ec07fb54a35b01baf2303
//
// By default it picks up 1.jpg and 2.jpg from the current directory and
// targets the signed-in account itself; pass --guid to hang the pair on
// a group/channel you administer instead. The first run asks for the
// SMS login code (and the 2FA password if enabled); after that the
// session file is reused.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	pr "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account"
)

func main() {
	phone := flag.String("phone", "", "phone number in international format, e.g. +989123456789")
	guid := flag.String("guid", "", "target group/channel GUID (empty = your own account)")
	photo1 := flag.String("photo1", "ParsRubika-v4.0.0 1.png", "first image file (jpg/png/gif/webp)")
	photo2 := flag.String("photo2", "ParsRubika-v4.0.0 2.png", "second image file (jpg/png/gif/webp)")
	rounds := flag.Int("rounds", 1, "dual-bind rounds (1 is enough; more only piles up avatars)")
	clean := flag.Bool("clean", false, "delete all existing avatars first")
	session := flag.String("session", "hangprofile", "session file name (reuses existing login)")
	flag.Parse()

	if *phone == "" {
		fmt.Fprintln(os.Stderr, "error: --phone is required")
		fmt.Fprintln(os.Stderr, "usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	opts := []account.HangiOption{
		account.WithHangiRounds(*rounds),
		account.WithHangiProgress(func(msg string) { fmt.Println(msg) }),
	}
	if *clean {
		opts = append(opts, account.WithHangiClean())
	}

	res, err := pr.Hangi(*session, *phone, *guid, *photo1, *photo2, opts...)
	if err != nil {
		log.Fatalf("hangi profile failed: %v", err)
	}

	fmt.Printf("done: %d avatar(s) on %s", res.FinalCount, res.ObjectGUID)
	for _, id := range res.AvatarIDs {
		fmt.Printf(" (new: %s)", id)
	}
	if res.Deleted > 0 {
		fmt.Printf(", %d leftover(s) removed", res.Deleted)
	}
	fmt.Println()
}
