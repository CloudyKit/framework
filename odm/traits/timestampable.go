package traits

import (
	"time"
)

type TimerAble struct {
	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

func (receiver *TimerAble) setTimeOnCreation() {
	now := time.Now()
	if receiver.CreatedAt.IsZero() {
		receiver.CreatedAt = now
		receiver.UpdatedAt = now
	} else {
		receiver.UpdatedAt = now
	}
}

type TimerAbleTrait struct {
}

func (trait TimerAbleTrait) OnDelete(filter interface{}) interface{} {
	return filter
}

func (trait TimerAbleTrait) OnFilter(filter interface{}) interface{} {
	return filter
}

func (trait TimerAbleTrait) OnInsert(doc interface{}) interface{} {
	if doc, is := doc.(interface{ setTimeOnCreation() }); is {
		doc.setTimeOnCreation()
	}
	return doc
}

func (trait TimerAbleTrait) OnUpdateOrReplace(filter, doc interface{}) interface{} {
	if doc, is := doc.(interface{ setTimeOnCreation() }); is {
		doc.setTimeOnCreation()
	}
	return doc
}
