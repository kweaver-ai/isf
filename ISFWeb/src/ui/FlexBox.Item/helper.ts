/**
 * 根据关键字获取对齐class
 * @param styles 样式表
 * @param keywords 对齐关键字
 */
export function getAlignCls(styles, keywords = 'left middle'): Array<string> {
    return keywords.split(/\s/).reduce((cls, align) => {
        switch (align) {
            case 'top':
                return cls.concat(styles['alignTop']);

            case 'right':
                return cls.concat(styles['alignRight']);

            case 'bottom':
                return cls.concat(styles['alignBottom']);

            case 'left':
                return cls.concat(styles['alignLeft']);

            case 'middle':
                return cls.concat(styles['alignMiddle']);

            case 'center':
                return cls.concat(styles['alignCenter']);
        }
    }, [])
}