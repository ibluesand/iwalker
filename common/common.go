package common

import (
    log "github.com/cihub/seelog"
)

//错误检查
func CheckError(err error, info string) (res bool) {
    if err != nil {
        log.Errorf("[%s] [%s]" , info , err.Error())
        return false
    }
    return true
}
