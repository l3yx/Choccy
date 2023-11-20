import request from "../utils/request";

export function getDatabases(page, pageSize,sortBy, sortOrder) {
    return request.get('/database', {
        params: {
            page,
            pageSize,
            sortBy,
            sortOrder
        }
    })
}