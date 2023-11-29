import request from "../utils/request";

export function getResults(page, pageSize,sortBy, sortOrder,filters) {
    return request.get('/result', {
        params: {
            page,
            pageSize,
            sortBy,
            sortOrder,
            filters
        }
    })
}

export function getResultSarif(fileName) {
    return request.get('/result/sarif', {
        params: {
            fileName
        }
    })
}

export function getResultSarifCodeFlows(fileName, id) {
    return request.get('/result/sarif/flows', {
        params: {
            fileName,
            id
        }
    })
}

export function deleteResult(ID) {
    return request.delete('/result', {
        params: {
            ID
        }
    })
}

export function setIsRead(idList, read) {
    return request.post('/result/read',{
        idList:idList,
        read:read
    })
}

export function getResultUnread() {
    return request.get('/result/unread')
}