import request from "../utils/request";
import {deleteProject} from "./project";

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

export function deleteDatabases(name) {
    return request.delete('/database', {
        params: {
            name
        }
    })
}