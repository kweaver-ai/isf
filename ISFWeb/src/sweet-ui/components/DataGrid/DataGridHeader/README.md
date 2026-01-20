# DataDridHeader

## Usage

```react
const columns = [
            {
                title: 'customizeTitle-name',
                key: 'name',
                width: '30%',
                renderCell: (name, record) => name
            },
            {
                title: 'customizeTitle-age',
                key: 'age',
                width: '70%',
                renderCell: (age, record) => age
            }
        ]
<DataGridHeader
    columns={columns}
    enableSelectAll={true}
    isSelectedAllChecked={true}
    onSelectAllChanged={checked=>{}}
/>
```

## Props

| 参数                                      | 说明                                     | 必须 | 类型     | 默认值 |
| ----------------------------------------- | ---------------------------------------- | ---- | -------- | ------ |
| [columns](#columns)                       | 数据列配置                               | Yes  | array    | /      |
| isSelectedAllChecked                      | 控制全选按钮是否选中                     | No   | boolean  | false  |
| enableSelectAll                           | 是否显示全选按钮                         | No   | boolean  | false  |
| [onSelectAllChanged](#onSelectAllChanged) | 全选按钮勾选状态发生改变时触发的处理函数 | No   | function | /      |

## Reference

### Props

#### <span id="columns">columns</span>

列配置对象数组，每一项列配置对象包含以下属性：

| 属性       | 说明                                                                                             | 必须 | 类型               |
| ---------- | ------------------------------------------------------------------------------------------------ | ---- | ------------------ |
| width      | 列宽度，多个列会按比例分配，建议使用百分比                                                       | No   | string&#124;number |
| key        | 用于从多层级的对象中查找特定属性，会将结果传递给 renderCell(value, record)                       | No   | string             |
| renderCell | 渲染单元格。`value` 会通过 `key` 属性在 `record` 对象中进行查找。`record` 是该行完整的数据对象。 | No   | function           |
| title      | 列标题                                                                                           | No   | string             |

#### <span id="onSelectAllChanged">onSelectAllChanged</span>

全选按钮勾选状态发生改变时触发的处理函数。

| 属性    | 说明             | 类型    |
| ------- | ---------------- | ------- |
| checked | 改变后的选中状态 | boolean |
