export function arrayTreeFilter<T>(
    data: T[],
    filterFn: (item: T, level: number) => boolean,
    options: { childrenKeyName: string } = { childrenKeyName: 'children' },
) {
    let children = data || [];
    let result: T[] = [];
    let level = 0;
    do {
        let foundItem: T = children.filter((item) => {
            return filterFn(item, level);
        })[0];

        if (!foundItem) {
            break;
        }
        result = [...result, foundItem]
        children = (foundItem as any)[options.childrenKeyName] || [];
        level += 1;
    } while (children.length > 0);

    return result;
}