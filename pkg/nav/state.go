package nav

import tea "github.com/charmbracelet/bubbletea"

var position = 0

type state struct {
	model    tea.Model
	next     *state
	previous *state
}

func NewState() *state {
	return &state{}
}

func (s *state) setCurrentCommand(m tea.Model) {
	if s.model == nil {
		*s = state{
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

func (s *state) NextCommand() tea.Model {
	position += 1
	for range position {
		if s.next != nil {
			s = s.next
		}
	}
	m := s.model
	return m
}

func (s *state) SetNextCommand(m tea.Model) *state {
	if s.next == nil {
		s.next = &state{
			model:    m,
			next:     nil,
			previous: s,
		}
		return s.next
	} else {
		return s.next.SetNextCommand(m)
	}
}

func (s *state) SetAndGetNextCommand(m tea.Model) tea.Model {
	s.setCurrentCommand(m)
	return s.NextCommand()
}

func (s *state) PreviousCommand() tea.Model {
	position -= 1
	for range position {
		if s.next != nil {
			s = s.next
		}
	}
	m := s.model
	return m
}

func (s *state) SetPreviousCommand(m tea.Model) *state {
	if s.previous == nil {
		s.previous = &state{
			model:    m,
			next:     s,
			previous: nil,
		}
		return s.previous
	} else {
		return s.previous.SetPreviousCommand(m)
	}
}

func (s *state) SetAndGetPreviousCommand(m tea.Model) tea.Model {
	s.setCurrentCommand(m)
	return s.PreviousCommand()
}
