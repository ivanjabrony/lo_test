package model

var EmptyFilter Filter = Filter{}

const Created TaskStatus = "created"
const InProgress TaskStatus = "inProgress"
const Done TaskStatus = "done"

type Filter struct {
	Status TaskStatus
}
