package models

// Ballot is a struct holding tezos ballot information
type Ballot struct {
	PKH          string
	Ballot       string
	ProposalHash string
	Rolls        int64
	IsTesting    bool
	Period       int

	// General statistic
	Quorum     float64
	TotalRolls float64
	Yay        int64
	Nay        int64
	Pass       int64
}

func (b *Ballot) PercentParticipation() float64 {
	if b.TotalRolls == 0 {
		return 0
	}

	return (b.Participations() / b.TotalRolls) * 100
}

func (b *Ballot) CountingPercentYay() float64 {
	if b.CountingParticipations() == 0 {
		return 0
	}

	return (float64(b.Yay) / b.CountingParticipations()) * 100
}

func (b *Ballot) CountingPercentNay() float64 {
	if b.CountingParticipations() == 0 {
		return 0
	}

	return (float64(b.Nay) / b.CountingParticipations()) * 100
}

func (b *Ballot) PercentYay() float64 {
	if b.Participations() == 0 {
		return 0
	}

	return (float64(b.Yay) / b.Participations()) * 100
}

func (b *Ballot) PercentNay() float64 {
	if b.Participations() == 0 {
		return 0
	}

	return (float64(b.Nay) / b.Participations()) * 100
}

func (b *Ballot) PercentPass() float64 {
	if b.Participations() == 0 {
		return 0
	}

	return (float64(b.Pass) / b.Participations()) * 100
}

func (b *Ballot) Participations() float64 {
	return float64(b.Yay + b.Nay + b.Pass)
}

func (b *Ballot) CountingParticipations() float64 {
	return float64(b.Yay + b.Nay)
}

func (b *Ballot) PercentTowardQuorum() float64 {
	return b.Quorum - b.PercentParticipation()
}

func (b *Ballot) Phase() string {
	if b.IsTesting {
		return "exploration"
	}
	return "promotion"
}
