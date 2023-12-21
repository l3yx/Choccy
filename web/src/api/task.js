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

export function addTask(database, suites, name) {
    return request.post('/task',{
        database: database,
        suites: suites,
        name: name
    })
}

export function addGithubBatchTasks(query,sort,order,number,offset,language,suites) {
    return request.post('/task/github',{
        query,
        sort,
        order,
        number,
        offset,
        language,
        suites
    })
}

export function getTaskUnread() {
    return request.get('/task/unread')
}

export function getGithubRepositoryQueryTotal(query) {
    return request.get('/task/github/query',{
        params:{
            query: query
        }
    })
}