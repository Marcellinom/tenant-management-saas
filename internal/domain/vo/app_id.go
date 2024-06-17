package vo

type AppId int

func (a AppId) Value() int {
	return int(a)
}
