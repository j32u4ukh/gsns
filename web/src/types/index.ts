export interface IPost {
    // 文章的唯一識別號
    id: number;
    // 作者的資訊，這裡使用空的物件表示，應該有更具體的型別
    author: IUser;
    // 文章標題
    title: String;
    // 文章描述，使用字串類型
    description: String;
    // 文章的圖片，使用字串陣列表示
    pictures: String[];
    // TODO: 之後可以新增一個影片
    // 貼文的分類，使用字串陣列
    categories: String[];
    // 喜歡貼文的人，使用空的陣列表示，應該有更具體的型別
    likes: IUser[];
    // 分享文章的人，使用空的陣列表示，應該有更具體的型別
    shares: IUser[];
    // 對文章的反應，使用空的陣列表示，應該有更具體的型別
    replys: IPost[];
}

export interface IUser{
    id: number,
    name: String,
}

export interface PostData{
    parentId: number,
    id: number,
    user_id: number,
    content: String,
}