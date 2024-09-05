package AchivmentNames

type AchivmentName string

var names []string

func InitNames() []string {
	names = append(names, "First blood", "Double kill", "Triple kill")
	return names
}
