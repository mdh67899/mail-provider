package sender

import (
	"github.com/mdh67899/go-utils/cron"
	List "github.com/mdh67899/go-utils/list"
	"github.com/mdh67899/mail-provider/model"
	"time"
)

type Store struct {
	SafeLinklist *List.Linklist
	Cron         *cron.CronScheduler
}

func NewStore() *Store {
	return &Store{
		SafeLinklist: List.NewLinklist(),
		Cron:         cron.NewCronScheduler(time.Second * 3),
	}
}

func (this *Store) SafePush(v interface{}) {
	this.SafeLinklist.PushFront(v)
}

func (this *Store) Len() int {
	return this.SafeLinklist.Len()
}

func (this *Store) PopBackByNum(num int) []*model.Mail {
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

func (this *Store) PopBack() *model.Mail {
	item := this.SafeLinklist.PopFront()
	return item.(*model.Mail)
}
