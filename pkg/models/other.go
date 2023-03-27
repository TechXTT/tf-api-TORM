package models

type VoteRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	SoftwareID  uint   `json:"software_id"`
	NetworksID  uint   `json:"networks_id"`
	EmbeddedID  uint   `json:"embedded_id"`
	BattleBotID uint   `json:"battlebot_id"`
}

type VerifyVoteRequest struct {
	Token string `json:"token"`
}
