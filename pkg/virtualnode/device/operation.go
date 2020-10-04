package device

import (
	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya/pkg/virtualnode/connectivity"
	"arhat.dev/pkg/log"
	"encoding/base64"
	"fmt"
	"io"
)

func (m *Manager) Operate(deviceName, operationName string, data []byte, out io.Writer) error {
	logger := m.Log.WithFields(
		log.String("device", deviceName),
		log.String("operation", operationName),
	)
	msgCh, _, err := m.ConnectivityManager.PostCmd(
		0,
		aranyagopb.CMD_DEVICE_OPERATE,
		aranyagopb.NewDeviceOperateCmd(deviceName, operationName, data),
	)

	if err != nil {
		logger.I("failed to post device ensure cmd", log.Error(err))
		return fmt.Errorf("failed to post device operate cmd")
	}

	connectivity.HandleMessages(msgCh, func(msg *aranyagopb.Msg) (exit bool) {
		if msgErr := msg.GetError(); msgErr != nil {
			err = msgErr
			logger.I("device operation failed", log.Error(msgErr))
			return true
		}

		dor := msg.GetDeviceOperationResult()
		if dor == nil {
			logger.I("unexpected non device operation result msg")
			return true
		}

		for _, data := range dor.Data {
			_, err = fmt.Fprintln(out, base64.StdEncoding.EncodeToString(data))
			if err != nil {
				logger.I("failed to write operation result", log.Error(err))
			}
		}

		return false
	}, nil, connectivity.HandleUnknownMessage(m.Log))
	return nil
}