package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ozym/dmc"

	"github.com/tbruyelle/hipchat-go/hipchat"
)

var (
	verbose bool
	room    string
	token   string
	base    string
)

func chat(msg, colour string, notify bool) error {

	if token == "" || room == "" {
		return nil
	}

	r := &hipchat.NotificationRequest{
		Message: msg,
		Color:   colour,
		Notify:  notify,
	}

	if _, err := hipchat.NewClient(token).Room.Notification(room, r); err != nil {
		return err
	}

	return nil
}

func store(d dmc.Device, s *dmc.State) error {

	// couldn't find a model ...
	if _, ok := s.Values["model"]; !ok {
		return nil
	}

	// cleanup fdqn ...
	n := strings.Split(d.Name, ".")
	if !(len(n) > 0) {
		return nil
	}
	p := strings.Split(n[0], "-")
	if !(len(p) > 0) {
		return nil
	}

	// per site directory
	b := base + "/" + p[len(p)-1]
	if err := os.MkdirAll(b, 0755); err != nil {
		return err
	}

	// output file name
	f := b + "/" + n[0] + ".json"
	file, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(s.Marshal())
	if err != nil {
		return err
	}

	_, err = file.Write([]byte{'\n'})

	return err
}

func notify(d dmc.Device, s *dmc.State) (bool, error) {

	// couldn't find a model ...
	if _, ok := s.Values["model"]; !ok {
		return false, nil
	}

	// hasn't really changed ...
	if d.Model == s.Values["model"].(string) {
		return true, nil
	}

	// cleanup fdqn ...
	n := strings.Split(d.Name, ".")
	if !(len(n) > 0) {
		return false, nil
	}

	// outgoing message ...
	items := []interface{}{d.Model, n[0], d.IP.String(), s.Values["model"].(string)}
	parts := []string{
		"---{{{... %s ...}}}---",
		"<b><em>%s</em></b> [%s]",
		"has been identified to be a",
		"<b>%s</b>",
	}

	if verbose {
		log.Printf("[%s] %s [%s] has been identified to be a %s\n", items...)
	}

	// hipchat ...
	msg := fmt.Sprintf(strings.Join(parts, "<br/>"), items...)
	if err := chat(msg, "red", true); err != nil {
		return true, err
	}

	return true, nil
}

func device(d dmc.Device, s *dmc.State) bool {

	err := store(d, s)
	if err != nil {
		log.Fatal(err)
	}

	ok, err := notify(d, s)
	if err != nil {
		log.Println(err)
	}

	return ok
}

func main() {

	flag.BoolVar(&verbose, "verbose", false, "make noise")
	flag.StringVar(&base, "base", ".", "base status storage directory")

	flag.StringVar(&room, "room", os.Getenv("HIPCHAT_ROOM_NAME"), "hipchat room name")
	flag.StringVar(&token, "token", os.Getenv("HIPCHAT_ROOM_TOKEN"), "hipchat room token")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Manage equipment identification and monitoring\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] <command> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  check  -- check for new equipment tagged as uninstalled\n")
		fmt.Fprintf(os.Stderr, "  status -- check existing equipment for state changes\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Use: \"%s <command> --help\" for more information about a specific command\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	args := flag.Args()
	if !(len(args) > 0) {
		flag.Usage()

		log.Fatalf("Missing command")
	}

	switch args[0] {
	case "check":
		check(args[1:])
	case "status":
		status(args[1:])
	default:
		flag.Usage()

		log.Fatalf("Unknown command: %s", args[0])
	}

}
