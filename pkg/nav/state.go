package nav

import tea "github.com/charmbracelet/bubbletea"

var position = 0

type State struct {
	model    tea.Model
	next     *State
	previous *State
}

func NewState() *State {
	return &State{}
}

func (s *State) setCurrentCommand(m tea.Model) {
	if s.model == nil {
		*s = State{
			model:    m,
			next:     s.next,
			previous: s.previous,
		}
	}
	for range position {
		if s.next != nil {
			s = s.next
		}
	}
	s.model = m
}

func (s *State) NextCommand() tea.Model {
	position++
	for range position {
		if s.next != nil {
			s = s.next
		}
	}
	m := s.model
	return m
}

func (s *State) SetNextCommand(m tea.Model) *State {
	if s.next == nil {
		s.next = &State{
			model:    m,
			next:     nil,
			previous: s,
		}
		return s.next
	}
	return s.next.SetNextCommand(m)
}

func (s *State) SetAndGetNextCommand(m tea.Model) tea.Model {
	s.setCurrentCommand(m)
	return s.NextCommand()
}

func (s *State) PreviousCommand() tea.Model {
	position--
	for range position {
		if s.next != nil {
			s = s.next
		}
	}
	m := s.model
	return m
}

func (s *State) SetPreviousCommand(m tea.Model) *State {
	if s.previous == nil {
		s.previous = &State{
			model:    m,
			next:     s,
			previous: nil,
		}
		return s.previous
	}
	return s.previous.SetPreviousCommand(m)
}

func (s *State) SetAndGetPreviousCommand(m tea.Model) tea.Model {
	s.setCurrentCommand(m)
	return s.PreviousCommand()
}
