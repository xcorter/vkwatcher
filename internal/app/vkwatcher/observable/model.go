package observable

type ObservableType struct {
	Value int
}

type Observable struct {
	Owner string
	ObservableType
	Value    string
	LastScan int
	ChatId   int64
}

func NewMusicObservable(owner string, value string, chatId int64) Observable {
	return Observable{
		Owner:          owner,
		ObservableType: ObservableType{Value: 1},
		Value:          value,
		ChatId:         chatId,
	}
}
