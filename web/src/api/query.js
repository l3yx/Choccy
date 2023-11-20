import request from "../utils/request";

export function getQueries(path) {
    return request.get('/query',{
        params: {
            path
        }
    })
}

export function getQueryContent(path) {
    return request.get('/query/content',{
        params: {
            path
        }
    })
}