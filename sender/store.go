package sender

import (
	List "github.com/mdh67899/go-utils/list"
	"github.com/mdh67899/mail-provider/model"
)

type safeLinklist struct {
	SafeLinklist *List.Linklist
}

func NewsafeLinklist() *safeLinklist {
	return &safeLinklist{SafeLinklist: List.NewLinklist()}
}

var Queue = NewsafeLinklist()

func (this *safeLinklist) SafePush(v interface{}) {
	this.SafeLinklist.PushFront(v)
}

func (this *safeLinklist) Len() int {
	return this.SafeLinklist.Len()
}

func (this *safeLinklist) PopBackByNum(num int) []*model.Mail {
	mails := this.SafeLinklist.BatchPopBack(num)
	length := len(mails)

	if length == 0 {
		return make([]*model.Mail, 0)
	}

	mail := make([]*model.Mail, length)

	for i := 0; i < length; i++ {
		mail[i] = mails[i].(*model.Mail)
	}
	return mail
}

func (this *safeLinklist) PopBack() *model.Mail {
	item := this.SafeLinklist.PopFront()
	return item.(*model.Mail)
}
