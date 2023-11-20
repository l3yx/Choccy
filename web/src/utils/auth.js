import Cookies from 'js-cookie';

export function setToken(token){
    Cookies.set('token', token, { expires: 30 });
}

export function getToken(){
    return Cookies.get('token');
}