package model

type Channel string

const (
	Email Channel = "email"
	SMS   Channel = "sms"
	Push  Channel = "push"
)
