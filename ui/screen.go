package ui

import "github.com/hajimehoshi/ebiten/v2"

type Screen interface {
	Update() error
	Draw(*ebiten.Image)
}

type Manager struct {
	current Screen
}

func NewManager(start Screen) *Manager {
	return &Manager{
		current: start,
	}
}

func (m *Manager) Set(screen Screen) {
	m.current = screen
}

func (m *Manager) Update() error {
	if m.current == nil {
		return nil
	}
	return m.current.Update()
}

func (m *Manager) Draw(screen *ebiten.Image) {
	if m.current != nil {
		m.current.Draw(screen)
	}
}