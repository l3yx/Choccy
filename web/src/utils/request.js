import axios from 'axios'
import { ElMessage } from 'element-plus'
import {getToken} from "./auth";


let baseURL
if(import.meta.env.DEV){
    baseURL = "http://localhost:80/api"
}else {
    baseURL = "/api"
}

const request = axios.create({
    baseURL: baseURL, timeout: 5*60*1000
})

request.interceptors.request.use(function (config) {
    config.headers.set("X-Token",getToken())
    return config;
}, function (error) {
    return Promise.reject(error);
});


request.interceptors.response.use(function (response) {
    let data = response.data
    if (data.err) {
        if (data.err!=="Unauthorized"){
            ElMessage.error(data.err)
        }
        return Promise.reject(data.err);
    }
    return data;
}, function (error) {
    ElMessage.error(error)
    return Promise.reject(error);
});

export default request