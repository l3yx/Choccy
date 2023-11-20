import request from "../utils/request";

export function runTaskByID(id) {
    return request.get('/task/run',{
        params: {
            id
        }
    })
}


export function getTasks(page, pageSize,sortBy, sortOrder, filters) {
    return request.get('/task',{
        params: {
            page,
            pageSize,
            sortBy,
            sortOrder,
            filters
        }
    })
}

export function setIsRead(idList, read) {
    return request.post('/task/read',{
        idList:idList,
        read:read
    })
}

export function getTaskUnread() {
    return request.get('/task/unread')
}