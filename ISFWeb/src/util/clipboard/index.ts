/**
 * 剪贴板写入
 * @param text 要写入剪贴板的文本
 */
export function copy(text: string): boolean {
    if (document.createRange) {
        // 1. 创建一个DOM
        // 2. 创建一个选区
        // 3. 使用选区选中DOM，并执行copy
        // 4. 移除DOM
        const textNode = document.createElement('span');
        const range = document.createRange();
        const selection = window.getSelection();

        // 隐藏DOM
        Object.assign(textNode.style, {
            'white-space': 'pre',
            position: 'fixed',
            top: '-9999px',
            left: '-9999px',
            '-moz-user-select': 'text',
            '-webkit-user-select': 'text',
            '-ms-user-select': 'text',
            'user-select': 'text',
        })

        textNode.innerText = text;

        document.body.appendChild(textNode);

        // 先移除所有选取，防止有文本选中的情况下复制失败
        selection.removeAllRanges();

        range.selectNode(textNode);
        window.getSelection().addRange(range);

        const result = document.execCommand('copy');

        document.body.removeChild(textNode);

        return result;
    } else if (window.clipboardData) {
        window.clipboardData.setData('Text', text);

        return true;
    }
}