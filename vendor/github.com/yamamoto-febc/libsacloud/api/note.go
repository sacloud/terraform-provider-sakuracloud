package api

type NoteAPI struct {
	*baseAPI
}

func NewNoteAPI(client *Client) *NoteAPI {
	return &NoteAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "note"
			},
		},
	}
}
