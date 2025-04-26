import requestEvent from '@/event/request'

import { message, Modal, notification } from 'ant-design-vue'

requestEvent.on('NetworkError', (statusResult: string | boolean) => {
  if (!statusResult) {
    Modal.error({
      title: 'Connection Lost',
      content: 'Please ensure you are properly connected to the internet'
    })
  }
  if (statusResult === 'browser') {
    Modal.error({
      title: 'Unable to Connect to Server',
      content: 'Please ensure your network environment is normal'
    })
  }
  if (statusResult === 'service') {
    Modal.error({
      title: 'Server Error',
      content: 'The server is temporarily experiencing issues, please try again later. If the problem persists, please contact technical support'
    })
  }
})

requestEvent.on('UnknownError', () => {
  notification.error({
    message: 'Backend Server Error',
    description: `Please refresh the page and try again. If the problem persists, please contact technical support`
  })
})

requestEvent.on('Message', (type: 'success' | 'error' | 'warn', msg: string) => {
  if (type == 'success') message.success(msg)
  if (type == 'error') message.error(msg)
  if (type == 'warn') message.warn(msg)
})