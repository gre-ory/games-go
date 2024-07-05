package model

type GameStatus int

const (
	GameStatus_JoinableNotStartable GameStatus = iota
	GameStatus_JoinableAndStartable
	GameStatus_NotJoinableAndStartable
	GameStatus_Started
	GameStatus_Stopped
	GameStatus_MarkedForDeletion
)

func (s GameStatus) IsJoinable() bool {
	switch s {
	case GameStatus_JoinableNotStartable,
		GameStatus_JoinableAndStartable:
		return true
	default:
		return false
	}
}

func (s GameStatus) IsNotJoinable() bool {
	return !s.IsJoinable()
}

func (s GameStatus) IsStartable() bool {
	switch s {
	case GameStatus_JoinableAndStartable,
		GameStatus_NotJoinableAndStartable:
		return true
	default:
		return false
	}
}

func (s GameStatus) IsNotStartable() bool {
	return !s.IsStartable()
}

func (s GameStatus) IsStarted() bool {
	return s == GameStatus_Started
}

func (s GameStatus) IsStopped() bool {
	return s == GameStatus_Stopped
}

func (s GameStatus) IsMarkedForDeletion() bool {
	return s == GameStatus_MarkedForDeletion
}

func (s GameStatus) CanJoin() error {
	switch s {
	case GameStatus_NotJoinableAndStartable:
		return ErrGameNotJoinable
	case GameStatus_Started:
		return ErrGameAlreadyStarted
	case GameStatus_Stopped:
		return ErrGameStopped
	case GameStatus_MarkedForDeletion:
		return ErrGameMarkedForDeletion
	default:
		return nil
	}
}

func (s GameStatus) CanStart() error {
	switch s {
	case GameStatus_JoinableNotStartable:
		return ErrGameNotStartable
	case GameStatus_Started:
		return ErrGameAlreadyStarted
	case GameStatus_Stopped:
		return ErrGameStopped
	case GameStatus_MarkedForDeletion:
		return ErrGameMarkedForDeletion
	default:
		return nil
	}
}

func (s GameStatus) CanStop() error {
	switch s {
	case GameStatus_JoinableNotStartable,
		GameStatus_JoinableAndStartable,
		GameStatus_NotJoinableAndStartable:
		return ErrGameNotStarted
	case GameStatus_Stopped:
		return ErrGameStopped
	case GameStatus_MarkedForDeletion:
		return ErrGameMarkedForDeletion
	default:
		return nil
	}
}

func (s GameStatus) CanPlay() error {
	switch s {
	case GameStatus_JoinableNotStartable,
		GameStatus_JoinableAndStartable,
		GameStatus_NotJoinableAndStartable:
		return ErrGameNotStarted
	case GameStatus_Stopped:
		return ErrGameStopped
	case GameStatus_MarkedForDeletion:
		return ErrGameMarkedForDeletion
	default:
		return nil
	}
}

func (s GameStatus) CanLeave() error {
	switch s {
	case GameStatus_MarkedForDeletion:
		return ErrGameMarkedForDeletion
	default:
		return nil
	}
}

func (s GameStatus) CanDelete() error {
	switch s {
	case GameStatus_Started:
		return ErrGameNotStopped
	case GameStatus_MarkedForDeletion:
		return ErrGameMarkedForDeletion
	default:
		return nil
	}
}

func (s GameStatus) IsValid() bool {
	switch s {
	case GameStatus_JoinableNotStartable,
		GameStatus_JoinableAndStartable,
		GameStatus_NotJoinableAndStartable,
		GameStatus_Started,
		GameStatus_Stopped:
		return true
	default:
		return false
	}
}

func (s GameStatus) String() string {
	switch s {
	case GameStatus_JoinableNotStartable:
		return "joinable"
	case GameStatus_JoinableAndStartable:
		return "joinable-startable"
	case GameStatus_NotJoinableAndStartable:
		return "not-joinable-startable"
	case GameStatus_Started:
		return "started"
	case GameStatus_Stopped:
		return "stopped"
	default:
		return ""
	}
}

func (s GameStatus) Labels() []string {
	var labels []string
	labels = append(labels, "game-status")
	if s.IsValid() {
		labels = append(labels, s.String())
	}
	return labels
}
