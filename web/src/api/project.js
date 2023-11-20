import request from '../utils/request'

export function getProjects(page, pageSize,sortBy, sortOrder) {
    return request.get('/project', {
        params: {
            page,
            pageSize,
            sortBy,
            sortOrder
        }
    })
}

export function saveProject(project) {
    return request.post('/project', project)
}

export function deleteProject(ID) {
    return request.delete('/project', {
        params: {
            ID
        }
    })
}