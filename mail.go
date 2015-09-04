package main

import (
	"fmt"
	"os/exec"
)

/**
 * simple mail func w/sendmail https://www.snip2code.com/Snippet/172752/golang-mail%28%29--because-porting-awful-PHP
 */
func mail(to, from, subject, message string) error {
	var err error
	cmd := exec.Command("/usr/sbin/sendmail", "-t", "-i")
	pipe, err := cmd.StdinPipe()

	if err != nil {
		return err
	}

	if from == "" {
		from = "MidDiff"
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	_, err = fmt.Fprintf(pipe, "To: %s\n", to)
	_, err = fmt.Fprintf(pipe, "From: %s\n", from)
	_, err = fmt.Fprintf(pipe, "Subject: %s\n", subject)
	_, err = fmt.Fprintf(pipe, "\n%s\n", message)

	err = pipe.Close()
	if err != nil {
		return err
	}

	return cmd.Wait()
}
