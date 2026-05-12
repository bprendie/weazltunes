package player

import (
	"errors"
	"os/exec"
)

type Player struct {
	cmd *exec.Cmd
}

func New() *Player {
	return &Player{}
}

func (p *Player) Play(url string) error {
	p.Stop()
	bin, err := exec.LookPath("mpv")
	if err != nil {
		return errors.New("mpv was not found; install mpv to play streams")
	}
	cmd := exec.Command(bin, "--no-video", "--force-window=no", url)
	if err := cmd.Start(); err != nil {
		return err
	}
	p.cmd = cmd
	go cmd.Wait()
	return nil
}

func (p *Player) Stop() {
	if p.cmd == nil || p.cmd.Process == nil {
		return
	}
	_ = p.cmd.Process.Kill()
	_ = p.cmd.Wait()
	p.cmd = nil
}
