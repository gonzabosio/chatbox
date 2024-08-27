import { axiosInstance } from "./axiosInstance";

//load chats
const getChats = () => axiosInstance({
    method: 'get',
    url: `/chat/load/${localStorage.getItem('user-id')}`,
    headers: {
        Authorization: 'Bearer ' + localStorage.getItem('access-token')
    }
}).then(res => {
    console.log(res.data.chats)
    console.log(res.data.message)
    return res.data.chats
}).catch(err => {
    console.log(err.message)
    console.log(err.error)
    return []
})

export default getChats
