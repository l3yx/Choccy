import request from '../utils/request'

export function getSetting() {
    return request.get('/setting')
}

export function saveSetting(setting) {
    return request.post('/setting', setting)
}

export function testSetting(key, value){
    return request.post('/setting/test',{
        key: key,
        value: value
    });
}