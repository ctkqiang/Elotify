export const PostPushNotifications = {
  pushToAll: '/push/all',
  pushToUser: '/push/user/{user_id}',
  pushToSegMent: '/push/group/{segment_id}',
  pushToTopic: '/push/topic/{topic_id}',
}

export const GetPushNotifications = {
    getHealth: '/health',
    getLatestNotification: '/notification/latest',
    getUserNotifications: '/notification/user/{user_id}',
    getSegmentNotifications: '/notification/group/{segment_id}',
    getTopicNotifications: '/notification/topic/{topic_id}',
}

export const PutPushNotifications = {
    updateNotificationById: '/notification/{notification_id}',
}

export const DeletePushNotifications = {
    deleteNotificationById: '/notification/{notification_id}',
}
