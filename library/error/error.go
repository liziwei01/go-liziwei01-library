package error

type (
	ErrMsg string
	ErrNo  int
)

type Err interface {
	ErrMsg() ErrMsg
	ErrNo() ErrNo
}

type err struct {
	errmsg ErrMsg
	errno  ErrNo
}

func New(errno ErrNo, errmsg ErrMsg) Err {
	return &err{
		errmsg: errmsg,
		errno:  errno,
	}
}

func NewDefault() Err {
	return &err{
		errmsg: "success",
		errno:  0,
	}
}

func (e *err) ErrMsg() ErrMsg {
	return e.errmsg
}

func (e *err) ErrNo() ErrNo {
	return e.errno
}

func Append(e Err, msg ErrMsg) Err {
	return &err{
		errmsg: e.ErrMsg() + msg,
		errno:  e.ErrNo(),
	}
}

func (e *err) Marshal(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data":   data,
		"errmsg": e.errmsg,
		"errno":  e.errno,
	}
}

var _ Err = (*err)(nil)
