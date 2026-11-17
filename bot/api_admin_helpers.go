package bot

import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"

// widgetsJoinChannel builds a JoinChannel button. It is a thin wrapper
// around models.Button so the api_admin file does not import the widgets
// package (which would create a cycle in some configurations).
func widgetsJoinChannel(id, text, username string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeLink,
		ButtonText: text,
		ButtonLink: &models.ButtonLink{
			Type: "joinchannel",
			JoinChannelData: &models.JoinChannelData{
				Username: username,
			},
		},
	}
}

// widgetsLink builds a plain URL link button.
func widgetsLink(id, text, url string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeLink,
		ButtonText: text,
		ButtonLink: &models.ButtonLink{
			Type:    "url",
			LinkURL: url,
		},
	}
}
