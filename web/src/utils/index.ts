// 定義一個函數 extractHashtags，用於從輸入的字符串中提取出所有的主題標簽（hashtags）
export function extractHashtags(str: string): string[] {

    // 使用正則表達式 rgx 來匹配字符串中的所有主題標簽，同時進行全局、不區分大小寫的匹配
    const rgx = /#(\w+)\b/gi;

    // 創建一個空數組 result 來存放提取出的主題標簽
    const result: string[] = [];
    let temp: RegExpExecArray | null;

    // 使用 while 循環進行正則匹配，將匹配的結果存入 temp 中
    while ((temp = rgx.exec(str)) !== null) {

        // 將 temp 中的第一個捕獲組（即主題標簽）加入到 result 數組中
        result.push(temp[1]);
    }
    return result;
}

export function formatNumber(number: number): string {
    if (number < 1000) {
      return `${number}`;
    } else if (number < 1000000) {
      const rounded = (number / 1000).toFixed(2);
      return `${parseFloat(rounded)}K`;
    } else {
      const rounded = (number / 1000000).toFixed(2);
      return `${parseFloat(rounded)}M`;
    }
  }