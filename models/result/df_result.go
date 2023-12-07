package result

import (
	"sync"

	"gitlab.galaxy123.cloud/base/public/constants/consts"
	"gitlab.galaxy123.cloud/base/public/constants/enums"
)

var (
	defaultErrCode = enums.ErrSysBadRequest
	defaultErr     = consts.ErrSysBadRequest
)

type ErrManger struct {
	sync.RWMutex
	errMap  map[string]map[int]error
	codeMap map[string]int
}

func NewErrManger() *ErrManger {

	errMap := make(map[string]map[int]error)
	codeMap := make(map[string]int)

	zhMap := map[int]error{
		enums.ErrSysBadRequest:   consts.ErrSysBadRequest,
		enums.ErrSysTokenExpired: consts.ErrSysTokenExpired,
		enums.ErrSysAuthFailed:   consts.ErrSysAuthFailed,
		enums.ErrRequestLimit:    consts.ErrRequestLimit,
		enums.ErrTimeout:         consts.ErrTimeout,
		enums.ErrIPLimit:         consts.ErrIPLimit,
		enums.ErrImageSizeLimit:  consts.ErrImageSizeLimit,
		enums.ErrImageSuffix:     consts.ErrImageSuffix,
	}

	errMap[consts.ZH] = zhMap

	for _, m := range errMap {
		for i, err := range m {
			codeMap[err.Error()] = i
		}
	}

	return &ErrManger{
		errMap:  errMap,
		codeMap: codeMap,
	}
}

func (e *ErrManger) GetCode(err error) int {
	e.RLock()
	defer e.RUnlock()

	code, ok := e.codeMap[err.Error()]
	if !ok {
		return defaultErrCode
	}

	return code
}

func (e *ErrManger) GetErr(code int, lang string) error {
	e.Lock()
	defer e.Unlock()

	if lang == "" {
		lang = consts.ZH
	}

	m, ok := e.errMap[lang]
	if !ok {
		return defaultErr
	}

	er, ok := m[code]
	if !ok {
		return defaultErr
	}

	return er
}
