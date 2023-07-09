/*
 * @Author: hiLin 123456
 * @Date: 2021-11-08 17:18:56
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-15 22:08:26
 * @FilePath: /arclinks-go/library/apollo_sdk/apollo_logger.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package apollo_sdk

import (
	"fmt"

	"arclinks-go/library/log"
)

type ApolloLogger struct {
	Logger log.LoggerContextIface
}

// Debug logs debug msg.
func (a *ApolloLogger) Debug(v ...interface{}) {
	a.Logger.Debug(fmt.Sprint(v...))
}

// Warn logs info msg.
func (a *ApolloLogger) Info(v ...interface{}) {
	a.Logger.Info(fmt.Sprint(v...))
}

// Warn logs warn msg.
func (a *ApolloLogger) Warn(v ...interface{}) error {
	a.Logger.Warn(fmt.Sprint(v...))
	return nil
}

// Error logs error msg.
func (a *ApolloLogger) Error(v ...interface{}) error {
	a.Logger.Error(fmt.Sprint(v...))
	return nil
}

// Debugf logs debuf msg.
func (a *ApolloLogger) Debugf(format string, params ...interface{}) {
	a.Logger.Debugf(format, params...)
}

// Infof logs info msg.
func (a *ApolloLogger) Infof(format string, params ...interface{}) {
	a.Logger.Infof(format, params...)
}

// Warnf logs warn msg.
func (a *ApolloLogger) Warnf(format string, params ...interface{}) error {
	a.Logger.Warnf(format, params...)
	return nil
}

// Errorf logs error msg.
func (a *ApolloLogger) Errorf(format string, params ...interface{}) error {
	a.Logger.Errorf(format, params...)
	return nil
}
