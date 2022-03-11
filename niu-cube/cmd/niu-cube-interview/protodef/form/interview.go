package form

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiniu/x/xlog"
	"github.com/solutions/niu-cube/cmd/niu-cube-interview/common"
	"github.com/solutions/niu-cube/cmd/niu-cube-interview/protodef"
	"github.com/solutions/niu-cube/internal/protodef/model"
	"regexp"
	"time"
)

var (
	defaultLogger = xlog.New("form validator")
)

type InterviewUpdateForm = InterviewCreateForm

type InterviewCreateForm struct {
	Title            string `form:"title"`
	StartTime        int64  `form:"startTime"`
	EndTime          int64  `form:"endTime"`
	Goverment        string `form:"goverment"`
	Career           string `form:"career"`
	IsAuth           bool   `form:"isAuth"`
	AuthCode         string `form:"authCode"`
	IsRecorded       bool   `form:"isRecorded"`
	CandidateName    string `form:"candidateName"`
	CandidatePhone   string `form:"candidatePhone"`
	InterviewerName  string `form:"interviewerName"`
	InterviewerPhone string `form:"interviewerPhone"`
}

var (
	RegPhone          = regexp.MustCompile("1[3-9][0-9]{9}")
	ErrPhoneCollision = fmt.Errorf("面试官手机号不可与面试者相同")
)

const (
	ErrPhoneMsg = "手机号校验失败"
	ErrTimeMsg  = "时间需至少大于当前时间"
	ErrTitleMsg = "标题过长"
)

func PhoneValidate(phone string) validation.Rule {
	if common.IsFixedPhone(phone) {
		return validation.Skip
	} else {
		return validation.Match(RegPhone).Error(ErrPhoneMsg)
	}
}

func (i *InterviewCreateForm) Validate() error {
	if i.CandidatePhone == i.InterviewerPhone {
		return ErrPhoneCollision
	}
	err := validation.ValidateStruct(i,
		validation.Field(&i.Title, validation.Required, validation.Length(0, 100).Error(ErrTitleMsg)),
		validation.Field(&i.InterviewerPhone, validation.Required, PhoneValidate(i.InterviewerPhone)),
		validation.Field(&i.CandidatePhone, validation.Required, PhoneValidate(i.InterviewerPhone)),
		validation.Field(&i.StartTime, validation.Required, validation.Min(time.Now().Unix()).Error(ErrTimeMsg)),
		validation.Field(&i.EndTime, validation.Required, validation.Min(time.Now().Unix()).Error(ErrTimeMsg)),
	)
	return err
}

func (i *InterviewCreateForm) FillDefault(c *gin.Context) {
	val, ok := c.Get(protodef.ContextUserKey)
	user := val.(model.AccountDo)
	if !ok {
		defaultLogger.Infof("error get user from context")
		return
	}
	if i.InterviewerName == "" {
		i.InterviewerName = user.Nickname
	}
	if i.InterviewerPhone == "" {
		i.InterviewerPhone = user.Phone
	}
	return
}

func (i *InterviewCreateForm) Map() map[string]interface{} {
	var res map[string]interface{}
	val, _ := json.Marshal(i)
	_ = json.Unmarshal(val, &res)
	return res
}
