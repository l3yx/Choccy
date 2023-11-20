import request from "../utils/request";

export function getNotifications() {
    return request.get('/notifications')
}
