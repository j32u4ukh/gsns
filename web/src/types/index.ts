export interface IPost {
    // 文章的唯一識別號
    id: number;
    // 文章標題（這裡應該是拼寫錯誤，正確應為 title）
    tittle: null;
    // 文章描述，使用字串類型
    description: string;
    // 作者的資訊，這裡使用空的物件表示，應該有更具體的型別
    author: {};
    // 作者的唯一識別號
    authorId: number;
    // 分享文章的人，使用空的陣列表示，應該有更具體的型別
    sharedBy: never[];
    // 文章的分類，使用字串陣列
    categories: string[];
    // 對文章的反應，使用空的陣列表示，應該有更具體的型別
    reactions: never[];
    // 文章的圖片，使用字串陣列表示
    pictures: string[];
}