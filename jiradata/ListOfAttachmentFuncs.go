package jiradata

func (l *ListOfAttachment) Len() int {
	return len(*l)
}

func (l *ListOfAttachment) Less(i, j int) bool {
	return (*l)[i].ID < (*l)[j].ID
}

func (l *ListOfAttachment) Swap(i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}
