package player

import (
	"errors"
	"os/exec"
	"syscall"
)

type Player struct {
	cmd    *exec.Cmd
	done   chan struct{}
	paused bool
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
	p.done = make(chan struct{})
	p.paused = false
	go func() {
		_ = cmd.Wait()
		close(p.done)
	}()
	return nil
}

func (p *Player) TogglePause() (bool, error) {
	if p.cmd == nil || p.cmd.Process == nil {
		return false, errors.New("nothing is playing")
	}
	signal := syscall.SIGSTOP
	if p.paused {
		signal = syscall.SIGCONT
	}
	if err := syscall.Kill(p.cmd.Process.Pid, signal); err != nil {
		return p.paused, err
	}
	p.paused = !p.paused
	return p.paused, nil
}

func (p *Player) Stop() {
	if p.cmd == nil || p.cmd.Process == nil {
		return
	}
	if p.paused {
		_ = syscall.Kill(p.cmd.Process.Pid, syscall.SIGCONT)
	}
	_ = p.cmd.Process.Kill()
	if p.done != nil {
		<-p.done
	}
	p.cmd = nil
	p.done = nil
	p.paused = false
}
