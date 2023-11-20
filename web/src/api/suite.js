import request from "../utils/request";

export function saveSuiteContent(name, content) {
    return request.post('/suite/content',{
        name,
        content
    })
}

export function getSuites(page, pageSize,sortBy, sortOrder) {
    return request.get('/suite', {
        params: {
            page,
            pageSize,
            sortBy,
            sortOrder
        }
    })
}

export function getSuiteContent(name) {
    return request.get('/suite/content',{
        params: {
            name
        }
    })
}

export function resolveSuite(path) {
    return request.get('/suite/resolve',{
        params: {
            path
        }
    })
}

export function deleteSuite(name) {
    return request.delete('/suite', {
        params: {
            name
        }
    })
}


export function createSuite(name) {
    return request.post('/suite',null,{
        params:{
            name:name
        }
    })
}

export function renameSuite(oldName, newName) {
    return request.post('/suite/rename',null,{
        params:{
            oldName,
            newName
        }
    })
}